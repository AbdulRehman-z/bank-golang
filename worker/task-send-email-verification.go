package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AbdulRehman-z/bank-golang/mail"
	"github.com/AbdulRehman-z/bank-golang/util"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	TaskSendVerificationEmail = "task:send_verification_email"
)

type PayloadSendVerificationEmail struct {
	Username string `json:"username"`
}

func (d *RedisTaskDistributor) TaskSendVerificationEmail(ctx context.Context, payload *PayloadSendVerificationEmail, options ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("cannot marshal payload: %w", err)
	}
	task := asynq.NewTask(TaskSendVerificationEmail, jsonPayload, options...)
	info, err := d.client.Enqueue(task)
	if err != nil {
		return fmt.Errorf("cannot enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendEmailVerify(ctx context.Context, task *asynq.Task) error {
	config, err := util.LoadConfig(".")
	if err != nil {
		return fmt.Errorf("cannot load config: %w", err)
	}

	var payload PayloadSendVerificationEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("cannot unmarshal payload: %w", err)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		return fmt.Errorf("cannot get user: %w", err)
	}

	// send email
	senderEmail := "yousafbhaikhan10@gmail.com"
	mailSender := mail.NewGmailSender(user.Username, senderEmail, config.APP_PASSWORD)
	receiverEmail := []string{user.Email}
	mailSender.SendEmail(receiverEmail, "Verify your email", "Please verify your email by clicking this link: http://localhost:8080/verify-email?token=123")

	log.Info().Str("type", task.Type()).Str("username", payload.Username).
		Str("email", user.Email).Msg("processed task")

	return nil
}
