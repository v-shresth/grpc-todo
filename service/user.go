package service

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"todo-grpc/models"
	"todo-grpc/pb"
)

type UserService interface {
	Register(ctx context.Context, user *models.User) (*pb.RegisterResponse, error)
	Login(ctx context.Context, user *models.User) (string, error)
	AddTodoIdToUser(ctx context.Context, userId, todoId primitive.ObjectID) error
}
