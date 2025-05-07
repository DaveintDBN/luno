package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
