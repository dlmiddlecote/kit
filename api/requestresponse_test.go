package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestProblemResponse(t *testing.T) {

	is := is.New(t)

	// Create a dummy request to pass to our problem response.
	r, err := newTestRequest("GET", "/teapot", nil, "/:path")
	is.NoErr(err)

	// Create a response recorder, which satisfies http.ResponseWriter, to record the response.
	rr := httptest.NewRecorder()

	// Respond with a problem.
	Problem(rr, r, "There's a problem", "The problem is that I'm a teapot", http.StatusTeapot)

	// Check things.
	is.Equal(rr.Code, http.StatusTeapot) // status code is teapot.

	is.Equal(rr.Header().Get("Content-Type"), "application/problem+json") // content-type is correct.

	type body struct {
		Type   string `json:"type"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
		Status int    `json:"status"`
	}

	expectedBody := body{
		Type:   "about:blank",
		Title:  "There's a problem",
		Detail: "The problem is that I'm a teapot",
		Status: http.StatusTeapot,
	}
	var actualBody body
	err = json.Unmarshal(rr.Body.Bytes(), &actualBody)
	is.NoErr(err)                      // actual body is json.
	is.Equal(actualBody, expectedBody) // response body is correct.

	d := getDetails(r)
	is.True(d != nil) // details exist.
	if d != nil {
		is.Equal(d.StatusCode, http.StatusTeapot) // status is set on details.
	}
}

func TestProblemResponseWithExtras(t *testing.T) {

	is := is.New(t)

	// Create a dummy request to pass to our problem response.
	r, err := newTestRequest("GET", "/teapot", nil, "/:path")
	is.NoErr(err)

	// Create a response recorder, which satisfies http.ResponseWriter, to record the response.
	rr := httptest.NewRecorder()

	// Respond with a problem.
	Problem(rr, r, "There's a problem", "The problem is that I'm a teapot", http.StatusTeapot, WithType("https://example.net/validation-error"), WithInstance("BROWN-BETTY"))

	// Check things.
	is.Equal(rr.Code, http.StatusTeapot) // status code is teapot.

	is.Equal(rr.Header().Get("Content-Type"), "application/problem+json") // content-type is correct.

	type body struct {
		Type     string `json:"type"`
		Title    string `json:"title"`
		Status   int    `json:"status"`
		Detail   string `json:"detail"`
		Instance string `json:"instance"`
	}

	expectedBody := body{
		Type:     "https://example.net/validation-error",
		Title:    "There's a problem",
		Detail:   "The problem is that I'm a teapot",
		Status:   http.StatusTeapot,
		Instance: "BROWN-BETTY",
	}
	var actualBody body
	err = json.Unmarshal(rr.Body.Bytes(), &actualBody)
	is.NoErr(err)                      // actual body is json.
	is.Equal(actualBody, expectedBody) // response body is correct.

	d := getDetails(r)
	is.True(d != nil) // details exist.
	if d != nil {
		is.Equal(d.StatusCode, http.StatusTeapot) // status is set on details.
	}
}

func TestProblemResponseWithFields(t *testing.T) {

	is := is.New(t)

	// Create a dummy request to pass to our problem response.
	r, err := newRequest("GET", "/teapot", nil)
	is.NoErr(err)

	// Create a response recorder, which satisfied http.ResponseWriter, to record the response.
	rr := httptest.NewRecorder()

	// Create the extra fields for the problem
	extras := map[string]interface{}{
		"validation_errors": []string{
			"no handle.",
			"no spout.",
		},
		"detail": "I won't be included because I clash with the problem response.",
	}

	// Respond with a problem.
	Problem(rr, r, "There's a problem", "The problem is that I'm a teapot", http.StatusTeapot, WithFields(extras))

	// Check things.
	is.Equal(rr.Code, http.StatusTeapot) // status code is teapot.

	is.Equal(rr.Header().Get("Content-Type"), "application/problem+json") // content-type is correct.

	type body struct {
		Type             string   `json:"type"`
		Title            string   `json:"title"`
		Detail           string   `json:"detail"`
		Status           int      `json:"status"`
		ValidationErrors []string `json:"validation_errors"`
	}

	expectedBody := body{
		Type:   "about:blank",
		Title:  "There's a problem",
		Detail: "The problem is that I'm a teapot",
		Status: http.StatusTeapot,
		ValidationErrors: []string{
			"no handle.",
			"no spout.",
		},
	}
	var actualBody body
	err = json.Unmarshal(rr.Body.Bytes(), &actualBody)
	is.NoErr(err)                      // actual body is json.
	is.Equal(actualBody, expectedBody) // response body is correct.

	d := getDetails(r)
	is.True(d != nil) // details exist.
	if d != nil {
		is.Equal(d.StatusCode, http.StatusTeapot) // status is set on details.
	}
}

func TestNotFoundResponse(t *testing.T) {

	is := is.New(t)

	// Create a dummy request to pass to our not found response.
	r, err := newTestRequest("GET", "/not-found", nil, "/:path")
	is.NoErr(err)

	// Create a response recorder, which satisfies http.ResponseWriter, to record the response.
	rr := httptest.NewRecorder()

	// Respond with not found.
	NotFound(rr, r)

	// Check things.
	is.Equal(rr.Code, 404) // status code is correct.

	is.Equal(rr.Header().Get("Content-Type"), "application/problem+json") // content-type is correct.

	type body struct {
		Type   string `json:"type"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
		Status int    `json:"status"`
	}

	expectedBody := body{
		Type:   "about:blank",
		Title:  "Not Found",
		Detail: "Not Found",
		Status: 404,
	}
	var actualBody body
	err = json.Unmarshal(rr.Body.Bytes(), &actualBody)
	is.NoErr(err)                      // actual body is json.
	is.Equal(actualBody, expectedBody) // response body is correct.

	d := getDetails(r)
	is.True(d != nil) // details exist.
	if d != nil {
		is.Equal(d.StatusCode, 404) // status is set on details.
	}
}

