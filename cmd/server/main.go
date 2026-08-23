package main

import (
	applogger "github.com/graned/go-service-template/internal/logger"
	"github.com/graned/go-service-template/internal/transport"
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

	restRuntime := transport.NewRuntime(
		":3000",
		logger,
	)

	restServer := rest.New(
		restRuntime,
		userService,
	)

	servers := []transport.Server{
		restServer,
	}
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer stop()

	serverErrors := make(chan error, len(servers))

	for _, server := range servers {
		go func(s transport.Server) {
			if err := s.Run(); err != nil {
				serverErrors <- err
			}
		}(server)
	}

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

	shutdownErrors := make(chan error, len(servers))

	for _, server := range servers {
		go func(s transport.Server) {
			shutdownErrors <- server.Shutdown(shutdownCtx)
		}(server)
	}

	for range servers {
		if err := <-shutdownErrors; err != nil {
			logger.Error(
				"Server shutdown failed",
				"error", err,
			)
		}
	}
	logger.Info("Application stopped")
}
