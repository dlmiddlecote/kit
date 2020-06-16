package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"

	"github.com/matryer/is"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type testAPI struct {
	logger *zap.SugaredLogger
}

func (a *testAPI) Endpoints() []Endpoint {
	return []Endpoint{
		{
			Method:  "GET",
			Path:    "/",
			Handler: a.handler(),
		},
		{
			Method:       "GET",
			Path:         "/no-logs",
			Handler:      a.handler(),
			SuppressLogs: true,
		},
		{
			Method:          "GET",
			Path:            "/no-metrics",
			Handler:         a.handler(),
			SuppressMetrics: true,
		},
	}
}

func (a *testAPI) handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.logger.Info("from handler")
		Respond(w, r, http.StatusOK, map[string]string{"status": "ok"})
	})
}

func TestServer(t *testing.T) {

	is := is.New(t)

	// Create logger, and captured logs.
	logger, logs := newTestLogger(zap.InfoLevel)

	// Start a test server that serves the prometheus metrics endpoint.
	p := httptest.NewServer(promhttp.Handler())
	defer p.Close()

	// scrape is a function that will return the exposed prometheus metrics endpoint as a string.
	scrape := func() string {
		resp, _ := http.Get(p.URL)
		buf, _ := ioutil.ReadAll(resp.Body)
		return string(buf)
	}

	// timeseriesMatches is a function that will return the count timeseries that matches the given labels
	timeseriesMatches := func(method, path, status string) []string {
		re := regexp.MustCompile(`http_request_duration_seconds_count{method="` + method + `",path="` + path + `",status="` + status + `"} ([0-9\.]+)`)
		return re.FindStringSubmatch(scrape())
	}

	// timeseriesValue is a function that will return the value of the count timeseries for the given labels.
	timeseriesValue := func(method, path, status string) float64 {
		matches := timeseriesMatches(method, path, status)
		f, _ := strconv.ParseFloat(matches[1], 64)
		return f
	}

	// timeseriesMissing is a function that will check whether the count timeseries for the given labels is missing.
	timeseriesMissing := func(method, path, status string) bool {
		return len(timeseriesMatches(method, path, status)) == 0
	}

	// Create test api
	a := &testAPI{logger}

	// Create server
	srv := NewServer(":0", logger, a)

	// Create test server from real server
	s := httptest.NewServer(srv.Handler)
	defer s.Close()

	// Check that all endpoints and statuses have been predeclared, if they are expected to be.
	for _, e := range a.Endpoints() {
		for _, status := range []string{"2XX", "3XX", "4XX", "5XX"} {
			if !e.SuppressMetrics {
				is.Equal(timeseriesValue(e.Method, e.Path, status), float64(0)) // timeseries created, as not suppressed, and is 0.
			} else {
				is.True(timeseriesMissing(e.Method, e.Path, status)) // timeseries is not created, as suppressed.
			}
		}
	}

	// Call all endpoints.
	for _, e := range a.Endpoints() {

		// construct url.
		u := fmt.Sprintf("%s%s", s.URL, e.Path)

		// make request to endpoint.
		resp, _ := http.Get(u)

		// Read response body.
		buf, _ := ioutil.ReadAll(resp.Body)
		body := string(buf)

		// Check response.
		is.Equal(resp.StatusCode, http.StatusOK) // response status code is as expected.
		is.Equal(body, `{"status":"ok"}`)        // response body is as expected.
	}

	// Check log line counts.
	is.Equal(logs.FilterMessage("from handler").Len(), 3) // log line from handlers are logged.
	is.Equal(logs.FilterMessage("request").Len(), 2)      // log line from desired requests are logged.

	// Check correct endpoints are logged.
	for _, ll := range logs.FilterMessage("request").All() {
		matches := false
		for _, p := range []string{"/", "/no-metrics"} {
			if ll.ContextMap()["path"].(string) == p {
				matches = true
				break
			}
		}
		is.True(matches) // request log line is from correct endpoint
	}

	// Check exposed metrics are as expected.
	for _, e := range a.Endpoints() {
		if !e.SuppressMetrics {
			is.Equal(timeseriesValue("GET", e.Path, "2XX"), float64(1)) // timeseries for request is now 1, as not suppressed.
		} else {
			is.True(timeseriesMissing("GET", e.Path, "2XX")) // timeseries still not created, as suppressed.
		}

		// check the remaining timeseries haven't changed
		for _, status := range []string{"3XX", "4XX", "5XX"} {
			if !e.SuppressMetrics {
				is.Equal(timeseriesValue(e.Method, e.Path, status), float64(0)) // timeseries is still 0.
			} else {
				is.True(timeseriesMissing(e.Method, e.Path, status)) // timeseries is not created, as suppressed.
			}
		}
	}
}
