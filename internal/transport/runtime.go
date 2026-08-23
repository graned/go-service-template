package transport

import "log/slog"

type Runtime struct {
	Address string
	Logger  *slog.Logger
}

func NewRuntime(
	address string,
	logger *slog.Logger,
) Runtime {
	return Runtime{
		Address: address,
		Logger:  logger,
	}
}
