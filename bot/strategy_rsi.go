package bot

// RSIStrategy implements an RSI-based trading signal.
type RSIStrategy struct {
	Period     int
	Overbought float64
	Oversold   float64
	prices     []float64
}

// NewRSIStrategy constructs an RSI strategy with the given parameters.
func NewRSIStrategy(period int, overbought, oversold float64) *RSIStrategy {
	if period <= 0 {
		panic("invalid RSI period")
	}
	return &RSIStrategy{Period: period, Overbought: overbought, Oversold: oversold}
}

// Next computes RSI over the last Period data points and returns a Signal.
func (r *RSIStrategy) Next(data MarketData, cfg Config) Signal {
	// mid-price
	price := (data.Bid + data.Ask) / 2
	r.prices = append(r.prices, price)
	if len(r.prices) <= r.Period {
		return SignalNone
	}
	// calculate gains and losses
	var gainSum, lossSum float64
	last := len(r.prices) - 1
	for i := last - r.Period; i < last; i++ {
		delta := r.prices[i+1] - r.prices[i]
		if delta > 0 {
			gainSum += delta
		} else {
			lossSum -= delta
		}
	}
	avgGain := gainSum / float64(r.Period)
	avgLoss := lossSum / float64(r.Period)
	if avgLoss == 0 {
		return SignalNone
	}
	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))
	// overbought => sell, oversold => buy
	if rsi >= r.Overbought {
		return SignalSell
	}
	if rsi <= r.Oversold {
		return SignalBuy
	}
	return SignalNone
}
