package ai

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

// ParamRange defines valid ranges for parameters
type ParamRange struct {
	Min          float64
	Max          float64
	StepSize     float64
	IsInteger    bool
	LogScale     bool   // Use logarithmic scaling for param search
	CurrentValue float64
}

// ParamSet represents a set of parameters to optimize
type ParamSet struct {
	Params       map[string]ParamRange
	FitnessScore float64
}

// StrategyPerformance contains metrics about strategy performance
type StrategyPerformance struct {
	ProfitLoss        float64  // Total P&L (%)
	SharpeRatio       float64  // Risk-adjusted returns
	MaxDrawdown       float64  // Maximum drawdown (%)
	WinRate           float64  // % of profitable trades
	ProfitFactor      float64  // Gross profits / gross losses
	RecoveryFactor    float64  // Profit / maximum drawdown
	ExpectedValue     float64  // Average return per trade
	NumTrades         int      // Total number of trades
	AvgHoldingPeriod  float64  // Average time in trade
	AvgProfitPerTrade float64  // Average profit per trade (%)
	CalmarRatio       float64  // Annual return / max drawdown
	SortinoRatio      float64  // Downside risk-adjusted
	PercentProfitable float64  // Percentage of profitable trades
	Alpha             float64  // Excess return over benchmark
	Beta              float64  // Volatility compared to market
}

// OptimizationResult contains results of parameter optimization
type OptimizationResult struct {
	BestParams        map[string]float64
	BestPerformance   StrategyPerformance
	AllTrials         []ParamSet
	CompletedTrials   int
	StartTime         time.Time
	EndTime           time.Time
	OptimizationMeta  map[string]interface{}
}

// Optimizer handles automatic parameter tuning
type Optimizer struct {
	pairs                []string
	timeframes           []string
	paramRanges          map[string]ParamRange
	optimizationHistory  []OptimizationResult
	backtest             func(map[string]float64) StrategyPerformance
	currentBest          ParamSet
	optimizationLock     sync.RWMutex
	backfilledData       map[string]map[string][]OHLCData // pair->timeframe->data
	dataLock             sync.RWMutex
	iterationCallback    func(trial int, params map[string]float64, perf StrategyPerformance)
}

// NewOptimizer creates a new optimization engine
func NewOptimizer(backtest func(map[string]float64) StrategyPerformance) *Optimizer {
	return &Optimizer{
		pairs:               []string{},
		timeframes:          []string{"1h", "4h", "1d"},
		paramRanges:         make(map[string]ParamRange),
		optimizationHistory: []OptimizationResult{},
		backtest:            backtest,
		currentBest: ParamSet{
			Params:       make(map[string]ParamRange),
			FitnessScore: 0,
		},
		backfilledData: make(map[string]map[string][]OHLCData),
	}
}

// SetPairs configures pairs to include in optimization
func (o *Optimizer) SetPairs(pairs []string) {
	o.pairs = pairs
}

// SetTimeframes configures timeframes to analyze
func (o *Optimizer) SetTimeframes(timeframes []string) {
	o.timeframes = timeframes
}

// AddParameterToOptimize adds a parameter to the optimization space
func (o *Optimizer) AddParameterToOptimize(name string, min, max, step float64, isInteger, logScale bool, currentValue float64) {
	o.paramRanges[name] = ParamRange{
		Min:          min,
		Max:          max,
		StepSize:     step,
		IsInteger:    isInteger,
		LogScale:     logScale,
		CurrentValue: currentValue,
	}
}

// SetIterationCallback sets a callback function that's called after each iteration
func (o *Optimizer) SetIterationCallback(callback func(trial int, params map[string]float64, perf StrategyPerformance)) {
	o.iterationCallback = callback
}

// AddHistoricalData adds backfill data for optimization
func (o *Optimizer) AddHistoricalData(pair, timeframe string, data []OHLCData) {
	o.dataLock.Lock()
	defer o.dataLock.Unlock()
	
	if _, exists := o.backfilledData[pair]; !exists {
		o.backfilledData[pair] = make(map[string][]OHLCData)
	}
	
	o.backfilledData[pair][timeframe] = data
}

