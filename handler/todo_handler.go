package handler

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/sharykhin/todoapp/entity"
	"github.com/sharykhin/todoapp/repository/todo"
	"github.com/sharykhin/todoapp/request"
)

type TodoHandler struct {
	Handler Handler
}

func (th TodoHandler) Index(w http.ResponseWriter, r *http.Request, repository todo.Repositier) {
	limit := th.Handler.queryParam(r, "limit", "10")
	offset := th.Handler.queryParam(r, "offset", "0")

	chTodos := make(chan []entity.Todo)
	chCount := make(chan *int)
	chErr := make(chan error)
	done := make(chan bool)

	var todos []entity.Todo
	var count *int
	var wg sync.WaitGroup
	wg.Add(2)

	go func(chTodos chan<- []entity.Todo, chErr chan<- error, wg *sync.WaitGroup) {
		todos, err := repository.Get(limit, offset)
		if err != nil {
			chErr <- err
		}
		chTodos <- todos
		wg.Done()
	}(chTodos, chErr, &wg)

	go func(chCount chan<- *int, chErr chan<- error, wg *sync.WaitGroup) {
		count, err := repository.Count()
		if err != nil {
			chErr <- err
		}
		chCount <- count
		wg.Done()
	}(chCount, chErr, &wg)

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		done <- true
	}(&wg)

	for {
		completed := false
		select {
		case t := <-chTodos:
			todos = t
		case c := <-chCount:
			count = c
		case err := <-chErr:
			th.Handler.serverError(w, err)
			return
		case <-done:
			close(chTodos)
			close(chCount)
			close(chErr)
			completed = true
			break
		}
		if completed == true {
			break
		}
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
