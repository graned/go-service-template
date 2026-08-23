package main

import (
	applogger "github.com/graned/go-service-template/internal/logger"
	"github.com/graned/go-service-template/internal/transport/rest"
	"github.com/graned/go-service-template/internal/user"

	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := applogger.New()

	userService := user.NewService()

	restServer := rest.New(
		":3000",
		logger,
		userService,
	)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer stop()

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info("Starting rest server")
		serverErrors <- restServer.Run()
	}()

	select {
	case err := <-serverErrors:
		if err != nil {
			logger.Error(
				"Server Failed",
				slog.Any("error", err),
			)
			os.Exit(1)
		}
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := restServer.Shutdown(shutdownCtx); err != nil {
		logger.Error(
			"Server shutdown failed",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
	logger.Info("Rest server stopped")
}
