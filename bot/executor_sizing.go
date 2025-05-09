package bot

import (
	"context"
)

// SizingExecutor wraps an Executor and applies position sizing.
type SizingExecutor struct {
	Inner Executor
	Sizer PositionSizer
}

// NewSizingExecutor constructs a SizingExecutor.
func NewSizingExecutor(inner Executor, sizer PositionSizer) *SizingExecutor {
	return &SizingExecutor{Inner: inner, Sizer: sizer}
}

// Execute computes stake size via the sizer, updates cfg, and delegates execution.
func (s *SizingExecutor) Execute(ctx context.Context, sig Signal, md MarketData, cfg Config) error {
	// Compute dynamic stake size
	size := s.Sizer.Size(cfg.InitialEquity, cfg)
	cfg.StakeSize = size
	return s.Inner.Execute(ctx, sig, md, cfg)
}

// CancelAll delegates cancellation.
func (s *SizingExecutor) CancelAll(ctx context.Context) error {
	return s.Inner.CancelAll(ctx)
}
