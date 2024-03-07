package todo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"todo-grpc/models"
	"todo-grpc/utils"
)

type repoClient struct {
	db     *mongo.Client
	todoC  *mongo.Collection
	logger *utils.Logger
}

type todosRepo interface {
	startSession() (mongo.Session, error)
	insertTodo(
		ctx context.Context, todo *models.Todo,
	) (primitive.ObjectID, error)
	fetchTodos(
		ctx context.Context, userId primitive.ObjectID, limitFilter *models.ListTodoFilter,
	) ([]models.Todo, error)
	countTodos(
		ctx context.Context, userId primitive.ObjectID,
	) (int64, error)
	fetchTodo(
		ctx context.Context, todoId, userId primitive.ObjectID,
	) (*models.Todo, error)
	updateTodo(
		ctx context.Context, taskId primitive.ObjectID, update bson.M,
	) (*models.Todo, error)
	fetchTodosWithNearbyDeadline(
		ctx context.Context,
	) ([]models.Todo, error)
}

func newRepoClient(
	db *mongo.Client, logger *utils.Logger,
) todosRepo {
	return &repoClient{
		db:     db,
		todoC:  utils.GetCollection(db, "todos"),
		logger: logger,
	}
}

func (r *repoClient) insertTodo(ctx context.Context, todo *models.Todo) (primitive.ObjectID, error) {
	insertedResp, err := r.todoC.InsertOne(
		ctx,
		todo,
	)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return insertedResp.InsertedID.(primitive.ObjectID), nil
}

func (r *repoClient) startSession() (mongo.Session, error) {
	return r.db.StartSession()
}

func (r *repoClient) fetchTodos(
	ctx context.Context, userId primitive.ObjectID, limitFilter *models.ListTodoFilter,
) ([]models.Todo, error) {
	filter := bson.M{
		"user_id": userId,
	}

	sortOpns := bson.M{
		"create_time": -1,
	}
	opns := options.Find().SetSort(sortOpns).SetLimit(int64(limitFilter.Limit)).SetSkip(int64(limitFilter.Page))

	cursor, err := r.todoC.Find(ctx, filter, opns)
	if err != nil {
		return nil, err
	}
	todos := make([]models.Todo, 0)
	if err = cursor.All(ctx, &todos); err != nil {
		return todos, err
	}

	return todos, nil
}

func (r *repoClient) countTodos(
	ctx context.Context, userId primitive.ObjectID,
) (int64, error) {
	filter := bson.M{
		"user_id": userId,
	}

	count, err := r.todoC.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *repoClient) fetchTodo(
	ctx context.Context, todoId, userId primitive.ObjectID,
) (*models.Todo, error) {
	filter := bson.M{
		"_id":     todoId,
		"user_id": userId,
	}

	var todo models.Todo
	err := r.todoC.FindOne(ctx, filter).Decode(&todo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			customErr := &utils.ReqInvalidArgumentError{
				GeneralError: &utils.GeneralError{
					DevInfo: err.Error(),
					Msg:     "todo not found",
				},
			}
			return nil, customErr
		}
		return nil, &utils.DBInternalError{
			GeneralError: &utils.GeneralError{
				DevInfo: err.Error(),
				Msg:     "unable to fetch todo",
			},
		}
	}

	return &todo, nil
}

func (r *repoClient) updateTodo(
	ctx context.Context, taskId primitive.ObjectID, update bson.M,
) (*models.Todo, error) {
	filter := bson.M{
		"_id": taskId,
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var todo models.Todo
	err := r.todoC.FindOneAndUpdate(ctx, filter, update, opts).Decode(&todo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			customErr := &utils.ReqInvalidArgumentError{
				GeneralError: &utils.GeneralError{
					DevInfo: err.Error(),
					Msg:     "todo not found",
				},
			}
			return nil, customErr
		}
		return nil, &utils.ReqInvalidArgumentError{
			GeneralError: &utils.GeneralError{
				DevInfo: err.Error(),
				Msg:     "unable to update todo",
			},
		}
	}

	return &todo, nil
}

func (r *repoClient) fetchTodosWithNearbyDeadline(
	ctx context.Context,
) ([]models.Todo, error) {
	var todos []models.Todo
	filter := bson.M{
		"deadline": bson.M{
			"$gte": primitive.NewDateTimeFromTime(time.Now().Add(-6 * time.Hour)),
		},
	}

	cursor, err := r.todoC.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &todos)
	if err != nil {
		return nil, err
	}

	return todos, nil
}
