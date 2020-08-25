package api

import (
	"net/http"
	"strings"
)

// API defines a HTTP API that can be exposed using a server
type API interface {
	// Endpoints must return all Endpoints of the HTTP API to register with a http router
	Endpoints() []Endpoint
}

// Endpoint defines an endpoint of a HTTP API
type Endpoint struct {
	// The HTTP Method of this endpoint. It is possible to register multiple
	// methods by separating them by '+', i.e. GET+OPTIONS.
	Method string
	// The URL Path of this endpoint. Should follow the format for
	// paths specified by https://github.com/julienschmidt/httprouter.
	Path string
	// The handler to invoke when a request for the given Method, Path is received
	Handler http.Handler
	// Any endpoint specific middlewares for this handler (i.e. access control)
	Middlewares []Middleware
	// Flag to suppress endpoint request/response information log line.
	SuppressLogs bool
	// Flag to suppress endpoint appearing in exposed Prometheus metrics.
	SuppressMetrics bool
}

// Methods returns all methods we should register this endpoint for.
func (e Endpoint) Methods() []string {
	return strings.Split(e.Method, "+")
}
