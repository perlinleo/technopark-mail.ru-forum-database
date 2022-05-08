package responses

import "github.com/prometheus/client_golang/prometheus"

type PromMetrics struct {
	Hits    *prometheus.CounterVec
	Timings *prometheus.HistogramVec
}