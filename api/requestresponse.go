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

	// If we have data to respond with, encode it into JSON, and set the correct
	// header. If we cannot encode, we'll return an Internal Server Error.
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		err := enc.Encode(data)
		if err != nil {
			enc.Encode(map[string]string{"msg": "Internal Server Error"})
			status = http.StatusInternalServerError
		}
	}

	// Set status code value on request details so other middlewares can access it
	if d := getDetails(r); d != nil {
		d.StatusCode = status
	}

	// Set the status code of the response. This should be the last header to be written.
	w.WriteHeader(status)
}
