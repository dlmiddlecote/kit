package api

import (
	"encoding/json"
	"net/http"

	"github.com/peterbourgon/mergemap"
)

// Decode should be used to convert the request's JSON body into the given v value.
func Decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	// FUTURE: Can handle different content types by looking at the request headers.
	return json.NewDecoder(r.Body).Decode(v)
}

// Respond should be used to respond to a http request within a http handler.
// Respond encodes any data passed in as JSON.
// Respond also sets the status code of the response on the request details, so
// middlewares can access this value.
func Respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {

	var jsonData []byte
	var err error

	// If we have data to respond with, encode it into JSON, and set the correct
	// header. If we cannot encode, we'll return an Internal Server Error.
	if data != nil {
		// Set the content-type if not already set
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "application/json")
		}

		// Marshal data into byte array
		jsonData, err = json.Marshal(data)
		if err != nil {
			// There was an error Marshalling, so return a server error
			Error(w, r, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	// Set status code value on request details so other middlewares can access it
	if d := getDetails(r); d != nil {
		d.StatusCode = code
	}

	// Set the status code of the response. This should be the last header to be written.
	w.WriteHeader(code)

	// Write the JSON body. This must be done last, otherwise we flush the response too quickly.
	if len(jsonData) > 0 {
		//nolint:errcheck
		w.Write(jsonData)
	}
}

// Redirect replies to the request with a redirect to url.
// The provided code should be in the 3xx range and is usually
// http.StatusMovedPermanently, http.StatusFound or http.StatusSeeOther.
func Redirect(w http.ResponseWriter, r *http.Request, url string, code int) {

	// Set status code value on request details so other middlewares can access it
	if d := getDetails(r); d != nil {
		d.StatusCode = code
	}

	// Use default http.Redirect function
	http.Redirect(w, r, url, code)
}

// NotFound replies to the request with an HTTP 404 not found error, using the problem response defined by RFC 7807.
func NotFound(w http.ResponseWriter, r *http.Request, extras ...ProblemExtra) {
	code := http.StatusNotFound
	title := http.StatusText(code)

	Problem(w, r, title, title, code, extras...)
}

// Error replies to the request with the specified error message and HTTP code,
// using the problem response defined by RFC 7807.
func Error(w http.ResponseWriter, r *http.Request, err string, code int, extras ...ProblemExtra) {
	Problem(w, r, http.StatusText(code), err, code, extras...)
}

//
// Problem
//

// problem is a wrapper around a map that is used to collect all fields required for a problem response in
// accordance to RFC 7807. It implements the json Marshaler interface to define a custom json marshalling.
type problem struct {
	fields map[string]interface{}
}

// MarshalJSON implements the json Marshaler interface.
func (p problem) MarshalJSON() ([]byte, error) {
	// Use the default marshaller for nested fields map.
	return json.Marshal(p.fields)
}

// ProblemExtra provides a way to add extra information into the problem response.
type ProblemExtra func(*problem)

// WithType allows the type field to be added to the problem response.
func WithType(t string) ProblemExtra {
	return func(p *problem) {
		p.fields["type"] = t
	}
}

// WithDetail allows the detail field to be added to the problem response.
func WithDetail(d string) ProblemExtra {
	return func(p *problem) {
		p.fields["detail"] = d
	}
}

// WithInstance allows the instance field to be added to the problem response.
func WithInstance(i string) ProblemExtra {
	return func(p *problem) {
		p.fields["instance"] = i
	}
}

// WithFields allows extra fields to be included in the problem response.
// Any field keys that clash with those expected in the problem response will not be used.
func WithFields(fields map[string]interface{}) ProblemExtra {
	return func(p *problem) {
		// merges extra fields into the existing fields,
		// where a conflict will choose the value from the existing fields.
		p.fields = mergemap.Merge(fields, p.fields)
	}
}

// Problem responds to HTTP request with a response that follows RFC 7807 (https://tools.ietf.org/html/rfc7807).
// It should be used for error responses.
func Problem(w http.ResponseWriter, r *http.Request, title, detail string, code int, extras ...ProblemExtra) {

	p := problem{
		fields: map[string]interface{}{
			"type":   "about:blank", // RFC 7807 defines this as the default type value.
			"title":  title,
			"detail": detail,
			"status": code,
		},
	}

	// Apply all extras to the response struct
	for _, e := range extras {
		e(&p)
	}

	// Set the correct content-type header, as defined by RFC 7807.
	w.Header().Set("Content-Type", "application/problem+json")

	Respond(w, r, code, p)
}
