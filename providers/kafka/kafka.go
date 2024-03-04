package kafkaQueueProvider

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"todo-grpc/models"
	"todo-grpc/utils"
)

// KafkaProvider is a Kafka provider that implements the Provider interface
type KafkaProvider struct {
	deadLineNearbyWriter *kafka.Writer
	premiumEndingWriter  *kafka.Writer
	envProvider          *utils.EnvConfig
	logger               *utils.Logger
}

type Provider interface {
	// Publish PublishData publishes data to a message queue
	Publish(topic models.QueueTopics, message []byte)
	Reconnect()
	Close()
}

func NewKafkaProvider(env *utils.EnvConfig, logger *utils.Logger) Provider {
	kafkaHost := env.GetKafkaHost()
	deadLineNearbyWriter := &kafka.Writer{
		Addr:        kafka.TCP(kafkaHost),
		Topic:       string(models.TopicDeadlineNearby),
		Async:       true,
		ErrorLogger: logger,
	}

	premiumEndingWriter := &kafka.Writer{
		Addr:        kafka.TCP(kafkaHost),
		Topic:       string(models.TopicPremiumEnding),
		Async:       true,
		ErrorLogger: logger,
	}

	return &KafkaProvider{
		deadLineNearbyWriter: deadLineNearbyWriter,
		premiumEndingWriter:  premiumEndingWriter,
		envProvider:          env,
		logger:               logger,
	}
}

// Publish publishes data to a Kafka topic
func (k *KafkaProvider) Publish(topic models.QueueTopics, message []byte) {
	switch topic {
	case models.TopicDeadlineNearby:
		err := k.deadLineNearbyWriter.WriteMessages(
			context.Background(),
			kafka.Message{
				Value: message,
			},
		)
		if err != nil {
			k.logger.Error(
				err, "Publish: failed to write unassigned ride request message on kafka: %v",
			)
			k.Reconnect()
		}
		k.logger.Info("Published kafka message: %v", message)

	case models.TopicPremiumEnding:
		err := k.premiumEndingWriter.WriteMessages(
			context.Background(),
			kafka.Message{
				Value: message,
			},
		)
		if err != nil {
			k.logger.Error(err, "Publish: failed to write push notification message on kafka: %v")
			k.Reconnect()
		}
		k.logger.Info("Published kafka message: %v", message)

	default:
		k.logger.Error(errors.New("trying to publish on wrong topic"), "Trying to publish on wrong topic")
	}
}

func (k *KafkaProvider) Reconnect() {
	kafkaHost := k.envProvider.GetKafkaHost()

	k.premiumEndingWriter = &kafka.Writer{
		Addr:  kafka.TCP(kafkaHost),
		Topic: string(models.TopicPremiumEnding),
	}

	k.deadLineNearbyWriter = &kafka.Writer{
		Addr:  kafka.TCP(kafkaHost),
		Topic: string(models.TopicDeadlineNearby),
	}

}

func (k *KafkaProvider) Close() {
	if err := k.premiumEndingWriter.Close(); err != nil {
		k.logger.Error(err, "error closing kafka connection")
	}
}
