package ai

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// PatternType represents recognized chart patterns
type PatternType string

const (
	// Bullish patterns
	PatternBullishFlag          PatternType = "bullish_flag"
	PatternCupAndHandle         PatternType = "cup_and_handle"
	PatternInverseHeadShoulder  PatternType = "inverse_head_shoulder"
	PatternBullishRectangle     PatternType = "bullish_rectangle"
	PatternAscendingTriangle    PatternType = "ascending_triangle"
	PatternBullishPennant       PatternType = "bullish_pennant"
	PatternMorningStarDoji      PatternType = "morning_star_doji"
	
	// Bearish patterns
	PatternBearishFlag          PatternType = "bearish_flag"
	PatternHeadAndShoulders     PatternType = "head_and_shoulders"
	PatternDoubleTop            PatternType = "double_top"
	PatternBearishRectangle     PatternType = "bearish_rectangle"
	PatternDescendingTriangle   PatternType = "descending_triangle"
	PatternBearishPennant       PatternType = "bearish_pennant"
	PatternEveningStarDoji      PatternType = "evening_star_doji"
	
	// Continuation patterns
	PatternSymmetricalTriangle  PatternType = "symmetrical_triangle"
	PatternRectangle            PatternType = "rectangle"
	
	// Reversal indicators
	PatternDoji                 PatternType = "doji"
	PatternHammer               PatternType = "hammer"
	PatternShootingStar         PatternType = "shooting_star"
	PatternEngulfing            PatternType = "engulfing"
	
	// Volatility patterns
	PatternVolatilityCompression PatternType = "volatility_compression"
	PatternVolatilityExpansion   PatternType = "volatility_expansion"
)

// PatternSignal represents detected chart patterns
type PatternSignal struct {
	Pattern       PatternType
	Confidence    float64        // 0.0-1.0
	Direction     float64        // -1.0 (bearish) to 1.0 (bullish)
	StartIndex    int            // Starting candle index
	EndIndex      int            // Ending candle index
	PredictedMove float64        // Expected price movement (in %)
	TimeFrame     string         // "1h", "4h", "1d", etc.
	Timestamp     time.Time
}

// Timeframe for analysis
type Timeframe string

const (
	Timeframe1Min  Timeframe = "1m"
	Timeframe5Min  Timeframe = "5m"
	Timeframe15Min Timeframe = "15m"
	Timeframe1Hour Timeframe = "1h"
	Timeframe4Hour Timeframe = "4h"
	Timeframe1Day  Timeframe = "1d"
	Timeframe1Week Timeframe = "1w"
)

// OHLCData represents price candle data
type OHLCData struct {
	Timestamp time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

// PatternRecognizer detects chart patterns in price data
type PatternRecognizer struct {
	signalThreshold float64
	patternCache    map[string][]PatternSignal
	cacheLock       sync.RWMutex
	lastUpdate      map[string]time.Time
}

// NewPatternRecognizer creates a new pattern recognition engine
func NewPatternRecognizer() *PatternRecognizer {
	return &PatternRecognizer{
		signalThreshold: 0.65, // Minimum confidence to report a pattern
		patternCache:    make(map[string][]PatternSignal),
		lastUpdate:      make(map[string]time.Time),
	}
}

// AnalyzePatterns searches for patterns in OHLC data
func (pr *PatternRecognizer) AnalyzePatterns(pair string, timeframe Timeframe, data []OHLCData) []PatternSignal {
	if len(data) < 30 {
		return nil // Need sufficient data for pattern recognition
	}
	
	var signals []PatternSignal
	
	// Run various pattern detection algorithms in parallel
	var wg sync.WaitGroup
	var mutex sync.Mutex // To protect signals slice
	
	// Head and Shoulders (bearish)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if signal := pr.detectHeadAndShoulders(data); signal != nil {
			mutex.Lock()
			signals = append(signals, *signal)
			mutex.Unlock()
		}
	}()
	
	// Inverse Head and Shoulders (bullish)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if signal := pr.detectInverseHeadAndShoulders(data); signal != nil {
			mutex.Lock()
			signals = append(signals, *signal)
			mutex.Unlock()
		}
	}()
	
	// Double Top (bearish)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if signal := pr.detectDoubleTop(data); signal != nil {
			mutex.Lock()
			signals = append(signals, *signal)
			mutex.Unlock()
		}
	}()
	
	// Ascending Triangle (bullish)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if signal := pr.detectAscendingTriangle(data); signal != nil {
			mutex.Lock()
			signals = append(signals, *signal)
			mutex.Unlock()
		}
	}()
	
	// Descending Triangle (bearish)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if signal := pr.detectDescendingTriangle(data); signal != nil {
			mutex.Lock()
			signals = append(signals, *signal)
			mutex.Unlock()
		}
	}()
	
	// Volatility Patterns
	wg.Add(1)
	go func() {
		defer wg.Done()
		if signal := pr.detectVolatilityPatterns(data); signal != nil {
			mutex.Lock()
			signals = append(signals, *signal)
			mutex.Unlock()
		}
	}()
	
	// Japanese Candlestick Patterns
	wg.Add(1)
	go func() {
		defer wg.Done()
		candleSignals := pr.detectCandlestickPatterns(data)
		if len(candleSignals) > 0 {
			mutex.Lock()
			signals = append(signals, candleSignals...)
			mutex.Unlock()
		}
	}()
	
	wg.Wait()
	
	// Update pattern cache
	cacheKey := fmt.Sprintf("%s-%s", pair, timeframe)
	pr.cacheLock.Lock()
	pr.patternCache[cacheKey] = signals
	pr.lastUpdate[cacheKey] = time.Now()
	pr.cacheLock.Unlock()
	
	return signals
}

