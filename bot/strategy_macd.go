package bot

import (
	// Standard library imports
)

// MACDStrategy uses the MACD indicator for buy/sell signals.
type MACDStrategy struct {
	FastPeriod   int
	SlowPeriod   int
	SignalPeriod int

	emaFast   float64
	emaSlow   float64
	emaSignal float64
	initialized bool
}

// NewMACDStrategy constructs a MACD strategy with given EMA periods.
func NewMACDStrategy(fast, slow, signal int) *MACDStrategy {
	if fast <= 0 || slow <= 0 || signal <= 0 {
		panic("invalid MACD periods")
	}
	return &MACDStrategy{FastPeriod: fast, SlowPeriod: slow, SignalPeriod: signal}
}

// Next updates EMA values and returns a signal: buy if MACD > signal, sell if MACD < signal.
func (m *MACDStrategy) Next(data MarketData, cfg Config) Signal {
	price := (data.Bid + data.Ask) / 2
	// Initialize EMAs on first iteration
	if !m.initialized {
		m.emaFast = price
		m.emaSlow = price
		m.emaSignal = 0
		m.initialized = true
		return SignalNone
	}
	// EMA smoothing constants
	alphaFast := 2.0 / float64(m.FastPeriod+1)
	alphaSlow := 2.0 / float64(m.SlowPeriod+1)
	alphaSignal := 2.0 / float64(m.SignalPeriod+1)
	// Update EMAs
	m.emaFast = alphaFast*price + (1-alphaFast)*m.emaFast
	m.emaSlow = alphaSlow*price + (1-alphaSlow)*m.emaSlow
	// MACD line
	macd := m.emaFast - m.emaSlow
	// Signal line
	m.emaSignal = alphaSignal*macd + (1-alphaSignal)*m.emaSignal
	// Generate signal
	if macd > m.emaSignal {
		return SignalBuy
	}
	if macd < m.emaSignal {
		return SignalSell
	}
	return SignalNone
}
