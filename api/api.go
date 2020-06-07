package api

import (
	"net/http"
)

type API interface {
	Endpoints() []Endpoint
}

type Endpoint struct {
	Method  string
	Path    string
	Handler http.Handler
}
