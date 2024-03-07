package service

import (
	"context"
	"todo-grpc/models"
)

type AlertService interface {
	CreateAlert(ctx context.Context, todo *models.Alert) (*models.Alert, error)
}