// RandomSearch performs random search optimization
func (o *Optimizer) RandomSearch(iterations int, seed int64) *OptimizationResult {
	if seed != 0 {
		rand.Seed(seed)
	} else {
		rand.Seed(time.Now().UnixNano())
	}
	
	result := &OptimizationResult{
		BestParams:      make(map[string]float64),
		AllTrials:       make([]ParamSet, 0, iterations),
		StartTime:       time.Now(),
		OptimizationMeta: map[string]interface{}{
			"method": "random_search",
			"iterations": iterations,
		},
	}
	
	var bestScore float64 = -math.MaxFloat64
	var bestParams map[string]float64
	var bestPerformance StrategyPerformance
	
	for i := 0; i < iterations; i++ {
		// Generate random parameters
		params := make(map[string]float64)
		for name, paramRange := range o.paramRanges {
			var value float64
			if paramRange.LogScale {
				// Log scale sampling
				logMin := math.Log(paramRange.Min)
				logMax := math.Log(paramRange.Max)
				value = math.Exp(logMin + rand.Float64()*(logMax-logMin))
			} else {
				// Linear scale sampling
				value = paramRange.Min + rand.Float64()*(paramRange.Max-paramRange.Min)
			}
			
			if paramRange.IsInteger {
				value = math.Round(value)
			} else if paramRange.StepSize > 0 {
				// Discretize to step size
				value = math.Round(value/paramRange.StepSize) * paramRange.StepSize
			}
			
			params[name] = value
		}
		
		// Run backtest with these parameters
		performance := o.backtest(params)
		
		// Calculate fitness score (can be customized based on goals)
		fitnessScore := calculateFitnessScore(performance)
		
		// Record this trial
		paramSet := ParamSet{
			Params:       makeCopyOfParamRanges(o.paramRanges),
			FitnessScore: fitnessScore,
		}
		// Update current values to what we just tested
		for name, value := range params {
			if p, exists := paramSet.Params[name]; exists {
				p.CurrentValue = value
				paramSet.Params[name] = p
			}
		}
		result.AllTrials = append(result.AllTrials, paramSet)
		
		// Call iteration callback if set
		if o.iterationCallback != nil {
			o.iterationCallback(i, params, performance)
		}
		
		// Check if this is the best so far
		if fitnessScore > bestScore {
			bestScore = fitnessScore
			bestParams = params
			bestPerformance = performance
			
			// Update current best
			o.optimizationLock.Lock()
			o.currentBest = paramSet
			o.optimizationLock.Unlock()
			
			fmt.Printf("New best parameters found (iteration %d/%d):\n", i+1, iterations)
			for name, value := range bestParams {
				fmt.Printf("  %s: %.6g\n", name, value)
			}
			fmt.Printf("Performance: PnL=%.2f%%, Sharpe=%.2f, Drawdown=%.2f%%, Win=%.2f%%\n",
				performance.ProfitLoss,
				performance.SharpeRatio,
				performance.MaxDrawdown,
				performance.WinRate*100)
		}
		
		// Update progress
		result.CompletedTrials = i + 1
	}
	
	// Save final results
	result.BestParams = bestParams
	result.BestPerformance = bestPerformance
	result.EndTime = time.Now()
	
	// Add to history
	o.optimizationHistory = append(o.optimizationHistory, *result)
	
	return result
}

