package api

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// LogMW returns a middleware that implements request + response detail logging.
// The middleware will log upon response.
func LogMW(logger *zap.SugaredLogger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				// Retrieve detail state of this request
				d := GetDetails(r)
				if d == nil {
					// There's nothing more to do if we can't find the details.
					return
				}

				logger.Infow("request",
					"request_id", d.RequestID,
					"method", d.Method,
					"path", r.URL.Path,
					"status", d.StatusCode,
					"duration", time.Since(d.Now).String(),
				)
			}()
			// Call the wrapped handler
			next.ServeHTTP(w, r)
		})
	}
}
