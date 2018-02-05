package middleware

import (
	"log"
	"net/http"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request income path: %s\n", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
