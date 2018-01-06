package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/takama/bit"
	"github.com/takama/k8sapp/pkg/config"
	"github.com/takama/k8sapp/pkg/logger"
	"github.com/takama/k8sapp/pkg/logger/standard"
)

func TestMetrics(t *testing.T) {
	h := New(standard.New(&logger.Config{}), new(config.Config))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Base(h.Root)(bit.NewControl(w, r))
	})

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}

	trw := httptest.NewRecorder()
	handler.ServeHTTP(trw, req)

	metricsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Base(h.MetricsFunc())(bit.NewControl(w, r))
	})

	req, err = http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}

	trw = httptest.NewRecorder()
	metricsHandler.ServeHTTP(trw, req)

	metrics := trw.Body.String()

	if !strings.Contains(metrics, "http_request_duration_seconds") {
		t.Fatalf("Cannot find metrics of request durations for service")
	}

	if !strings.Contains(metrics, "http_requests_total") {
		t.Fatalf("Cannot find metrics of response statuses for service")
	}
}
