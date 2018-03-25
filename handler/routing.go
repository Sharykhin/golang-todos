package handler

import (
	"net/http"

	"github.com/Sharykhin/golang-todos/middleware"
	"github.com/gorilla/mux"
)

// Handler holds project routing
func Handler() http.Handler {
	r := mux.NewRouter()
	r.Handle("/", middleware.Chain(http.HandlerFunc(index), middleware.Logger)).Methods("GET")
	r.Handle("/", middleware.Chain(http.HandlerFunc(create), middleware.Logger)).Methods("POST")
	return r
}
