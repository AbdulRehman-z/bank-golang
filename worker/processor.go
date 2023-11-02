package worker

import (
	"context"

	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/hibiken/asynq"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendEmailVerify(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(redisOpts asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(redisOpts, asynq.Config{})
	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerificationEmail, processor.ProcessTaskSendEmailVerify)
	return processor.server.Start(mux)
}
