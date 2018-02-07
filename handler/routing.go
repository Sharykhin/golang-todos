package handler

import (
	"net/http"

	"github.com/Sharykhin/golang-todos/middleware"
	"github.com/gorilla/mux"
)

// Handler holds project routing
func Handler() http.Handler {
	r := mux.NewRouter()
	// TODO: routes can be moved to somewhere else I guess
	r.Handle("/", middleware.Chain(http.HandlerFunc(index), middleware.Logger))

	// TODO: need somehow handle method request, I guess that could be reach through middleware
	r.Handle("/create", middleware.Chain(http.HandlerFunc(create), middleware.Logger))
	return r
}
