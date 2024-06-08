package main

import (
	"auth_service/internal/app"
	"auth_service/internal/config"
	"auth_service/internal/lib/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	application := app.New(log, cfg.GRPC.Port, cfg.TokenTTL)

	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	// initiate  shutdown
	application.GRPCServer.Stop()
	application.Storage.Stop()
	log.Info("Gracefully stopped")

}
