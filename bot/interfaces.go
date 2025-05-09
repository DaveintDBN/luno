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
	// GetBalances retrieves account balances from Luno API
	GetBalances(ctx context.Context, req *luno.GetBalancesRequest) (*luno.GetBalancesResponse, error)
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
	Pair             string        // e.g. "XBTZAR"
	EntryThreshold   float64       // threshold to enter
	ExitThreshold    float64       // threshold to exit
	StakeSize        float64       // amount per trade
	Cooldown         time.Duration // min time between trades
	PositionLimit    float64       // max allowed position size
	MaxDrawdown      float64       // max permitted drawdown
	ShortWindow      int           // SMA short window
	LongWindow       int           // SMA long window
	BaseAccountId    int64         // base currency account ID for trades
	CounterAccountId int64         // counter currency account ID for trades
	// RSI indicator parameters
	RSIPeriod       int           // number of periods for RSI
	RSIOverBought   float64       // RSI level above which to sell
	RSIOverSold     float64       // RSI level below which to buy
	// MACD indicator parameters
	MACDFastPeriod   int          // fast EMA period for MACD
	MACDSlowPeriod   int          // slow EMA period for MACD
	MACDSignalPeriod int          // signal line EMA period for MACD
	// Bollinger Bands parameters
	BBPeriod         int          // window size for Bollinger Bands
	BBMultiplier     float64      // stddev multiplier for Bollinger Bands
	// Risk & execution parameters
	InitialEquity       float64      // starting capital for sizing
	PositionSizerType   string       // "fixed" or "kelly"
	KellyWinProb        float64      // win probability for Kelly sizing
	KellyWinLossRatio   float64      // average win/loss ratio for Kelly sizing
	TWAPSlices          int          // number of slices for TWAP execution
	TWAPIntervalSeconds int          // seconds between TWAP slices
	// VWAP parameters
	VWAPSource               string      // VWAP source: "historical", "orderbook", or "hybrid"
	VWAPHistoryWindowMinutes int         // window in minutes for historical VWAP
	VWAPOrderbookDepthLevels int         // depth levels for orderbook VWAP
	VWAPHybridWeight         float64     // weight factor for hybrid VWAP combination
}

// MarketData packages latest market metrics.
type MarketData struct {
	Bid       float64
	Ask       float64
	Timestamp time.Time
}

// LunoClient implements the Client interface by wrapping luno-go.