// GetCachedPatterns returns previously detected patterns
func (pr *PatternRecognizer) GetCachedPatterns(pair string, timeframe Timeframe) []PatternSignal {
	cacheKey := fmt.Sprintf("%s-%s", pair, timeframe)
	
	pr.cacheLock.RLock()
	defer pr.cacheLock.RUnlock()
	
	if signals, ok := pr.patternCache[cacheKey]; ok {
		return signals
	}
	
	return nil
}

// PatternToSignalFeatures converts pattern signals to ML model features
func (pr *PatternRecognizer) PatternToSignalFeatures(patterns []PatternSignal) []SignalFeature {
	if len(patterns) == 0 {
		return nil
	}
	
	// Aggregate pattern signals
	var bullishScore, bearishScore float64
	var patternCount int
	
	for _, p := range patterns {
		weight := p.Confidence
		if p.Direction > 0 {
			bullishScore += p.Direction * weight
		} else {
			bearishScore += math.Abs(p.Direction) * weight
		}
		patternCount++
	}
	
	// Normalize scores
	if patternCount > 0 {
		bullishScore /= float64(patternCount)
		bearishScore /= float64(patternCount)
	}
	
	// Net score (-1 to 1)
	netScore := bullishScore - bearishScore
	
	// Convert to 0-1 scale for ML model
	normalizedScore := (netScore + 1.0) / 2.0
	
	// Pattern strength based on number and confidence
	patternStrength := math.Min(1.0, float64(patternCount) / 5.0) // Cap at 5 patterns
	
	features := []SignalFeature{
		{Name: "patternSignal", Value: normalizedScore},
		{Name: "patternStrength", Value: patternStrength},
	}
	
	return features
}

// Implementation of pattern detection algorithms
// In a production system, these would use more sophisticated algorithms

func (pr *PatternRecognizer) detectHeadAndShoulders(data []OHLCData) *PatternSignal {
	// Simplified detection algorithm for demonstration
	// In production, would use peak detection and correlation
	
	if len(data) < 20 {
		return nil
	}
	
	// For demonstration purposes, we're using a simplified detection
	// This should be replaced with actual technical analysis
	
	// Randomly detect a pattern with 10% probability for demo
	if time.Now().UnixNano()%10 == 0 {
		return &PatternSignal{
			Pattern:       PatternHeadAndShoulders,
			Confidence:    0.65 + (float64(time.Now().UnixNano()%20) / 100.0),
			Direction:     -0.8, // Bearish
			StartIndex:    len(data) - 20,
			EndIndex:      len(data) - 1,
			PredictedMove: -2.5 - (float64(time.Now().UnixNano()%10) / 10.0),
			TimeFrame:     "1d",
			Timestamp:     time.Now(),
		}
	}
	
	return nil
}

func (pr *PatternRecognizer) detectInverseHeadAndShoulders(data []OHLCData) *PatternSignal {
	if len(data) < 20 {
		return nil
	}
	
	// Simplified detection for demonstration
	if time.Now().UnixNano()%10 == 1 {
		return &PatternSignal{
			Pattern:       PatternInverseHeadShoulder,
			Confidence:    0.70 + (float64(time.Now().UnixNano()%15) / 100.0),
			Direction:     0.8, // Bullish
			StartIndex:    len(data) - 20,
			EndIndex:      len(data) - 1,
			PredictedMove: 2.5 + (float64(time.Now().UnixNano()%10) / 10.0),
			TimeFrame:     "1d",
			Timestamp:     time.Now(),
		}
	}
	
	return nil
}

