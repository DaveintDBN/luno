package config

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

// Config matches the JSON schema for trading parameters.
type Config struct {
	Pair             string        `json:"pair"`
	EntryThreshold   float64       `json:"entry_threshold"`
	ExitThreshold    float64       `json:"exit_threshold"`
	StakeSize        float64       `json:"stake_size"`
	Cooldown         time.Duration `json:"cooldown"`
	PositionLimit    float64       `json:"position_limit"`
	MaxDrawdown      float64       `json:"max_drawdown"`
	ShortWindow      int           `json:"short_window"`
	LongWindow       int           `json:"long_window"`
	BaseAccountId    int64         `json:"base_account_id"`
	CounterAccountId int64         `json:"counter_account_id"`
	RSIPeriod        int           `json:"rsi_period"`
	RSIOverBought    float64       `json:"rsi_overbought"`
	RSIOverSold      float64       `json:"rsi_oversold"`
	// MACD indicator params
	MACDFastPeriod   int `json:"macd_fast_period"`
	MACDSlowPeriod   int `json:"macd_slow_period"`
	MACDSignalPeriod int `json:"macd_signal_period"`
	// Bollinger Bands params
	BBPeriod     int     `json:"bb_period"`
	BBMultiplier float64 `json:"bb_multiplier"`
	// Risk & execution parameters
	InitialEquity       float64 `json:"initial_equity"`
	PositionSizerType   string  `json:"position_sizer_type"`
	KellyWinProb        float64 `json:"kelly_win_prob"`
	KellyWinLossRatio   float64 `json:"kelly_win_loss_ratio"`
	TWAPSlices          int     `json:"twap_slices"`
	TWAPIntervalSeconds int     `json:"twap_interval_seconds"`
	// VWAP parameters
	VWAPSource               string  `json:"vwap_source"`
	VWAPHistoryWindowMinutes int     `json:"vwap_history_window_minutes"`
	VWAPOrderbookDepthLevels int     `json:"vwap_orderbook_depth_levels"`
	VWAPHybridWeight         float64 `json:"vwap_hybrid_weight"`
	DBPath                   string  `json:"db_path"`
}

// StateStore persists and retrieves bot configuration.
type StateStore interface {
	LoadConfig() (*Config, error)
	SaveConfig(*Config) error
}

// JSONStateStore implements persistence of Config to a JSON file.
type JSONStateStore struct {
	Path string
}

// NewStateStore returns a StateStore backed by the given file path.
func NewStateStore(path string) *JSONStateStore {
	return &JSONStateStore{Path: path}
}

// LoadConfig reads and unmarshals the config JSON file.
func (s *JSONStateStore) LoadConfig() (*Config, error) {
	data, err := ioutil.ReadFile(s.Path)
	if err != nil {
		return nil, err
	}
	// intermediate to parse duration as string
	type raw struct {
		Pair                     string  `json:"pair"`
		EntryThreshold           float64 `json:"entry_threshold"`
		ExitThreshold            float64 `json:"exit_threshold"`
		StakeSize                float64 `json:"stake_size"`
		Cooldown                 string  `json:"cooldown"`
		PositionLimit            float64 `json:"position_limit"`
		MaxDrawdown              float64 `json:"max_drawdown"`
		ShortWindow              int     `json:"short_window"`
		LongWindow               int     `json:"long_window"`
		BaseAccountId            int64   `json:"base_account_id"`
		CounterAccountId         int64   `json:"counter_account_id"`
		RSIPeriod                int     `json:"rsi_period"`
		RSIOverBought            float64 `json:"rsi_overbought"`
		RSIOverSold              float64 `json:"rsi_oversold"`
		MACDFastPeriod           int     `json:"macd_fast_period"`
		MACDSlowPeriod           int     `json:"macd_slow_period"`
		MACDSignalPeriod         int     `json:"macd_signal_period"`
		BBPeriod                 int     `json:"bb_period"`
		BBMultiplier             float64 `json:"bb_multiplier"`
		InitialEquity            float64 `json:"initial_equity"`
		PositionSizerType        string  `json:"position_sizer_type"`
		KellyWinProb             float64 `json:"kelly_win_prob"`
		KellyWinLossRatio        float64 `json:"kelly_win_loss_ratio"`
		TWAPSlices               int     `json:"twap_slices"`
		TWAPIntervalSeconds      int     `json:"twap_interval_seconds"`
		VWAPSource               string  `json:"vwap_source"`
		VWAPHistoryWindowMinutes int     `json:"vwap_history_window_minutes"`
		VWAPOrderbookDepthLevels int     `json:"vwap_orderbook_depth_levels"`
		VWAPHybridWeight         float64 `json:"vwap_hybrid_weight"`
		DBPath                   string  `json:"db_path"`
	}
	var r raw
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	dur, err := time.ParseDuration(r.Cooldown)
	if err != nil {
		return nil, err
	}
	cfg := &Config{
		Pair:                     r.Pair,
		EntryThreshold:           r.EntryThreshold,
		ExitThreshold:            r.ExitThreshold,
		StakeSize:                r.StakeSize,
		Cooldown:                 dur,
		PositionLimit:            r.PositionLimit,
		MaxDrawdown:              r.MaxDrawdown,
		ShortWindow:              r.ShortWindow,
		LongWindow:               r.LongWindow,
		BaseAccountId:            r.BaseAccountId,
		CounterAccountId:         r.CounterAccountId,
		RSIPeriod:                r.RSIPeriod,
		RSIOverBought:            r.RSIOverBought,
		RSIOverSold:              r.RSIOverSold,
		MACDFastPeriod:           r.MACDFastPeriod,
		MACDSlowPeriod:           r.MACDSlowPeriod,
		MACDSignalPeriod:         r.MACDSignalPeriod,
		BBPeriod:                 r.BBPeriod,
		BBMultiplier:             r.BBMultiplier,
		InitialEquity:            r.InitialEquity,
		PositionSizerType:        r.PositionSizerType,
		KellyWinProb:             r.KellyWinProb,
		KellyWinLossRatio:        r.KellyWinLossRatio,
		TWAPSlices:               r.TWAPSlices,
		TWAPIntervalSeconds:      r.TWAPIntervalSeconds,
		VWAPSource:               r.VWAPSource,
		VWAPHistoryWindowMinutes: r.VWAPHistoryWindowMinutes,
		VWAPOrderbookDepthLevels: r.VWAPOrderbookDepthLevels,
		VWAPHybridWeight:         r.VWAPHybridWeight,
		DBPath:                   r.DBPath,
	}
	return cfg, nil
}

