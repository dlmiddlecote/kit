package api

import (
	"net/http"

	"github.com/rs/cors"
)

// DefaultCorsMW returns a middleware that adds Cross Origin Resource Sharing
// support for all origins for GET and POST methods.
// This middleware must be used in combination with an OPTIONS method endpoint.
func DefaultCorsMW() Middleware {
	return func(next http.Handler) http.Handler {
		return cors.Default().Handler(next)
	}
}

// AllowAllCorsMW returns a middleware that adds Cross Origin Resource Sharing
// support for all origins, for all methods, with any header or credential.
// This middleware must be used in combination with an OPTIONS method endpoint.
func AllowAllCorsMW() Middleware {
	return func(next http.Handler) http.Handler {
		return cors.AllowAll().Handler(next)
	}
}