// BayesianOptimization performs intelligent parameter search
func (o *Optimizer) BayesianOptimization(maxIterations int, explorationFactor float64) *OptimizationResult {
	// In a real implementation, this would use proper Bayesian optimization
	// For now, we'll implement a simplified version with adaptive random search
	
	result := &OptimizationResult{
		BestParams:      make(map[string]float64),
		AllTrials:       make([]ParamSet, 0, maxIterations),
		StartTime:       time.Now(),
		OptimizationMeta: map[string]interface{}{
			"method": "bayesian_optimization",
			"iterations": maxIterations,
			"exploration_factor": explorationFactor,
		},
	}
	
	// Start with the current values as a baseline
	currentParams := make(map[string]float64)
	for name, paramRange := range o.paramRanges {
		currentParams[name] = paramRange.CurrentValue
	}
	
	// Run initial evaluation
	initialPerformance := o.backtest(currentParams)
	bestScore := calculateFitnessScore(initialPerformance)
	bestParams := currentParams
	bestPerformance := initialPerformance
	
	// Record initial trial
	paramSet := ParamSet{
		Params:       makeCopyOfParamRanges(o.paramRanges),
		FitnessScore: bestScore,
	}
	result.AllTrials = append(result.AllTrials, paramSet)
	
	// Adaptive search radius (starts large, shrinks as we converge)
	searchRadiusFactor := 1.0
	
	for i := 1; i < maxIterations; i++ {
		// Generate parameters in the neighborhood of the best parameters
		params := make(map[string]float64)
		for name, paramRange := range o.paramRanges {
			var value float64
			
			// Get the current best value for this parameter
			baseValue := bestParams[name]
			
			// Calculate search radius based on parameter range
			searchRadius := (paramRange.Max - paramRange.Min) * searchRadiusFactor * explorationFactor
			
			// Sample from neighborhood of best value, clamped to valid range
			minVal := math.Max(paramRange.Min, baseValue - searchRadius)
			maxVal := math.Min(paramRange.Max, baseValue + searchRadius)
			
			if paramRange.LogScale {
				// Log scale sampling
				logMin := math.Log(minVal)
				logMax := math.Log(maxVal)
				value = math.Exp(logMin + rand.Float64()*(logMax-logMin))
			} else {
				// Linear scale sampling
				value = minVal + rand.Float64()*(maxVal-minVal)
			}
			
			if paramRange.IsInteger {
				value = math.Round(value)
			} else if paramRange.StepSize > 0 {
				value = math.Round(value/paramRange.StepSize) * paramRange.StepSize
			}
			
			params[name] = value
		}
		
		// Run backtest with these parameters
		performance := o.backtest(params)
		
		// Calculate fitness score
		fitnessScore := calculateFitnessScore(performance)
		
		// Record this trial
		paramSet := ParamSet{
			Params:       makeCopyOfParamRanges(o.paramRanges),
			FitnessScore: fitnessScore,
		}
		// Update current values to what we just tested
		for name, value := range params {
			if p, exists := paramSet.Params[name]; exists {
				p.CurrentValue = value
				paramSet.Params[name] = p
			}
		}
		result.AllTrials = append(result.AllTrials, paramSet)
		
		// Call iteration callback if set
		if o.iterationCallback != nil {
			o.iterationCallback(i, params, performance)
		}
		
		// Check if this is the best so far
		if fitnessScore > bestScore {
			bestScore = fitnessScore
			bestParams = params
			bestPerformance = performance
			
			// Update current best
			o.optimizationLock.Lock()
			o.currentBest = paramSet
			o.optimizationLock.Unlock()
			
			// Reset search radius when we find a better solution
			searchRadiusFactor = 1.0
			
			fmt.Printf("New best parameters found (iteration %d/%d):\n", i+1, maxIterations)
			for name, value := range bestParams {
				fmt.Printf("  %s: %.6g\n", name, value)
			}
			fmt.Printf("Performance: PnL=%.2f%%, Sharpe=%.2f, Drawdown=%.2f%%, Win=%.2f%%\n",
				performance.ProfitLoss,
				performance.SharpeRatio,
				performance.MaxDrawdown,
				performance.WinRate*100)
		} else {
			// Shrink search radius gradually if we're not finding better solutions
			searchRadiusFactor *= 0.95
			if searchRadiusFactor < 0.1 {
				// Occasionally reset to avoid getting stuck
				if rand.Float64() < 0.2 {
					searchRadiusFactor = 1.0
				}
			}
		}
		
		// Update progress
		result.CompletedTrials = i + 1
	}
	
	// Save final results
	result.BestParams = bestParams
	result.BestPerformance = bestPerformance
	result.EndTime = time.Now()
	
	// Add to history
	o.optimizationHistory = append(o.optimizationHistory, *result)
	
	return result
}

