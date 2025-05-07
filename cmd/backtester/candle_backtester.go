package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	luno "github.com/luno/luno-go"
	"github.com/luno/luno-bot/bot"
)

func main() {
	apiKeyID := flag.String("api_key_id", "", "Luno API key ID")
	apiKeySecret := flag.String("api_key_secret", "", "Luno API key secret")
	pair := flag.String("pair", "", "Market pair, e.g. XBTZAR")
	sinceMin := flag.Int("since_minutes", 60, "Minutes back to fetch 1m candles")
	shortW := flag.Int("short", 5, "Short SMA window")
	longW := flag.Int("long", 10, "Long SMA window")
	// Fee rate per trade side
	feeRate := flag.Float64("fee_rate", 0.001, "Trading fee rate per trade side (e.g. 0.001)")
	flag.Parse()

	if *apiKeyID == "" || *apiKeySecret == "" || *pair == "" {
		fmt.Println("Usage: candle_backtester --api_key_id <id> --api_key_secret <secret> --pair <pair> [--since_minutes <min>] [--short <n>] [--long <n>] [--fee_rate <rate>]")
		return
	}

	// Setup client
	lc := bot.NewLunoClient()
	if err := lc.SetAuth(*apiKeyID, *apiKeySecret); err != nil {
		fmt.Println("Error setting auth:", err)
		return
	}
	ctx := context.Background()

	// Fetch 1m candles
	since := time.Now().Add(-time.Duration(*sinceMin) * time.Minute)
	req := &luno.GetCandlesRequest{
		Pair:     *pair,
		Duration: 60,
		Since:    luno.Time(since),
	}
	res, err := lc.GetCandles(ctx, req)
	if err != nil {
		fmt.Println("Error fetching candles:", err)
		return
	}

	// Extract closes and timestamps
	n := len(res.Candles)
	closes := make([]float64, n)
	times := make([]time.Time, n)
	for i, c := range res.Candles {
		p := c.Close.Float64()
		closes[i] = p
		times[i] = time.Time(c.Timestamp)
	}

	// Backtest SMA
	strat := bot.NewSMAStrategy(*shortW, *longW)
	var cfg bot.Config
	cfg.EntryThreshold = 0
	cfg.ExitThreshold = 0
	cfg.StakeSize = 1

	inPos := false
	var entry float64
	var trades, wins, losses int
	var pnlTotal float64
	for i := 0; i < n; i++ {
		md := bot.MarketData{Bid: closes[i], Ask: closes[i], Timestamp: times[i]}
		sig := strat.Next(md, cfg)
		if sig == bot.SignalBuy && !inPos {
			entry = closes[i]
			inPos = true
		} else if sig == bot.SignalSell && inPos {
			gross := (closes[i] - entry) * cfg.StakeSize
			fee := (*feeRate) * (closes[i] + entry) * cfg.StakeSize
			profit := gross - fee
			pnlTotal += profit
			trades++
			if profit > 0 {
				wins++
			} else {
				losses++
			}
			inPos = false
		}
	}

	// Summary
	fmt.Printf("Candle Backtest (%dm): Trades=%d, Wins=%d, Losses=%d, Win rate=%.2f%%, Total PnL=%.2f\n",
		*sinceMin, trades, wins, losses, float64(wins)/float64(trades)*100, pnlTotal)
	if trades > 0 {
		fmt.Printf("Avg PnL per trade: %.2f\n", pnlTotal/float64(trades))
	}
}
