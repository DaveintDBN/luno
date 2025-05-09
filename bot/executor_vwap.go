package bot

import (
    "context"
    "fmt"
    "math"
    "strconv"
    "time"
    luno "github.com/luno/luno-go"
    "github.com/luno/luno-bot/storage"
)

// VWAPExecutor slices large orders into smaller chunks based on volume-weighted logic.
type VWAPExecutor struct {
    Inner    Executor
    Client   Client
    Slices   int
    Interval time.Duration
    Store    *storage.SQLiteStore
}

// NewVWAPExecutor constructs a VWAP executor that distributes execution over given slices and interval.
func NewVWAPExecutor(inner Executor, client Client, slices int, interval time.Duration, store *storage.SQLiteStore) *VWAPExecutor {
    if slices <= 1 {
        slices = 1
    }
    return &VWAPExecutor{Inner: inner, Client: client, Slices: slices, Interval: interval, Store: store}
}

// Execute slices execution based on VWAP; currently evenly weighted as placeholder.
func (v *VWAPExecutor) Execute(ctx context.Context, sig Signal, md MarketData, cfg Config) error {
    if sig == SignalNone {
        return nil
    }
    fmt.Printf("VWAPExecutor: executing %d slices every %s based on VWAP\n", v.Slices, v.Interval)
    // Determine slice weights based on VWAP source
    var weights []float64
    switch cfg.VWAPSource {
    case "historical":
        weights = v.computeHistoricalWeights(ctx, cfg)
    case "orderbook":
        weights = v.computeOrderbookWeights(ctx, cfg, sig)
    case "hybrid":
        hist := v.computeHistoricalWeights(ctx, cfg)
        book := v.computeOrderbookWeights(ctx, cfg, sig)
        weights = make([]float64, v.Slices)
        for i := 0; i < v.Slices; i++ {
            weights[i] = cfg.VWAPHybridWeight*hist[i] + (1-cfg.VWAPHybridWeight)*book[i]
        }
    default:
        weights = make([]float64, v.Slices)
        for i := range weights {
            weights[i] = 1.0 / float64(v.Slices)
        }
    }
    // Persist trade record
    price := (md.Bid + md.Ask) / 2
    var tradeID int64
    if v.Store != nil {
        var side string
        if sig == SignalBuy {
            side = "buy"
        } else if sig == SignalSell {
            side = "sell"
        } else {
            side = "none"
        }
        id, err := v.Store.SaveTrade(md.Timestamp, cfg.Pair, side, price, cfg.StakeSize)
        if err != nil {
            return fmt.Errorf("save trade: %w", err)
        }
        tradeID = id
    }
    for i := 0; i < v.Slices; i++ {
        sliceCfg := cfg
        sliceCfg.StakeSize = cfg.StakeSize * weights[i]
        if err := v.Inner.Execute(ctx, sig, md, sliceCfg); err != nil {
            return err
        }
        // Persist slice
        if v.Store != nil {
            if err := v.Store.SaveSlice(tradeID, i, sliceCfg.StakeSize, weights[i]); err != nil {
                return fmt.Errorf("save slice: %w", err)
            }
        }
        if i < v.Slices-1 {
            select {
            case <-ctx.Done():
                return ctx.Err()
            case <-time.After(v.Interval):
            }
        }
    }
    return nil
}

// CancelAll delegates cancellation to inner executor.
func (v *VWAPExecutor) CancelAll(ctx context.Context) error {
    return v.Inner.CancelAll(ctx)
}

// computeHistoricalWeights calculates weights from historical volume data.
func (v *VWAPExecutor) computeHistoricalWeights(ctx context.Context, cfg Config) []float64 {
    since := time.Now().Add(-time.Duration(cfg.VWAPHistoryWindowMinutes) * time.Minute)
    req := &luno.GetCandlesRequest{Pair: cfg.Pair, Duration: 60, Since: luno.Time(since)}
    resp, err := v.Client.GetCandles(ctx, req)
    weights := make([]float64, v.Slices)
    if err != nil {
        for i := range weights {
            weights[i] = 1.0 / float64(v.Slices)
        }
        return weights
    }
    candles := resp.Candles
    n := len(candles)
    if n == 0 {
        for i := range weights {
            weights[i] = 1.0 / float64(v.Slices)
        }
        return weights
    }
    vols := make([]float64, n)
    var totalVol float64
    for i, c := range candles {
        vol, err := strconv.ParseFloat(c.Volume.String(), 64)
        if err != nil {
            vol = 0
        }
        vols[i] = vol
        totalVol += vol
    }
    bucket := float64(n) / float64(v.Slices)
    for i := 0; i < v.Slices; i++ {
        start := int(math.Floor(float64(i) * bucket))
        end := int(math.Floor(float64(i+1) * bucket))
        if i == v.Slices-1 {
            end = n
        }
        var sumVol float64
        for j := start; j < end; j++ {
            sumVol += vols[j]
        }
        if totalVol > 0 {
            weights[i] = sumVol / totalVol
        } else {
            weights[i] = 1.0 / float64(v.Slices)
        }
    }
    return weights
}

// computeOrderbookWeights calculates weights from orderbook depth data.
func (v *VWAPExecutor) computeOrderbookWeights(ctx context.Context, cfg Config, sig Signal) []float64 {
    req := &luno.GetOrderBookRequest{Pair: cfg.Pair}
    resp, err := v.Client.GetOrderBook(ctx, req)
    weights := make([]float64, v.Slices)
    if err != nil {
        for i := range weights {
            weights[i] = 1.0 / float64(v.Slices)
        }
        return weights
    }
    items := resp.Asks
    if sig == SignalSell {
        items = resp.Bids
    }
    depth := len(items)
    if depth == 0 {
        for i := range weights {
            weights[i] = 1.0 / float64(v.Slices)
        }
        return weights
    }
    vols := make([]float64, depth)
    var totalVol float64
    for i, item := range items {
        vol, err := strconv.ParseFloat(item.Volume.String(), 64)
        if err != nil {
            vol = 0
        }
        vols[i] = vol
        totalVol += vol
    }
    bucket := float64(depth) / float64(v.Slices)
    for i := 0; i < v.Slices; i++ {
        start := int(math.Floor(float64(i) * bucket))
        end := int(math.Floor(float64(i+1) * bucket))
        if i == v.Slices-1 {
            end = depth
        }
        var sumVol float64
        for j := start; j < end; j++ {
            sumVol += vols[j]
        }
        if totalVol > 0 {
            weights[i] = sumVol / totalVol
        } else {
            weights[i] = 1.0 / float64(v.Slices)
        }
    }
    return weights
}
