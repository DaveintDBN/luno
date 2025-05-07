package bot

import (
	"context"
	"time"

	"github.com/luno/luno-go"
)

// Client wraps Luno API calls used by the bot.
type Client interface {
	// SetAuth configures API credentials.
	SetAuth(id, secret string) error
	GetTickers(ctx context.Context, req *luno.GetTickersRequest) (*luno.GetTickersResponse, error)
	GetOrderBook(ctx context.Context, req *luno.GetOrderBookRequest) (*luno.GetOrderBookResponse, error)
	PostLimitOrder(ctx context.Context, req *luno.PostLimitOrderRequest) (*luno.PostLimitOrderResponse, error)
	// Fetch historical trades for backtesting
	ListTrades(ctx context.Context, req *luno.ListTradesRequest) (*luno.ListTradesResponse, error)
	// Fetch historical candles for backtesting
	GetCandles(ctx context.Context, req *luno.GetCandlesRequest) (*luno.GetCandlesResponse, error)
}

// Strategy generates trading signals.
type Strategy interface {
	// Next returns a Signal based on market data and config.
	Next(data MarketData, cfg Config) Signal
}

// Executor places and manages orders based on signals.
type Executor interface {
	Execute(ctx context.Context, sig Signal, md MarketData, cfg Config) error
	CancelAll(ctx context.Context) error
}

// Signal indicates trading actions.
type Signal int

const (
	SignalNone Signal = iota
	SignalBuy
	SignalSell
)

// Config holds adjustable parameters for a strategy.
type Config struct {
	Pair           string        // e.g. "XBTZAR"
	EntryThreshold float64       // threshold to enter
	ExitThreshold  float64       // threshold to exit
	StakeSize      float64       // amount per trade
	Cooldown       time.Duration // min time between trades
	PositionLimit  float64       // max allowed position size
	MaxDrawdown    float64       // max permitted drawdown
	ShortWindow    int           // SMA short window
	LongWindow     int           // SMA long window
	BaseAccountId    int64         // base currency account ID for trades
	CounterAccountId int64         // counter currency account ID for trades
}

// MarketData packages latest market metrics.
type MarketData struct {
	Bid       float64
	Ask       float64
	Timestamp time.Time
}

// LunoClient implements the Client interface by wrapping luno-go.
