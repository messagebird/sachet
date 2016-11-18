package main

import "github.com/prometheus/client_golang/prometheus"

var (
	requestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sachet_requests_total",
			Help: "How many requests processed, partitioned by status code and provider.",
		},
		[]string{"code", "provider"},
	)
)

func init() {
	prometheus.MustRegister(requestTotal)
}
