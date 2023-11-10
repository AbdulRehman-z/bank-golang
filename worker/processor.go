package worker

import (
	"context"

	db "github.com/AbdulRehman-z/bank-golang/db/sqlc"
	"github.com/AbdulRehman-z/bank-golang/mail"
	"github.com/rs/zerolog/log"

	"github.com/hibiken/asynq"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendEmailVerify(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
	mailer mail.Mail
}

func NewRedisTaskProcessor(redisOpts asynq.RedisClientOpt, store db.Store, mailer mail.Mail) TaskProcessor {
	server := asynq.NewServer(redisOpts, asynq.Config{
		Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Error().Err(err).Str("type", task.Type()).Bytes("payload", task.Payload()).Msg("failed to process task")
		}),
		Logger: NewLogger(),
	})

	return &RedisTaskProcessor{
		server: server,
		store:  store,
		mailer: mailer,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerificationEmail, processor.ProcessTaskSendEmailVerify)
	return processor.server.Start(mux)
}
