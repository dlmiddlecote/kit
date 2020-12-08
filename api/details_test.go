package api

import (
	"net/http"
	"testing"

	"github.com/matryer/is"
)

func TestRequestIDNotEmpty(t *testing.T) {
	is := is.New(t)

	// Create request.
	r, err := http.NewRequest("GET", "/foo", nil)
	is.NoErr(err)

	// Set Details on request.
	r = SetDetails(r, "/foo", map[string]string{})

	// Retrieve details from request.
	d := getDetails(r)
	is.True(d != nil) // details found

	// Check request id is not empty.
	if d != nil {
		is.True(d.RequestID != "") // request id is not empty
	}
}

func TestRequestIDFromHeader(t *testing.T) {
	is := is.New(t)

	// Create request.
	r, err := http.NewRequest("GET", "/foo", nil)
	is.NoErr(err)

	// Add request id header.
	r.Header.Set("X-Request-ID", "request-id")

	// Set Details on request.
	r = SetDetails(r, "/foo", map[string]string{})

	// Retrieve details from request.
	d := getDetails(r)
	is.True(d != nil) // details found

	// Check request id.
	if d != nil {
		is.Equal(d.RequestID, "request-id") // request id is as expected
	}
}
