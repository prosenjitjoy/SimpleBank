package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"main/database/db"
	"os"

	"github.com/hibiken/asynq"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				slogAttrs := []slog.Attr{
					slog.String("type", task.Type()),
					slog.String("payload", string(task.Payload())),
				}

				var logger *slog.Logger

				if os.Getenv("ENVIRONMENT") == "dev" {
					logger = slog.New(slog.NewTextHandler(os.Stdout, nil).WithAttrs(slogAttrs))
				} else {
					logger = slog.New(slog.NewJSONHandler(os.Stdout, nil).WithAttrs(slogAttrs))
				}

				logger.Error("process task failed", slog.String("error", err.Error()))
			}),
			Logger: NewLogger(),
		},
	)

	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(mux)
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
