package ai

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"
)

// AnalysisResult contains AI-enhanced market analysis
type AnalysisResult struct {
	Pair             string
	Timeframe        string
	Score            float64          // 0-1 opportunity score 
	Signal           string           // "buy", "sell", "hold"
	Confidence       float64          // 0-1 confidence level
	MLFeatures       []SignalFeature  // ML model features
	PatternSignals   []PatternSignal  // Detected patterns
	SentimentData    *SentimentData   // Sentiment analysis
	PredictedMove    float64          // Expected price movement (%)
	RecommendedSize  float64          // Suggested position size (0-1)
	Timestamp        time.Time
	AnalysisDuration time.Duration
}

// AIEngine coordinates all AI components
type AIEngine struct {
	// Core AI components
	mlModel           *MLModel
	sentimentAnalyzer *SentimentAnalyzer
	patternRecognizer *PatternRecognizer
	optimizer         *Optimizer

	// Configuration
	enabledComponents map[string]bool
	pairs             []string
	timeframes        []string
	scanInterval      time.Duration
	autoOptimize      bool
	optimizeInterval  time.Duration

	// State
	running             bool
	lastScan            time.Time
	lastOptimization    time.Time
	scanResults         map[string]map[string]*AnalysisResult // pair -> timeframe -> result
	runningAvgScore     map[string]float64                   // pair -> running average score
	scanLock            sync.RWMutex
	
	// Integration points
	onNewOpportunity    func(result *AnalysisResult)
	fetchCandles        func(pair string, timeframe string, limit int) ([]OHLCData, error)
	executeBacktest     func(params map[string]float64) StrategyPerformance
	executeOrder        func(pair string, side string, volume float64, price float64) error

	// Logging
	logger              *log.Logger
}

// NewAIEngine creates a new AI analysis engine
func NewAIEngine() *AIEngine {
	engine := &AIEngine{
		mlModel:           NewMLModel(),
		sentimentAnalyzer: NewSentimentAnalyzer(),
		patternRecognizer: NewPatternRecognizer(),
		enabledComponents: map[string]bool{
			"ml":        true,
			"sentiment": true,
			"patterns":  true,
			"optimize":  true,
		},
		pairs:              []string{},
		timeframes:         []string{"1h", "4h", "1d"},
		scanInterval:       5 * time.Minute,
		autoOptimize:       true,
		optimizeInterval:   24 * time.Hour,
		running:            false,
		scanResults:        make(map[string]map[string]*AnalysisResult),
		runningAvgScore:    make(map[string]float64),
	}
	
	// Initialize optimizer with a dummy backtest function
	// (will be replaced with real backtest integration)
	engine.optimizer = NewOptimizer(func(params map[string]float64) StrategyPerformance {
		if engine.executeBacktest != nil {
			return engine.executeBacktest(params)
		}
		return StrategyPerformance{}
	})
	
	return engine
}

// Configure sets up the AI engine
func (e *AIEngine) Configure(pairs []string, timeframes []string, scanInterval time.Duration) {
	e.pairs = pairs
	e.timeframes = timeframes
	e.scanInterval = scanInterval
	
	// Configure optimizer
	e.optimizer.SetPairs(pairs)
	e.optimizer.SetTimeframes(timeframes)
}

// SetSentimentAPIKeys configures API keys for sentiment data sources
func (e *AIEngine) SetSentimentAPIKeys(lunarCrushKey, newsAPIKey string) {
	if lunarCrushKey != "" {
		e.sentimentAnalyzer.SetAPIKey(SourceLunarCrush, lunarCrushKey)
	}
	if newsAPIKey != "" {
		e.sentimentAnalyzer.SetAPIKey(SourceNewsAPI, newsAPIKey)
	}
}

// SetCandleDataProvider sets a function to fetch candle data
func (e *AIEngine) SetCandleDataProvider(provider func(pair string, timeframe string, limit int) ([]OHLCData, error)) {
	e.fetchCandles = provider
}

// SetBacktestFunction sets a function to execute backtests
func (e *AIEngine) SetBacktestFunction(backtest func(params map[string]float64) StrategyPerformance) {
	e.executeBacktest = backtest
	e.optimizer = NewOptimizer(backtest)
}

// SetOrderExecutor sets a function to execute trades
func (e *AIEngine) SetOrderExecutor(executor func(pair string, side string, volume float64, price float64) error) {
	e.executeOrder = executor
}

