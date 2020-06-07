package api

import (
	"net/http"
	"time"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyValues is how request values or stored/retrieved.
const KeyValues ctxKey = 1

// Values represent state for each request.
type Values struct {
	Now         time.Time
	TraceID     string
	Method      string
	RequestPath string
	StatusCode  int
}

func getValues(r *http.Request) *Values {
	v, ok := r.Context().Value(KeyValues).(*Values)
	if !ok {
		return nil
	}
	return v
}
