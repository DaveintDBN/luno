package bot

import (
  "math"
)

// BBandsStrategy uses Bollinger Bands for trading signals.
type BBandsStrategy struct {
  Period     int
  Multiplier float64
  prices     []float64
}

// NewBBandsStrategy constructs a BBandsStrategy.
func NewBBandsStrategy(period int, multiplier float64) *BBandsStrategy {
  if period <= 0 || multiplier <= 0 {
    panic("invalid Bollinger Bands parameters")
  }
  return &BBandsStrategy{Period: period, Multiplier: multiplier}
}

// Next calculates bands over the last Period prices and signals based on price
func (b *BBandsStrategy) Next(data MarketData, cfg Config) Signal {
  price := (data.Bid + data.Ask) / 2
  b.prices = append(b.prices, price)
  n := len(b.prices)
  if n < b.Period {
    return SignalNone
  }
  // compute mean and variance
  var sum, sumSq float64
  for i := n - b.Period; i < n; i++ {
    v := b.prices[i]
    sum += v
    sumSq += v * v
  }
  mean := sum / float64(b.Period)
  variance := sumSq/float64(b.Period) - mean*mean
  stddev := math.Sqrt(variance)
  upper := mean + b.Multiplier*stddev
  lower := mean - b.Multiplier*stddev
  if price > upper {
    return SignalSell
  }
  if price < lower {
    return SignalBuy
  }
  return SignalNone
}
