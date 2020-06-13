package api

import (
	"context"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/segmentio/ksuid"
)

// ctxKey represents the type of value for the context key
type ctxKey int

// KeyDetails is how request details are stored and retrieved
const KeyDetails ctxKey = 1

// Details represent state for each request
type Details struct {
	Now         time.Time
	RequestID   string
	Method      string
	RequestPath string
	Params      httprouter.Params
	StatusCode  int
}

// setDetails adds the required Details into the given request's context. The returned request should then be used.
func setDetails(r *http.Request, path string, params httprouter.Params) *http.Request {

	d := Details{
		Now:         time.Now(),
		RequestID:   ksuid.New().String(), // TODO: Get from request to add some request tracing.
		Method:      r.Method,
		RequestPath: path,
		Params:      params,
	}

	// Add details to the context, so other functions can access them.
	ctx := context.WithValue(r.Context(), KeyDetails, &d)

	return r.WithContext(ctx)
}

// getDetails returns any Details found within the http.Request, or nil
func getDetails(r *http.Request) *Details {
	v, ok := r.Context().Value(KeyDetails).(*Details)
	if !ok {
		return nil
	}
	return v
}

// URLParam returns the named parameter from the request's URL path.
func URLParam(r *http.Request, name string) string {
	d := getDetails(r)
	if d == nil {
		return ""
	}
	return d.Params.ByName(name)
}
