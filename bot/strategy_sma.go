package bot

// SMAStrategy implements a simple moving average crossover strategy.
type SMAStrategy struct {
	ShortWindow int
	LongWindow  int
	shortBuf    []float64
	longBuf     []float64
	shortSum    float64
	longSum     float64
}

// NewSMAStrategy returns a new SMAStrategy. shortWindow must be < longWindow.
func NewSMAStrategy(shortWindow, longWindow int) *SMAStrategy {
	if shortWindow <= 0 || longWindow <= 0 || shortWindow >= longWindow {
		panic("invalid SMA window sizes")
	}
	return &SMAStrategy{ShortWindow: shortWindow, LongWindow: longWindow}
}

// Next processes a new MarketData and returns a Signal.
func (s *SMAStrategy) Next(data MarketData, cfg Config) Signal {
	// Use mid-price
	price := (data.Bid + data.Ask) / 2

	// Update short SMA buffer
	s.shortBuf = append(s.shortBuf, price)
	s.shortSum += price
	if len(s.shortBuf) > s.ShortWindow {
		s.shortSum -= s.shortBuf[0]
		s.shortBuf = s.shortBuf[1:]
	}

	// Update long SMA buffer
	s.longBuf = append(s.longBuf, price)
	s.longSum += price
	if len(s.longBuf) > s.LongWindow {
		s.longSum -= s.longBuf[0]
		s.longBuf = s.longBuf[1:]
	}

	// Not enough data yet
	if len(s.longBuf) < s.LongWindow {
		return SignalNone
	}

	shortAvg := s.shortSum / float64(s.ShortWindow)
	longAvg := s.longSum / float64(s.LongWindow)

	// Entry
	if shortAvg > longAvg+cfg.EntryThreshold {
		return SignalBuy
	}
	// Exit
	if shortAvg < longAvg-cfg.ExitThreshold {
		return SignalSell
	}
	return SignalNone
}
