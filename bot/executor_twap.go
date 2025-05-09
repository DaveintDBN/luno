package bot

import (
	"context"
	"fmt"
	"time"
)

// TWAPExecutor slices large orders into multiple smaller orders.
type TWAPExecutor struct {
	Inner    Executor
	Slices   int
	Interval time.Duration
}

// NewTWAPExecutor creates a TWAP executor that executes orders in Slices over Interval durations.
func NewTWAPExecutor(inner Executor, slices int, interval time.Duration) *TWAPExecutor {
	if slices <= 1 {
		slices = 1
	}
	return &TWAPExecutor{Inner: inner, Slices: slices, Interval: interval}
}

// Execute slices the execution into smaller timed chunks.
func (t *TWAPExecutor) Execute(ctx context.Context, sig Signal, md MarketData, cfg Config) error {
	// No action if no signal
	if sig == SignalNone {
		return nil
	}
	fmt.Printf("TWAPExecutor: executing %d slices every %s\n", t.Slices, t.Interval)
	// Divide stake size across slices
	sliceSize := cfg.StakeSize / float64(t.Slices)
	for i := 0; i < t.Slices; i++ {
		// configure this slice
		sliceCfg := cfg
		sliceCfg.StakeSize = sliceSize
		if err := t.Inner.Execute(ctx, sig, md, sliceCfg); err != nil {
			return err
		}
		// wait for next slice if not last
		if i < t.Slices-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(t.Interval):
			}
		}
	}
	return nil
}

// CancelAll delegates cancellation.
func (t *TWAPExecutor) CancelAll(ctx context.Context) error {
	return t.Inner.CancelAll(ctx)
}
