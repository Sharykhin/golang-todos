package handler

import (
	"log"
	"net/http"

	"github.com/sharykhin/todoapp/repository/todo"
	"github.com/sharykhin/todoapp/service/response"
)

func Index(w http.ResponseWriter, r *http.Request, repository todo.Repository) {
	todos, err := repository.Get(10, 0)

	if err != nil {
		log.Println(err)
		res, _ := response.NewJson(w, http.StatusInternalServerError, response.Response{false, nil, err})
		w.Write(res)
		return
	}
	res, _ := response.NewJson(w, http.StatusOK, response.Response{false, todos, nil})
	w.Write(res)
}
