package user

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"todo-grpc/models"
	"todo-grpc/utils"
)

type repoClient struct {
	db     *mongo.Client
	usersC *mongo.Collection
	logger *utils.Logger
}

type userRepo interface {
	insertUser(ctx context.Context, user *models.User) (primitive.ObjectID, error)
	fetchUserByEmail(ctx context.Context, email string) (*models.User, error)
	addTodoIdToUser(ctx context.Context, todoId, userId primitive.ObjectID) error
	startSession() (mongo.Session, error)
}

func newRepoClient(
	db *mongo.Client, logger *utils.Logger,
) userRepo {
	return &repoClient{
		db:     db,
		usersC: utils.GetCollection(db, "users"),
		logger: logger,
	}
}

func (r *repoClient) insertUser(ctx context.Context, user *models.User) (primitive.ObjectID, error) {
	insertedRes, err := r.usersC.InsertOne(ctx, user)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return insertedRes.InsertedID.(primitive.ObjectID), err
}

func (r *repoClient) fetchUserByEmail(ctx context.Context, email string) (*models.User, error) {
	filter := bson.M{
		"email": email,
	}
	var user models.User
	err := r.usersC.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			customErr := &utils.ReqInvalidArgumentError{
				GeneralError: &utils.GeneralError{
					DevInfo: err.Error(),
					Msg:     "user not found",
				},
			}
			return nil, customErr
		}
		return nil, &utils.DBInternalError{
			GeneralError: &utils.GeneralError{
				DevInfo: err.Error(),
				Msg:     "unable to fetch user",
			},
		}
	}
	return &user, err
}

func (r *repoClient) addTodoIdToUser(ctx context.Context, todoId, userId primitive.ObjectID) error {
	filter := bson.M{
		"_id": userId,
	}
	update := bson.M{
		"$addToSet": bson.M{
			"todos": todoId,
		},
	}
	_, err := r.usersC.UpdateOne(
		ctx,
		filter,
		update,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *repoClient) startSession() (mongo.Session, error) {
	return r.db.StartSession()
}
