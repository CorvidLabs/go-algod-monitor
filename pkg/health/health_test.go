package health

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestStatus_String(t *testing.T) {
	tests := []struct {
		s    Status
		want string
	}{
		{StatusHealthy, "healthy"},
		{StatusDegraded, "degraded"},
		{StatusDown, "down"},
		{Status(99), "unknown"},
	}
	for _, tt := range tests {
		if got := tt.s.String(); got != tt.want {
			t.Errorf("Status(%d).String() = %q, want %q", tt.s, got, tt.want)
		}
	}
}

func makeAlgodHandler(status algodStatusResponse, httpCode int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/status" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(httpCode)
		json.NewEncoder(w).Encode(status)
	}
}

func TestCheck_Healthy(t *testing.T) {
	srv := httptest.NewServer(makeAlgodHandler(algodStatusResponse{
		LastRound:          42000000,
		TimeSinceLastRound: uint64(3 * time.Second),
		CatchupTime:        0,
		LastVersion:        "v2.0",
	}, http.StatusOK))
	defer srv.Close()

	c := NewChecker(srv.Client(), 30*time.Second)
	result := c.Check(context.Background(), srv.URL, "")

	if result.Status != StatusHealthy {
		t.Errorf("expected healthy, got %s (error: %s)", result.StatusText, result.Error)
	}
	if result.LastRound != 42000000 {
		t.Errorf("expected round 42000000, got %d", result.LastRound)
	}
	if result.Latency == 0 {
		t.Error("expected non-zero latency")
	}
	if result.Version != "v2.0" {
		t.Errorf("expected version v2.0, got %s", result.Version)
	}
}

func TestCheck_Degraded_Catchup(t *testing.T) {
	srv := httptest.NewServer(makeAlgodHandler(algodStatusResponse{
		LastRound:   100,
		CatchupTime: 5000000000, // 5s in ns
	}, http.StatusOK))
	defer srv.Close()

	c := NewChecker(srv.Client(), 30*time.Second)
	result := c.Check(context.Background(), srv.URL, "")

	if result.Status != StatusDegraded {
		t.Errorf("expected degraded, got %s", result.StatusText)
	}
}

func TestCheck_Degraded_RoundLag(t *testing.T) {
	srv := httptest.NewServer(makeAlgodHandler(algodStatusResponse{
		LastRound:          100,
		TimeSinceLastRound: uint64(60 * time.Second), // 60s > 30s max lag
	}, http.StatusOK))
	defer srv.Close()

	c := NewChecker(srv.Client(), 30*time.Second)
	result := c.Check(context.Background(), srv.URL, "")

	if result.Status != StatusDegraded {
		t.Errorf("expected degraded, got %s", result.StatusText)
	}
}

func TestCheck_Down_StoppedAtUnsupported(t *testing.T) {
	srv := httptest.NewServer(makeAlgodHandler(algodStatusResponse{
		LastRound:                 100,
		StoppedAtUnsupportedRound: true,
	}, http.StatusOK))
	defer srv.Close()

	c := NewChecker(srv.Client(), 30*time.Second)
	result := c.Check(context.Background(), srv.URL, "")

	if result.Status != StatusDown {
		t.Errorf("expected down, got %s", result.StatusText)
	}
	if result.Error != "stopped at unsupported round" {
		t.Errorf("unexpected error: %s", result.Error)
	}
}

func TestCheck_Down_BadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := NewChecker(srv.Client(), 30*time.Second)
	result := c.Check(context.Background(), srv.URL, "")

	if result.Status != StatusDown {
		t.Errorf("expected down, got %s", result.StatusText)
	}
}

func TestCheck_Down_Unreachable(t *testing.T) {
	c := NewChecker(&http.Client{Timeout: 100 * time.Millisecond}, 30*time.Second)
	result := c.Check(context.Background(), "http://192.0.2.1:9999", "")

	if result.Status != StatusDown {
		t.Errorf("expected down, got %s", result.StatusText)
	}
}

func TestCheck_Down_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not json"))
	}))
	defer srv.Close()

	c := NewChecker(srv.Client(), 30*time.Second)
	result := c.Check(context.Background(), srv.URL, "")

	if result.Status != StatusDown {
		t.Errorf("expected down, got %s", result.StatusText)
	}
}

func TestCheck_TokenHeader(t *testing.T) {
	var gotToken string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotToken = r.Header.Get("X-Algo-API-Token")
		json.NewEncoder(w).Encode(algodStatusResponse{LastRound: 1})
	}))
	defer srv.Close()

	c := NewChecker(srv.Client(), 30*time.Second)
	c.Check(context.Background(), srv.URL, "my-secret-token")

	if gotToken != "my-secret-token" {
		t.Errorf("expected token header 'my-secret-token', got %q", gotToken)
	}
}

func TestCheck_NoTokenHeader(t *testing.T) {
	var gotToken string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotToken = r.Header.Get("X-Algo-API-Token")
		json.NewEncoder(w).Encode(algodStatusResponse{LastRound: 1})
	}))
	defer srv.Close()

	c := NewChecker(srv.Client(), 30*time.Second)
	c.Check(context.Background(), srv.URL, "")

	if gotToken != "" {
		t.Errorf("expected no token header, got %q", gotToken)
	}
}

func TestCheckMultiple(t *testing.T) {
	healthy := httptest.NewServer(makeAlgodHandler(algodStatusResponse{LastRound: 100}, http.StatusOK))
	defer healthy.Close()

	down := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer down.Close()

	c := NewChecker(nil, 30*time.Second)
	nodes := []NodeConfig{
		{Address: healthy.URL, Name: "healthy"},
		{Address: down.URL, Name: "down"},
	}

	results := c.CheckMultiple(context.Background(), nodes)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Status != StatusHealthy {
		t.Errorf("node 0: expected healthy, got %s", results[0].StatusText)
	}
	if results[1].Status != StatusDown {
		t.Errorf("node 1: expected down, got %s", results[1].StatusText)
	}
}

func TestNewChecker_Defaults(t *testing.T) {
	c := NewChecker(nil, 0)
	if c.client == nil {
		t.Error("expected non-nil default client")
	}
	if c.maxRoundLag != 30*time.Second {
		t.Errorf("expected 30s default max lag, got %v", c.maxRoundLag)
	}
}

func TestCheck_Down_BadURL(t *testing.T) {
	c := NewChecker(nil, 30*time.Second)
	result := c.Check(context.Background(), "://invalid", "")

	if result.Status != StatusDown {
		t.Errorf("expected down, got %s", result.StatusText)
	}
}

func TestCheck_ContextCancelled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	c := NewChecker(srv.Client(), 30*time.Second)
	result := c.Check(ctx, srv.URL, "")

	if result.Status != StatusDown {
		t.Errorf("expected down on cancelled context, got %s", result.StatusText)
	}
}
