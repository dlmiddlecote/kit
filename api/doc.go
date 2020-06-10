/*
Package api provides a minimal framework for APIs.

There is one main interface and one main function that are used to interact with this package.

API

This interface defines an API, and concrete implementations of an API should be registered with a server, which is returned by;

NewServer

This function takes an API, and returns a `http.Server`.

These two should be used in conjunction to provide a conformant experience across many HTTP APIs.

Example

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
		a := MyAPI{logger}
    	srv := api.NewServer(":8080", logger, a)
    	srv.ListenAndServe()
	}

*/
package api
