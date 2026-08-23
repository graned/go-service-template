package app

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type Component interface {
	Name() string
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
				a.logger.Error(
					"application component failed",
					"error", err,
					"component", component.Name(),
				)
				errors <- err
			}
		}(component)
	}

	var runErr error

	select {
	case <-ctx.Done():
		a.logger.Info("application shutdown requested")

	case err := <-errors:
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

	var wg sync.WaitGroup

	for _, component := range a.components {
		wg.Add(1)

		go func(c Component) {
			defer wg.Done()

			a.logger.Info(
				"Shutting down component",
				"timeout", a.shutdownTimeout.String(),
				"component", c.Name(),
			)
			if error := c.Shutdown(ctx); error != nil {
				a.logger.Error(
					"Component shutdown failed",
					"error", error,
					"component", c.Name(),
				)
			}
		}(component)
	}

	wg.Wait()

	a.logger.Info("application stopped")
}
