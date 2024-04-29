package main

import (
	"auth_service/internal/app"
	"auth_service/internal/config"
	"auth_service/internal/lib/logger"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	if log == nil {
	}

	application := app.New(log, cfg.GRPC.Port, cfg.TokenTTL)

	//
	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)
	defer cancelCtx() // ← make sure we always call this, even on panic
	_, err := application.Storage.IsAdmin(ctx, 33)
	if err != nil {
		panic(err)
	}
	//

	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	// initiate  shutdown
	application.GRPCServer.Stop()
	application.Storage.Stop()
	log.Info("Gracefully stopped")

}
