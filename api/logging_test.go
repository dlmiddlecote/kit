package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
	"go.uber.org/zap"
)

func TestLogMW(t *testing.T) {

	is := is.New(t)

	// Create a dummy request to pass to our handler.
	r, err := newTestRequest("GET", "/status", nil)
	is.NoErr(err) // http request created ok.

	// Create a response recorder, which satisfies http.ResponseWriter, to record the response.
	rr := httptest.NewRecorder()

	// Create logger, and captured logs.
	logger, logs := newTestLogger(zap.InfoLevel)

	// Create logging middleware.
	mw := LogMW(logger)

	// Create dummy handler.
	var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log some message.
		logger.Info("from handler")
		// Use Respond function so that response details are captured.
		Respond(w, r, http.StatusAccepted, map[string]string{"status": "ok"})
	})

	// Wrap handler in logging middleware.
	h = mw(h)

	// Invoke our handler.
	h.ServeHTTP(rr, r)

	// Check response is as expected.
	is.Equal(rr.Code, http.StatusAccepted)                        // response status code is 202.
	is.Equal(rr.Header().Get("Content-Type"), "application/json") // response content-type is application/json.

	type body struct {
		Status string `json:"status"`
	}

	expectedBody := body{"ok"}
	var actualBody body
	err = json.Unmarshal(rr.Body.Bytes(), &actualBody)
	is.NoErr(err)                      // actual response body is JSON.
	is.Equal(actualBody, expectedBody) // response body is as expected.

	// Check logs are as expected.
	is.Equal(logs.Len(), 2) // 2 log lines are logged.

	// Check log line from the handler is logged.
	is.Equal(logs.FilterMessage("from handler").Len(), 1) // log line from handler is logged.

	// Check middleware log line is as expected, and is logged last.
	mwll := logs.All()[logs.Len()-1]
	is.Equal(mwll.Message, "request")                         // log line message is 'request'.
	is.Equal(mwll.ContextMap()["method"].(string), "GET")     // log line method field is 'GET'.
	is.Equal(mwll.ContextMap()["path"].(string), "/status")   // log line path field is '/status'.
	is.Equal(mwll.ContextMap()["status"].(int64), int64(202)) // log line status field is '202'.
	is.True(mwll.ContextMap()["request_id"].(string) != "")   // log line request_id field isn't empty.
	is.True(mwll.ContextMap()["duration"].(string) != "")     // log line duration field isn't empty.
}
