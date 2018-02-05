package main

import (
	"net/http"

	"github.com/sharykhin/todoapp/handler"
	"github.com/sharykhin/todoapp/middleware"
	"github.com/sharykhin/todoapp/provider"
	"github.com/sharykhin/todoapp/repository/todo/sql"
)

func main() {

	p := provider.Register()
	defer p.Storage.Close()

	// TODO: routes can be moved to somewhere else I guess
	http.Handle("/", middleware.Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sqlRepository := sql.TodoRepository{DB: p.Storage}
		handler.Index(w, r, sqlRepository)
	})))

	// TODO: need somehow hanlde method request, I guess that could be reach through middleware
	http.Handle("/create", middleware.Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sqlRepository := sql.TodoRepository{DB: p.Storage}
		handler.Create(w, r, sqlRepository)
	})))
	// TODO: this one might be taken from env variable.
	http.ListenAndServe("localhost:8082", nil)
}
