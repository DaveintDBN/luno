package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/luno/luno-go"
	"github.com/luno/luno-go/decimal"
)

func TestHealthzEndpoint(t *testing.T) {
	r := SetupRouter(nil, nil, nil, nil, nil)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Healthz returned %d, expected %d", w.Code, http.StatusOK)
	}
}

func TestMetricsEndpoint(t *testing.T) {
	r := SetupRouter(nil, nil, nil, nil, nil)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/metrics", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Metrics returned %d, expected %d", w.Code, http.StatusOK)
	}
}

// --- Fake client and new endpoint tests ---
// fakeClient implements bot.Client for tests.
type fakeClient struct{}

func (f *fakeClient) SetAuth(id, secret string) error { return nil }
func (f *fakeClient) GetTickers(ctx context.Context, req *luno.GetTickersRequest) (*luno.GetTickersResponse, error) {
	return &luno.GetTickersResponse{Tickers: []luno.Ticker{{Pair: "XBTZAR", Bid: decimal.NewFromInt64(100), Ask: decimal.NewFromInt64(110), Rolling24HourVolume: decimal.NewFromInt64(1000000)}}}, nil
}
func (f *fakeClient) GetOrderBook(ctx context.Context, req *luno.GetOrderBookRequest) (*luno.GetOrderBookResponse, error) {
	return &luno.GetOrderBookResponse{}, nil
}
func (f *fakeClient) PostLimitOrder(ctx context.Context, req *luno.PostLimitOrderRequest) (*luno.PostLimitOrderResponse, error) {
	return &luno.PostLimitOrderResponse{}, nil
}
func (f *fakeClient) ListTrades(ctx context.Context, req *luno.ListTradesRequest) (*luno.ListTradesResponse, error) {
	return &luno.ListTradesResponse{}, nil
}
func (f *fakeClient) GetCandles(ctx context.Context, req *luno.GetCandlesRequest) (*luno.GetCandlesResponse, error) {
	return &luno.GetCandlesResponse{}, nil
}
func (f *fakeClient) GetBalances(ctx context.Context, req *luno.GetBalancesRequest) (*luno.GetBalancesResponse, error) {
	return &luno.GetBalancesResponse{}, nil
}

func TestPairsEndpoint(t *testing.T) {
	fc := &fakeClient{}
	r := SetupRouter(nil, fc, nil, nil, nil)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/pairs", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Pairs returned %d, expected %d", w.Code, http.StatusOK)
	}
	var pairs []string
	if err := json.Unmarshal(w.Body.Bytes(), &pairs); err != nil {
		t.Fatalf("Invalid JSON: %v", err)
	}
	if len(pairs) != 1 || pairs[0] != "XBTZAR" {
		t.Errorf("Unexpected pairs: %v", pairs)
	}
}

func TestScanEndpoint(t *testing.T) {
	fc := &fakeClient{}
	r := SetupRouter(nil, fc, nil, nil, nil)
	body := map[string]interface{}{"pairs": []string{"XBTZAR"}, "min_volume": 0, "entry_threshold": 0.05, "exit_threshold": 0.01}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/scan", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Scan returned %d, expected %d", w.Code, http.StatusOK)
	}
	var results []struct {
		Pair   string  `json:"pair"`
		Bid    float64 `json:"bid"`
		Ask    float64 `json:"ask"`
		Signal string  `json:"signal"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &results); err != nil {
		t.Fatalf("Invalid JSON: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("Expected one result, got %d", len(results))
	}
	if results[0].Signal != "buy" {
		t.Errorf("Expected signal buy, got %s", results[0].Signal)
	}
}

// Test continuous auto-scan endpoints
func TestAutoScanEndpoints(t *testing.T) {
	fc := &fakeClient{}
	r := SetupRouter(nil, fc, nil, nil, nil)

	// Start auto-scan
	body := map[string]interface{}{"pairs": []string{"XBTZAR"}, "min_volume": 0, "entry_threshold": 0, "exit_threshold": 0, "interval_seconds": 1, "auto_execute": false}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/autoscan/start", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("StartAutoScan returned %d, expected %d", w.Code, http.StatusOK)
	}
	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Invalid JSON: %v", err)
	}
	if resp["status"] != "auto-scan started" {
		t.Errorf("Expected status 'auto-scan started', got '%s'", resp["status"])
	}

	// Starting again should error
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected BadRequest on double start, got %d", w.Code)
	}

	// Stop auto-scan
	w = httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/autoscan/stop", nil)
	r.ServeHTTP(w, req2)
	if w.Code != http.StatusOK {
		t.Fatalf("StopAutoScan returned %d, expected %d", w.Code, http.StatusOK)
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Invalid JSON: %v", err)
	}
	if resp["status"] != "auto-scan stopped" {
		t.Errorf("Expected status 'auto-scan stopped', got '%s'", resp["status"])
	}

	// Stopping again should error
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req2)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected BadRequest on double stop, got %d", w.Code)
	}
}
