package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"main/database/db"
	"main/util"
	"os"

	"github.com/hibiken/asynq"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	slogAttrs := []slog.Attr{
		slog.String("type", task.Type()),
		slog.String("payload", string(task.Payload())),
		slog.String("queue", info.Queue),
		slog.Int("max_retry", info.MaxRetry),
	}

	var logger *slog.Logger

	if os.Getenv("ENVIRONMENT") == "dev" {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil).WithAttrs(slogAttrs))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil).WithAttrs(slogAttrs))
	}

	logger.Info("enqueued task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		// if err == pgx.ErrNoRows {
		// 	return fmt.Errorf("user doesn't exist: %w", asynq.SkipRetry)
		// }

		return fmt.Errorf("failed to get user: %w", err)
	}

	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, &db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	})
	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	subject := "Welcome to Simple Bank"
	verifyUrl := fmt.Sprintf("http://localhost:3000/v1/verify_email?email_id=%d&secret_code=%s", verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(`
	<h1>Hello %s</h1>
	<p>Thank you for registering with us!</p>
	<p>Please <a href="%s">click here</a> to verify your email address</p>
	`, user.FullName, verifyUrl)
	to := []string{user.Email}

	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	slogAttrs := []slog.Attr{
		slog.String("type", task.Type()),
		slog.String("payload", string(task.Payload())),
		slog.String("email", user.Email),
	}

	var logger *slog.Logger

	if os.Getenv("ENVIRONMENT") == "dev" {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil).WithAttrs(slogAttrs))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil).WithAttrs(slogAttrs))
	}

	logger.Info("processed task")
	return nil
}
