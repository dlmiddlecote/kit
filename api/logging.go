package api

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func LogMW(logger *zap.SugaredLogger) Middleware {
	return func(next http.Handler) http.Handler {
		var h http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				v := getValues(r)
				if v == nil {
					return
				}

				logger.Infow("request",
					"request_id", v.TraceID,
					"method", v.Method,
					"path", v.RequestPath,
					"status", v.StatusCode,
					"duration", time.Since(v.Now),
				)

			}()

			next.ServeHTTP(w, r)

		}

		return h
	}
}
