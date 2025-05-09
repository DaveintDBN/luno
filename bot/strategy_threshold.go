package bot

// ThresholdStrategy uses bid/ask spread thresholds to generate signals.
type ThresholdStrategy struct{}

// NewThresholdStrategy creates a threshold-based strategy.
func NewThresholdStrategy() *ThresholdStrategy {
	return &ThresholdStrategy{}
}

// Next returns buy/sell/none based on entry/exit thresholds in cfg.
func (s *ThresholdStrategy) Next(data MarketData, cfg Config) Signal {
	bid := data.Bid
	ask := data.Ask
	// Entry
	if cfg.EntryThreshold > 0 && ask > bid*(1+cfg.EntryThreshold) {
		return SignalBuy
	}
	// Exit
	if cfg.ExitThreshold > 0 && bid < ask*(1-cfg.ExitThreshold) {
		return SignalSell
	}
	return SignalNone
}
