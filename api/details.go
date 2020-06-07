package api

import (
	"net/http"
	"time"
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
	StatusCode  int
}

// getDetails returns any Details found within the http.Request, or nil
func getDetails(r *http.Request) *Details {
	v, ok := r.Context().Value(KeyDetails).(*Details)
	if !ok {
		return nil
	}
	return v
}
