package api

import (
	"net/http"
)

type Endpointer interface {
	Endpoints() []Endpoint
}

type Endpoint struct {
	Method  string
	Path    string
	Handler http.Handler
}
