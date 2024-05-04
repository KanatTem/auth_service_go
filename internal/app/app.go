package app

import (
	"auth_service/internal/app/grpc"
	"auth_service/internal/config"
	"auth_service/internal/services/auth"
	"auth_service/internal/storage/postgress"
	"fmt"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpc.App
	Storage    *postgress.Storage
}

func New(
	log *slog.Logger,
	grpcPort int,
	tokenTTL time.Duration,
) *App {
	cfg := config.MustLoad()

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
	)
	storage, err := postgress.New(connStr)

	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpc.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
		Storage:    storage,
	}
}
