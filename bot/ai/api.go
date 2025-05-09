package ai

import (
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// AIOpportunityResponse represents the API response for AI-enhanced opportunity detection
type AIOpportunityResponse struct {
	Pair             string            `json:"pair"`
	Timeframe        string            `json:"timeframe"`
	Score            float64           `json:"score"`
	Signal           string            `json:"signal"`
	Confidence       float64           `json:"confidence"`
	PredictedMove    float64           `json:"predicted_move"`
	RecommendedSize  float64           `json:"recommended_size"`
	PatternSignals   []string          `json:"pattern_signals,omitempty"`
	SentimentScore   float64           `json:"sentiment_score,omitempty"`
	TopFeatures      map[string]float64 `json:"top_features,omitempty"`
	LastUpdated      string            `json:"last_updated"`
}

// AIOpportunitiesRequest represents the API request for fetching opportunities
type AIOpportunitiesRequest struct {
	Pairs         []string `json:"pairs"`
	MinScore      float64  `json:"min_score"`
	Limit         int      `json:"limit"`
	IncludeDetail bool     `json:"include_detail"`
}

// AIBacktestRequest represents a request to run a backtest with AI optimization
type AIBacktestRequest struct {
	Pairs          []string          `json:"pairs"`
	Timeframe      string            `json:"timeframe"`
	StartDate      string            `json:"start_date"`
	EndDate        string            `json:"end_date"`
	OptimizeParams bool              `json:"optimize_params"`
	Parameters     map[string]float64 `json:"parameters,omitempty"`
}

// AIOptimizeRequest represents a parameter optimization request
type AIOptimizeRequest struct {
	Method         string   `json:"method"` // "random", "bayesian", "walkforward"
	Iterations     int      `json:"iterations"`
	Pairs          []string `json:"pairs"`
	TimeframeStart string   `json:"timeframe_start"`
	TimeframeEnd   string   `json:"timeframe_end"`
}

// AIModelInfoResponse contains information about the AI model
type AIModelInfoResponse struct {
	LastUpdated     string                 `json:"last_updated"`
	LastOptimized   string                 `json:"last_optimized"`
	FeatureWeights  map[string]float64     `json:"feature_weights"`
	RunningScores   map[string]float64     `json:"running_scores"`
	ModelParameters map[string]interface{} `json:"model_parameters"`
}

// RegisterAIRoutes adds AI-related API endpoints to a Gin router
func RegisterAIRoutes(router *gin.RouterGroup, engine *AIEngine) {
	if router == nil || engine == nil {
		log.Println("Cannot register AI routes: router or engine is nil")
		return
	}

	// GET /ai/status - Get AI engine status
	router.GET("/status", func(c *gin.Context) {
		modelSummary := engine.GetModelSummary()
		
		response := gin.H{
			"ai_enabled":        true,
			"components":        engine.enabledComponents,
			"last_scan":         engine.GetLastScanTime().Format(time.RFC3339),
			"last_optimization": engine.GetLastOptimizationTime().Format(time.RFC3339),
			"model_summary":     modelSummary,
		}
		
		c.JSON(http.StatusOK, response)
	})

	// POST /ai/opportunities - Get AI-ranked trading opportunities
	router.POST("/opportunities", func(c *gin.Context) {
		var req AIOpportunitiesRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		// Use provided minimum score or default to 0.6
		minScore := 0.6
		if req.MinScore > 0 {
			minScore = req.MinScore
		}
		
		// Use provided limit or default to 10
		limit := 10
		if req.Limit > 0 {
			limit = req.Limit
		}
		
		// Filter by requested pairs if specified
		if len(req.Pairs) > 0 {
			// Force a scan of the requested pairs
			var filteredPairs []string
			for _, p := range req.Pairs {
				for _, supportedPair := range engine.pairs {
					if p == supportedPair {
						filteredPairs = append(filteredPairs, p)
						break
					}
				}
			}
			
			// Set pairs temporarily and restore after scan
			originalPairs := engine.pairs
			engine.pairs = filteredPairs
			engine.ScanAllMarkets()
			engine.pairs = originalPairs
		}
		
		// Get top opportunities
		opportunities := engine.GetOpportunityRanking(minScore, limit)
		
		// Convert to response format
		response := make([]AIOpportunityResponse, 0, len(opportunities))
		for _, opp := range opportunities {
			item := AIOpportunityResponse{
				Pair:            opp.Pair,
				Timeframe:       opp.Timeframe,
				Score:           opp.Score,
				Signal:          opp.Signal,
				Confidence:      opp.Confidence,
				PredictedMove:   opp.PredictedMove,
				RecommendedSize: opp.RecommendedSize,
				LastUpdated:     opp.Timestamp.Format(time.RFC3339),
			}
			
			// Include additional details if requested
			if req.IncludeDetail {
				// Add pattern signals
				patterns := make([]string, 0)
				for _, p := range opp.PatternSignals {
					patterns = append(patterns, string(p.Pattern))
				}
				item.PatternSignals = patterns
				
				// Add sentiment if available
				if opp.SentimentData != nil {
					item.SentimentScore = opp.SentimentData.Score
				}
				
				// Add top 5 features by weight
				item.TopFeatures = make(map[string]float64)
				featsByWeight := make([]SignalFeature, len(opp.MLFeatures))
				copy(featsByWeight, opp.MLFeatures)
				
				// Sort by weight
				sort.Slice(featsByWeight, func(i, j int) bool {
					return featsByWeight[i].Weight > featsByWeight[j].Weight
				})
				
				// Take top 5
				for i := 0; i < len(featsByWeight) && i < 5; i++ {
					item.TopFeatures[featsByWeight[i].Name] = featsByWeight[i].Value
				}
			}
			
			response = append(response, item)
		}
		
		c.JSON(http.StatusOK, response)
	})

	// POST /ai/analyze - Analyze a specific pair
	router.POST("/analyze", func(c *gin.Context) {
		pair := c.PostForm("pair")
		timeframe := c.PostForm("timeframe")
		
		if pair == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "pair is required"})
			return
		}
		
		if timeframe == "" {
			timeframe = "1h" // Default timeframe
		}
		
		// Perform analysis
		result := engine.AnalyzeMarket(pair, timeframe)
		
		// Convert to response format
		patterns := make([]string, 0)
		for _, p := range result.PatternSignals {
			patterns = append(patterns, string(p.Pattern))
		}
		
		topFeatures := make(map[string]float64)
		featsByWeight := make([]SignalFeature, len(result.MLFeatures))
		copy(featsByWeight, result.MLFeatures)
		
		// Sort by weight
		sort.Slice(featsByWeight, func(i, j int) bool {
			return featsByWeight[i].Weight > featsByWeight[j].Weight
		})
		
		// Take top 5
		for i := 0; i < len(featsByWeight) && i < 5; i++ {
			topFeatures[featsByWeight[i].Name] = featsByWeight[i].Value
		}
		
		var sentimentScore float64
		if result.SentimentData != nil {
			sentimentScore = result.SentimentData.Score
		}
		
		response := AIOpportunityResponse{
			Pair:            result.Pair,
			Timeframe:       result.Timeframe,
			Score:           result.Score,
			Signal:          result.Signal,
			Confidence:      result.Confidence,
			PredictedMove:   result.PredictedMove,
			RecommendedSize: result.RecommendedSize,
			PatternSignals:  patterns,
			SentimentScore:  sentimentScore,
			TopFeatures:     topFeatures,
			LastUpdated:     result.Timestamp.Format(time.RFC3339),
		}
		
		c.JSON(http.StatusOK, response)
	})

	// POST /ai/backtest - Backtest with AI optimization
	router.POST("/backtest", func(c *gin.Context) {
		var req AIBacktestRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		// This endpoint would run a backtest with the provided parameters
		// or perform optimization if requested
		if req.OptimizeParams {
			// For now, just return a simple response (full implementation would be complex)
			c.JSON(http.StatusOK, gin.H{
				"message": "Optimization and backtesting initiated",
				"status": "running",
				"job_id": "backtest-" + strconv.FormatInt(time.Now().Unix(), 10),
			})
			
			// In a real implementation, this would start a background job
			go func() {
				log.Println("Starting AI-optimized backtest...")
				// Here we would start the optimizer and run backtests
			}()
			
			return
		}
		
		// Run a single backtest with provided parameters
		performance, err := engine.RunSingleBacktest()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"performance": performance,
			"parameters": req.Parameters,
		})
	})

	// GET /ai/model - Get model information
	router.GET("/model", func(c *gin.Context) {
		modelSummary := engine.GetModelSummary()
		runningScores := engine.GetRunningAverageScores()
		
		response := AIModelInfoResponse{
			LastUpdated:     engine.GetLastScanTime().Format(time.RFC3339),
			LastOptimized:   engine.GetLastOptimizationTime().Format(time.RFC3339),
			RunningScores:   runningScores,
			ModelParameters: modelSummary,
		}
		
		if weights, ok := modelSummary["weights"].(map[string]float64); ok {
			response.FeatureWeights = weights
		}
		
		c.JSON(http.StatusOK, response)
	})

	// POST /ai/optimize - Trigger parameter optimization
	router.POST("/optimize", func(c *gin.Context) {
		var req AIOptimizeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		// Set iterations with reasonable defaults and limits
		iterations := 100
		if req.Iterations > 0 {
			iterations = req.Iterations
		}
		if iterations > 500 {
			iterations = 500 // Cap for reasonable runtime
		}
		
		// This would start optimization in the background
		c.JSON(http.StatusOK, gin.H{
			"message": "Optimization initiated",
			"method": req.Method,
			"iterations": iterations,
			"status": "running",
			"job_id": "optimize-" + strconv.FormatInt(time.Now().Unix(), 10),
		})
		
		// Start optimization in background
		go func() {
			log.Printf("Starting %s optimization with %d iterations...", req.Method, iterations)
			
			switch req.Method {
			case "random":
				engine.optimizer.RandomSearch(iterations, 0)
			case "bayesian":
				engine.optimizer.BayesianOptimization(iterations, 0.5)
			case "walkforward":
				windowSize := 30 * 24 * time.Hour // 30 days
				stepSize := 7 * 24 * time.Hour   // 7 days
				engine.optimizer.WalkForwardOptimization(windowSize, stepSize, iterations)
			default:
				// Default to random search
				engine.optimizer.RandomSearch(iterations, 0)
			}
			
			log.Println("Optimization complete")
		}()
	})
}
