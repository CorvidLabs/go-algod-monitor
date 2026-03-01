package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/CorvidLabs/go-algod-monitor/pkg/health"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if len(cfg.Nodes) != 2 {
		t.Fatalf("expected 2 default nodes, got %d", len(cfg.Nodes))
	}
	if cfg.IntervalSec != 30 {
		t.Errorf("expected interval 30, got %d", cfg.IntervalSec)
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("default config should be valid: %v", err)
	}
}

func TestLoad_Valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	data := `{"nodes":[{"address":"http://localhost:4001","token":"abc"}],"interval_sec":10,"max_lag_sec":60,"output":"json"}`
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(cfg.Nodes))
	}
	if cfg.Nodes[0].Token != "abc" {
		t.Errorf("expected token 'abc', got %q", cfg.Nodes[0].Token)
	}
	if cfg.Output != "json" {
		t.Errorf("expected output 'json', got %q", cfg.Output)
	}
}

func TestLoad_NoFile(t *testing.T) {
	_, err := Load("/nonexistent/file.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("{bad json"), 0644)

	_, err := Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLoad_NoNodes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.json")
	os.WriteFile(path, []byte(`{"nodes":[],"interval_sec":10}`), 0644)

	_, err := Load(path)
	if err == nil {
		t.Error("expected error for empty nodes")
	}
}

func TestValidate_EmptyAddress(t *testing.T) {
	cfg := Config{
		Nodes:       []health.NodeConfig{{Address: ""}},
		IntervalSec: 10,
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for empty address")
	}
}

func TestValidate_NoNodes(t *testing.T) {
	cfg := Config{IntervalSec: 10}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for no nodes")
	}
}

func TestValidate_BadInterval(t *testing.T) {
	cfg := Config{
		Nodes:       []health.NodeConfig{{Address: "http://localhost"}},
		IntervalSec: 0,
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero interval")
	}
}
