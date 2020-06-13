package api

import (
	"encoding/json"
	"net/http"
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
			jsonData = []byte(`{"msg": "Internal Server Error"}`)
			code = http.StatusInternalServerError
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

type problemResponse struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

// ProblemExtra provides a way to add extra information into the problem response.
type ProblemExtra func(*problemResponse)

// WithType allows the type field to be added to the problem response.
func WithType(t string) ProblemExtra {
	return func(r *problemResponse) {
		r.Type = t
	}
}

// WithDetail allows the detail field to be added to the problem response.
func WithDetail(detail string) ProblemExtra {
	return func(r *problemResponse) {
		r.Detail = detail
	}
}

// WithInstance allows the instance field to be added to the problem response.
func WithInstance(instance string) ProblemExtra {
	return func(r *problemResponse) {
		r.Instance = instance
	}
}

// Problem responds to HTTP request with a response that follows RFC 7807 (https://tools.ietf.org/html/rfc7807).
// It should be used for error responses.
func Problem(w http.ResponseWriter, r *http.Request, title string, code int, extras ...ProblemExtra) {

	response := problemResponse{
		Type:   "about:blank", // RFC 7807 defines this as the default type value.
		Title:  title,
		Status: code,
	}

	// Apply all extras to the response struct
	for _, e := range extras {
		e(&response)
	}

	// Set the correct content-type header, as defined by RFC 7807.
	w.Header().Set("Content-Type", "application/problem+json")

	Respond(w, r, code, response)
}

// NotFound replies to the request with an HTTP 404 not found error, using the problem response defined by RFC 7807.
func NotFound(w http.ResponseWriter, r *http.Request, extras ...ProblemExtra) {
	code := http.StatusNotFound
	title := http.StatusText(code)

	Problem(w, r, title, code, extras...)
}

// Error replies to the request with the specified error message and HTTP code,
// using the problem response defined by RFC 7807.
func Error(w http.ResponseWriter, r *http.Request, err string, code int, extras ...ProblemExtra) {
	Problem(w, r, err, code, extras...)
}
