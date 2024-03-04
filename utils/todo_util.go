package utils

import (
	"errors"
	"fmt"
	"github.com/go-logr/logr"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"slices"
	"strings"
	"todo-grpc/models"
	"todo-grpc/pb"
)

var todoUpdatableFields = []string{
	"name", "description", "status", "priority",
}

func ValidateCreateTodoReq(req *pb.CreateTodoReq) error {
	if req == nil {
		return errors.New("req not present")
	}

	if strings.TrimSpace(req.Todo.Name) == "" {
		return errors.New("name can't be empty")
	}

	return nil
}

func ValidateUpdateTodoFieldMask(fm *fieldmaskpb.FieldMask) ([]string, error) {
	if fm == nil || len(fm.Paths) == 0 {
		return todoUpdatableFields, nil
	}

	for _, path := range fm.Paths {
		if !slices.Contains(todoUpdatableFields, path) {
			return nil, fmt.Errorf(
				"invalid field mask: only support updating the following fields: %s, found: %s",
				strings.Join(todoUpdatableFields, ","),
				fm.String(),
			)
		}
	}

	return fm.Paths, nil
}

func ConvertApiTodoDbToto(apiTodo *pb.Todo) (*models.Todo, error) {
	if apiTodo == nil {
		return nil, errors.New("todo not present")
	}

	dbTodo := &models.Todo{
		Name:        strings.TrimSpace(apiTodo.Name),
		Description: strings.TrimSpace(apiTodo.Description),
		Status:      apiTodo.Status,
		DeadLine:    primitive.NewDateTimeFromTime(apiTodo.Deadline.AsTime()),
		Priority:    apiTodo.Priority.String(),
	}

	if apiTodo.Id != "" {
		todoId, err := primitive.ObjectIDFromHex(apiTodo.Id)
		if err != nil {
			return nil, err
		}
		dbTodo.ID = todoId
	}

	if apiTodo.UserId != "" {
		userId, err := primitive.ObjectIDFromHex(apiTodo.UserId)
		if err != nil {
			return nil, err
		}
		dbTodo.UserID = userId
	}
	return dbTodo, nil
}

func ConvertDbTodoApiToto(dbTodo *models.Todo) *pb.Todo {
	if dbTodo == nil {
		return nil
	}

	apiTodo := &pb.Todo{
		Id:          dbTodo.ID.Hex(),
		UserId:      dbTodo.UserID.Hex(),
		Name:        dbTodo.Name,
		Description: dbTodo.Description,
		Status:      dbTodo.Status,
		Priority:    pb.Todo_Priority(pb.Todo_Priority_value[dbTodo.Priority]),
	}

	if dbTodo.CreateTime != 0 {
		apiTodo.CreatedAt = timestamppb.New(dbTodo.CreateTime.Time())
	}
	if dbTodo.DeadLine != 0 {
		apiTodo.Deadline = timestamppb.New(dbTodo.DeadLine.Time())
	}
	if dbTodo.UpdateTime != 0 {
		apiTodo.UpdatedAt = timestamppb.New(dbTodo.UpdateTime.Time())
	}
	return apiTodo
}

func ParseListTodoReq(req *pb.ListTodoReq, logger *logr.Logger) (*models.ListTodoFilter, error) {
	if req.GetLimit() < 1 || req.GetLimit() > 20 {
		if req.GetLimit() > 20 {
			return nil, CreateStatusErrorFromError(errors.New("limit cannot exceed 20"), logger)
		}
		req.Limit = 20
	}
	if req.GetPage() == 0 {
		req.Page = 1
	}
	return &models.ListTodoFilter{
		Limit: req.GetLimit(),
		Page:  req.GetPage(),
	}, nil
}
