package api

import (
	"net/http"
)

// API defines a HTTP API that can be exposed using a server
type API interface {
	// Endpoints must return all Endpoints of the HTTP API to register with a http router
	Endpoints() []Endpoint
}

// Endpoint defines an endpoint of a HTTP API
type Endpoint struct {
	// The HTTP Method of this endpoint
	Method string
	// The URL Path of this endpoint. Should follow the format for
	// paths specified by https://github.com/julienschmidt/httprouter.
	Path string
	// The handler to invoke when a request for the given Method, Path is received
	Handler http.Handler
	// Any endpoint specific middlewares for this handler (i.e. access control)
	Middlewares []Middleware
}
