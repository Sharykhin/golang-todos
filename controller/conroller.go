package controller

import (
	"context"
	"fmt"
	"sync"

	"github.com/Sharykhin/golang-todos/contract"
	db "github.com/Sharykhin/golang-todos/database"
	"github.com/Sharykhin/golang-todos/entity"
)

var (
	// TODOCtrl provides references to a private struct that handles all request around items
	TODOCtrl = todo{storage: db.Storage}
)

type (
	todo struct {
		storage contract.TodoProvider
	}

	listResult struct {
		err  error
		list []entity.Todo
	}

	countResult struct {
		err   error
		count int
	}
)

func (t todo) Index3(ctx context.Context, limit, offset int) ([]entity.Todo, int, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	chList := t.runList(ctx, limit, offset)
	chCount := t.runCount(ctx)

	var todosList []entity.Todo
	var count int

	for {
		if chList == nil && chCount == nil {
			return todosList, count, nil
		}
		select {
		case listResult, ok := <-chList:
			if !ok {
				chList = nil
				continue
			}
			if listResult.err != nil {
				cancel()
				return nil, 0, listResult.err
			}
			todosList = listResult.list
		case countResult, ok := <-chCount:
			if !ok {
				chCount = nil
				continue
			}
			if countResult.err != nil {
				cancel()
				return nil, 0, countResult.err
			}
			count = countResult.count
		}
	}

}

func (t todo) runList(ctx context.Context, limit, offset int) <-chan listResult {
	chListResult := make(chan listResult)
	go func() {
		defer close(chListResult)
		var lr listResult
		list, err := t.storage.Get(ctx, limit, offset)
		if err != nil {
			lr.err = err
		}
		lr.list = list
		select {
		case <-ctx.Done():
			return
		case chListResult <- lr:
		}
	}()
	return chListResult
}

func (t todo) runCount(ctx context.Context) <-chan countResult {
	chCountResult := make(chan countResult)
	go func() {
		defer close(chCountResult)
		var cr countResult
		count, err := t.storage.Count(ctx)
		if err != nil {
			cr.err = err
		}
		cr.count = count
		select {
		case <-ctx.Done():
			return
		case chCountResult <- cr:
		}
	}()
	return chCountResult
}

func (t todo) Index2(ctx context.Context, limit, offset int) ([]entity.Todo, int, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	chTodos, chErr := getList(ctx, limit, offset, t.storage)
	chTotal, chTotalErr := getTotal(ctx, t.storage)

	var todos []entity.Todo
	var total int

	for {
		if chTodos == nil && chTotal == nil {
			return todos, total, nil
		}
		select {
		case gotTodos, ok := <-chTodos:
			if !ok {
				chTodos = nil
				continue
			}
			todos = gotTodos
		case gotErr, ok := <-chErr:
			if ok {
				cancel()
				return nil, 0, gotErr
			}
		case gotTotal, ok := <-chTotal:
			if !ok {
				chTotal = nil
				continue
			}
			total = gotTotal
		case gotErr, ok := <-chTotalErr:
			if ok {
				cancel()
				return nil, 0, gotErr
			}
		}
	}
}

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

func getList(ctx context.Context, limit, offset int, storage contract.TodoProvider) (<-chan []entity.Todo, <-chan error) {
	chTodos := make(chan []entity.Todo)
	chErr := make(chan error)

	go func(ctx context.Context, chTodos chan<- []entity.Todo, chErr chan<- error, storage contract.TodoProvider) {
		defer close(chTodos)
		defer close(chErr)
		todos, err := storage.Get(ctx, limit, offset)
		if err != nil {
			chErr <- err
			return
		}
		chTodos <- todos
	}(ctx, chTodos, chErr, storage)

	return chTodos, chErr
}

func getTotal(ctx context.Context, storage contract.TodoProvider) (<-chan int, <-chan error) {
	chTotal := make(chan int)
	chErr := make(chan error)

	go func(ctx context.Context, chTotal chan<- int, chErr chan<- error, storage contract.TodoProvider) {
		defer close(chTotal)
		defer close(chErr)
		total, err := storage.Count(ctx)
		if err != nil {
			chErr <- err
			return
		}
		chTotal <- total

	}(ctx, chTotal, chErr, storage)

	return chTotal, chErr
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
