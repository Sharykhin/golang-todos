package main

import (
	"net/http"

	"github.com/sharykhin/todoapp/handler"
	"github.com/sharykhin/todoapp/middleware"
	"github.com/sharykhin/todoapp/provider"
	"github.com/sharykhin/todoapp/repository/todo/sql"
)

// @QUESTION:
// Is a structure of project ok?[describe based on repository]?
// What the best place for storing interfaces types? [describe based on repository]
// Do we need in some sort of config package?

func main() {
	// @QUESTION:
	// Should we register some global dependencies somewhere. Currently I mean databases connection?
	p := provider.Register()
	defer p.Storage.Close()

	// TODO: routes can be moved to somewhere else I guess
	http.Handle("/", middleware.Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// @QUESTION:
		// Is is a good approach for resolving dependencies and putting them into a handler?

		//Resolve dependencies
		sqlRepository := sql.TodoRepository{DB: p.Storage}
		// Initialize appropriate handler
		th := handler.TodoHandler{Handler: handler.Handler{}}
		// @QUESTION:
		// In general, what the best approach of getting dependencies. From my experience all the dependencies should be
		// passed as parameters of methods, of some sort of construct. And we have to get interfaces for better testing
		// and maintain?
		th.Index(w, r, sqlRepository)
	})))

	// TODO: need somehow hanlde method request, I guess that could be reach through middleware
	http.Handle("/create", middleware.Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sqlRepository := sql.TodoRepository{DB: p.Storage}
		handler.Create(w, r, sqlRepository)
	})))
	// TODO: this one might be taken from env variable.
	http.ListenAndServe(":8082", nil)
}
