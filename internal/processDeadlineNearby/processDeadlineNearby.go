package processDeadlineNearby

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"todo-grpc/models"
	"todo-grpc/utils"
)

const (
	defaultPartitions        = 1
	defaultReplicationFactor = 1
)

func SubscribeDeadlineNearbyPartitions(logger *utils.Logger, config utils.EnvConfig, controllerConn *kafka.Conn) {
	err := controllerConn.CreateTopics(
		kafka.TopicConfig{
			Topic:             string(models.TopicDeadlineNearby),
			NumPartitions:     defaultPartitions,
			ReplicationFactor: defaultReplicationFactor,
		},
	)
	if err != nil {
		logger.Error(err, "SubscribeAllPartitions: error creating topic %v")
		return
	}

	go SubscribeProcessPaymentStatus(config, logger)
}

func SubscribeProcessPaymentStatus(env utils.EnvConfig, logger *utils.Logger) {
	kafkaHost := env.GetKafkaHost()
	groupID := fmt.Sprintf("processed-dealine-nearby")

	r := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:     []string{kafkaHost},
			Topic:       string(models.TopicDeadlineNearby),
			GroupID:     groupID,
			Logger:      logger,
			ErrorLogger: logger,
		},
	)
	_ = r.SetOffset(kafka.LastOffset)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			logger.Error(err, "error reading message %v\n")
			continue
		}

		var kafkaMsg models.Todo
		err = json.Unmarshal(m.Value, &kafkaMsg)
		if err != nil {
			logger.Error(err, "error unmarshalling message %v\n")
			continue
		}

		// get ctx with request id same as that of original request
		ctx := context.WithValue(context.Background(), utils.RequestIdKey, kafkaMsg.ID.Hex())

		err = r.CommitMessages(ctx, m)
		if err != nil {
			fmt.Println("Failed while committing message", err)
		}

		fmt.Println("New message:", kafkaMsg)
	}
}
