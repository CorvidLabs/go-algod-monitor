// Package health provides types and functions for checking Algorand node health.
package health

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Status represents the overall health of a node.
type Status int

const (
	StatusHealthy  Status = iota
	StatusDegraded        // node responds but is behind
	StatusDown            // node unreachable or erroring
)

func (s Status) String() string {
	switch s {
	case StatusHealthy:
		return "healthy"
	case StatusDegraded:
		return "degraded"
	case StatusDown:
		return "down"
	default:
		return "unknown"
	}
}

// NodeStatus is the result of a health check against a single algod node.
type NodeStatus struct {
	Address        string        `json:"address"`
	Status         Status        `json:"status"`
	StatusText     string        `json:"status_text"`
	LastRound      uint64        `json:"last_round"`
	TimeSinceRound time.Duration `json:"time_since_round_ns"`
	CatchupTime    time.Duration `json:"catchup_time_ns"`
	Latency        time.Duration `json:"latency_ns"`
	Version        string        `json:"version,omitempty"`
	Error          string        `json:"error,omitempty"`
	CheckedAt      time.Time     `json:"checked_at"`
}

// algodStatusResponse mirrors the /v2/status response from algod.
type algodStatusResponse struct {
	LastRound                 uint64 `json:"last-round"`
	TimeSinceLastRound        uint64 `json:"time-since-last-round"`
	CatchupTime               uint64 `json:"catchup-time"`
	LastVersion               string `json:"last-version"`
	StoppedAtUnsupportedRound bool   `json:"stopped-at-unsupported-round"`
}

// algodVersionsResponse mirrors the /versions response from algod.
type algodVersionsResponse struct {
	GenesisID string `json:"genesis_id"`
	Build     struct {
		Major int `json:"major"`
		Minor int `json:"minor"`
		Build int `json:"build_number"`
	} `json:"build"`
}

// Checker performs health checks against algod nodes.
type Checker struct {
	client      *http.Client
	maxRoundLag time.Duration
}

// NewChecker creates a Checker with the given HTTP client and max acceptable round lag.
// If client is nil, a default client with a 10s timeout is used.
// If maxRoundLag is 0, defaults to 30 seconds.
func NewChecker(client *http.Client, maxRoundLag time.Duration) *Checker {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	if maxRoundLag == 0 {
		maxRoundLag = 30 * time.Second
	}
	return &Checker{client: client, maxRoundLag: maxRoundLag}
}

// Check performs a health check against the given algod node.
func (c *Checker) Check(ctx context.Context, address, token string) NodeStatus {
	ns := NodeStatus{
		Address:   address,
		CheckedAt: time.Now(),
	}

	start := time.Now()

	// Call /v2/status
	statusURL := address + "/v2/status"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, statusURL, nil)
	if err != nil {
		ns.Status = StatusDown
		ns.StatusText = StatusDown.String()
		ns.Error = fmt.Sprintf("creating request: %v", err)
		return ns
	}
	if token != "" {
		req.Header.Set("X-Algo-API-Token", token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		ns.Status = StatusDown
		ns.StatusText = StatusDown.String()
		ns.Error = fmt.Sprintf("connecting: %v", err)
		return ns
	}
	defer resp.Body.Close()
	ns.Latency = time.Since(start)

	if resp.StatusCode != http.StatusOK {
		ns.Status = StatusDown
		ns.StatusText = StatusDown.String()
		ns.Error = fmt.Sprintf("unexpected status code: %d", resp.StatusCode)
		return ns
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ns.Status = StatusDown
		ns.StatusText = StatusDown.String()
		ns.Error = fmt.Sprintf("reading body: %v", err)
		return ns
	}

	var sr algodStatusResponse
	if err := json.Unmarshal(body, &sr); err != nil {
		ns.Status = StatusDown
		ns.StatusText = StatusDown.String()
		ns.Error = fmt.Sprintf("parsing status: %v", err)
		return ns
	}

	ns.LastRound = sr.LastRound
	ns.TimeSinceRound = time.Duration(sr.TimeSinceLastRound) * time.Nanosecond
	ns.CatchupTime = time.Duration(sr.CatchupTime) * time.Nanosecond
	ns.Version = sr.LastVersion

	if sr.StoppedAtUnsupportedRound {
		ns.Status = StatusDown
		ns.StatusText = StatusDown.String()
		ns.Error = "stopped at unsupported round"
		return ns
	}

	if sr.CatchupTime > 0 || ns.TimeSinceRound > c.maxRoundLag {
		ns.Status = StatusDegraded
		ns.StatusText = StatusDegraded.String()
		return ns
	}

	ns.Status = StatusHealthy
	ns.StatusText = StatusHealthy.String()
	return ns
}

// CheckMultiple checks several nodes concurrently and returns all results.
func (c *Checker) CheckMultiple(ctx context.Context, nodes []NodeConfig) []NodeStatus {
	results := make([]NodeStatus, len(nodes))
	done := make(chan int, len(nodes))

	for i, node := range nodes {
		go func(idx int, n NodeConfig) {
			results[idx] = c.Check(ctx, n.Address, n.Token)
			done <- idx
		}(i, node)
	}

	for range nodes {
		<-done
	}

	return results
}

// NodeConfig holds connection info for a single algod node.
type NodeConfig struct {
	Address string `json:"address"`
	Token   string `json:"token"`
	Name    string `json:"name,omitempty"`
}