func TestNotFoundWithDetailResponse(t *testing.T) {

	is := is.New(t)

	// Create a dummy request to pass to our not found response.
	r, err := newRequest("GET", "/not-found", nil)
	is.NoErr(err)

	// Create a response recorder, which satisfied http.ResponseWriter, to record the response.
	rr := httptest.NewRecorder()

	// Respond with not found.
	NotFound(rr, r, WithDetail("teapot '1' not found"))

	// Check things.
	is.Equal(rr.Code, 404) // status code is correct.

	is.Equal(rr.Header().Get("Content-Type"), "application/problem+json") // content-type is correct.

	type body struct {
		Type   string `json:"type"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
		Status int    `json:"status"`
	}

	expectedBody := body{
		Type:   "about:blank",
		Title:  "Not Found",
		Detail: "teapot '1' not found",
		Status: 404,
	}
	var actualBody body
	err = json.Unmarshal(rr.Body.Bytes(), &actualBody)
	is.NoErr(err)                      // actual body is json.
	is.Equal(actualBody, expectedBody) // response body is correct.

	d := getDetails(r)
	is.True(d != nil) // details exist.
	if d != nil {
		is.Equal(d.StatusCode, 404) // status is set on details.
	}
}

func TestErrorResponse(t *testing.T) {

	is := is.New(t)

	// Create a dummy request to pass to our error response.
	r, err := newTestRequest("GET", "/error", nil, "/:path")
	is.NoErr(err)

	// Create a response recorder, which satisfies http.ResponseWriter, to record the response.
	rr := httptest.NewRecorder()

	// Respond with an error.
	Error(rr, r, "Oops...", http.StatusInternalServerError)

	// Check things.
	is.Equal(rr.Code, http.StatusInternalServerError) // status code is correct.

	is.Equal(rr.Header().Get("Content-Type"), "application/problem+json") // content-type is correct.

	type body struct {
		Type   string `json:"type"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
		Status int    `json:"status"`
	}

	expectedBody := body{
		Type:   "about:blank",
		Title:  "Internal Server Error",
		Detail: "Oops...",
		Status: http.StatusInternalServerError,
	}
	var actualBody body
	err = json.Unmarshal(rr.Body.Bytes(), &actualBody)
	is.NoErr(err)                      // actual body is json.
	is.Equal(actualBody, expectedBody) // response body is correct.

	d := getDetails(r)
	is.True(d != nil) // details exist.
	if d != nil {
		is.Equal(d.StatusCode, http.StatusInternalServerError) // status is set on details.
	}
}

