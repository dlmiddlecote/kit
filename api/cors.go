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
// This middleware must be used in combination with an OPTIONS method endpoint.
func DefaultCorsMW() CorsMiddleware {
	return func(next http.Handler) http.Handler {
		return cors.Default().Handler(next)
	}
}

// AllowAllCorsMW returns a cors middleware that adds support for all origins,
// for all methods, with any header or credential.
// This middleware must be used in combination with an OPTIONS method endpoint.
func AllowAllCorsMW() CorsMiddleware {
	return func(next http.Handler) http.Handler {
		return cors.AllowAll().Handler(next)
	}
}
