package util

import (
	"context"

	"go.uber.org/zap"
)

type contextKey struct {
	name string
}

var (
	loggerCtxKey = &contextKey{"Logger"}
)

// ContextWithLogger context with zap logger
func ContextWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}

// LoggerFromContext zap logger from context
func LoggerFromContext(ctx context.Context) (*zap.Logger, bool) {
	logger, ok := ctx.Value(loggerCtxKey).(*zap.Logger)
	return logger, ok
}
