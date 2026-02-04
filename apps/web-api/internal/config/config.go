package config

import (
	"context"
	"log/slog"
	"os"
)

func Init(ctx context.Context) Config {
	cfg := Config{ // TODO: init this from env and/or config file
		Address:   ":8080",
		DebugMode: true,
	}

	slogger := NewSlogger(cfg.DebugMode)
	slog.SetDefault(slogger)

	return cfg
}

func NewSlogger(verbose bool) *slog.Logger {
	var handlerOpt *slog.HandlerOptions
	if verbose {
		handlerOpt = &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}
	}

	return slog.New(slog.NewTextHandler(os.Stderr, handlerOpt))
}

type Config struct {
	Address   string // e.g. ":8080"
	DebugMode bool   // e.g.  false
}
