# Kit
[![Test](https://github.com/dlmiddlecote/kit/workflows/Test/badge.svg)](https://github.com/dlmiddlecote/kit/actions?query=workflow%3ATest)
[![codecov](https://codecov.io/gh/dlmiddlecote/kit/branch/master/graph/badge.svg)](https://codecov.io/gh/dlmiddlecote/kit)
[![Documentation](https://godoc.org/github.com/dlmiddlecote/kit?status.svg)](http://godoc.org/github.com/dlmiddlecote/kit)

A collection of building blocks for Go apps 🧱.

## Packages

### package api

This package provides building blocks for HTTP APIs. There is one main interface and one main function that are used
to interact with this package.

- `API`: This interface defines an API, and concrete implementations of an API should be registered with a server, which is returned by;
- `NewServer`: This function takes an API, and returns a `http.Server`.

These two should be used in conjunction to provide a conformant experience across many HTTP APIs.

#### Example

```go
...

type MyAPI struct {
    logger *zap.SugaredLogger
}

func (a *MyAPI) Endpoints() []api.Endpoint {
    return []api.Endpoint{
        {"GET", "/:id", a.handleGet(), []api.Middleware{}},
    }
}

func (a *MyAPI) handleGet() http.Handler {
    var h http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
        api.Respond(w, r, http.StatusOK, nil)
    }
    return h
}

func main() {
    ...
    a := MyAPI{logger}
    srv := api.NewServer(":8080", logger, a)
    srv.ListenAndServe()
    ...
}
```
