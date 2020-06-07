package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// MetricsMW returns a middleware that implements counting + timing of requests
// using a Prometheus Histogram
func MetricsMW() Middleware {
	// Create Histogram that will observe request latency.
	// This Histogram will also expose a 'count' metric that can be used
	// to rate requests.
	duration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "api_http_latency_seconds",
		Help:    "HTTP Latency distributions",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path", "status"})

	// Register the Histogram to be exposed via the Prometheus metrics handler
	prometheus.MustRegister(duration)

	return func(next http.Handler) http.Handler {
		var h http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				// Retrieve detail state of this request
				d := getDetails(r)
				if d == nil {
					// There's nothing more to do if we can't find the details.
					return
				}

				// Calculate the 'group' of the status code, i.e. 2XX, 3XX etc.
				statusGroup := fmt.Sprintf("%dXX", d.StatusCode/100)

				// Observe latency of request
				duration.WithLabelValues(d.Method, d.RequestPath, statusGroup).Observe(time.Since(d.Now).Seconds())
			}()
			// Call the wrapped handler
			next.ServeHTTP(w, r)
		}
		return h
	}
}
