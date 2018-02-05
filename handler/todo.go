package handler

import (
	"log"
	"net/http"

	"encoding/json"

	"time"

	"github.com/sharykhin/todoapp/repository/todo"
	"github.com/sharykhin/todoapp/request"
	"github.com/sharykhin/todoapp/service/response"
)

// Get list of todos
func Index(w http.ResponseWriter, r *http.Request, repository todo.Repository) {
	var limit string = "10"
	var offset string = "0"

	if rL := r.FormValue("limit"); rL != "" {
		limit = rL
	}
	if rO := r.FormValue("offset"); rO != "" {
		offset = rO
	}

	todos, err := repository.Get(limit, offset)

	if err != nil {
		log.Println(err)
		res, _ := response.NewJson(w, http.StatusInternalServerError, response.Response{false, nil, err})
		w.Write(res)
		return
	}
	res, _ := response.NewJson(w, http.StatusOK, response.Response{true, todos, nil})
	w.Write(res)
}

// Create a new todo
func Create(w http.ResponseWriter, r *http.Request, repository todo.Repository) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var rt request.Todo
	err := decoder.Decode(&rt)
	if err != nil {
		log.Println(err)
		res, _ := response.NewJson(w, http.StatusInternalServerError, response.Response{false, nil, err.Error()})
		w.Write(res)
		return
	}
	rt.Completed = false
	rt.Created = time.Now().UTC()
	t, err := repository.Create(rt)
	if err != nil {
		log.Println(err)
		res, _ := response.NewJson(w, http.StatusInternalServerError, response.Response{false, nil, err.Error()})
		w.Write(res)
		return
	}
	res, _ := response.NewJson(w, http.StatusCreated, response.Response{true, t, nil})
	w.Write(res)
}
