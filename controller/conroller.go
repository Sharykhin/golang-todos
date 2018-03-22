package controller

import (
	"context"
	"fmt"
	"sync"

	db "github.com/Sharykhin/golang-todos/database"
	"github.com/Sharykhin/golang-todos/entity"
)

var (
	TODO = todo{storage: db.Storage}
)

type (
	// TodoCreator interface describes creation method
	TodoProvider interface {
		Create(ctx context.Context, rt entity.CreateParams) (*entity.Todo, error)
		Get(ctx context.Context, limit, offset int) ([]entity.Todo, error)
		Count(ctx context.Context) (int, error)
	}

	todo struct {
		storage TodoProvider
	}
)

//func init() {
//	TODO.storage = db.Storage
//}

// Index returns list of todos
func (t todo) Index(ctx context.Context, limit, offset int) ([]entity.Todo, int, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// So closing in defer way may create not obvious error of sending value to a closed channel. If error occurs
	// somewhere select statement get it and make a return hence all defer's are called and while gorouting a still running
	// they may try to send a value to a closed channel and panic will be thrown
	chTodos := make(chan []entity.Todo)
	//defer close(chTodos)
	chCount := make(chan int)
	//defer close(chCount)
	chErr := make(chan error)
	// Since we have a few suppliers of error we can't close channel in a specific goroutine, that's why we can close it
	// in defer by rely on context.Canceled property of ctx.Err() method
	defer close(chErr)
	done := make(chan struct{})
	//defer close(done)

	var todos []entity.Todo
	var count int
	var wg sync.WaitGroup

	wg.Add(1)
	// TODO: here we can apply patter when method create a channel run gorouting and return read-only channel
	go t.getList(ctx, limit, offset, chTodos, chErr, &wg)

	wg.Add(1)
	// TODO: here we can apply patter when method create a channel run gorouting and return read-only channel
	go t.getCount(ctx, chCount, chErr, &wg)
	// TODO: don't you think that it's better to check nil channel in for loop?
	go t.wait(&wg, done)

	for {
		select {
		case t, ok := <-chTodos:
			if !ok {
				chTodos = nil
				continue
			}
			todos = t
		case c, ok := <-chCount:
			if !ok {
				chCount = nil
				continue
			}
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
func (t todo) Create(ctx context.Context, rt entity.CreateParams) (*entity.Todo, error) {
	// TODO: narrow case, how to provide the exact utc time
	//rt.Created = time.Now().UTC()
	return t.storage.Create(ctx, rt)
}

func (t *todo) getList(ctx context.Context, limit, offset int, chTodos chan<- []entity.Todo, chErr chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(chTodos)
	todos, err := t.storage.Get(ctx, limit, offset)
	if err != nil {
		if ctx.Err() == context.Canceled {
			return
		}
		chErr <- fmt.Errorf("could not get all todos: %s", err)
	} else {
		chTodos <- todos
	}
}

func (t *todo) getCount(ctx context.Context, chCount chan<- int, chErr chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(chCount)
	count, err := t.storage.Count(ctx)

	if err != nil {
		if ctx.Err() == context.Canceled {
			return
		}
		chErr <- fmt.Errorf("could not get count of todos: %s", err)
	} else {
		chCount <- count
	}
}

func (t *todo) wait(wg *sync.WaitGroup, done chan<- struct{}) {
	wg.Wait()
	done <- struct{}{}
	close(done)
}