// WalkForwardOptimization performs optimization on rolling time windows
func (o *Optimizer) WalkForwardOptimization(windowSize, stepSize time.Duration, iterations int) *OptimizationResult {
	// In a full implementation, this would:
	// 1. Split historical data into windows
	// 2. Optimize on each training window
	// 3. Validate on out-of-sample data
	// 4. Roll forward and repeat
	// 5. Analyze parameter stability and performance consistency
	
	// This is a simplified implementation for the framework
	
	result := &OptimizationResult{
		BestParams:      make(map[string]float64),
		AllTrials:       make([]ParamSet, 0),
		StartTime:       time.Now(),
		OptimizationMeta: map[string]interface{}{
			"method": "walk_forward",
			"window_size": windowSize.String(),
			"step_size": stepSize.String(),
			"iterations": iterations,
		},
	}
	
	fmt.Println("Starting walk-forward optimization...")
	fmt.Printf("Window size: %s, Step size: %s\n", windowSize, stepSize)
	
	// For each pair and timeframe, determine the date range
	var startDate, endDate time.Time
	
	o.dataLock.RLock()
	for _, timeframeMap := range o.backfilledData {
		for _, data := range timeframeMap {
			if len(data) > 0 {
				if startDate.IsZero() || data[0].Timestamp.Before(startDate) {
					startDate = data[0].Timestamp
				}
				
				lastIdx := len(data) - 1
				if endDate.IsZero() || data[lastIdx].Timestamp.After(endDate) {
					endDate = data[lastIdx].Timestamp
				}
			}
		}
	}
	o.dataLock.RUnlock()
	
	if startDate.IsZero() || endDate.IsZero() {
		fmt.Println("Error: No data available for walk-forward optimization")
		return result
	}
	
	// Calculate the number of windows
	totalDuration := endDate.Sub(startDate)
	numWindows := int(totalDuration / stepSize)
	
	// Parameters from each window optimization
	windowParameters := make([]map[string]float64, 0, numWindows)
	
	// For each window, perform optimization
	for i := 0; i < numWindows; i++ {
		windowStart := startDate.Add(time.Duration(i) * stepSize)
		windowEnd := windowStart.Add(windowSize)
		
		if windowEnd.After(endDate) {
			windowEnd = endDate
		}
		
		fmt.Printf("\nOptimizing window %d/%d (%s to %s)...\n", 
			i+1, numWindows, windowStart.Format("2006-01-02"), windowEnd.Format("2006-01-02"))
		
		// In a real implementation, we would filter the data for this time window
		// and pass it to a backtest function that's constrained to the window
		
		// Perform a small bayesian optimization for this window
		tempOptimizer := *o
		// Temporary backtest function that would be constrained to window
		tempOptimizer.backtest = func(params map[string]float64) StrategyPerformance {
			// This simulates a backtest on just the window data
			// In a real implementation, we'd filter the historical data
			return o.backtest(params)
		}
		
		windowResult := tempOptimizer.BayesianOptimization(iterations/numWindows, 0.5)
		
		// Save the best parameters from this window
		windowParameters = append(windowParameters, windowResult.BestParams)
		
		// Append trials to the overall result
		result.AllTrials = append(result.AllTrials, windowResult.AllTrials...)
		result.CompletedTrials += windowResult.CompletedTrials
	}
	
	// Analyze parameter stability across windows
	stableParams := make(map[string]float64)
	
	// Only process if we have at least one window result
	if len(windowParameters) > 0 {
		// For each parameter, calculate statistics
		paramStats := make(map[string]struct {
			Mean   float64
			StdDev float64
			Values []float64
		})
		
		// Initialize stats
		for paramName := range o.paramRanges {
			paramStats[paramName] = struct {
				Mean   float64
				StdDev float64
				Values []float64
			}{
				Values: make([]float64, 0, len(windowParameters)),
			}
		}
		
		// Collect all values
		for _, params := range windowParameters {
			for paramName, value := range params {
				stats := paramStats[paramName]
				stats.Values = append(stats.Values, value)
				paramStats[paramName] = stats
			}
		}
		
		// Calculate mean and standard deviation
		for paramName, stats := range paramStats {
			// Calculate mean
			sum := 0.0
			for _, v := range stats.Values {
				sum += v
			}
			mean := sum / float64(len(stats.Values))
			
			// Calculate standard deviation
			variance := 0.0
			for _, v := range stats.Values {
				variance += math.Pow(v-mean, 2)
			}
			stdDev := math.Sqrt(variance / float64(len(stats.Values)))
			
			// Update stats
			stats.Mean = mean
			stats.StdDev = stdDev
			paramStats[paramName] = stats
			
			// Use the mean as the stable parameter value
			stableParams[paramName] = mean
			
			fmt.Printf("Parameter %s: Mean=%.6g, StdDev=%.6g, CV=%.2f%%\n", 
				paramName, mean, stdDev, (stdDev/mean)*100)
		}
	}
	
	// Final evaluation of the stable parameters
	finalPerf := o.backtest(stableParams)
	
	// Save final results
	result.BestParams = stableParams
	result.BestPerformance = finalPerf
	result.EndTime = time.Now()
	
	fmt.Println("\nWalk-forward optimization complete.")
	fmt.Println("Final stable parameters:")
	for name, value := range stableParams {
		fmt.Printf("  %s: %.6g\n", name, value)
	}
	fmt.Printf("Overall Performance: PnL=%.2f%%, Sharpe=%.2f, Drawdown=%.2f%%, Win=%.2f%%\n",
		finalPerf.ProfitLoss,
		finalPerf.SharpeRatio,
		finalPerf.MaxDrawdown,
		finalPerf.WinRate*100)
	
	// Add to history
	o.optimizationHistory = append(o.optimizationHistory, *result)
	
	return result
}

