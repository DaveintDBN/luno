package bot

import (
	"context"
	"fmt"
	"time"
)

// SimulatedExecutor enforces risk controls and simulates order execution.
type SimulatedExecutor struct {
	Position            float64   // current position size
	EntryPrice          float64   // price at entry
	TotalPnL            float64   // cumulative PnL
	PeakPnL             float64   // highest PnL
	MaxDrawdownExceeded bool      // flag if drawdown breached
	LastTradeTime       time.Time // last execution timestamp
}

// NewSimulatedExecutor constructs a new Simulation executor.
func NewSimulatedExecutor() *SimulatedExecutor {
	return &SimulatedExecutor{}
}

// Execute processes a trading signal using market data and config, enforcing limits.
func (e *SimulatedExecutor) Execute(ctx context.Context, sig Signal, md MarketData, cfg Config) error {
	// cooldown enforcement
	if !e.LastTradeTime.IsZero() && md.Timestamp.Sub(e.LastTradeTime) < cfg.Cooldown {
		return nil
	}
	e.LastTradeTime = md.Timestamp

	price := (md.Bid + md.Ask) / 2

	switch sig {
	case SignalBuy:
		// only enter if no position
		if e.Position != 0 {
			return nil
		}
		if cfg.StakeSize > cfg.PositionLimit {
			return fmt.Errorf("stake size %.2f > position limit %.2f", cfg.StakeSize, cfg.PositionLimit)
		}
		e.Position = cfg.StakeSize
		e.EntryPrice = price
	case SignalSell:
		// only exit if in position
		if e.Position == 0 {
			return nil
		}
		profit := (price - e.EntryPrice) * e.Position
		e.TotalPnL += profit
		// update peak for drawdown
		if e.TotalPnL > e.PeakPnL {
			e.PeakPnL = e.TotalPnL
		}
		drawdown := e.PeakPnL - e.TotalPnL
		if drawdown > cfg.MaxDrawdown {
			e.MaxDrawdownExceeded = true
			return fmt.Errorf("max drawdown %.2f exceeded", cfg.MaxDrawdown)
		}
		e.Position = 0
	}
	return nil
}

// CancelAll resets any open position.
func (e *SimulatedExecutor) CancelAll(ctx context.Context) error {
	if e.Position != 0 {
		e.Position = 0
	}
	return nil
}
