package env

import (
	"log/slog"
	"os"
)

const (
	Development = "development"
	Production  = "production"
)

const Env = Development

var Log = newLogger()

func newLogger() *slog.Logger {
	if Env != Development {
		return slog.New(slog.DiscardHandler)
	}

	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}
