package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

// TODO: refactor to use viper config
func Init(ctx context.Context) Config {
	cfg := Config{
		Address:   ":8080",
		DebugMode: true,
	}

	err := godotenv.Load(".env")
	if err != nil {
		err = fmt.Errorf("failed to load environment variables: %w", err)
		panic(err)
	}

	cfg.PostgresDSN = os.Getenv("DB_URL")

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
	Address     string // e.g. ":8080"
	DebugMode   bool   // e.g.  false
	PostgresDSN string
}
