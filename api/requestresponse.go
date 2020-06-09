package api

import (
	"encoding/json"
	"net/http"
)

// Respond should be used to respond to a http request within a http handler.
// Respond encodes any data passed in as JSON.
// Respond also sets the status code of the response on the request details, so
// middlewares can access this value.
func Respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {

	var jsonData []byte
	var err error

	// If we have data to respond with, encode it into JSON, and set the correct
	// header. If we cannot encode, we'll return an Internal Server Error.
	if data != nil {
		// Set the correct header
		w.Header().Set("Content-Type", "application/json")

		// Marshal data into byte array
		jsonData, err = json.Marshal(data)
		if err != nil {
			// There was an error Marshalling, so return a server error
			jsonData = []byte(`{"msg": "Internal Server Error"}`)
			status = http.StatusInternalServerError
		}
	}

	// Set status code value on request details so other middlewares can access it
	if d := getDetails(r); d != nil {
		d.StatusCode = status
	}

	// Set the status code of the response. This should be the last header to be written.
	w.WriteHeader(status)

	// Write the JSON body. This must be done last, otherwise we flush the response too quickly.
	if len(jsonData) > 0 {
		w.Write(jsonData)
	}
}

// Decode should be used to convert the request's JSON body into the given v value.
func Decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	// FUTURE: Can handle different content types by looking at the request headers.
	return json.NewDecoder(r.Body).Decode(v)
}
