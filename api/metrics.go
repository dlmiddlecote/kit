package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func MetricsMW() Middleware {
	// Create service level metrics
	duration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "api_http_latency_seconds",
		Help:    "HTTP Latency distributions",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path", "status"})
	prometheus.MustRegister(duration)

	return func(next http.Handler) http.Handler {
		var h http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				v := getValues(r)
				if v == nil {
					return
				}

				// Observe latency
				duration.WithLabelValues(v.Method, v.RequestPath, fmt.Sprintf("%dXX", v.StatusCode/100)).Observe(time.Since(v.Now).Seconds())
			}()

			next.ServeHTTP(w, r)
		}
		return h
	}
}
