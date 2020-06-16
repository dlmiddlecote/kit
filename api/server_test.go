package api

import (
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
		{"GET", "/", a.handler(), []Middleware{}},
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

	// timeseriesValue is a function that will return the value of the count timeseries for the given labels.
	timeseriesValue := func(method, path, status string) float64 {
		re := regexp.MustCompile(`http_request_duration_seconds_count{method="` + method + `",path="` + path + `",status="` + status + `"} ([0-9\.]+)`)
		matches := re.FindStringSubmatch(scrape())
		f, _ := strconv.ParseFloat(matches[1], 64)
		return f
	}

	// Create test api
	a := &testAPI{logger}

	// Create server
	srv := NewServer(":0", logger, a)

	// Create test server from real server
	s := httptest.NewServer(srv.Handler)
	defer s.Close()

	// Check that all endpoints and statuses have been predeclared.
	for _, e := range a.Endpoints() {
		for _, status := range []string{"2XX", "3XX", "4XX", "5XX"} {
			is.Equal(timeseriesValue(e.Method, e.Path, status), float64(0)) // timeseries created, and is 0
		}
	}

	// Call server
	resp, _ := http.Get(s.URL)

	// Get response body
	buf, _ := ioutil.ReadAll(resp.Body)
	body := string(buf)

	// Check response
	is.Equal(resp.StatusCode, http.StatusOK) // response status code is as expected.
	is.Equal(body, `{"status":"ok"}`)        // response body is as expected.

	// Check log line from the handler is logged.
	is.Equal(logs.FilterMessage("from handler").Len(), 1) // log line from handler is logged.

	// Check request log line is logged.
	is.Equal(logs.FilterMessage("request").Len(), 1) // log line from handler is logged.

	// Check metric is as expected.
	is.Equal(timeseriesValue("GET", "/", "2XX"), float64(1)) // timeseries for request is now 1.

	for _, e := range a.Endpoints() {
		for _, status := range []string{"2XX", "3XX", "4XX", "5XX"} {
			// Skip timeseries we've just checked.
			if e.Method == "GET" && e.Path == "/" && status == "2XX" {
				continue
			}
			is.Equal(timeseriesValue(e.Method, e.Path, status), float64(0)) // timeseries is still 0
		}
	}
}
