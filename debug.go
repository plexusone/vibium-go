package vibium

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

// debugKey is the context key for the debug logger.
type debugKey struct{}

// Debug returns true if debug logging is enabled via VIBIUM_DEBUG environment variable.
func Debug() bool {
	val := os.Getenv("VIBIUM_DEBUG")
	return val == "1" || strings.EqualFold(val, "true")
}

// NewDebugLogger creates a new debug logger that writes to stderr.
// Returns nil if debug logging is disabled.
func NewDebugLogger() *slog.Logger {
	if !Debug() {
		return nil
	}
	return slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

// ContextWithLogger returns a new context with the logger attached.
func ContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	if logger == nil {
		return ctx
	}
	return context.WithValue(ctx, debugKey{}, logger)
}

// LoggerFromContext returns the logger from the context, or nil if not present.
func LoggerFromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(debugKey{}).(*slog.Logger); ok {
		return logger
	}
	return nil
}

// debugLog logs a debug message if a logger is present in the context.
func debugLog(ctx context.Context, msg string, args ...any) {
	if logger := LoggerFromContext(ctx); logger != nil {
		logger.Debug(msg, args...)
	}
}
