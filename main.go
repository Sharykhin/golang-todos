package main

import (
	"net/http"

	"log"

	"github.com/sharykhin/todoapp/handler"
	"github.com/sharykhin/todoapp/middleware"
	"github.com/sharykhin/todoapp/repository/todo/sql"
)

func main() {
	// TODO: routes can be moved to somewhere else I guess
	http.Handle("/", middleware.Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: Dependencies should be resolved somewhere else I guess
		db, err := sql.New()
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		sqlRepository := sql.TodoRepository{db}
		handler.Index(w, r, sqlRepository)
	})))
	http.ListenAndServe("localhost:8082", nil)
}
