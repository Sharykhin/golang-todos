package main

import (
	"net/http"

	"log"

	"github.com/sharykhin/todoapp/handler"
	"github.com/sharykhin/todoapp/middleware"
	"github.com/sharykhin/todoapp/repository/todo/sql"
)

func main() {
	//TODO: is it okay that we can permamment connection to sqlite database?
	db, err := sql.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// TODO: routes can be moved to somewhere else I guess
	http.Handle("/", middleware.Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sqlRepository := sql.TodoRepository{db}
		handler.Index(w, r, sqlRepository)
	})))
	// TODO: this one might be taken from env variable.
	http.ListenAndServe("localhost:8082", nil)
}