func (pr *PatternRecognizer) detectDoubleTop(data []OHLCData) *PatternSignal {
	if len(data) < 15 {
		return nil
	}
	
	// Simplified detection for demonstration
	if time.Now().UnixNano()%10 == 2 {
		return &PatternSignal{
			Pattern:       PatternDoubleTop,
			Confidence:    0.68 + (float64(time.Now().UnixNano()%20) / 100.0),
			Direction:     -0.7, // Bearish
			StartIndex:    len(data) - 15,
			EndIndex:      len(data) - 1,
			PredictedMove: -2.0 - (float64(time.Now().UnixNano()%10) / 10.0),
			TimeFrame:     "1d",
			Timestamp:     time.Now(),
		}
	}
	
	return nil
}

func (pr *PatternRecognizer) detectAscendingTriangle(data []OHLCData) *PatternSignal {
	if len(data) < 15 {
		return nil
	}
	
	// Simplified detection for demonstration
	if time.Now().UnixNano()%10 == 3 {
		return &PatternSignal{
			Pattern:       PatternAscendingTriangle,
			Confidence:    0.72 + (float64(time.Now().UnixNano()%15) / 100.0),
			Direction:     0.75, // Bullish
			StartIndex:    len(data) - 15,
			EndIndex:      len(data) - 1,
			PredictedMove: 2.2 + (float64(time.Now().UnixNano()%10) / 10.0),
			TimeFrame:     "1d",
			Timestamp:     time.Now(),
		}
	}
	
	return nil
}

func (pr *PatternRecognizer) detectDescendingTriangle(data []OHLCData) *PatternSignal {
	if len(data) < 15 {
		return nil
	}
	
	// Simplified detection for demonstration
	if time.Now().UnixNano()%10 == 4 {
		return &PatternSignal{
			Pattern:       PatternDescendingTriangle,
			Confidence:    0.71 + (float64(time.Now().UnixNano()%15) / 100.0),
			Direction:     -0.75, // Bearish
			StartIndex:    len(data) - 15,
			EndIndex:      len(data) - 1,
			PredictedMove: -2.2 - (float64(time.Now().UnixNano()%10) / 10.0),
			TimeFrame:     "1d",
			Timestamp:     time.Now(),
		}
	}
	
	return nil
}

func (pr *PatternRecognizer) detectVolatilityPatterns(data []OHLCData) *PatternSignal {
	if len(data) < 10 {
		return nil
	}
	
	// Calculate recent volatility
	var recentVolatility, previousVolatility float64
	
	for i := len(data) - 5; i < len(data); i++ {
		range1 := data[i].High - data[i].Low
		recentVolatility += range1 / data[i].Open
	}
	recentVolatility /= 5.0
	
	for i := len(data) - 10; i < len(data) - 5; i++ {
		range1 := data[i].High - data[i].Low
		previousVolatility += range1 / data[i].Open
	}
	previousVolatility /= 5.0
	
	// Check for volatility compression (leads to breakouts)
	if recentVolatility < previousVolatility*0.7 {
		return &PatternSignal{
			Pattern:       PatternVolatilityCompression,
			Confidence:    0.70 + (0.3 * (1.0 - (recentVolatility / previousVolatility))),
			Direction:     0, // Neutral until breakout
			StartIndex:    len(data) - 10,
			EndIndex:      len(data) - 1,
			PredictedMove: 0, // Direction unclear
			TimeFrame:     "1d",
			Timestamp:     time.Now(),
		}
	}
	
	// Check for volatility expansion (often breakouts in progress)
	if recentVolatility > previousVolatility*1.5 {
		// Determine direction of breakout
		var direction float64
		if data[len(data)-1].Close > data[len(data)-6].Close {
			direction = 0.8 // Bullish breakout
		} else {
			direction = -0.8 // Bearish breakout
		}
		
		return &PatternSignal{
			Pattern:       PatternVolatilityExpansion,
			Confidence:    0.75,
			Direction:     direction,
			StartIndex:    len(data) - 10,
			EndIndex:      len(data) - 1,
			PredictedMove: direction * 3.0, // 3% in the breakout direction
			TimeFrame:     "1d",
			Timestamp:     time.Now(),
		}
	}
	
	return nil
}

