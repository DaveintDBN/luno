package config

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

// Config matches the JSON schema for trading parameters.
type Config struct {
	Pair           string        `json:"pair"`
	EntryThreshold float64       `json:"entry_threshold"`
	ExitThreshold  float64       `json:"exit_threshold"`
	StakeSize      float64       `json:"stake_size"`
	Cooldown       time.Duration `json:"cooldown"`
	PositionLimit  float64       `json:"position_limit"`
	MaxDrawdown    float64       `json:"max_drawdown"`
	ShortWindow    int           `json:"short_window"`
	LongWindow     int           `json:"long_window"`
	BaseAccountId    int64         `json:"base_account_id"`
	CounterAccountId int64         `json:"counter_account_id"`
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
		Pair           string  `json:"pair"`
		EntryThreshold float64 `json:"entry_threshold"`
		ExitThreshold  float64 `json:"exit_threshold"`
		StakeSize      float64 `json:"stake_size"`
		Cooldown       string  `json:"cooldown"`
		PositionLimit  float64 `json:"position_limit"`
		MaxDrawdown    float64 `json:"max_drawdown"`
		ShortWindow    int     `json:"short_window"`
		LongWindow     int     `json:"long_window"`
		BaseAccountId    int64   `json:"base_account_id"`
		CounterAccountId int64   `json:"counter_account_id"`
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
		Pair:           r.Pair,
		EntryThreshold: r.EntryThreshold,
		ExitThreshold:  r.ExitThreshold,
		StakeSize:      r.StakeSize,
		Cooldown:       dur,
		PositionLimit:  r.PositionLimit,
		MaxDrawdown:    r.MaxDrawdown,
		ShortWindow:    r.ShortWindow,
		LongWindow:     r.LongWindow,
		BaseAccountId:    r.BaseAccountId,
		CounterAccountId: r.CounterAccountId,
	}
	return cfg, nil
}

// SaveConfig marshals and writes the Config back to the JSON file.
func (s *JSONStateStore) SaveConfig(cfg *Config) error {
	type raw struct {
		Pair           string  `json:"pair"`
		EntryThreshold float64 `json:"entry_threshold"`
		ExitThreshold  float64 `json:"exit_threshold"`
		StakeSize      float64 `json:"stake_size"`
		Cooldown       string  `json:"cooldown"`
		PositionLimit  float64 `json:"position_limit"`
		MaxDrawdown    float64 `json:"max_drawdown"`
		ShortWindow    int     `json:"short_window"`
		LongWindow     int     `json:"long_window"`
		BaseAccountId    int64   `json:"base_account_id"`
		CounterAccountId int64   `json:"counter_account_id"`
	}
	r := raw{
		Pair:           cfg.Pair,
		EntryThreshold: cfg.EntryThreshold,
		ExitThreshold:  cfg.ExitThreshold,
		StakeSize:      cfg.StakeSize,
		Cooldown:       cfg.Cooldown.String(),
		PositionLimit:  cfg.PositionLimit,
		MaxDrawdown:    cfg.MaxDrawdown,
		ShortWindow:    cfg.ShortWindow,
		LongWindow:     cfg.LongWindow,
		BaseAccountId:    cfg.BaseAccountId,
		CounterAccountId: cfg.CounterAccountId,
	}
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(s.Path, data, 0644)
}
