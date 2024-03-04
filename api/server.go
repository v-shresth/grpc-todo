package api

import (
	"todo-grpc/pb"
	kafkaQueueProvider "todo-grpc/providers/kafka"
	"todo-grpc/service"
	"todo-grpc/utils"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	pb.UnimplementedTodoServiceServer

	TodoSvc       service.TodoService
	UserSvc       service.UserService
	Config        *utils.EnvConfig
	Logger        *utils.Logger
	KafkaProvider kafkaQueueProvider.Provider
}
