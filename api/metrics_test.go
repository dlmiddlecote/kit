package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"

	"github.com/matryer/is"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestMetricsMW(t *testing.T) {

	is := is.New(t)

	// Start a test server that serves the prometheus metrics endpoint.
	s := httptest.NewServer(promhttp.Handler())
	defer s.Close()

	// scrape is a function that will return the exposed prometheus metrics endpoint as a string.
	scrape := func() string {
		resp, _ := http.Get(s.URL)
		buf, _ := ioutil.ReadAll(resp.Body)
		return string(buf)
	}

	// timeseriesValue is a function that will return the value of the count timeseries for the given labels.
	timeseriesValue := func(method, path, status string) float64 {
		re := regexp.MustCompile(`http_request_duration_seconds_count{method="` + method + `",path="` + path + `",status="` + status + `"} ([0-9\.]+)`)
		matches := re.FindStringSubmatch(scrape())
		f, _ := strconv.ParseFloat(matches[1], 64)
		return f
	}

	// Create a dummy handler.
	var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use Respond function so that response details are captured.
		Respond(w, r, http.StatusOK, map[string]string{"status": "ripe"})
	})

	// Create endpoints that'll be registered in the metrics middleware.
	endpoints := []Endpoint{
		{
			Method:      "GET",
			Path:        "/fruits/:fruit",
			Handler:     h,
			Middlewares: []Middleware{},
		},
		{
			Method:      "GET",
			Path:        "/vegetables/:vegetable",
			Handler:     h,
			Middlewares: []Middleware{},
		},
	}

	// Create the metrics middleware.
	mw := MetricsMW(endpoints)

	// Wrap handler in metrics middleware.
	h = mw(h)

	// Check that all endpoints and statuses have been predeclared.
	for _, e := range endpoints {
		for _, status := range []string{"2XX", "3XX", "4XX", "5XX"} {
			is.Equal(timeseriesValue(e.Method, e.Path, status), float64(0)) // timeseries created, and is 0
		}
	}

	// Create a dummy request to pass to our handler.
	r, err := newTestRequest("GET", "/fruits/grapes", nil, "/fruits/:fruit")
	is.NoErr(err) // http request created ok.

	// Create a response recorder, which satisfies http.ResponseWriter, to record the response.
	rr := httptest.NewRecorder()

	// Invoke our handler.
	h.ServeHTTP(rr, r)

	// Check response is as expected.
	is.Equal(rr.Code, http.StatusOK)                              // response status code is 200.
	is.Equal(rr.Header().Get("Content-Type"), "application/json") // response content-type is application/json.

	type body struct {
		Status string `json:"status"`
	}

	expectedBody := body{"ripe"}
	var actualBody body
	err = json.Unmarshal(rr.Body.Bytes(), &actualBody)
	is.NoErr(err)                      // actual response body is JSON.
	is.Equal(actualBody, expectedBody) // response body is as expected.

	// Check metric is as expected.
	is.Equal(timeseriesValue("GET", "/fruits/:fruit", "2XX"), float64(1)) // timeseries for request is now 1.

	for _, e := range endpoints {
		for _, status := range []string{"2XX", "3XX", "4XX", "5XX"} {
			// Skip timeseries we've just checked.
			if e.Method == "GET" && e.Path == "/fruits/:fruit" && status == "2XX" {
				continue
			}
			is.Equal(timeseriesValue(e.Method, e.Path, status), float64(0)) // timeseries is still 0
		}
	}
}
