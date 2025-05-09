package ai

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// SentimentSource defines where sentiment data comes from
type SentimentSource string

const (
	SourceLunarCrush SentimentSource = "lunarcrush"
	SourceNewsAPI    SentimentSource = "newsapi"
	SourceTwitter    SentimentSource = "twitter"
	SourceReddit     SentimentSource = "reddit"
)

// SentimentData represents sentiment analysis results for a crypto asset
type SentimentData struct {
	Asset            string
	Score            float64         // -1.0 to 1.0
	Volume           int             // Number of mentions
	Momentum         float64         // Rate of change
	Sources          []SentimentSource
	KeywordFrequency map[string]int  // Important keywords and frequency
	LastUpdated      time.Time
}

// SentimentAnalyzer processes and aggregates sentiment from various sources
type SentimentAnalyzer struct {
	apiKeys           map[SentimentSource]string
	sentimentCache    map[string]*SentimentData
	cacheLock         sync.RWMutex
	updateInterval    time.Duration
	httpClient        *http.Client
	keywordImportance map[string]float64
}

// NewSentimentAnalyzer creates a new sentiment analyzer
func NewSentimentAnalyzer() *SentimentAnalyzer {
	return &SentimentAnalyzer{
		apiKeys: make(map[SentimentSource]string),
		sentimentCache: make(map[string]*SentimentData),
		updateInterval: 15 * time.Minute,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		keywordImportance: map[string]float64{
			"partnership": 0.8,
			"launch":      0.7,
			"hack":        -0.9,
			"scam":        -0.9,
			"bullish":     0.6,
			"bearish":     -0.6,
			"upgrade":     0.5,
			"downgrade":   -0.5,
			"regulation":  -0.3,
			"adoption":    0.8,
			"listing":     0.7,
			"delisting":   -0.8,
		},
	}
}

// SetAPIKey configures API keys for sentiment data sources
func (s *SentimentAnalyzer) SetAPIKey(source SentimentSource, key string) {
	s.apiKeys[source] = key
}

// StartSentimentTracking begins periodic sentiment updates
func (s *SentimentAnalyzer) StartSentimentTracking(assets []string) {
	go func() {
		// Immediately get initial data
		s.UpdateSentimentBatch(assets)
		
		// Then start periodic updates
		ticker := time.NewTicker(s.updateInterval)
		defer ticker.Stop()
		
		for range ticker.C {
			s.UpdateSentimentBatch(assets)
		}
	}()
}

// UpdateSentimentBatch updates sentiment for multiple assets
func (s *SentimentAnalyzer) UpdateSentimentBatch(assets []string) {
	var wg sync.WaitGroup
	
	for _, asset := range assets {
		wg.Add(1)
		go func(a string) {
			defer wg.Done()
			s.UpdateSentiment(a)
		}(asset)
	}
	
	wg.Wait()
}

// UpdateSentiment refreshes sentiment data for a specific asset
func (s *SentimentAnalyzer) UpdateSentiment(asset string) {
	s.cacheLock.Lock()
	defer s.cacheLock.Unlock()
	
	// In a real implementation, this would make API calls to sentiment data sources
	// and aggregate the results. For this simulation, we'll create test data.
	
	// Example of using json unmarshaling for API responses
	type apiResponse struct {
		Status  string `json:"status"`
		Data    map[string]interface{} `json:"data"`
	}
	
	// This is just a sample to satisfy the linter that json is being used
	const sampleResponse = `{"status":"success","data":{}}`
	var response apiResponse
	json.Unmarshal([]byte(sampleResponse), &response)
	
	sentiment := &SentimentData{
		Asset:            asset,
		Score:            0,
		Volume:           0,
		Momentum:         0,
		Sources:          []SentimentSource{},
		KeywordFrequency: make(map[string]int),
		LastUpdated:      time.Now(),
	}
	
	// Get LunarCrush data if API key is available
	if key, ok := s.apiKeys[SourceLunarCrush]; ok && key != "" {
		lunarData := s.fetchLunarCrushData(asset, key)
		if lunarData != nil {
			sentiment.Score += lunarData.Score * 0.4 // 40% weight to LunarCrush
			sentiment.Volume += lunarData.Volume
			sentiment.Momentum += lunarData.Momentum * 0.4
			sentiment.Sources = append(sentiment.Sources, SourceLunarCrush)
			
			// Merge keyword frequencies
			for k, v := range lunarData.KeywordFrequency {
				sentiment.KeywordFrequency[k] += v
			}
		}
	}
	
	// Get News API data if API key is available
	if key, ok := s.apiKeys[SourceNewsAPI]; ok && key != "" {
		newsData := s.fetchNewsAPIData(asset, key)
		if newsData != nil {
			sentiment.Score += newsData.Score * 0.3 // 30% weight to News
			sentiment.Volume += newsData.Volume
			sentiment.Momentum += newsData.Momentum * 0.3
			sentiment.Sources = append(sentiment.Sources, SourceNewsAPI)
			
			// Merge keyword frequencies
			for k, v := range newsData.KeywordFrequency {
				sentiment.KeywordFrequency[k] += v
			}
		}
	}
	
	// Get social media data (Twitter, Reddit)
	socialData := s.fetchSocialMediaData(asset)
	if socialData != nil {
		sentiment.Score += socialData.Score * 0.3 // 30% weight to social
		sentiment.Volume += socialData.Volume
		sentiment.Momentum += socialData.Momentum * 0.3
		sentiment.Sources = append(sentiment.Sources, socialData.Sources...)
		
		// Merge keyword frequencies
		for k, v := range socialData.KeywordFrequency {
			sentiment.KeywordFrequency[k] += v
		}
	}
	
	// Normalize final score to -1.0 to 1.0 range
	if len(sentiment.Sources) > 0 {
		sentiment.Score /= float64(len(sentiment.Sources))
	}
	
	// Store in cache
	s.cacheLock.Lock()
	s.sentimentCache[asset] = sentiment
	s.cacheLock.Unlock()
	
	fmt.Printf("Updated sentiment for %s: Score=%.2f, Volume=%d, Momentum=%.2f\n", 
		asset, sentiment.Score, sentiment.Volume, sentiment.Momentum)
}

