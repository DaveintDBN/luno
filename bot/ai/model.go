package ai

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

// SignalFeature represents a feature used for ML signal generation
type SignalFeature struct {
	Name  string
	Value float64
	Weight float64
}

// OpportunityScore represents a ranked trading opportunity
type OpportunityScore struct {
	Pair           string
	Score          float64
	Confidence     float64
	Features       []SignalFeature
	Timestamp      time.Time
	RecommendedAction string // "buy", "sell", "hold"
	PredictedMovement float64 // % change predicted
}

// MLModel represents the core machine learning model for signal reinforcement
type MLModel struct {
	modelWeights    map[string]float64
	featureHistory  map[string][]float64
	modelLock       sync.RWMutex
	lastUpdate      time.Time
	optimizationRun bool
}

// NewMLModel creates a new machine learning model
func NewMLModel() *MLModel {
	return &MLModel{
		modelWeights:   initializeDefaultWeights(),
		featureHistory: make(map[string][]float64),
		lastUpdate:     time.Now(),
	}
}

// initializeDefaultWeights sets up initial feature weights
func initializeDefaultWeights() map[string]float64 {
	return map[string]float64{
		"rsi":              0.8,
		"macd":             0.7,
		"bollinger":        0.6,
		"volume":           0.5,
		"priceAction":      0.9,
		"volatility":       0.4,
		"momentum":         0.7,
		"trendStrength":    0.8,
		"marketSentiment":  0.5,
		"newsImpact":       0.3,
		"whaleMovement":    0.4,
		"exchangeInflows":  0.3,
		"fundingRate":      0.5,
		"onChainActivity": 0.4,
	}
}

// ScoreOpportunity evaluates a potential trading opportunity
func (m *MLModel) ScoreOpportunity(pair string, features []SignalFeature) *OpportunityScore {
	m.modelLock.RLock()
	defer m.modelLock.RUnlock()

	var totalScore float64
	var totalWeight float64
	
	// Apply model weights to features
	for i := range features {
		if weight, exists := m.modelWeights[features[i].Name]; exists {
			features[i].Weight = weight
			weightedValue := features[i].Value * weight
			totalScore += weightedValue
			totalWeight += weight
		}
	}
	
	// Normalize score
	var normalizedScore float64
	if totalWeight > 0 {
		normalizedScore = totalScore / totalWeight
	}
	
	// Calculate confidence based on feature alignment
	var featureVariance float64
	for i := range features {
		featureVariance += math.Pow(features[i].Value - normalizedScore, 2)
	}
	
	// Lower variance = higher confidence
	confidence := 1.0
	if len(features) > 1 {
		confidence = 1.0 - math.Min(1.0, math.Sqrt(featureVariance/float64(len(features))))
	}
	
	// Determine recommended action
	action := "hold"
	if normalizedScore > 0.7 {
		action = "buy"
	} else if normalizedScore < 0.3 {
		action = "sell"
	}
	
	// Predict movement magnitude (simplified formula, would be replaced with actual ML prediction)
	predictedMovement := (normalizedScore - 0.5) * 5.0 // Scale to ±2.5%
	
	return &OpportunityScore{
		Pair:              pair,
		Score:             normalizedScore,
		Confidence:        confidence,
		Features:          features,
		Timestamp:         time.Now(),
		RecommendedAction: action,
		PredictedMovement: predictedMovement,
	}
}

// RankOpportunities ranks multiple opportunities by score and confidence
func (m *MLModel) RankOpportunities(opportunities []*OpportunityScore) []*OpportunityScore {
	// Create a copy to avoid modifying original
	rankedOpps := make([]*OpportunityScore, len(opportunities))
	copy(rankedOpps, opportunities)
	
	// Sort by combined score and confidence
	sort.Slice(rankedOpps, func(i, j int) bool {
		// Use both score and confidence as ranking factors
		iRank := rankedOpps[i].Score * (0.7 + 0.3*rankedOpps[i].Confidence)
		jRank := rankedOpps[j].Score * (0.7 + 0.3*rankedOpps[j].Confidence)
		return iRank > jRank
	})
	
	return rankedOpps
}

// UpdateModelWeight updates a specific feature weight based on performance feedback
func (m *MLModel) UpdateModelWeight(featureName string, performanceFeedback float64) {
	m.modelLock.Lock()
	defer m.modelLock.Unlock()
	
	if weight, exists := m.modelWeights[featureName]; exists {
		// Adjust weight (limit between 0.1 and 1.0)
		newWeight := weight + 0.1*performanceFeedback
		m.modelWeights[featureName] = math.Max(0.1, math.Min(1.0, newWeight))
	}
}

// AddFeatureObservation records feature values for historical analysis
func (m *MLModel) AddFeatureObservation(featureName string, value float64) {
	m.modelLock.Lock()
	defer m.modelLock.Unlock()
	
	if _, exists := m.featureHistory[featureName]; !exists {
		m.featureHistory[featureName] = make([]float64, 0, 1000)
	}
	
	// Add to history, maintaining a maximum size
	history := m.featureHistory[featureName]
	if len(history) >= 1000 {
		history = history[1:] // Remove oldest
	}
	history = append(history, value)
	m.featureHistory[featureName] = history
}

// GetModelSummary returns current model weights and performance statistics
func (m *MLModel) GetModelSummary() map[string]interface{} {
	m.modelLock.RLock()
	defer m.modelLock.RUnlock()
	
	// Deep copy weights
	weights := make(map[string]float64)
	for k, v := range m.modelWeights {
		weights[k] = v
	}
	
	// Calculate feature importance
	featureImportance := make(map[string]float64)
	totalWeight := 0.0
	for _, weight := range weights {
		totalWeight += weight
	}
	
	for featureName, weight := range weights {
		if totalWeight > 0 {
			featureImportance[featureName] = weight / totalWeight
		}
	}
	
	return map[string]interface{}{
		"weights":           weights,
		"featureImportance": featureImportance,
		"lastUpdate":        m.lastUpdate,
		"optimizationRun":   m.optimizationRun,
	}
}

// ScheduleOptimization triggers model optimization on a schedule
func (m *MLModel) ScheduleOptimization(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		
		for range ticker.C {
			m.OptimizeModel()
		}
	}()
}

// OptimizeModel performs weight optimization based on historical performance
func (m *MLModel) OptimizeModel() {
	m.modelLock.Lock()
	defer m.modelLock.Unlock()
	
	// This would be replaced with actual optimization logic
	// such as Bayesian optimization or gradient-based methods
	fmt.Println("Optimizing AI model weights based on historical performance...")
	
	// Track optimization run
	m.optimizationRun = true
	m.lastUpdate = time.Now()
	
	// Simple simulation of optimization (in real implementation this would
	// analyze feature correlations with successful trades)
	for feature := range m.modelWeights {
		if len(m.featureHistory[feature]) > 100 {
			// Small random adjustment for demonstration (would be replaced with ML)
			adjustment := (math.Sin(float64(time.Now().Unix())) * 0.05)
			m.modelWeights[feature] = math.Max(0.1, math.Min(1.0, m.modelWeights[feature]+adjustment))
		}
	}
	
	fmt.Println("Model optimization complete.")
}
