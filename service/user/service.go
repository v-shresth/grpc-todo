package user

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"todo-grpc/models"
	"todo-grpc/pb"
	"todo-grpc/service"
	"todo-grpc/utils"
)

type serviceClient struct {
	userRepo userRepo
	logger   *utils.Logger
	config   utils.EnvConfig
}

func NewUserService(
	db *mongo.Client,
	logger *utils.Logger,
	config utils.EnvConfig,
) service.UserService {
	return &serviceClient{
		userRepo: newRepoClient(db, logger),
		logger:   logger,
		config:   config,
	}
}

func (s *serviceClient) Register(ctx context.Context, user *models.User) (*pb.RegisterResponse, error) {
	var err error
	user.ID = primitive.NewObjectID()
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	userId, err := s.userRepo.insertUser(ctx, user)
	if err != nil {
		return nil, err
	}
	token, err := utils.GenerateToken(userId.Hex(), s.config)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterResponse{
		Token:  token,
		UserId: userId.Hex(),
	}, nil
}

func (s *serviceClient) Login(ctx context.Context, user *models.User) (string, error) {
	dbUser, err := s.userRepo.fetchUserByEmail(ctx, user.Email)
	if err != nil {
		return "", err
	}

	if ok := utils.CheckPassword(user.Password, dbUser.Password); !ok {
		customErr := &utils.ReqInvalidArgumentError{
			GeneralError: &utils.GeneralError{
				DevInfo: err.Error(),
				Msg:     "invalid todo id",
			},
		}
		return "", customErr
	}

	token, err := utils.GenerateToken(dbUser.ID.Hex(), s.config)
	if err != nil {
		customErr := &utils.SystemInternalError{
			GeneralError: &utils.GeneralError{
				DevInfo: err.Error(),
				Msg:     "invalid todo id",
			},
		}
		return "", customErr
	}
	return token, nil
}

func (s *serviceClient) AddTodoIdToUser(ctx context.Context, userId, todoId primitive.ObjectID) error {
	return s.userRepo.addTodoIdToUser(ctx, todoId, userId)
}