// SetOpportunityHandler sets a callback for new opportunities
func (e *AIEngine) SetOpportunityHandler(handler func(result *AnalysisResult)) {
	e.onNewOpportunity = handler
}

// SetLogger configures logging
func (e *AIEngine) SetLogger(logger *log.Logger) {
	e.logger = logger
}

// ToggleComponent enables or disables an AI component
func (e *AIEngine) ToggleComponent(component string, enabled bool) {
	e.enabledComponents[component] = enabled
}

// Start begins continuous market analysis
func (e *AIEngine) Start() {
	if e.running {
		return
	}
	
	e.running = true
	
	// Begin tracking sentiment for all pairs
	if e.enabledComponents["sentiment"] {
		e.sentimentAnalyzer.StartSentimentTracking(e.pairs)
	}
	
	// Begin continuous scanning for opportunities
	go e.continuousScan()
	
	// Begin auto-optimization if enabled
	if e.enabledComponents["optimize"] && e.autoOptimize {
		go e.scheduleOptimization()
	}
	
	// Schedule model optimization
	if e.enabledComponents["ml"] {
		e.mlModel.ScheduleOptimization(12 * time.Hour)
	}
	
	e.log("AI engine started")
}

// Stop halts continuous analysis
func (e *AIEngine) Stop() {
	e.running = false
	e.log("AI engine stopped")
}

// continuousScan periodically scans for trading opportunities
func (e *AIEngine) continuousScan() {
	ticker := time.NewTicker(e.scanInterval)
	defer ticker.Stop()
	
	// Do an initial scan
	e.ScanAllMarkets()
	
	for e.running {
		select {
		case <-ticker.C:
			e.ScanAllMarkets()
		}
	}
}

// scheduleOptimization periodically optimizes trading parameters
func (e *AIEngine) scheduleOptimization() {
	ticker := time.NewTicker(e.optimizeInterval)
	defer ticker.Stop()
	
	for e.running {
		select {
		case <-ticker.C:
			if e.enabledComponents["optimize"] && e.autoOptimize {
				e.log("Starting scheduled parameter optimization...")
				
				// Run walk-forward optimization
				window := 30 * 24 * time.Hour // 30 days
				step := 7 * 24 * time.Hour    // 7 days
				result := e.optimizer.WalkForwardOptimization(window, step, 100)
				
				e.lastOptimization = time.Now()
				
				e.log(fmt.Sprintf("Optimization complete. Best params: %v", result.BestParams))
			}
		}
	}
}

// ScanAllMarkets analyzes all configured markets
func (e *AIEngine) ScanAllMarkets() {
	var wg sync.WaitGroup
	
	e.log(fmt.Sprintf("Scanning %d pairs across %d timeframes...", len(e.pairs), len(e.timeframes)))
	
	// Process each pair in a separate goroutine
	for _, pair := range e.pairs {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			for _, tf := range e.timeframes {
				result := e.AnalyzeMarket(p, tf)
				
				// Store result
				e.scanLock.Lock()
				if _, exists := e.scanResults[p]; !exists {
					e.scanResults[p] = make(map[string]*AnalysisResult)
				}
				e.scanResults[p][tf] = result
				
				// Update running average score
				alpha := 0.1 // Weight for new score in EMA
				if oldAvg, exists := e.runningAvgScore[p]; exists {
					e.runningAvgScore[p] = oldAvg*(1-alpha) + result.Score*alpha
				} else {
					e.runningAvgScore[p] = result.Score
				}
				e.scanLock.Unlock()
				
				// Trigger callback if significant opportunity detected
				if result.Score > 0.7 && result.Confidence > 0.65 && e.onNewOpportunity != nil {
					e.onNewOpportunity(result)
				}
			}
		}(pair)
	}
	
	wg.Wait()
	e.lastScan = time.Now()
	e.log("Market scan completed")
}

