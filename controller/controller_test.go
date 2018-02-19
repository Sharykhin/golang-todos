package controller

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/Sharykhin/golang-todos/entity"
	"context"
	"time"
	"github.com/Sharykhin/golang-todos/utils"
	"errors"
)

type mockCreate struct {
	mock.Mock
}

func (m mockCreate) successCreate(ctx context.Context, rt entity.CreateParams) (*entity.Todo, error) {
	ret := m.Called(ctx, rt)
	return ret.Get(0).(*entity.Todo), nil
}

func (m mockCreate) errorCreate(ctx context.Context, rt entity.CreateParams) (*entity.Todo, error) {
	ret := m.Called(ctx, rt)
	return nil, ret.Get(1).(error)
}

func TestCreate(t *testing.T) {
	t.Run("success creation", func(t *testing.T) {
		var oldCreate = create
		defer func(){
			create = oldCreate
		}()

		m := new(mockCreate)

		ctx := context.Background()
		rt := entity.CreateParams{
			Title: "test title",
			Description: "test desc",
			Completed: false,
		}

		newTodo := &entity.Todo{
			ID: 18,
			Title: "test title",
			Description: "test desc",
			Completed: false,
			Created: utils.JSONTime(time.Now().UTC()),
		}

		m.On("successCreate", ctx, rt).Return(newTodo, nil).Once()
		create = m.successCreate

		todo, err := Create(ctx, rt)
		if err  != nil {
			t.Errorf("unexpected error: %v", err)
		}
		m.AssertExpectations(t)

		assert.Equal(t, rt.Title, todo.Title)
		assert.Equal(t, rt.Description, todo.Description)
		assert.Equal(t, rt.Completed, todo.Completed)
	})

	t.Run("error creation", func(t *testing.T) {
		var oldCreate = create
		defer func(){
			create = oldCreate
		}()

		m := new(mockCreate)
		var errExpect = errors.New("something went wrong")
		ctx := context.Background()
		rt := entity.CreateParams{
			Title: "test title",
			Description: "test desc",
			Completed: false,
		}

		m.On("errorCreate", ctx, rt).Return(nil, errExpect).Once()
		create = m.errorCreate

		todo, err := Create(ctx, rt)
		if err  == nil {
			t.Error("expected error but got nil", err)
		}
		m.AssertExpectations(t)

		assert.Nil(t, todo)
		assert.Equal(t, err.Error(), errExpect.Error())
	})
}
