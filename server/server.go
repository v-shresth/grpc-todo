package server

import (
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"todo-grpc/api"
	"todo-grpc/cron"
	"todo-grpc/middleware"
	"todo-grpc/pb"
	kafkaQueueProvider "todo-grpc/providers/kafka"
	"todo-grpc/service/alerts"
	"todo-grpc/service/todo"
	"todo-grpc/service/user"
	"todo-grpc/utils"
)

func NewServer(
	db *mongo.Client,
	logger *utils.Logger,
	config utils.EnvConfig,
	kafkaProvider kafkaQueueProvider.Provider,
) *grpc.Server {
	srv := &api.Server{
		TodoSvc:       todo.NewTodoService(db, logger, config),
		UserSvc:       user.NewUserService(db, logger, config),
		AlertSvc:      alerts.NewAlertService(db, logger, config),
		Config:        config,
		Logger:        logger,
		KafkaProvider: kafkaProvider,
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(middleware.AuthMiddleware),
		grpc.ChainStreamInterceptor(middleware.AuthStreamInterceptor),
	)

	pb.RegisterUserServiceServer(server, srv)
	pb.RegisterTodoServiceServer(server, srv)

	reflection.Register(server)

	go cron.InitCron(srv.TodoSvc, kafkaProvider)

	return server
}
