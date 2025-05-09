package bot

import "github.com/luno/luno-bot/config"

// MultiTimeframeStrategy wraps fast and slow composite strategies.
type MultiTimeframeStrategy struct {
	Fast Strategy
	Slow Strategy
}

// NewMultiTimeframeStrategy builds two composites (fast and slow timeframes) from cfg.
func NewMultiTimeframeStrategy(cfg *config.Config) *MultiTimeframeStrategy {
	// Fast timeframe strategies
	fastStrats := []Strategy{
		NewSMAStrategy(cfg.ShortWindow, cfg.LongWindow),
		NewThresholdStrategy(),
		NewRSIStrategy(cfg.RSIPeriod, cfg.RSIOverBought, cfg.RSIOverSold),
		NewMACDStrategy(cfg.MACDFastPeriod, cfg.MACDSlowPeriod, cfg.MACDSignalPeriod),
		NewBBandsStrategy(cfg.BBPeriod, cfg.BBMultiplier),
	}
	// Slow timeframe: periods doubled, same thresholds
	slowStrats := []Strategy{
		NewSMAStrategy(cfg.ShortWindow*2, cfg.LongWindow*2),
		NewThresholdStrategy(),
		NewRSIStrategy(cfg.RSIPeriod*2, cfg.RSIOverBought, cfg.RSIOverSold),
		NewMACDStrategy(cfg.MACDFastPeriod*2, cfg.MACDSlowPeriod*2, cfg.MACDSignalPeriod*2),
		NewBBandsStrategy(cfg.BBPeriod*2, cfg.BBMultiplier),
	}
	fast := NewCompositeStrategy(fastStrats...)
	slow := NewCompositeStrategy(slowStrats...)
	return &MultiTimeframeStrategy{Fast: fast, Slow: slow}
}

// Next returns a signal only if fast and slow agree, else none.
func (m *MultiTimeframeStrategy) Next(data MarketData, cfg Config) Signal {
	sigFast := m.Fast.Next(data, cfg)
	sigSlow := m.Slow.Next(data, cfg)
	if sigFast == sigSlow {
		return sigFast
	}
	return SignalNone
}
