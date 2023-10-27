package main

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"main/api"
	"main/database/db"
	"main/gapi"
	"main/mail"
	"main/pb"
	"main/util"
	"main/worker"
	"net"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hibiken/asynq"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

//go:embed swagger/*
var content embed.FS

func main() {
	cfg, err := util.LoadConfig(".env")
	if err != nil {
		slog.Error("cannot load config:", slog.String("error", err.Error()))
		return
	}

	if cfg.Environment == "dev" {
		var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
		slog.SetDefault(logger)
	} else {
		var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
		slog.SetDefault(logger)
	}

	conn, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		slog.Error("cannot connect to db:", slog.String("error", err.Error()))
		return
	}

	// run db migration
	runMigration(cfg.MigrationURL, cfg.DatabaseURL)

	store := db.NewStore(conn)

	// redis
	redisOpt := asynq.RedisClientOpt{Addr: cfg.RedisAddress}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	go runTaskProcessor(cfg, redisOpt, store)

	go runGatewayServer(store, cfg, taskDistributor)
	runGrpcServer(store, cfg, taskDistributor)
}

func runTaskProcessor(cfg *util.ConfigDatabase, redisOpt asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewGmailSender(cfg.EmailSenderName, cfg.EmailSenderAddress, cfg.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)
	slog.Info("start task processor")

	err := taskProcessor.Start()
	if err != nil {
		slog.Error("Failed to start task processor", slog.String("error", err.Error()))
		return
	}
}

func runMigration(migrationURL string, databaseURL string) {
	migration, err := migrate.New(migrationURL, databaseURL)
	if err != nil {
		slog.Error("cannot create new migrate instance:", slog.String("error", err.Error()))
		return
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("failed to run migrate up:", slog.String("error", err.Error()))
		return
	}

	slog.Info("db migrated successfully")
}

func runHttpServer(store db.Store, cfg *util.ConfigDatabase) {
	server, err := api.NewServer(store, cfg)
	if err != nil {
		slog.Error("cannot initialize server:", slog.String("error", err.Error()))
		return
	}

	err = server.Start(cfg.HTTPServerAddress)
	if err != nil {
		slog.Error("cannot start HTTP server:", slog.String("error", err.Error()))
		return
	}
}

func runGrpcServer(store db.Store, cfg *util.ConfigDatabase, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(store, cfg, taskDistributor)
	if err != nil {
		slog.Error("cannot initialize server:", slog.String("error", err.Error()))
		return
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)

	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", cfg.GRPCServerAddress)
	if err != nil {
		slog.Error("cannot create listener:", slog.String("error", err.Error()))
		return
	}

	slog.Info(fmt.Sprintf("starting gRPC server at %s", listener.Addr().String()))
	err = grpcServer.Serve(listener)
	if err != nil {
		slog.Error("cannot start gRPC server:", slog.String("error", err.Error()))
		return
	}
}

func runGatewayServer(store db.Store, cfg *util.ConfigDatabase, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(store, cfg, taskDistributor)
	if err != nil {
		slog.Error("cannot initialize server:", slog.String("error", err.Error()))
		return
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		slog.Error("cannot register handler server:", slog.String("error", err.Error()))
		return
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	fs := http.FileServer(http.FS(content))
	mux.Handle("/doc/", http.StripPrefix("/doc/", fs))

	listener, err := net.Listen("tcp", cfg.HTTPServerAddress)
	if err != nil {
		slog.Error("cannot create listener:", slog.String("error", err.Error()))
		return
	}

	slog.Info(fmt.Sprintf("starting HTTP gateway server at %s", listener.Addr().String()))

	handler := gapi.HttpLogger(mux)

	err = http.Serve(listener, handler)
	if err != nil {
		slog.Error("cannot start HTTP gateway server:", err)
	}
}
