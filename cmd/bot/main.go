package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	"github.com/joho/godotenv"
	"github.com/luno/luno-bot/bot"
	"github.com/luno/luno-bot/bot/ai"
	api "github.com/luno/luno-bot/cmd/bot/api"
	"github.com/luno/luno-bot/config"
	"github.com/luno/luno-bot/storage"
	luno "github.com/luno/luno-go"
)

func main() {
	// Load .env from root or parent dirs
	for _, envFile := range []string{".env", "../.env", "../../.env"} {
		if err := godotenv.Load(envFile); err == nil {
			break
		}
	}
	// Use env vars as default flag values
	apiKeyID := flag.String("api_key_id", "", "Luno API key ID")
	apiKeySecret := flag.String("api_key_secret", "", "Luno API key secret")
	configPath := flag.String("config", "../../config/config.json", "Path to config file")
	flag.Parse()

	// Fallback to environment variables if flags not provided
	if *apiKeyID == "" {
		*apiKeyID = os.Getenv("API_KEY_ID")
	}
	if *apiKeySecret == "" {
		*apiKeySecret = os.Getenv("API_KEY_SECRET")
	}

	fmt.Printf("Using API_KEY_ID: %s\n", *apiKeyID)
	// Require both credentials
	if *apiKeyID == "" || *apiKeySecret == "" {
		fmt.Println("Missing Luno API credentials (set API_KEY_ID and API_KEY_SECRET in .env or via flags)")
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
	strat := bot.NewMultiTimeframeStrategy(cfg)
	// Setup position sizing and TWAP executor chain
	var sizer bot.PositionSizer
	switch cfg.PositionSizerType {
	case "kelly":
		sizer = &bot.KellySizer{WinProb: cfg.KellyWinProb, WinLoss: cfg.KellyWinLossRatio}
	default:
		sizer = &bot.FixedSizer{}
	}
	simInner := bot.NewSimulatedExecutor()
	simSizing := bot.NewSizingExecutor(simInner, sizer)
	// Setup VWAP executor for simulation
	// Initialize SQLite store
	sqlStore, err := storage.NewSQLiteStore(cfg.DBPath)
	if err != nil {
		fmt.Println("Error opening SQLite DB:", err)
		return
	}
	defer sqlStore.Close()
	simVWAP := bot.NewVWAPExecutor(simSizing, lc, cfg.TWAPSlices, time.Duration(cfg.TWAPIntervalSeconds)*time.Second, sqlStore)
	// Initialize live VWAP executor
	liveInner := bot.NewLunoExecutor(lc)
	liveSizing := bot.NewSizingExecutor(liveInner, sizer)
	var liveExec bot.Executor = bot.NewVWAPExecutor(liveSizing, lc, cfg.TWAPSlices, time.Duration(cfg.TWAPIntervalSeconds)*time.Second, sqlStore)
	// Wrap live executor with logging
	actFile, err := os.OpenFile("live_activity.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error opening live activity log:", err)
		return
	}
	defer actFile.Close()
	errFile, err := os.OpenFile("live_errors.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error opening live error log:", err)
		return
	}
	defer errFile.Close()
	actLogger := log.New(actFile, "", log.LstdFlags)
	errLogger := log.New(errFile, "", log.LstdFlags|log.Lshortfile)
	liveExec = bot.NewLoggingExecutor(liveExec, actLogger, errLogger)
	
	// Initialize AI controller
	aiController := ai.NewAIController(lc, sqlStore, cfg, strat, liveExec)
	aiController.Start()
	
	// Launch REST API server with simulation and live execution
	r := api.SetupRouter(store, lc, strat, simVWAP, liveExec)
	
	// Register AI routes
	aiGroup := r.Group("/api/ai")
	ai.RegisterAIRoutes(aiGroup, aiController.Engine)
	
	fmt.Println("AI enhancements activated")
	fmt.Println("Starting server on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		fmt.Println("Server error:", err)
	}
}