// AnalyzeMarket performs full AI analysis on one market
func (e *AIEngine) AnalyzeMarket(pair string, timeframe string) *AnalysisResult {
	startTime := time.Now()
	
	// Prepare result container
	result := &AnalysisResult{
		Pair:             pair,
		Timeframe:        timeframe,
		Timestamp:        startTime,
		MLFeatures:       []SignalFeature{},
		PatternSignals:   []PatternSignal{},
	}
	
	// 1. Fetch candle data
	var candles []OHLCData
	var err error
	
	if e.fetchCandles != nil {
		candles, err = e.fetchCandles(pair, timeframe, 200) // Get 200 candles
		if err != nil {
			e.log(fmt.Sprintf("Error fetching candles for %s-%s: %v", pair, timeframe, err))
			return result
		}
	} else {
		// Mock data for demo
		candles = generateMockCandles(pair, timeframe, 200)
	}
	
	// 2. Collect features from various sources
	var allFeatures []SignalFeature
	
	// 2.1 Technical analysis features
	taFeatures := generateTAFeatures(candles)
	allFeatures = append(allFeatures, taFeatures...)
	
	// 2.2 Pattern recognition
	if e.enabledComponents["patterns"] {
		patterns := e.patternRecognizer.AnalyzePatterns(pair, Timeframe(timeframe), candles)
		patternFeatures := e.patternRecognizer.PatternToSignalFeatures(patterns)
		allFeatures = append(allFeatures, patternFeatures...)
		result.PatternSignals = patterns
	}
	
	// 2.3 Sentiment analysis
	if e.enabledComponents["sentiment"] {
		// Extract base asset from pair (e.g., "BTCZAR" -> "BTC")
		baseAsset := extractBaseAsset(pair)
		sentFeatures := e.sentimentAnalyzer.SentimentToSignalFeature(baseAsset)
		allFeatures = append(allFeatures, sentFeatures...)
		result.SentimentData = e.sentimentAnalyzer.GetSentiment(baseAsset)
	}
	
	// Store all features
	result.MLFeatures = allFeatures
	
	// 3. Score opportunity using ML model
	if e.enabledComponents["ml"] && len(allFeatures) > 0 {
		opportunityScore := e.mlModel.ScoreOpportunity(pair, allFeatures)
		
		result.Score = opportunityScore.Score
		result.Signal = opportunityScore.RecommendedAction
		result.Confidence = opportunityScore.Confidence
		result.PredictedMove = opportunityScore.PredictedMovement
		
		// Calculate recommended position size (0 to 1)
		// Higher for strong signals with high confidence
		result.RecommendedSize = calculatePositionSize(
			opportunityScore.Score,
			opportunityScore.Confidence,
			math.Abs(opportunityScore.PredictedMovement),
		)
		
		// Add to ML model observation history
		for _, feature := range allFeatures {
			e.mlModel.AddFeatureObservation(feature.Name, feature.Value)
		}
	}
	
	// Record analysis duration
	result.AnalysisDuration = time.Since(startTime)
	
	return result
}

// GetOpportunityRanking returns top opportunities sorted by score
func (e *AIEngine) GetOpportunityRanking(minScore float64, limit int) []*AnalysisResult {
	e.scanLock.RLock()
	defer e.scanLock.RUnlock()
	
	// Flatten results into a slice
	var allResults []*AnalysisResult
	for _, timeframeMap := range e.scanResults {
		for _, result := range timeframeMap {
			if result.Score >= minScore {
				allResults = append(allResults, result)
			}
		}
	}
	
	// Sort by score
	sort.Slice(allResults, func(i, j int) bool {
		// Primary sort by score
		if allResults[i].Score != allResults[j].Score {
			return allResults[i].Score > allResults[j].Score
		}
		// Secondary sort by confidence
		return allResults[i].Confidence > allResults[j].Confidence
	})
	
	// Apply limit
	if limit > 0 && len(allResults) > limit {
		allResults = allResults[:limit]
	}
	
	return allResults
}

// GetAnalysisForPair returns all analyses for a specific pair
func (e *AIEngine) GetAnalysisForPair(pair string) map[string]*AnalysisResult {
	e.scanLock.RLock()
	defer e.scanLock.RUnlock()
	
	if results, exists := e.scanResults[pair]; exists {
		// Create a copy to avoid race conditions
		resultsCopy := make(map[string]*AnalysisResult)
		for tf, result := range results {
			resultsCopy[tf] = result
		}
		return resultsCopy
	}
	
	return nil
}

// GetRunningAverageScores returns EMA scores for all pairs
func (e *AIEngine) GetRunningAverageScores() map[string]float64 {
	e.scanLock.RLock()
	defer e.scanLock.RUnlock()
	
	// Create a copy to avoid race conditions
	scoresCopy := make(map[string]float64)
	for pair, score := range e.runningAvgScore {
		scoresCopy[pair] = score
	}
	
	return scoresCopy
}

