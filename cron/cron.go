package cron

import (
	"context"
	"encoding/json"
	"github.com/robfig/cron"
	"todo-grpc/models"
	kafkaQueueProvider "todo-grpc/providers/kafka"
	"todo-grpc/service"
)

func InitCron(ctx context.Context, ts service.TodoService, kfp kafkaQueueProvider.Provider) {
	c := cron.New()
	c.AddFunc(
		"@every 3600m", func() {
			CheckDeadlineIsNear(ctx, ts, kfp)
		},
	)
	c.Start()
}

func CheckDeadlineIsNear(ctx context.Context, ts service.TodoService, kfp kafkaQueueProvider.Provider) {
	todos, err := ts.FetchTodosWithNearbyDeadline(ctx)
	if err != nil {
		return
	}

	for _, todo := range todos {
		data, err := json.Marshal(&todo)
		if err != nil {
			//err = ("unable to marshall message to bytes for kafka queue", err)
			return
		}
		kfp.Publish(models.TopicDeadlineNearby, data)
	}
}
