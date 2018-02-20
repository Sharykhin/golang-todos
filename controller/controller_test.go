package controller

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/Sharykhin/golang-todos/entity"
	"context"
	"github.com/pkg/errors"
)

type mockStorage struct {
	mock.Mock
}

func (m mockStorage) Create(ctx context.Context, rt entity.CreateParams) (*entity.Todo, error) {
	ret := m.Called(ctx, rt)
	t, err := ret.Get(0), ret.Get(1)
	if err != nil {
		return nil, err.(error)
	}
	return t.(*entity.Todo), nil
}

func TestCreate(t *testing.T) {
	t.Run("success creation", func(t *testing.T) {
		ctx := context.Background()
		rt := entity.CreateParams{
			Title: "test title",
			Description: "test desc",
			Completed: false,
		}

		var returnErr error = errors.New("something went wrong")

		m := new(mockStorage)
		m.On("Create", ctx, rt).Return(nil, returnErr).Once()

		todo, err := Create(ctx, rt, m)
		if err  == nil {
			t.Error("expected error but got nil")
		}
		m.AssertExpectations(t)

		assert.Nil(t, todo)
		assert.Equal(t, returnErr.Error(), err.Error())
	})

	t.Run("error creation", func(t *testing.T) {
		var oldCreate = create
		defer func(){
			create = oldCreate
		}()

		m := new(mockStorage)
		var errExpect = errors.New("something went wrong")
		ctx := context.Background()
		rt := entity.CreateParams{
			Title: "test title",
			Description: "test desc",
			Completed: false,
		}

		m.On("Create", ctx, rt).Return(nil, errExpect).Once()

		todo, err := Create(ctx, rt, m)
		if err  == nil {
			t.Error("expected error but got nil", err)
		}
		m.AssertExpectations(t)

		assert.Nil(t, todo)
		assert.Equal(t, err.Error(), errExpect.Error())
	})
}
