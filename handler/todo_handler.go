package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sharykhin/golang-todos/controller"
	"github.com/Sharykhin/golang-todos/entity"
	"github.com/Sharykhin/golang-todos/database"
)

var (
	todoIndex  = controller.Index
	todoCreate = controller.Create
)

func index(w http.ResponseWriter, r *http.Request) {
	limit, err := queryParamInt(r, "limit", 10)
	if err != nil {
		badRequest(w, fmt.Sprintf("could not parse limit param: %s", err))
		return
	}
	offset, err := queryParamInt(r, "offset", 0)
	if err != nil {
		badRequest(w, fmt.Sprintf("could not parse offset param: %s", err))
		return
	}
	todos, count, err := todoIndex(r.Context(), limit, offset, database.Storage{})
	if err != nil {
		serverError(w, err)
		return
	}
	success(w, todos, map[string]int{"total": count, "count": len(todos)})
}

func create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var rt entity.CreateParams
	err := decoder.Decode(&rt)
	if err != nil {
		serverError(w, err)
		return
	}

	// Pass the storage as parameter
	t, err := todoCreate(r.Context(), rt, database.Storage{})
	if err != nil {
		serverError(w, err)
		return
	}
	successCreated(w, t)
}
