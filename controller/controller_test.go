package controller

import (
	"context"
	"testing"

	"github.com/Sharykhin/golang-todos/entity"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func (m mockStorage) Get(ctx context.Context, limit, offset int) ([]entity.Todo, error) {
	ret := m.Called(ctx, limit, offset)
	t, err := ret.Get(0), ret.Get(1)
	if err != nil {
		return nil, err.(error)
	}
	return t.([]entity.Todo), nil
}

func (m mockStorage) Count(ctx context.Context) (int, error) {
	ret := m.Called(ctx)
	c, err := ret.Get(0), ret.Get(1)
	if err != nil {
		return 0, err.(error)
	}
	return c.(int), nil
}

func TestCreate(t *testing.T) {
	t.Run("success creation", func(t *testing.T) {
		ctx := context.Background()
		rt := entity.CreateParams{
			Title:       "test title",
			Description: "test desc",
			Completed:   false,
		}

		var returnErr error = errors.New("something went wrong")

		m := new(mockStorage)
		m.On("Create", ctx, rt).Return(nil, returnErr).Once()

		to := &todo{
			storage: m,
		}

		todo, err := to.Create(ctx, rt)
		if err == nil {
			t.Error("expected error but got nil")
		}
		m.AssertExpectations(t)

		assert.Nil(t, todo)
		assert.Equal(t, returnErr.Error(), err.Error())
	})

	t.Run("error creation", func(t *testing.T) {

		var errExpect = errors.New("something went wrong")
		ctx := context.Background()
		rt := entity.CreateParams{
			Title:       "test title",
			Description: "test desc",
			Completed:   false,
		}

		m := new(mockStorage)
		m.On("Create", ctx, rt).Return(nil, errExpect).Once()

		to := &todo{
			storage: m,
		}

		todo, err := to.Create(ctx, rt)
		if err == nil {
			t.Error("expected error but got nil", err)
		}
		m.AssertExpectations(t)

		assert.Nil(t, todo)
		assert.Equal(t, err.Error(), errExpect.Error())
	})
}

func TestIndex(t *testing.T) {
	t.Run("success index", func(t *testing.T) {

		ctx := context.Background()
		cc, cancel := context.WithCancel(ctx)
		defer cancel()

		tt := []entity.Todo{
			{
				ID:          19,
				Title:       "test title",
				Description: "test description",
				Completed:   false,
			},
			{
				ID:          20,
				Title:       "test title",
				Description: "test description",
				Completed:   true,
			},
		}

		m := new(mockStorage)
		m.On("Get", cc, 10, 0).Return(tt, nil).Once()
		m.On("Count", cc).Return(10, nil).Once()

		to := &todo{
			storage: m,
		}

		ts, c, err := to.Index(ctx, 10, 0)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		m.AssertExpectations(t)

		assert.Equal(t, 10, c)
		assert.Equal(t, 2, len(ts))
	})

	t.Run("error getting list", func(t *testing.T) {
		ctx := context.Background()
		cc, cancel := context.WithCancel(ctx)
		defer cancel()

		tt := []entity.Todo{
			{
				ID:          19,
				Title:       "test title",
				Description: "test description",
				Completed:   false,
			},
			{
				ID:          20,
				Title:       "test title",
				Description: "test description",
				Completed:   true,
			},
		}

		exErr := errors.New("something went wrong")

		m := new(mockStorage)
		m.On("Get", cc, 10, 0).Return(tt, nil).Maybe()
		m.On("Count", cc).Return(0, exErr).Once()

		to := &todo{
			storage: m,
		}
		ts, c, err := to.Index(ctx, 10, 0)
		m.AssertExpectations(t)

		assert.Nil(t, ts)
		assert.Equal(t, 0, c)
		assert.NotNil(t, err)
		assert.Equal(t, "could not get count of todos: something went wrong", err.Error())
	})
	t.Run("error on count", func(t *testing.T) {
		ctx := context.Background()
		cc, cancel := context.WithCancel(ctx)
		defer cancel()

		exErr := errors.New("something went wrong")

		m := new(mockStorage)
		m.On("Get", cc, 10, 0).Return(nil, exErr).Once()
		m.On("Count", cc).Return(10, nil).Maybe()

		to := &todo{
			storage: m,
		}

		ts, c, err := to.Index(ctx, 10, 0)
		m.AssertExpectations(t)

		assert.Nil(t, ts)
		assert.Equal(t, 0, c)
		assert.NotNil(t, err)
		assert.Equal(t, "could not get all todos: something went wrong", err.Error())
	})
}
