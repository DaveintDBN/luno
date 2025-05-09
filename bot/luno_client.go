package bot

import (
	"context"

	"github.com/luno/luno-go"
)

// LunoClient implements the Client interface by wrapping luno-go.
type LunoClient struct {
	cli *luno.Client
}

// NewLunoClient constructs a new LunoClient.
func NewLunoClient() *LunoClient {
	return &LunoClient{cli: luno.NewClient()}
}

// SetAuth configures API credentials.
func (c *LunoClient) SetAuth(id, secret string) error {
	return c.cli.SetAuth(id, secret)
}

// GetTickers fetches market tickers.
func (c *LunoClient) GetTickers(ctx context.Context, req *luno.GetTickersRequest) (*luno.GetTickersResponse, error) {
	return c.cli.GetTickers(ctx, req)
}

// GetOrderBook retrieves the order book.
func (c *LunoClient) GetOrderBook(ctx context.Context, req *luno.GetOrderBookRequest) (*luno.GetOrderBookResponse, error) {
	return c.cli.GetOrderBook(ctx, req)
}

// PostLimitOrder places a new limit order.
func (c *LunoClient) PostLimitOrder(ctx context.Context, req *luno.PostLimitOrderRequest) (*luno.PostLimitOrderResponse, error) {
	return c.cli.PostLimitOrder(ctx, req)
}

// ListTrades fetches recent trades for backtesting.
func (c *LunoClient) ListTrades(ctx context.Context, req *luno.ListTradesRequest) (*luno.ListTradesResponse, error) {
	return c.cli.ListTrades(ctx, req)
}

// GetCandles fetches historical candles for backtesting.
func (c *LunoClient) GetCandles(ctx context.Context, req *luno.GetCandlesRequest) (*luno.GetCandlesResponse, error) {
	return c.cli.GetCandles(ctx, req)
}

// GetBalances fetches account balances
func (c *LunoClient) GetBalances(ctx context.Context, req *luno.GetBalancesRequest) (*luno.GetBalancesResponse, error) {
	return c.cli.GetBalances(ctx, req)
}
