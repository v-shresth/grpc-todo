package cron

import (
	"context"
	"encoding/json"
	"github.com/robfig/cron"
	"todo-grpc/models"
	kafkaQueueProvider "todo-grpc/providers/kafka"
	"todo-grpc/service"
)

func InitCron(ts service.TodoService, kfp kafkaQueueProvider.Provider) {
	c := cron.New()
	c.AddFunc(
		"@every 3m", func() {
			CheckDeadlineIsNear(ts, kfp)
		},
	)
	c.Start()
}

func CheckDeadlineIsNear(ts service.TodoService, kfp kafkaQueueProvider.Provider) {
	todos, err := ts.FetchTodosWithNearbyDeadline(context.Background())
	if err != nil {
		return
	}

	for _, todo := range todos {
		data, err := json.Marshal(&todo)
		if err != nil {
			return
		}
		kfp.Publish(models.TopicDeadlineNearby, data)
	}
}
