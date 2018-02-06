package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sharykhin/todoapp/repository/todo"
	"github.com/sharykhin/todoapp/request"
)

type TodoHandler struct {
	Handler Handler
}

func (th TodoHandler) Index(w http.ResponseWriter, r *http.Request, repository todo.Repositier) {
	limit := th.Handler.queryParam(r, "limit", "10")
	offset := th.Handler.queryParam(r, "offset", "0")
	//TODO: so we have two functions that should be run in a separate goroutines
	todos, err := repository.Get(limit, offset)

	if err != nil {
		th.Handler.serverError(w, err)
		return
	}

	count, err := repository.Count()
	if err != nil {
		th.Handler.serverError(w, err)
		return
	}

	th.Handler.success(w, todos, map[string]interface{}{"total": count, "count": len(todos)})
}

func (th TodoHandler) Create(w http.ResponseWriter, r *http.Request, repository todo.Repositier) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	rt := request.Todo{
		Completed: false,
		Created:   time.Now().UTC(),
	}
	err := decoder.Decode(&rt)

	if err != nil {
		th.Handler.serverError(w, err)
		return
	}

	t, err := repository.Create(rt)
	if err != nil {
		th.Handler.serverError(w, err)
		return
	}
	th.Handler.success(w, t, map[string]interface{}{})
}
