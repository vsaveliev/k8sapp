// Copyright 2017 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	"net/http"
	"runtime"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/takama/bit"
	"github.com/takama/k8sapp/pkg/version"
)

var (
	totalDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of http requests duration in seconds.",
			Buckets: []float64{0.0001, 0.001, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 7.5, 10},
		}, 
		[]string{"status"},
	)

	totalCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of http requests.",
		},
		[]string{"status"},
	)

	buildTimestamp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "build_timestamp",
			Help: "Build timestamp with additional info.",
		},
		[]string{"go_version", "version", "commit"},
	)
)

// MetricsFunc returns a func for work with Prometheus
func (h *Handler) MetricsFunc() func(c bit.Control) {
	handler := promhttp.Handler()

	return func(c bit.Control) {
		c.Code(http.StatusOK)
		handler.ServeHTTP(c, c.Request())
	}
}

func setConstMetricsValues() {
	time, err := strconv.ParseFloat(version.BUILD_TIMESTAMP, 64)
    if err != nil {
		return
	}
	
	buildTimestamp.With(prometheus.Labels{
		"go_version": runtime.Version(),
		"version": version.RELEASE,
		"commit": version.COMMIT,
	}).Set(time)
}

func init() {
	prometheus.MustRegister(totalDuration)
	prometheus.MustRegister(totalCounter)
	prometheus.MustRegister(buildTimestamp)

	setConstMetricsValues()
}