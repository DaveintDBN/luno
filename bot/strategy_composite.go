package bot

// CompositeStrategy combines multiple strategies and signals only when all agree.
type CompositeStrategy struct {
	strategies []Strategy
}

// NewCompositeStrategy constructs a CompositeStrategy from given sub-strategies.
func NewCompositeStrategy(strats ...Strategy) *CompositeStrategy {
	return &CompositeStrategy{strategies: strats}
}

// Next returns SignalBuy if all sub-strategies return buy, SignalSell if all return sell, else SignalNone.
func (c *CompositeStrategy) Next(data MarketData, cfg Config) Signal {
	buyCount, sellCount := 0, 0
	n := len(c.strategies)
	for _, strat := range c.strategies {
		sig := strat.Next(data, cfg)
		if sig == SignalBuy {
			buyCount++
		} else if sig == SignalSell {
			sellCount++
		}
	}
	if buyCount == n {
		return SignalBuy
	}
	if sellCount == n {
		return SignalSell
	}
	return SignalNone
}