// SaveConfig marshals and writes the Config back to the JSON file.
func (s *JSONStateStore) SaveConfig(cfg *Config) error {
	type raw struct {
		Pair                     string  `json:"pair"`
		EntryThreshold           float64 `json:"entry_threshold"`
		ExitThreshold            float64 `json:"exit_threshold"`
		StakeSize                float64 `json:"stake_size"`
		Cooldown                 string  `json:"cooldown"`
		PositionLimit            float64 `json:"position_limit"`
		MaxDrawdown              float64 `json:"max_drawdown"`
		ShortWindow              int     `json:"short_window"`
		LongWindow               int     `json:"long_window"`
		BaseAccountId            int64   `json:"base_account_id"`
		CounterAccountId         int64   `json:"counter_account_id"`
		RSIPeriod                int     `json:"rsi_period"`
		RSIOverBought            float64 `json:"rsi_overbought"`
		RSIOverSold              float64 `json:"rsi_oversold"`
		MACDFastPeriod           int     `json:"macd_fast_period"`
		MACDSlowPeriod           int     `json:"macd_slow_period"`
		MACDSignalPeriod         int     `json:"macd_signal_period"`
		BBPeriod                 int     `json:"bb_period"`
		BBMultiplier             float64 `json:"bb_multiplier"`
		InitialEquity            float64 `json:"initial_equity"`
		PositionSizerType        string  `json:"position_sizer_type"`
		KellyWinProb             float64 `json:"kelly_win_prob"`
		KellyWinLossRatio        float64 `json:"kelly_win_loss_ratio"`
		TWAPSlices               int     `json:"twap_slices"`
		TWAPIntervalSeconds      int     `json:"twap_interval_seconds"`
		VWAPSource               string  `json:"vwap_source"`
		VWAPHistoryWindowMinutes int     `json:"vwap_history_window_minutes"`
		VWAPOrderbookDepthLevels int     `json:"vwap_orderbook_depth_levels"`
		VWAPHybridWeight         float64 `json:"vwap_hybrid_weight"`
		DBPath                   string  `json:"db_path"`
	}
	r := raw{
		Pair:                     cfg.Pair,
		EntryThreshold:           cfg.EntryThreshold,
		ExitThreshold:            cfg.ExitThreshold,
		StakeSize:                cfg.StakeSize,
		Cooldown:                 cfg.Cooldown.String(),
		PositionLimit:            cfg.PositionLimit,
		MaxDrawdown:              cfg.MaxDrawdown,
		ShortWindow:              cfg.ShortWindow,
		LongWindow:               cfg.LongWindow,
		BaseAccountId:            cfg.BaseAccountId,
		CounterAccountId:         cfg.CounterAccountId,
		RSIPeriod:                cfg.RSIPeriod,
		RSIOverBought:            cfg.RSIOverBought,
		RSIOverSold:              cfg.RSIOverSold,
		MACDFastPeriod:           cfg.MACDFastPeriod,
		MACDSlowPeriod:           cfg.MACDSlowPeriod,
		MACDSignalPeriod:         cfg.MACDSignalPeriod,
		BBPeriod:                 cfg.BBPeriod,
		BBMultiplier:             cfg.BBMultiplier,
		InitialEquity:            cfg.InitialEquity,
		PositionSizerType:        cfg.PositionSizerType,
		KellyWinProb:             cfg.KellyWinProb,
		KellyWinLossRatio:        cfg.KellyWinLossRatio,
		TWAPSlices:               cfg.TWAPSlices,
		TWAPIntervalSeconds:      cfg.TWAPIntervalSeconds,
		VWAPSource:               cfg.VWAPSource,
		VWAPHistoryWindowMinutes: cfg.VWAPHistoryWindowMinutes,
		VWAPOrderbookDepthLevels: cfg.VWAPOrderbookDepthLevels,
		VWAPHybridWeight:         cfg.VWAPHybridWeight,
		DBPath:                   cfg.DBPath,
	}
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(s.Path, data, 0644)
}
