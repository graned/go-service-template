package main

import (
	"github.com/graned/go-service-template/internal/app"
	applogger "github.com/graned/go-service-template/internal/logger"
	"github.com/graned/go-service-template/internal/transport"
	restserver "github.com/graned/go-service-template/internal/transport/rest"
	"github.com/graned/go-service-template/internal/user"

	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := applogger.New()
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer stop()

	userService := user.NewService()

	restRuntime := transport.NewRuntime(
		":3000",
		logger,
	)

	restServer := restserver.New(
		restRuntime,
		userService,
	)

	application := app.New(
		logger,
		10*time.Second,
		restServer,
	)

	if err := application.Run(ctx); err != nil {
		logger.Error(
			"application exited with error",
			"error", err,
		)

		os.Exit(1)
	}
}