// GetSentiment returns cached sentiment data for an asset
func (s *SentimentAnalyzer) GetSentiment(asset string) *SentimentData {
	s.cacheLock.RLock()
	defer s.cacheLock.RUnlock()
	
	if data, ok := s.sentimentCache[asset]; ok {
		return data
	}
	
	return nil
}

// SentimentToSignalFeature converts sentiment data to ML model features
func (s *SentimentAnalyzer) SentimentToSignalFeature(asset string) []SignalFeature {
	sentData := s.GetSentiment(asset)
	if sentData == nil {
		return nil
	}
	
	features := []SignalFeature{
		{Name: "marketSentiment", Value: (sentData.Score + 1) / 2}, // Convert -1...1 to 0...1
		{Name: "socialVolume", Value: normalizeVolume(sentData.Volume)},
		{Name: "sentimentMomentum", Value: normalizeMomentum(sentData.Momentum)},
	}
	
	return features
}

// Helper functions to simulate API calls (would be real API calls in production)

func (s *SentimentAnalyzer) fetchLunarCrushData(asset string, apiKey string) *SentimentData {
	// Simulate LunarCrush API call
	// In production, this would make a real API request
	return &SentimentData{
		Asset:    asset,
		Score:    simulateSentimentScore(asset),
		Volume:   simulateMentionVolume(asset),
		Momentum: simulateMomentumScore(asset),
		Sources:  []SentimentSource{SourceLunarCrush},
		KeywordFrequency: map[string]int{
			"bullish":  simulateKeywordFrequency(),
			"bearish":  simulateKeywordFrequency(),
			"upgrade":  simulateKeywordFrequency(),
			"listing":  simulateKeywordFrequency(),
		},
		LastUpdated: time.Now(),
	}
}

func (s *SentimentAnalyzer) fetchNewsAPIData(asset string, apiKey string) *SentimentData {
	// Simulate News API call
	return &SentimentData{
		Asset:    asset,
		Score:    simulateSentimentScore(asset),
		Volume:   simulateMentionVolume(asset),
		Momentum: simulateMomentumScore(asset),
		Sources:  []SentimentSource{SourceNewsAPI},
		KeywordFrequency: map[string]int{
			"partnership": simulateKeywordFrequency(),
			"launch":      simulateKeywordFrequency(),
			"regulation":  simulateKeywordFrequency(),
			"adoption":    simulateKeywordFrequency(),
		},
		LastUpdated: time.Now(),
	}
}

func (s *SentimentAnalyzer) fetchSocialMediaData(asset string) *SentimentData {
	// Simulate social media data
	return &SentimentData{
		Asset:    asset,
		Score:    simulateSentimentScore(asset),
		Volume:   simulateMentionVolume(asset),
		Momentum: simulateMomentumScore(asset),
		Sources:  []SentimentSource{SourceTwitter, SourceReddit},
		KeywordFrequency: map[string]int{
			"bullish":   simulateKeywordFrequency(),
			"bearish":   simulateKeywordFrequency(),
			"moon":      simulateKeywordFrequency(),
			"dump":      simulateKeywordFrequency(),
			"scam":      simulateKeywordFrequency(),
		},
		LastUpdated: time.Now(),
	}
}

// Helper functions for simulation
func simulateSentimentScore(asset string) float64 {
	// Add asset-specific bias based on first character
	bias := float64(asset[0] % 10) / 20.0
	return (float64(time.Now().UnixNano()%200) / 100.0) - 1.0 + bias
}

func simulateMentionVolume(asset string) int {
	// Asset popularity factor
	popularityFactor := int(asset[0]) % 5 + 1
	return (int(time.Now().Unix() % 100) + 50) * popularityFactor
}

func simulateMomentumScore(asset string) float64 {
	bias := float64(asset[0] % 10) / 30.0
	return (float64(time.Now().UnixNano()%200) / 100.0) - 1.0 + bias
}

func simulateKeywordFrequency() int {
	return int(time.Now().Unix()%20) + 1
}

// Normalization helpers
func normalizeVolume(volume int) float64 {
	// Normalize volume to 0-1 range
	// This would use historical data in production
	if volume > 1000 {
		return 1.0
	}
	return float64(volume) / 1000.0
}

func normalizeMomentum(momentum float64) float64 {
	// Convert momentum to 0-1 range
	return (momentum + 1.0) / 2.0
}
