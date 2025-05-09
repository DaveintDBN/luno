//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	luno "github.com/luno/luno-go"
	"github.com/luno/luno-bot/bot"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: backtester <apiKeyID> <apiKeySecret> <pair>")
		return
	}
	apiKeyID, apiKeySecret, pair := os.Args[1], os.Args[2], os.Args[3]

	// Setup client
	lc := bot.NewLunoClient()
	lc.SetAuth(apiKeyID, apiKeySecret)

	// Fetch recent trades
	req := &luno.ListTradesRequest{Pair: pair}
	res, err := lc.ListTrades(context.Background(), req)
	if err != nil {
		fmt.Println("Error fetching trades:", err)
		return
	}

	// Reverse trades for chronological order
	for i, j := 0, len(res.Trades)-1; i < j; i, j = i+1, j-1 {
		res.Trades[i], res.Trades[j] = res.Trades[j], res.Trades[i]
	}

	prices := make([]float64, len(res.Trades))
	for i, t := range res.Trades {
		p := t.Price.Float64()
		prices[i] = p
	}

	// Simulate backtest trades with PnL and metrics
	strat := bot.NewSMAStrategy(5, 10)
	var cfg bot.Config
	cfg.EntryThreshold = 0
	cfg.ExitThreshold = 0
	cfg.StakeSize = 1

	inPosition := false
	var entryPrice float64
	var totalTrades, wins, losses int
	var totalPnL float64
	for i, t := range res.Trades {
		price := prices[i]
		md := bot.MarketData{Bid: price, Ask: price, Timestamp: time.Time(t.Timestamp)}
		sig := strat.Next(md, cfg)
		if sig == bot.SignalBuy && !inPosition {
			entryPrice = price
			inPosition = true
		} else if sig == bot.SignalSell && inPosition {
			profit := (price - entryPrice) * cfg.StakeSize
			totalPnL += profit
			totalTrades++
			if profit > 0 {
				wins++
			} else {
				losses++
			}
			inPosition = false
		}
	}

	// Print summary metrics
	fmt.Printf("Backtest Summary: Trades=%d, Wins=%d, Losses=%d, Win rate=%.2f%%, Total PnL=%.2f\n",
		totalTrades, wins, losses, float64(wins)/float64(totalTrades)*100, totalPnL)
	if totalTrades > 0 {
		avg := totalPnL / float64(totalTrades)
		fmt.Printf("Avg PnL per trade: %.2f\n", avg)
	}
}