func TestDecode(t *testing.T) {

	tests := []struct {
		Name              string
		Body              string
		IsErr             bool
		ExpectedName      string
		ExpectedHasHandle bool
		ExpectedHasSpout  bool
	}{
		{
			Name:              "correct json is decoded correctly",
			Body:              `{"name": "Brown Betty", "hasHandle": true, "hasSpout": false}`,
			IsErr:             false,
			ExpectedName:      "Brown Betty",
			ExpectedHasHandle: true,
			ExpectedHasSpout:  false,
		},
		{
			Name:  "incorrect json is decoded, but empty",
			Body:  `{"nom": "Brown Betty"}`,
			IsErr: false,
		},
		{
			Name:  "not json isn't decoded",
			Body:  `yaml: yes`,
			IsErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			is := is.New(t)

			// Create a dummy request to pass to our decoder.
			r, err := newTestRequest("GET", "/teapot", strings.NewReader(tt.Body), "/:path")
			is.NoErr(err)

			// Create a response recorder, which satisfies http.ResponseWriter, to record the response.
			rr := httptest.NewRecorder()

			type body struct {
				Name      string `json:"name"`
				HasHandle bool   `json:"hasHandle"`
				HasSpout  bool   `json:"hasSpout"`
			}

			var teapot body
			err = Decode(rr, r, &teapot)
			is.Equal(err != nil, tt.IsErr) // decode exits as expected.

			if !tt.IsErr {
				// check teapot created correctly
				is.Equal(teapot.Name, tt.ExpectedName)           // teapot name is as expected.
				is.Equal(teapot.HasHandle, tt.ExpectedHasHandle) // teapot handle is as expected.
				is.Equal(teapot.HasSpout, tt.ExpectedHasSpout)   // teapot spout is as expected.
			}
		})
	}

}

func TestRespond(t *testing.T) {

	tests := []struct {
		Name                string
		Response            interface{}
		Code                int
		ExpectedCode        int
		ExpectedContentType string
		ExpectedBody        string
	}{
		{
			Name:                "Empty body",
			Response:            nil,
			Code:                http.StatusAccepted,
			ExpectedCode:        http.StatusAccepted,
			ExpectedContentType: "",
			ExpectedBody:        "",
		},
		{
			Name:                "JSON Body",
			Response:            map[string]int{"status": 200},
			Code:                http.StatusOK,
			ExpectedCode:        http.StatusOK,
			ExpectedContentType: "application/json",
			ExpectedBody:        `{"status":200}`,
		},
		{
			Name:                "Non-JSON Marshallable Body",
			Response:            func() {},
			Code:                http.StatusOK,
			ExpectedCode:        http.StatusInternalServerError,
			ExpectedContentType: "application/problem+json",
			ExpectedBody:        `{"type":"about:blank","title":"Internal Server Error","status":500}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			is := is.New(t)

			// Create a dummy request.
			r, err := newTestRequest("GET", "/teapot", nil, "/:path")
			is.NoErr(err)

			// Create a response recorder, which satisfies http.ResponseWriter, to record the response.
			rr := httptest.NewRecorder()

			// Invoke Respond.
			Respond(rr, r, tt.Code, tt.Response)

			// Check response is as expected.
			is.Equal(rr.Code, tt.ExpectedCode)                                // response code is as expected.
			is.Equal(rr.Header().Get("Content-Type"), tt.ExpectedContentType) // content-type is as expected.
			is.Equal(rr.Body.String(), tt.ExpectedBody)                       // response body is as expected
		})
	}
}

func TestRedirect(t *testing.T) {

	is := is.New(t)

	// Create a dummy request.
	r, err := newTestRequest("GET", "/teapot", nil, "/:path")
	is.NoErr(err)

	// Create a response recorder, which satisfies http.ResponseWriter, to record the response.
	rr := httptest.NewRecorder()

	// Invoke Redirect.
	Redirect(rr, r, "https://example.com", http.StatusPermanentRedirect)

	// Check that the status is set in request details.
	d := getDetails(r)
	is.True(d != nil) // details exists.
	if d != nil {
		is.Equal(d.StatusCode, http.StatusPermanentRedirect) // details contains correct status code.
	}

	// Check response is as expected.
	is.Equal(rr.Code, http.StatusPermanentRedirect)              // response status is correct.
	is.Equal(rr.Header().Get("Location"), "https://example.com") // redirect location header is set correctly.
}
