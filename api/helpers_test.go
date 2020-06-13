package api

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// newTestLogger returns a logger usable in tests, and also a struct that captures log lines
// logged via the returned logger. It is possible to change the returned loggers level with the
// available level argument.
func newTestLogger(level zapcore.LevelEnabler) (*zap.SugaredLogger, *observer.ObservedLogs) {
	core, recorded := observer.New(level)
	return zap.New(core).Sugar(), recorded
}

// newTestRequest creates a new *http.Request, that has the kit framework request details
// inside the requests context. This request should be used for all tests to mimic the
// real execution path.
func newTestRequest(method, path string, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	r = setDetails(r, "/:path", httprouter.Params{})
	return r, nil
}
