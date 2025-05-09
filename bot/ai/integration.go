package ai

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/luno/luno-bot/bot"
)

// CandleData represents OHLC price data
type CandleData struct {
	Timestamp time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

// BotConfig represents the trading bot configuration
type BotConfig struct {
	Pair                 string
	EntryThreshold       float64
	ExitThreshold        float64
	StakeSize            float64
	Cooldown             int
	ShortWindow          int
	LongWindow           int
	PositionSizerType    string
	KellyWinProb         float64
	KellyWinLossRatio    float64
	TWAPSlices           int
	TWAPIntervalSeconds  int
	DBPath               string
	ScanIntervalMinutes  int
	MaxPositionSize      float64
	AutoExecute          *bool
	EnableSentiment      bool
	SentimentWeight      float64
	LunarCrushAPIKey     string
	NewsAPIKey           string
	EnablePatterns       bool
	EnableOptimization   bool
	OptimizationInterval int
	EnabledPairs         []string
}

// AIController coordinates the AI engine with the bot's core components
type AIController struct {
	Engine     *AIEngine
	LunoClient *bot.LunoClient
	Store      interface{}
	Config     *BotConfig
	Strategy   bot.Strategy
	Executor   bot.Executor
	Logger     *log.Logger
}

// NewAIController creates a new AI controller
func NewAIController(
	lunoClient *bot.LunoClient,
	store interface{},
	cfg interface{},
	strategy bot.Strategy,
	executor bot.Executor,
) *AIController {
	// Create logger
	logFile, err := os.OpenFile("ai_engine.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Println("Error opening AI log file:", err)
		logFile = os.Stdout
	}
	
	aiLogger := log.New(logFile, "[AI] ", log.LstdFlags|log.Lshortfile)
	
	// Create AI engine
	engine := NewAIEngine()
	engine.SetLogger(aiLogger)
	
	// Convert config to our internal format
	var botConfig *BotConfig
	if cfg != nil {
		// Try to convert from the bot's config format
		botConfig = &BotConfig{
			Pair:                "BTCZAR", // Default
			EntryThreshold:      0.5,
			ExitThreshold:       0.5,
			StakeSize:           0.1,
			ShortWindow:         10,
			LongWindow:          30,
			AutoExecute:         new(bool),
			ScanIntervalMinutes: 5,
			MaxPositionSize:     0.1,
			EnabledPairs:        []string{"BTCZAR", "ETHZAR", "XRPZAR"},
		}
		
		// Try to extract values if they exist in the config
		// This is implementation-specific and would be customized based on your config structure
	}

	// Create controller
	controller := &AIController{
		Engine:     engine,
		LunoClient: lunoClient,
		Store:      store,
		Config:     botConfig,
		Strategy:   strategy,
		Executor:   executor,
		Logger:     aiLogger,
	}
	
	// Set up integration points
	controller.setupIntegration()
	
	return controller
}

// setupIntegration connects AI engine with bot components
func (c *AIController) setupIntegration() {
	// Set up candle data provider
	c.Engine.SetCandleDataProvider(c.fetchCandles)
	
	// Set up backtest function
	c.Engine.SetBacktestFunction(c.executeBacktest)
	
	// Set up order executor
	c.Engine.SetOrderExecutor(c.executeOrder)
	
	// Set up opportunity handler
	c.Engine.SetOpportunityHandler(c.handleOpportunity)
	
	// Configure engine with pairs from config
	var pairs []string
	if c.Config != nil && c.Config.EnabledPairs != nil && len(c.Config.EnabledPairs) > 0 {
		pairs = c.Config.EnabledPairs
	} else {
		// Default to a few common pairs if not specified
		pairs = []string{"BTCZAR", "ETHZAR", "XRPZAR", "LTCZAR"}
	}
	
	// Configure timeframes
	timeframes := []string{"15m", "1h", "4h", "1d"}
	
	// Set scan interval
	scanInterval := 5 * time.Minute
	if c.Config != nil && c.Config.ScanIntervalMinutes > 0 {
		scanInterval = time.Duration(c.Config.ScanIntervalMinutes) * time.Minute
	}
	
	// Configure engine
	c.Engine.Configure(pairs, timeframes, scanInterval)
	
	// Add parameters to optimize
	c.setupParamsToOptimize()
	
	// Log configuration
	c.Logger.Printf("AI controller initialized with %d pairs and %d timeframes", 
		len(pairs), len(timeframes))
}

// setupParamsToOptimize configures which parameters to optimize
func (c *AIController) setupParamsToOptimize() {
	optimizer := c.Engine.optimizer
	
	// RSI parameters
	optimizer.AddParameterToOptimize("rsi_period", 7, 21, 1, true, false, 14)
	optimizer.AddParameterToOptimize("rsi_overbought", 65, 85, 1, true, false, 70)
	optimizer.AddParameterToOptimize("rsi_oversold", 15, 35, 1, true, false, 30)
	
	// MACD parameters
	optimizer.AddParameterToOptimize("macd_fast_period", 8, 20, 1, true, false, 12)
	optimizer.AddParameterToOptimize("macd_slow_period", 20, 40, 1, true, false, 26)
	optimizer.AddParameterToOptimize("macd_signal_period", 6, 14, 1, true, false, 9)
	
	// Bollinger Bands parameters
	optimizer.AddParameterToOptimize("bb_period", 10, 30, 1, true, false, 20)
	optimizer.AddParameterToOptimize("bb_deviation", 1.5, 3.0, 0.1, false, false, 2.0)
	
	// Moving Averages parameters
	optimizer.AddParameterToOptimize("sma_short_period", 5, 20, 1, true, false, 10)
	optimizer.AddParameterToOptimize("sma_long_period", 20, 50, 1, true, false, 30)
	
	// Risk parameters
	optimizer.AddParameterToOptimize("risk_per_trade", 0.01, 0.05, 0.005, false, false, 0.02)
}

// Start activates the AI engine
func (c *AIController) Start() {
	c.Engine.Start()
	c.Logger.Println("AI controller started")
}

// Stop deactivates the AI engine
func (c *AIController) Stop() {
	c.Engine.Stop()
	c.Logger.Println("AI controller stopped")
}

// fetchCandles retrieves historical candle data
func (c *AIController) fetchCandles(pair string, timeframe string, limit int) ([]OHLCData, error) {
	// Convert timeframe to duration
	var duration time.Duration
	switch timeframe {
	case "1m":
		duration = time.Minute
	case "5m":
		duration = 5 * time.Minute
	case "15m":
		duration = 15 * time.Minute
	case "1h":
		duration = time.Hour
	case "4h":
		duration = 4 * time.Hour
	case "1d":
		duration = 24 * time.Hour
	default:
		return nil, fmt.Errorf("unsupported timeframe: %s", timeframe)
	}
	
	// Calculate time range
	now := time.Now()
	since := now.Add(-duration * time.Duration(limit))
	
	// Generate mock data for now
	c.Logger.Printf("Generating mock candle data for %s-%s", pair, timeframe)
	mockCandles := generateMockCandleData(pair, duration, since, now)
	
	// Convert to OHLC format
	result := make([]OHLCData, len(mockCandles))
	for i, candle := range mockCandles {
		result[i] = OHLCData{
			Timestamp: candle.Timestamp,
			Open:      candle.Open,
			High:      candle.High,
			Low:       candle.Low,
			Close:     candle.Close,
			Volume:    candle.Volume,
		}
	}
	
	return result, nil
}

// fetchCandlesFromAPI is not needed anymore as we use mock data directly

// executeBacktest runs a backtest with parameters
func (c *AIController) executeBacktest(params map[string]float64) StrategyPerformance {
	// This would run a backtest with the specified parameters
	// For now, return simulated results
	
	// Calculate a deterministic but varied result based on parameters
	profitLoss := 10.0
	
	// Adjust based on RSI parameters - prefer middle periods
	rsiPeriod := params["rsi_period"]
	optimalRSI := 14.0
	rsiAdjustment := 1.0 - (math.Abs(rsiPeriod-optimalRSI) / 7.0) * 0.2
	profitLoss *= rsiAdjustment
	
	// Adjust based on MACD parameters - prefer standard settings
	macdFast := params["macd_fast_period"]
	macdSlow := params["macd_slow_period"]
	optimalFastSlow := 14.0 // optimal gap between fast and slow
	macdAdjustment := 1.0 - (math.Abs((macdSlow-macdFast)-optimalFastSlow) / 10.0) * 0.2
	profitLoss *= macdAdjustment
	
	// Adjust based on risk - higher risk can mean higher returns but worse drawdown
	risk := params["risk_per_trade"]
	drawdown := risk * 10 * (1.0 + (math.Sin(risk*100) * 0.3))
	
	// Simulate other metrics
	winRate := 0.55 + (profitLoss / 100.0) * 0.1
	
	return StrategyPerformance{
		ProfitLoss:        profitLoss,
		SharpeRatio:       profitLoss / (drawdown * 2),
		MaxDrawdown:       drawdown,
		WinRate:           winRate,
		ProfitFactor:      1.5 + (profitLoss / 20.0),
		RecoveryFactor:    profitLoss / drawdown,
		ExpectedValue:     profitLoss / 20, // per trade
		NumTrades:         20,
		AvgHoldingPeriod:  18, // hours
		AvgProfitPerTrade: profitLoss / 20,
		CalmarRatio:       profitLoss / drawdown,
		SortinoRatio:      (profitLoss / 100) / (drawdown / 200),
		PercentProfitable: winRate * 100,
		Alpha:             0.2,
		Beta:              0.8,
	}
}

// executeOrder places an order via the Executor
func (c *AIController) executeOrder(pair string, side string, volume float64, price float64) error {
	if c.Executor == nil {
		return fmt.Errorf("executor not available")
	}
	
	// Format volume to Luno precision
	formattedVolume := fmt.Sprintf("%.6f", volume)
	
	// Log the order (simulated execution for now)
	c.Logger.Printf("AI order requested: %s %s %s @ %.2f", 
		side, formattedVolume, pair, price)
	
	// In a real implementation, would call the executor with proper parameters
	// This is a simplified version that just logs and returns success
	return nil
}

// handleOpportunity processes new trading opportunities
func (c *AIController) handleOpportunity(result *AnalysisResult) {
	c.Logger.Printf("New opportunity detected: %s %s (Score: %.2f, Confidence: %.2f)",
		result.Pair, result.Signal, result.Score, result.Confidence)
	
	// Check if auto-execution is enabled
	autoExecute := false
	if c.Config.AutoExecute != nil {
		autoExecute = *c.Config.AutoExecute
	}
	
	// Check confidence and score thresholds
	if autoExecute && result.Score >= 0.75 && result.Confidence >= 0.7 {
		// Determine maximum position size (from config or default)
		maxPositionSize := 0.1 // 10% of available funds
		if c.Config.MaxPositionSize > 0 {
			maxPositionSize = c.Config.MaxPositionSize
		}
		
		// Execute the trade
		if err := c.Engine.ExecuteTrade(result, maxPositionSize); err != nil {
			c.Logger.Printf("Error executing auto-trade: %v", err)
		} else {
			c.Logger.Printf("Auto-executed %s trade for %s (size: %.2f%%)",
				result.Signal, result.Pair, result.RecommendedSize*maxPositionSize*100)
		}
	} else {
		c.Logger.Printf("Trade opportunity below auto-execution threshold")
	}
}

// generateMockCandleData creates dummy candle data for testing
func generateMockCandleData(pair string, duration time.Duration, since, until time.Time) []CandleData {
	// Base price varies by pair
	basePrice := 100.0
	switch pair[:3] {
	case "BTC":
		basePrice = 50000.0
	case "ETH":
		basePrice = 2500.0
	case "XRP":
		basePrice = 0.5
	case "LTC":
		basePrice = 150.0
	}
	
	// Number of candles to generate
	intervals := int(until.Sub(since) / duration)
	if intervals <= 0 {
		intervals = 1
	}
	
	candles := make([]CandleData, intervals)
	
	// Generate price data
	price := basePrice
	for i := 0; i < intervals; i++ {
		timestamp := since.Add(duration * time.Duration(i))
		
		// Random walk
		change := (rand.Float64() - 0.48) * 0.02 * price
		price += change
		
		// OHLC values
		open := price
		high := price * (1 + 0.005*rand.Float64())
		low := price * (1 - 0.005*rand.Float64())
		close := price * (1 + (rand.Float64()-0.5)*0.01)
		volume := basePrice * rand.Float64() * 100
		
		candles[i] = CandleData{
			Timestamp: timestamp,
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    volume,
		}
	}
	
	return candles
}
