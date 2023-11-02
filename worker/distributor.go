package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	TaskSendVerificationEmail(ctx context.Context, payload *PayloadSendVerificationEmail, options ...asynq.Option) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewTaskDistributor(options asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(options)
	return &RedisTaskDistributor{
		client: client,
	}
}
