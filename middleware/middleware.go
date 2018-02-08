package middleware

import (
	"log"
	"net/http"
)

// Logger logs request
func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request income path: %s\n", r.URL.Path)
		h.ServeHTTP(w, r)
	})
}
