package app

import (
	"context"
	"log/slog"
	"os"
)

func NewLogger(ctx context.Context) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
