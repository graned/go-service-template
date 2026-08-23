package app

import (
	"context"
	"log/slog"
	"time"
)

type Component interface {
	Run() error
	Shutdown(context.Context) error
}

type App struct {
	logger          *slog.Logger
	components      []Component
	shutdownTimeout time.Duration
}

func New(
	logger *slog.Logger,
	shutdownTimeout time.Duration,
	components ...Component,
) *App {
	return &App{
		logger:          logger,
		components:      components,
		shutdownTimeout: shutdownTimeout,
	}
}

func (a *App) Run(ctx context.Context) error {
	errors := make(chan error, len(a.components))

	for _, component := range a.components {
		go func(c Component) {
			if err := c.Run(); err != nil {
				errors <- err
			}
		}(component)
	}

	var runErr error

	select {
	case <-ctx.Done():
		a.logger.Info("application shutdown requested")

	case err := <-errors:
		a.logger.Error(
			"application component failed",
			"error", err,
		)

		runErr = err
	}

	a.shutdown()

	return runErr
}

func (a *App) shutdown() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		a.shutdownTimeout,
	)
	defer cancel()

	a.logger.Info(
		"shutting down application",
		"timeout", a.shutdownTimeout.String(),
	)

	errors := make(chan error, len(a.components))

	for _, component := range a.components {
		go func(c Component) {
			errors <- c.Shutdown(ctx)
		}(component)
	}

	for range a.components {
		if err := <-errors; err != nil {
			a.logger.Error(
				"component shutdown failed",
				"error", err,
			)
		}
	}

	a.logger.Info("application stopped")
}
