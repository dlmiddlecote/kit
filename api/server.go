package api

import (
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type server struct {
	router *httptreemux.TreeMux
	logger *zap.SugaredLogger
	mw     []Middleware
}

// NewServer returns a HTTP server for accessing the the given API.
func NewServer(addr string, logger *zap.SugaredLogger, a API) http.Server {

	// Create our server
	s := server{
		router: httptreemux.New(),
		logger: logger,
		mw:     make([]Middleware, 0), // Place for default middleware.
	}

	// Gather endpoints to register with metrics middleware.
	// Some endpoints may not wish to be instrumented.
	metricEndpoints := make([]Endpoint, 0)
	for _, e := range a.Endpoints() {
		if !e.SuppressMetrics {
			metricEndpoints = append(metricEndpoints, e)
		}
	}

	// Create metrics middleware.
	metricsmw := MetricsMW(prometheus.DefaultRegisterer, metricEndpoints)

	// Create logging middleware.
	logmw := LogMW(logger)

	// Add all endpoints to the server's router.
	for _, e := range a.Endpoints() {

		methods := []string{e.Method}

		// Gather list of middleware to wrap this endpoint in.
		mws := make([]Middleware, 0)
		if !e.SuppressMetrics {
			// Add metrics middleware if metrics should not be suppressed.
			mws = append(mws, metricsmw)
		}
		if !e.SuppressLogs {
			// Add logging middleware if logs should not be suppressed.
			mws = append(mws, logmw)
		}
		if e.CorsMiddleware != nil {
			// Add cors middleware.
			mws = append(mws, Middleware(*e.CorsMiddleware))
			// Add OPTIONS method to be registered.
			methods = append(methods, "OPTIONS")
		}

		// Add all of the endpoint specific middleware.
		mws = append(mws, e.Middlewares...)

		for _, method := range methods {
			// Register endpoint with the server.
			s.handle(method, e.Path, e.Handler, mws...)
		}
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
	h := func(w http.ResponseWriter, r *http.Request, params map[string]string) {

		// Update request context with the required details to process the request
		r = SetDetails(r, path, params)

		// Call the wrapped handler
		handler.ServeHTTP(w, r)
	}

	// Register the handler to the router
	s.router.Handle(method, path, h)
}

// ServeHTTP implements http.Handler
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
