package alerts

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"todo-grpc/models"
	"todo-grpc/service"
	"todo-grpc/utils"
)

type serviceClient struct {
	alertsRepo alertsRepo
	logger     *utils.Logger
}

func NewAlertService(
	db *mongo.Client,
	logger *utils.Logger,
	config utils.EnvConfig,
) service.AlertService {
	return &serviceClient{
		alertsRepo: newRepoClient(db, logger),
		logger:     logger,
	}
}

func (s *serviceClient) CreateAlert(ctx context.Context, alert *models.Alert) (*models.Alert, error) {
	alert.ID = primitive.NewObjectID()
	_, err := s.alertsRepo.insertAlert(ctx, alert)
	if err != nil {
		return nil, err
	}

	return alert, nil
}
