package worker

import (
	"context"
	"log/slog"
	"main/database/db"
	"main/mail"
	"os"

	"github.com/hibiken/asynq"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	Shutdown()
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
	mailer mail.EmailSender
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, mailer mail.EmailSender) TaskProcessor {
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
		mailer: mailer,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(mux)
}

func (processor *RedisTaskProcessor) Shutdown() {
	processor.server.Shutdown()
}
