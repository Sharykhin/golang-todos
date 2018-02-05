package middleware

import (
	"log"
	"net/http"
)

//TODO: this middleware is totally independend, so is a good candidate to move to a separate github repository
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request income path: %s\n", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
