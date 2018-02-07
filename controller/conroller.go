package controller

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Sharykhin/golang-todos/entity"
	"github.com/Sharykhin/golang-todos/provider"
	"github.com/Sharykhin/golang-todos/repository/todo"
	"github.com/Sharykhin/golang-todos/repository/todo/sql"
)

var (
	repository todo.Repositier
)

func init() {
	// @QUESTION:
	// Should we register some global dependencies somewhere. Currently I mean databases connection?
	p := provider.Register()
	repository = sql.TodoRepository{DB: p.Storage}
}

// Index returns list of todos
func Index(ctx context.Context, limit, offset int) (*[]entity.Todo, int, error) {
	ctx, cancel := context.WithCancel(ctx)
	chTodos := make(chan *[]entity.Todo)
	defer close(chTodos)
	chCount := make(chan int)
	defer close(chCount)
	chErr := make(chan error)
	defer close(chErr)
	done := make(chan bool)
	defer close(done)

	var todos *[]entity.Todo
	var count int
	var wg sync.WaitGroup

	// @QUESTION:
	// I am concerned whether this code looks ok, since from the first view it make not so good impression?
	wg.Add(1)
	go getList(ctx, limit, offset, chTodos, chErr, &wg)

	wg.Add(1)
	go getCount(ctx, chCount, chErr, &wg)

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		done <- true
	}(&wg)

	for {
		select {
		case t := <-chTodos:
			todos = t
		case c := <-chCount:
			count = c
		case err := <-chErr:
			cancel()
			return nil, 0, err
		case <-done:
			return todos, count, nil
		}
	}
}

// Create creates new todo
func Create(ctx context.Context, rt *entity.CreateParams) (*entity.Todo, error) {
	rt.Created = time.Now().UTC()
	return repository.Create(ctx, rt)
}

func getList(ctx context.Context, limit, offset int, chTodos chan<- *[]entity.Todo, chErr chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	todos, err := repository.Get(ctx, limit, offset)
	if err != nil {
		chErr <- fmt.Errorf("could not get all todos: %s", err)
	}
	chTodos <- todos
}

func getCount(ctx context.Context, chCount chan<- int, chErr chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	count, err := repository.Count(ctx)
	if err != nil {
		chErr <- fmt.Errorf("could not get count of todos: %s", err)
	}
	chCount <- count
}