package service

import (
	"context"
	"todo-grpc/models"
)

type TodoService interface {
	CreateTodo(ctx context.Context, todo *models.Todo) (*models.Todo, error)
	ListTodos(ctx context.Context, userId string, filter *models.ListTodoFilter) (*models.ListTodoRes, error)
	FetchTodo(ctx context.Context, todoId, userId string) (*models.Todo, error)
	UpdateTodo(ctx context.Context, todo *models.Todo, fieldMasks []string) (*models.Todo, error)
	FetchTodosWithNearbyDeadline(ctx context.Context) ([]models.Todo, error)
}