// ExecuteTrade places an order based on AI analysis
func (e *AIEngine) ExecuteTrade(result *AnalysisResult, maxPositionSize float64) error {
	if e.executeOrder == nil {
		return fmt.Errorf("no order executor configured")
	}
	
	// Don't trade if signal is "hold"
	if result.Signal == "hold" {
		return nil
	}
	
	// Calculate position size
	positionSize := result.RecommendedSize * maxPositionSize
	
	// Determine side (buy/sell)
	side := "buy"
	if result.Signal == "sell" {
		side = "sell"
	}
	
	// Log the trade
	e.log(fmt.Sprintf("AI-generated trade: %s %s, size: %.4f, score: %.2f, confidence: %.2f", 
		side, result.Pair, positionSize, result.Score, result.Confidence))
	
	// Execute the order (price = 0 for market order)
	return e.executeOrder(result.Pair, side, positionSize, 0)
}

// GetLastScanTime returns when the last market scan was performed
func (e *AIEngine) GetLastScanTime() time.Time {
	return e.lastScan
}

// GetLastOptimizationTime returns when parameters were last optimized
func (e *AIEngine) GetLastOptimizationTime() time.Time {
	return e.lastOptimization
}

// GetModelSummary returns ML model status
func (e *AIEngine) GetModelSummary() map[string]interface{} {
	if e.enabledComponents["ml"] {
		return e.mlModel.GetModelSummary()
	}
	return nil
}

// RunSingleBacktest executes one backtest with current parameters
func (e *AIEngine) RunSingleBacktest() (*StrategyPerformance, error) {
	if e.executeBacktest == nil {
		return nil, fmt.Errorf("backtest function not configured")
	}
	
	// Get current parameters from optimizer
	params := e.optimizer.GetBestParameters()
	if len(params) == 0 {
		// Use default parameters
		params = getDefaultParameters()
	}
	
	// Run backtest
	performance := e.executeBacktest(params)
	
	return &performance, nil
}

// Helper logging function
func (e *AIEngine) log(message string) {
	if e.logger != nil {
		e.logger.Println(message)
	} else {
		fmt.Println(message)
	}
}

// Helper function to extract base asset from pair
func extractBaseAsset(pair string) string {
	// Simplified extraction for common quote currencies
	for _, quote := range []string{"ZAR", "USD", "USDT", "BTC", "ETH"} {
		if len(pair) > len(quote) && pair[len(pair)-len(quote):] == quote {
			return pair[:len(pair)-len(quote)]
		}
	}
	
	// Default fallback: take first 3 characters
	if len(pair) >= 3 {
		return pair[:3]
	}
	
	return pair
}

// Helper to calculate position size
func calculatePositionSize(score, confidence, expectedMove float64) float64 {
	// Base size on score and confidence
	baseSize := score * math.Pow(confidence, 0.5)
	
	// Adjust for expected move magnitude (higher for larger expected moves)
	moveMultiplier := math.Min(1.5, 1.0 + expectedMove/5.0)
	
	// Combine factors and cap at 1.0
	return math.Min(1.0, baseSize * moveMultiplier)
}

// Helper to generate default parameters
func getDefaultParameters() map[string]float64 {
	return map[string]float64{
		"rsi_period":         14,
		"rsi_overbought":     70,
		"rsi_oversold":       30,
		"macd_fast_period":   12,
		"macd_slow_period":   26,
		"macd_signal_period": 9,
		"bb_period":          20,
		"bb_deviation":       2.0,
		"sma_short_period":   10,
		"sma_long_period":    30,
		"risk_per_trade":     0.02,
	}
}

