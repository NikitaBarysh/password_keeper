package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Listen(address string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(address, mux)
}

var httpRequestTotal = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "password_keeper",
	Subsystem: "server",
	Name:      "http_request_total",
	Help:      "count request to project",
}, []string{"path", "status"})

var httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: "password_keeper",
	Subsystem: "server",
	Name:      "http_request_duration_seconds",
	Help:      "duration of http request",
}, []string{"path"})

var dbRequestsTotal = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace:  "password_keeper",
	Subsystem:  "server",
	Name:       "db",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
}, []string{"type"})

func IncRequestStatus(path string, status int) {
	httpRequestTotal.WithLabelValues(path, strconv.Itoa(status)).Inc()
}

func IncRequestDuration(path string, d time.Duration) {
	httpRequestDuration.WithLabelValues(path).Observe(d.Seconds())
}

func IncDBRequestStatus(action string, duration time.Duration) {
	dbRequestsTotal.WithLabelValues(action).Observe(duration.Seconds())
}
