package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/luno/luno-bot/bot"
	"github.com/luno/luno-bot/config"
	luno "github.com/luno/luno-go"
	api "github.com/luno/luno-bot/cmd/bot/api"
)

func main() {
	apiKeyID := flag.String("api_key_id", "", "Luno API key ID")
	apiKeySecret := flag.String("api_key_secret", "", "Luno API key secret")
	configPath := flag.String("config", "config/config.json", "Path to config file")
	flag.Parse()

	if *apiKeyID == "" || *apiKeySecret == "" {
		fmt.Println("Usage: cmd/bot --api_key_id <id> --api_key_secret <secret> [--config <path>]")
		return
	}

	store := config.NewStateStore(*configPath)
	cfg, err := store.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	fmt.Printf("Loaded config: %+v\n", cfg)
	fmt.Println("Phase 2 scaffolding: interfaces defined and config loaded.")

	// Ensure SMA windows default if not set in config
	if cfg.ShortWindow <= 0 || cfg.LongWindow <= 0 {
		cfg.ShortWindow = 5
		cfg.LongWindow = 10
		fmt.Println("Defaulting SMA windows to short=5, long=10")
	}

	// Initialize Luno client and fetch order book
	lc := bot.NewLunoClient()
	if err := lc.SetAuth(*apiKeyID, *apiKeySecret); err != nil {
		fmt.Println("Error setting auth:", err)
		return
	}
	ctx := context.Background()
	ob, err := lc.GetOrderBook(ctx, &luno.GetOrderBookRequest{Pair: cfg.Pair})
	if err != nil {
		fmt.Println("Error fetching order book:", err)
		return
	}
	fmt.Printf("Order Book (%s): Bids: %+v\nAsks: %+v\n", cfg.Pair, ob.Bids, ob.Asks)

	// Initialize strategy and simulated executor
	strat := bot.NewSMAStrategy(cfg.ShortWindow, cfg.LongWindow)
	simExec := bot.NewSimulatedExecutor()
	// Initialize live executor
	liveExec := bot.NewLunoExecutor(lc)
	// Launch REST API server with simulation and live execution
	r := api.SetupRouter(store, lc, strat, simExec, liveExec)
	fmt.Println("Starting server on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		fmt.Println("Server error:", err)
	}
}
