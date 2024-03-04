package todo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"golang.org/x/sync/errgroup"
	"slices"
	"time"
	"todo-grpc/models"
	"todo-grpc/service"
	"todo-grpc/service/user"
	"todo-grpc/utils"
)

type serviceClient struct {
	todoRepo    todosRepo
	logger      *utils.Logger
	userService service.UserService
}

func NewTodoService(
	db *mongo.Client,
	logger *utils.Logger,
	config *utils.EnvConfig,
) service.TodoService {
	return &serviceClient{
		todoRepo:    newRepoClient(db, logger),
		logger:      logger,
		userService: user.NewUserService(db, logger, config),
	}
}

func (s *serviceClient) CreateTodo(ctx context.Context, todo *models.Todo) (*models.Todo, error) {
	session, err := s.todoRepo.startSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	// ACID database transactions
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
	callback := func(ctx mongo.SessionContext) (interface{}, error) {
		todo.ID = primitive.NewObjectID()
		todo.CreateTime = primitive.NewDateTimeFromTime(time.Now())
		todoId, err := s.todoRepo.insertTodo(ctx, todo)
		if err != nil {
			return nil, err
		}
		todo.ID = todoId
		err = s.userService.AddTodoIdToUser(ctx, todo.UserID, todoId)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	_, err = session.WithTransaction(ctx, callback, txnOpts)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *serviceClient) ListTodos(
	ctx context.Context, userId string, filter *models.ListTodoFilter,
) (*models.ListTodoRes, error) {
	var todoRes models.ListTodoRes
	userID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	erg, _ := errgroup.WithContext(ctx)
	erg.Go(
		func() error {
			var todoErr error
			todoRes.Todos, todoErr = s.todoRepo.fetchTodos(ctx, userID, filter)
			return todoErr
		},
	)

	erg.Go(
		func() error {
			var countErr error
			todoRes.Count, countErr = s.todoRepo.countTodos(ctx, userID)
			return countErr
		},
	)
	if err = erg.Wait(); err != nil {
		return nil, err
	}

	return &todoRes, err
}

func (s *serviceClient) FetchTodo(ctx context.Context, todoId, userId string) (*models.Todo, error) {
	todoID, err := primitive.ObjectIDFromHex(todoId)
	if err != nil {
		customErr := &utils.ReqInvalidArgumentError{
			GeneralError: &utils.GeneralError{
				DevInfo: err.Error(),
				Msg:     "invalid todo id",
			},
		}
		return nil, customErr
	}
	userID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		customErr := &utils.ReqInvalidArgumentError{
			GeneralError: &utils.GeneralError{
				DevInfo: err.Error(),
				Msg:     "invalid user id",
			},
		}
		return nil, customErr
	}

	todo, err := s.todoRepo.fetchTodo(ctx, todoID, userID)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *serviceClient) UpdateTodo(ctx context.Context, todo *models.Todo, fieldMasks []string) (*models.Todo, error) {
	update := bson.M{}
	if slices.Contains(fieldMasks, "name") {
		if todo.Name == "" {
			customErr := &utils.ReqInvalidArgumentError{
				GeneralError: &utils.GeneralError{
					Msg: "invalid request: contains field mask name but name is empty",
				},
			}
			return nil, customErr
		}
		update["name"] = todo.Name
	}
	if slices.Contains(fieldMasks, "description") {
		update["description"] = todo.Description
	}
	if slices.Contains(fieldMasks, "status") {
		update["status"] = todo.Status
	}
	if slices.Contains(fieldMasks, "priority") {
		update["priority"] = todo.Priority
	}
	if slices.Contains(fieldMasks, "deadline") {
		update["deadline"] = todo.DeadLine
	}
	update["update_time"] = primitive.NewDateTimeFromTime(time.Now())
	updateTodo, err := s.todoRepo.updateTodo(ctx, todo.ID, update)
	if err != nil {
		return nil, err
	}

	return updateTodo, nil
}

func (s *serviceClient) FetchTodosWithNearbyDeadline(
	ctx context.Context,
) ([]models.Todo, error) {
	todos, err := s.todoRepo.fetchTodosWithNearbyDeadline(ctx)
	if err != nil {
		return nil, err
	}

	return todos, nil
}
