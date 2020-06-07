# Kit

A collection of building blocks for Go apps.

## Packages

### package api

This package provides building blocks for HTTP APIs. There is one main interface and one main function that are used
to interact with this package.

- `API`: This interface defines an API, and concrete implementations of an API should be registered with a server, which is returned by:
- `NewServer`: This function takes an API, and returns a `http.Server`.

These two should be used in conjunction to provide a conformant experience across many HTTP APIs.
