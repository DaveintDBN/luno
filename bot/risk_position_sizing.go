package bot

import (
	"math"
)

// PositionSizer defines how much to stake per trade.
type PositionSizer interface {
	Size(equity float64, cfg Config) float64
}

// KellySizer uses the Kelly criterion for position sizing.
type KellySizer struct {
	WinProb float64 // probability of a winning trade
	WinLoss float64 // average win/loss ratio
}

// Size computes optimal fraction of equity per Kelly.
func (k *KellySizer) Size(equity float64, cfg Config) float64 {
	f := k.WinProb - (1-k.WinProb)/k.WinLoss
	// cap by stake size
	return math.Max(0, math.Min(f*equity, cfg.StakeSize))
}

// FixedSizer always uses cfg.StakeSize.
type FixedSizer struct{}

func (f *FixedSizer) Size(equity float64, cfg Config) float64 {
	return cfg.StakeSize
}
