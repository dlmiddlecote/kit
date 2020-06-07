package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type server struct {
	router *httprouter.Router
	logger *zap.SugaredLogger
	mw     []Middleware
}

// NewServer returns a HTTP server for accessing the the given API.
func NewServer(logger *zap.SugaredLogger, a API, addr string) http.Server {
	router := httprouter.New()
	router.HandleOPTIONS = true

	s := server{
		router: router,
		logger: logger,
		mw:     []Middleware{LogMW(logger), MetricsMW()},
	}

	// initialise server router
	for _, e := range a.Endpoints() {
		s.handle(e.Method, e.Path, e.Handler)
	}

	return http.Server{
		Addr:    addr,
		Handler: &s,
	}
}

func (s *server) handle(method, path string, handler http.Handler) {

	handler = wrapMiddleware(s.mw, handler)

	// Create the function to execute for each request
	h := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Set the context with the required values to process the request
		v := Values{
			Now:         time.Now(),
			TraceID:     uuid.New().String(),
			Method:      method,
			RequestPath: path,
		}

		ctx = context.WithValue(ctx, KeyValues, &v)

		handler.ServeHTTP(w, r.WithContext(ctx))
	}

	s.router.HandlerFunc(method, path, h)
}

// ServeHTTP implements http.Handler
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

//
// helpers
//

func Respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	if v := getValues(r); v != nil {
		v.StatusCode = status
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		enc := json.NewEncoder(w)
		err := enc.Encode(data)
		if err == nil {
			enc.Encode(map[string]string{"msg": "Internal Server Error"})
		}
	}
}
