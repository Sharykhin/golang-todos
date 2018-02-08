package middleware

import "net/http"

type (
	// Middleware defines new type for http.Handler wrapper
	Middleware func(http.Handler) http.Handler
)

// Chain is a helper function that returns new http.Handler
// made of wrapping h onto middlewares
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}
