package bot

import (
    "context"
    "log"
    "time"
)

// LoggingExecutor wraps an Executor to log each action and errors.
type LoggingExecutor struct {
    inner          Executor
    activityLogger *log.Logger
    errorLogger    *log.Logger
}

// NewLoggingExecutor constructs a new LoggingExecutor.
func NewLoggingExecutor(inner Executor, activityLogger, errorLogger *log.Logger) *LoggingExecutor {
    return &LoggingExecutor{inner: inner, activityLogger: activityLogger, errorLogger: errorLogger}
}

// Execute logs the signal, market data, config, and any execution errors.
func (l *LoggingExecutor) Execute(ctx context.Context, sig Signal, md MarketData, cfg Config) error {
    l.activityLogger.Printf("Execute: signal=%v, bid=%.8f, ask=%.8f, time=%s, cfg=%+v", sig, md.Bid, md.Ask, md.Timestamp.Format(time.RFC3339), cfg)
    err := l.inner.Execute(ctx, sig, md, cfg)
    if err != nil {
        l.errorLogger.Printf("Execute error: %v", err)
    }
    return err
}

// CancelAll logs cancel events and any errors.
func (l *LoggingExecutor) CancelAll(ctx context.Context) error {
    l.activityLogger.Printf("CancelAll")
    err := l.inner.CancelAll(ctx)
    if err != nil {
        l.errorLogger.Printf("CancelAll error: %v", err)
    }
    return err
}
