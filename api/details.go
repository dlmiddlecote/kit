package api

import (
	"context"
	"net/http"
	"time"

	"github.com/segmentio/ksuid"
)

// ctxKey represents the type of value for the context key
type ctxKey int

// keyDetails is how request details are stored and retrieved
const keyDetails ctxKey = 1

// Details represent state for each request
type Details struct {
	Now         time.Time
	RequestID   string
	Method      string
	RequestPath string
	Params      map[string]string
	StatusCode  int
}

// SetDetails adds the required setails into the given request's context. The returned request should then be used.
func SetDetails(r *http.Request, path string, params map[string]string) *http.Request {

	d := Details{
		Now:         time.Now(),
		RequestID:   ksuid.New().String(), // TODO: Get from request to add some request tracing.
		Method:      r.Method,
		RequestPath: path,
		Params:      params,
	}

	// Add details to the context, so other functions can access them.
	ctx := context.WithValue(r.Context(), keyDetails, &d)

	return r.WithContext(ctx)
}

// GetDetails returns any details found within the http.Request, or nil
func GetDetails(r *http.Request) *Details {
	v, ok := r.Context().Value(keyDetails).(*Details)
	if !ok {
		return nil
	}
	return v
}

// URLParam returns the named parameter from the request's URL path.
func URLParam(r *http.Request, name string) string {
	d := GetDetails(r)
	if d == nil {
		return ""
	}
	return d.Params[name]
}
