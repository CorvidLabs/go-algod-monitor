// Package config handles loading configuration for the monitor.
package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/CorvidLabs/go-algod-monitor/pkg/health"
)

// Config is the top-level monitor configuration.
type Config struct {
	Nodes       []health.NodeConfig `json:"nodes"`
	IntervalSec int                 `json:"interval_sec"`
	MaxLagSec   int                 `json:"max_lag_sec"`
	Output      string              `json:"output"` // "text" or "json"
}

// DefaultConfig returns a config pointing at the public AlgoNode endpoints.
func DefaultConfig() Config {
	return Config{
		Nodes: []health.NodeConfig{
			{Address: "https://mainnet-api.algonode.cloud", Name: "algonode-mainnet"},
			{Address: "https://testnet-api.algonode.cloud", Name: "algonode-testnet"},
		},
		IntervalSec: 30,
		MaxLagSec:   30,
		Output:      "text",
	}
}

// Load reads a config from a JSON file.
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("reading config %s: %w", path, err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parsing config %s: %w", path, err)
	}
	if len(cfg.Nodes) == 0 {
		return Config{}, fmt.Errorf("config %s: no nodes defined", path)
	}
	return cfg, nil
}

// Validate checks the config for basic issues.
func (c Config) Validate() error {
	if len(c.Nodes) == 0 {
		return fmt.Errorf("no nodes configured")
	}
	for i, n := range c.Nodes {
		if n.Address == "" {
			return fmt.Errorf("node %d: empty address", i)
		}
	}
	if c.IntervalSec < 1 {
		return fmt.Errorf("interval must be >= 1 second")
	}
	return nil
}
