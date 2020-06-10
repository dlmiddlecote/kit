package api

import (
	"context"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/segmentio/ksuid"
	"go.uber.org/zap"
)

type server struct {
	router *httprouter.Router
	logger *zap.SugaredLogger
	mw     []Middleware
}

// NewServer returns a HTTP server for accessing the the given API.
func NewServer(addr string, logger *zap.SugaredLogger, a API) http.Server {
	// Create our server, with default middlewares
	s := server{
		router: httprouter.New(),
		logger: logger,
		mw: []Middleware{
			LogMW(logger),
			MetricsMW(a.Endpoints()),
		},
	}

	// Add all endpoints to the server's router
	for _, e := range a.Endpoints() {
		s.handle(e.Method, e.Path, e.Handler, e.Middlewares...)
	}

	// Convert our server into a http.Server
	return http.Server{
		Addr:    addr,
		Handler: &s,
	}
}

// handle registers handlers with the given middleware to the server's router
func (s *server) handle(method, path string, handler http.Handler, mw ...Middleware) {

	// First wrap the handler with its specific middleware
	handler = wrapMiddleware(mw, handler)

	// Then wrap the handler in the server's middleware
	handler = wrapMiddleware(s.mw, handler)

	// Create the function to execute for each request
	h := func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx := r.Context()

		// Set the context with the required details to process the request
		d := Details{
			Now:         time.Now(),
			RequestID:   ksuid.New().String(), // TODO: Get from request to add some request tracing.
			Method:      method,
			RequestPath: path,
			Params:      params,
		}

		// Add details to the context, so other functions can access them.
		ctx = context.WithValue(ctx, KeyDetails, &d)

		// Call the wrapped handler
		handler.ServeHTTP(w, r.WithContext(ctx))
	}

	// Register the handler to the router
	s.router.Handle(method, path, h)
}

// ServeHTTP implements http.Handler
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
