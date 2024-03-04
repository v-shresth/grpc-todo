package main

import (
	"log"
	"net"
	kafkaQueueProvider "todo-grpc/providers/kafka"
	"todo-grpc/server"
	"todo-grpc/utils"
)

func main() {
	logger := utils.NewLogger()
	config, err := utils.NewEnvConfig()
	if err != nil {
		logger.Error(err, "enable to load env config")
		log.Fatal()
	}

	kafkaProvider := kafkaQueueProvider.NewKafkaProvider(config, logger)

	db, err := utils.ConnectMongoDB(config.GetMongoURI())
	if err != nil {
		logger.Error(err, "enable to connect database")
		log.Fatal()
	}

	grpcServer := server.NewServer(db, logger, config, kafkaProvider)

	lis, err := net.Listen("tcp", config.GetServerPort())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