// Helper to generate mock candles for demonstration
func generateMockCandles(pair string, timeframe string, count int) []OHLCData {
	candles := make([]OHLCData, count)
	
	// Base price varies by pair to make different results
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
	
	// Timeframe duration
	var interval time.Duration
	switch timeframe {
	case "1m":
		interval = time.Minute
	case "5m":
		interval = 5 * time.Minute
	case "15m":
		interval = 15 * time.Minute
	case "1h":
		interval = time.Hour
	case "4h":
		interval = 4 * time.Hour
	case "1d":
		interval = 24 * time.Hour
	default:
		interval = time.Hour
	}
	
	// Generate mock price data
	now := time.Now()
	price := basePrice
	
	for i := count - 1; i >= 0; i-- {
		timestamp := now.Add(-interval * time.Duration(i))
		
		// Random walk with momentum
		change := (rand.Float64() - 0.48) * 0.02 * price
		price += change
		
		// Generate OHLC values
		open := price
		high := price * (1 + 0.005*rand.Float64())
		low := price * (1 - 0.005*rand.Float64())
		close := price * (1 + (rand.Float64()-0.5)*0.01)
		volume := basePrice * rand.Float64() * 100
		
		candles[count-1-i] = OHLCData{
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

// Helper to generate technical analysis features
func generateTAFeatures(candles []OHLCData) []SignalFeature {
	if len(candles) < 30 {
		return nil
	}
	
	// Calculate recent price action
	priceAction := calculatePriceAction(candles)
	
	// Calculate volume trends
	volumeTrend := calculateVolumeTrend(candles)
	
	// Calculate momentum
	momentum := calculateMomentum(candles)
	
	// Calculate volatility
	volatility := calculateVolatility(candles)
	
	// Calculate trend strength
	trendStrength := calculateTrendStrength(candles)
	
	return []SignalFeature{
		{Name: "priceAction", Value: priceAction},
		{Name: "volume", Value: volumeTrend},
		{Name: "momentum", Value: momentum},
		{Name: "volatility", Value: volatility},
		{Name: "trendStrength", Value: trendStrength},
	}
}

// Price action helper (simplified calculation)
func calculatePriceAction(candles []OHLCData) float64 {
	n := len(candles)
	if n < 10 {
		return 0.5 // Neutral if not enough data
	}
	
	// Recent price movement
	recentClose := candles[n-1].Close
	prevClose := candles[n-10].Close
	
	percentChange := (recentClose - prevClose) / prevClose
	
	// Normalize to 0-1 range
	return math.Max(0, math.Min(1, (percentChange+0.1)/0.2))
}

// Volume trend helper
func calculateVolumeTrend(candles []OHLCData) float64 {
	n := len(candles)
	if n < 10 {
		return 0.5
	}
	
	// Compare recent volume to historical average
	var recentVolume, historicalVolume float64
	
	for i := n - 5; i < n; i++ {
		recentVolume += candles[i].Volume
	}
	recentVolume /= 5
	
	for i := n - 20; i < n - 5; i++ {
		historicalVolume += candles[i].Volume
	}
	historicalVolume /= 15
	
	if historicalVolume == 0 {
		return 0.5
	}
	
	ratio := recentVolume / historicalVolume
	
	// Normalize to 0-1 range
	return math.Max(0, math.Min(1, ratio/2))
}

// Momentum helper
func calculateMomentum(candles []OHLCData) float64 {
	n := len(candles)
	if n < 14 {
		return 0.5
	}
	
	// Calculate rate of change
	roc := (candles[n-1].Close - candles[n-14].Close) / candles[n-14].Close
	
	// Normalize to 0-1 range
	return math.Max(0, math.Min(1, (roc+0.1)/0.2))
}

// Volatility helper
func calculateVolatility(candles []OHLCData) float64 {
	n := len(candles)
	if n < 14 {
		return 0.5
	}
	
	// Calculate average true range (ATR) as volatility measure
	var atr float64
	
	for i := n - 14; i < n; i++ {
		trueRange := math.Max(
			candles[i].High-candles[i].Low,
			math.Max(
				math.Abs(candles[i].High-candles[i-1].Close),
				math.Abs(candles[i].Low-candles[i-1].Close),
			),
		)
		atr += trueRange / candles[i].Close // Normalize by price
	}
	atr /= 14
	
	// Normalize to 0-1 range (ATR of 2% -> 0.5)
	return math.Max(0, math.Min(1, atr/0.04))
}

// Trend strength helper
func calculateTrendStrength(candles []OHLCData) float64 {
	n := len(candles)
	if n < 20 {
		return 0.5
	}
	
	// Simple moving averages
	var sma10, sma20 float64
	
	for i := n - 10; i < n; i++ {
		sma10 += candles[i].Close
	}
	sma10 /= 10
	
	for i := n - 20; i < n; i++ {
		sma20 += candles[i].Close
	}
	sma20 /= 20
	
	// Measure alignment of SMAs
	diff := (sma10 - sma20) / sma20
	
	// Normalize to 0-1 range
	return math.Max(0, math.Min(1, (diff+0.03)/0.06))
}