func (pr *PatternRecognizer) detectCandlestickPatterns(data []OHLCData) []PatternSignal {
	if len(data) < 3 {
		return nil
	}
	
	var signals []PatternSignal
	idx := len(data) - 1 // Most recent candle
	prev1 := len(data) - 2
	prev2 := len(data) - 3
	
	current := data[idx]
	prev := data[prev1]
	prevPrev := data[prev2]
	
	// Detect doji (open and close are nearly equal)
	bodySize := math.Abs(current.Open - current.Close)
	totalRange := current.High - current.Low
	if totalRange > 0 && bodySize/totalRange < 0.1 {
		signals = append(signals, PatternSignal{
			Pattern:       PatternDoji,
			Confidence:    0.7,
			Direction:     0, // Neutral until context is clear
			StartIndex:    idx,
			EndIndex:      idx,
			PredictedMove: 0,
			TimeFrame:     "1d",
			Timestamp:     time.Now(),
		})
	}
	
	// Detect hammer (bullish reversal)
	if current.Close > current.Open { // Green candle
		body := current.Close - current.Open
		lowerWick := current.Open - current.Low
		upperWick := current.High - current.Close
		
		if lowerWick > body*2 && upperWick < body*0.5 && prev.Close < prev.Open {
			signals = append(signals, PatternSignal{
				Pattern:       PatternHammer,
				Confidence:    0.65 + (lowerWick/body)/10,
				Direction:     0.7, // Bullish
				StartIndex:    idx-1,
				EndIndex:      idx,
				PredictedMove: 1.5,
				TimeFrame:     "1d",
				Timestamp:     time.Now(),
			})
		}
	}
	
	// Detect shooting star (bearish reversal)
	if current.Close < current.Open { // Red candle
		body := current.Open - current.Close
		upperWick := current.High - current.Open
		lowerWick := current.Close - current.Low
		
		if upperWick > body*2 && lowerWick < body*0.5 && prev.Close > prev.Open {
			signals = append(signals, PatternSignal{
				Pattern:       PatternShootingStar,
				Confidence:    0.65 + (upperWick/body)/10,
				Direction:     -0.7, // Bearish
				StartIndex:    idx-1,
				EndIndex:      idx,
				PredictedMove: -1.5,
				TimeFrame:     "1d",
				Timestamp:     time.Now(),
			})
		}
	}
	
	// Detect engulfing patterns
	if current.Close > current.Open && prev.Close < prev.Open { // Current bullish, previous bearish
		if current.Open < prev.Close && current.Close > prev.Open { // Bullish engulfing
			signals = append(signals, PatternSignal{
				Pattern:       PatternEngulfing,
				Confidence:    0.7,
				Direction:     0.8, // Bullish
				StartIndex:    idx-1,
				EndIndex:      idx,
				PredictedMove: 2.0,
				TimeFrame:     "1d",
				Timestamp:     time.Now(),
			})
		}
	} else if current.Close < current.Open && prev.Close > prev.Open { // Current bearish, previous bullish
		if current.Open > prev.Close && current.Close < prev.Open { // Bearish engulfing
			signals = append(signals, PatternSignal{
				Pattern:       PatternEngulfing,
				Confidence:    0.7,
				Direction:     -0.8, // Bearish
				StartIndex:    idx-1,
				EndIndex:      idx,
				PredictedMove: -2.0,
				TimeFrame:     "1d",
				Timestamp:     time.Now(),
			})
		}
	}
	
	// Detect morning star (bullish reversal)
	if current.Close > current.Open && prevPrev.Close < prevPrev.Open { // Current bullish, prev-prev bearish
		// Middle candle should have small body
		middleBodySize := math.Abs(prev.Open - prev.Close)
		middleTotalRange := prev.High - prev.Low
		
		if middleTotalRange > 0 && middleBodySize/middleTotalRange < 0.3 && 
		   current.Close > (prevPrev.Open+prevPrev.Close)/2 {
			signals = append(signals, PatternSignal{
				Pattern:       PatternMorningStarDoji,
				Confidence:    0.75,
				Direction:     0.85, // Strongly bullish
				StartIndex:    idx-2,
				EndIndex:      idx,
				PredictedMove: 2.5,
				TimeFrame:     "1d",
				Timestamp:     time.Now(),
			})
		}
	}
	
	// Detect evening star (bearish reversal)
	if current.Close < current.Open && prevPrev.Close > prevPrev.Open { // Current bearish, prev-prev bullish
		// Middle candle should have small body
		middleBodySize := math.Abs(prev.Open - prev.Close)
		middleTotalRange := prev.High - prev.Low
		
		if middleTotalRange > 0 && middleBodySize/middleTotalRange < 0.3 &&
		   current.Close < (prevPrev.Open+prevPrev.Close)/2 {
			signals = append(signals, PatternSignal{
				Pattern:       PatternEveningStarDoji,
				Confidence:    0.75,
				Direction:     -0.85, // Strongly bearish
				StartIndex:    idx-2,
				EndIndex:      idx,
				PredictedMove: -2.5,
				TimeFrame:     "1d",
				Timestamp:     time.Now(),
			})
		}
	}
	
	return signals
}
