package alerts

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"todo-grpc/models"
	"todo-grpc/utils"
)

type repoClient struct {
	db      *mongo.Client
	alertsC *mongo.Collection
	logger  *utils.Logger
}

type alertsRepo interface {
	insertAlert(ctx context.Context, alert *models.Alert) (primitive.ObjectID, error)
}

func newRepoClient(
	db *mongo.Client, logger *utils.Logger,
) alertsRepo {
	return &repoClient{
		db:      db,
		alertsC: utils.GetCollection(db, "alerts"),
		logger:  logger,
	}
}

func (r *repoClient) insertAlert(ctx context.Context, alert *models.Alert) (primitive.ObjectID, error) {
	insertedResp, err := r.alertsC.InsertOne(
		ctx,
		alert,
	)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return insertedResp.InsertedID.(primitive.ObjectID), nil
}
