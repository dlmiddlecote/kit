package api

import (
	"net/http"

	"github.com/rs/cors"
)

// CorsMiddleware is a function designed to run some code before and/or after
// another Handler, related to Cross Origin Resource Sharing.
type CorsMiddleware func(http.Handler) http.Handler

// DefaultCorsMW returns a cors middleware that adds support for all origins
// for GET and POST methods.
func DefaultCorsMW() *CorsMiddleware {
	return CorsMW(cors.Default())
}

// AllowAllCorsMW returns a cors middleware that adds support for all origins,
// for all methods, with any header or credential.
func AllowAllCorsMW() *CorsMiddleware {
	return CorsMW(cors.AllowAll())
}

// CorsMW returns a cors middleware that adds cors support to wrapped handlers,
// as defined by the given cors options.
func CorsMW(c *cors.Cors) *CorsMiddleware {
	var mw CorsMiddleware = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Wrap the response writer with another that captures the status code.
			scrw := &statusCodeResponseWriter{w, http.StatusOK}

			defer func() {
				// Set status code value on request details so other middlewares can access it.
				if d := getDetails(r); d != nil {
					d.StatusCode = scrw.statusCode
				}
			}()

			// Create the CORS handler.
			h := c.Handler(next)

			// Call the CORS handler with the wrapped response writer.
			h.ServeHTTP(scrw, r)
		})
	}
	return &mw
}

// statusCodeResponseWriter is a http.ResponseWriter that can capture the status code written to it.
type statusCodeResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader overrides the underlying ResponseWriter to capture the status code written.
func (w *statusCodeResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}
