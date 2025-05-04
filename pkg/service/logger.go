package service

import (
	"context"
	"log/slog"
	"os"
	"sync"
)

type LogHandler struct {
	slog.Handler
}

type contextKey string

const (
	ToolNameKey contextKey = "toolName"
	ApiNameKey  contextKey = "apiName"
)

var syncOnceLogger sync.Once

var logger *slog.Logger

// Custom log handler to add toolName and apiName attributes to each log
func (l *LogHandler) Handle(ctx context.Context, r slog.Record) error {
	if toolName, ok := ctx.Value(ToolNameKey).(string); ok {
		r.AddAttrs(slog.String("toolName", toolName))
	}
	if apiName, ok := ctx.Value(ApiNameKey).(string); ok {
		r.AddAttrs(slog.String("apiName", apiName))
	}
	return l.Handler.Handle(ctx, r)
}

func GetLogger() *slog.Logger {
	syncOnceLogger.Do(func() {
		if logger == nil {
			attributes := []slog.Attr{}
			baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: false}).WithAttrs(attributes)
			customHandler := &LogHandler{Handler: baseHandler}
			logger = slog.New(customHandler)
		}
	})
	return logger
}