// GetBestParameters returns the current best parameter set
func (o *Optimizer) GetBestParameters() map[string]float64 {
	o.optimizationLock.RLock()
	defer o.optimizationLock.RUnlock()
	
	bestParams := make(map[string]float64)
	for name, param := range o.currentBest.Params {
		bestParams[name] = param.CurrentValue
	}
	
	return bestParams
}

// GetOptimizationHistory returns all historical optimization runs
func (o *Optimizer) GetOptimizationHistory() []OptimizationResult {
	return o.optimizationHistory
}

// Helper functions

// calculateFitnessScore combines multiple performance metrics into a single score
func calculateFitnessScore(perf StrategyPerformance) float64 {
	// Customize this formula based on trading goals
	// Here we combine multiple factors with weights
	
	// Prevent division by zero
	if perf.MaxDrawdown <= 0 {
		perf.MaxDrawdown = 0.01
	}
	
	// Base score on profit/loss
	score := perf.ProfitLoss
	
	// Adjust based on risk metrics
	score *= (1.0 + 0.5*perf.SharpeRatio)      // Reward good Sharpe ratio
	score *= (1.0 + 0.3*perf.ProfitFactor)      // Reward good profit factor
	score *= (1.0 + 0.2*perf.RecoveryFactor)    // Reward recovery speed
	score *= (1.0 + 0.2*perf.WinRate)           // Reward consistent wins
	
	// Penalize for drawdown
	score *= math.Max(0.1, 1.0 - 0.5*(perf.MaxDrawdown/100.0))
	
	// Reward for decent number of trades (at least 20)
	tradeScaleFactor := math.Min(1.0, float64(perf.NumTrades) / 20.0)
	score *= tradeScaleFactor
	
	return score
}

// makeCopyOfParamRanges creates a deep copy of parameter ranges
func makeCopyOfParamRanges(original map[string]ParamRange) map[string]ParamRange {
	copy := make(map[string]ParamRange)
	for k, v := range original {
		copy[k] = v
	}
	return copy
}
