package internal

import (
	"github.com/segmentio/kafka-go"
	"net"
	"strconv"
	"todo-grpc/internal/processDeadlineNearby"
	"todo-grpc/utils"
)

func SubscribeAllPartitions(logger *utils.Logger, config utils.EnvConfig) {
	kafkaHost := config.GetKafkaHost()

	conn, err := kafka.Dial("tcp", kafkaHost)
	if err != nil {
		logger.Error(err, "error while connecting to kafka:")
		return
	}

	var controllerConn *kafka.Conn
	defer func() {
		_ = conn.Close()
		_ = controllerConn.Close()
	}()

	controller, err := conn.Controller()
	if err != nil {
		logger.Error(err, "error while connecting to kafka:")
		return
	}

	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		logger.Error(err, "error while connecting to kafka:")
		panic(err.Error())
	}

	processDeadlineNearby.SubscribeDeadlineNearbyPartitions(logger, config, controllerConn)
}
