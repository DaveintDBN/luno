package bot

import (
	"context"
	"fmt"
	luno "github.com/luno/luno-go"
	dec "github.com/luno/luno-go/decimal"
	"github.com/google/uuid"
)

// LunoExecutor places real orders via Luno API with simple risk checks.
// Uses local Client interface from this package.
type LunoExecutor struct {
	client     Client
	position   float64
	entryPrice float64
}

// NewLunoExecutor constructs a live executor using the given client.
func NewLunoExecutor(client Client) *LunoExecutor {
	return &LunoExecutor{client: client}
}

// Execute sends a limit order based on signal, tracking position and enforcing limits.
func (e *LunoExecutor) Execute(ctx context.Context, sig Signal, md MarketData, cfg Config) error {
	price := (md.Bid + md.Ask) / 2
	switch sig {
	case SignalBuy:
		if e.position != 0 {
			return nil // already in position
		}
		if cfg.StakeSize > cfg.PositionLimit {
			return fmt.Errorf("stake %.2f exceeds position limit %.2f", cfg.StakeSize, cfg.PositionLimit)
		}
		req := &luno.PostLimitOrderRequest{
			Pair:            cfg.Pair,
			Price:           dec.NewFromFloat64(price, 8),
			Type:            luno.OrderTypeBid,
			Volume:          dec.NewFromFloat64(cfg.StakeSize, 8),
			BaseAccountId:   cfg.BaseAccountId,
			CounterAccountId:cfg.CounterAccountId,
			ClientOrderId:   uuid.New().String(),
		}
		if _, err := e.client.PostLimitOrder(ctx, req); err != nil {
			return err
		}
		e.position = cfg.StakeSize
		e.entryPrice = price
	case SignalSell:
		if e.position == 0 {
			return nil // no position to exit
		}
		req := &luno.PostLimitOrderRequest{
			Pair:            cfg.Pair,
			Price:           dec.NewFromFloat64(price, 8),
			Type:            luno.OrderTypeAsk,
			Volume:          dec.NewFromFloat64(e.position, 8),
			BaseAccountId:   cfg.BaseAccountId,
			CounterAccountId:cfg.CounterAccountId,
			ClientOrderId:   uuid.New().String(),
		}
		if _, err := e.client.PostLimitOrder(ctx, req); err != nil {
			return err
		}
		e.position = 0
	}
	return nil
}

// CancelAll does nothing for live executor.
func (e *LunoExecutor) CancelAll(ctx context.Context) error {
	// implement if needed to cancel outstanding orders
	return nil
}
