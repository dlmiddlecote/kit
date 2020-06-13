package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/matryer/is"
)

func newRequest(method, path string, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	r = setDetails(r, path, httprouter.Params{})
	return r, nil
}

func TestProblemResponse(t *testing.T) {

	is := is.New(t)

	// Create a dummy request to pass to our problem response.
	r, err := newRequest("GET", "/teapot", nil)
	is.NoErr(err)

	// Create a response recorder, which satisfied http.ResponseWriter, to record the response.
	rr := httptest.NewRecorder()

	// Respond with a problem.
	Problem(rr, r, "The problem is that I'm a teapot", http.StatusTeapot)

	// Check things.
	is.Equal(rr.Code, http.StatusTeapot) // status code is teapot.

	is.Equal(rr.Header().Get("Content-Type"), "application/problem+json") // content-type is correct.

	type body struct {
		Type   string `json:"type"`
		Title  string `json:"title"`
		Status int    `json:"status"`
	}

	expectedBody := body{
		Type:   "about:blank",
		Title:  "The problem is that I'm a teapot",
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
	r, err := newRequest("GET", "/teapot", nil)
	is.NoErr(err)

	// Create a response recorder, which satisfied http.ResponseWriter, to record the response.
	rr := httptest.NewRecorder()

	// Respond with a problem.
	Problem(rr, r, "The problem is that I'm a teapot", http.StatusTeapot, WithType("https://example.net/validation-error"), WithDetail("I need a handle."), WithInstance("BROWN-BETTY"))

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
		Title:    "The problem is that I'm a teapot",
		Status:   http.StatusTeapot,
		Detail:   "I need a handle.",
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

func TestNotFoundResponse(t *testing.T) {

	is := is.New(t)

	// Create a dummy request to pass to our problem response.
	r, err := newRequest("GET", "/not-found", nil)
	is.NoErr(err)

	// Create a response recorder, which satisfied http.ResponseWriter, to record the response.
	rr := httptest.NewRecorder()

	// Respond with not found.
	NotFound(rr, r)

	// Check things.
	is.Equal(rr.Code, 404) // status code is correct.

	is.Equal(rr.Header().Get("Content-Type"), "application/problem+json") // content-type is correct.

	type body struct {
		Type   string `json:"type"`
		Title  string `json:"title"`
		Status int    `json:"status"`
	}

	expectedBody := body{
		Type:   "about:blank",
		Title:  "Not Found",
		Status: 404,
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

func TestErrorResponse(t *testing.T) {

	is := is.New(t)

	// Create a dummy request to pass to our problem response.
	r, err := newRequest("GET", "/error", nil)
	is.NoErr(err)

	// Create a response recorder, which satisfied http.ResponseWriter, to record the response.
	rr := httptest.NewRecorder()

	// Respond with an error.
	Error(rr, r, "Oops...", http.StatusInternalServerError)

	// Check things.
	is.Equal(rr.Code, http.StatusInternalServerError) // status code is correct.

	is.Equal(rr.Header().Get("Content-Type"), "application/problem+json") // content-type is correct.

	type body struct {
		Type   string `json:"type"`
		Title  string `json:"title"`
		Status int    `json:"status"`
	}

	expectedBody := body{
		Type:   "about:blank",
		Title:  "Oops...",
		Status: http.StatusInternalServerError,
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
