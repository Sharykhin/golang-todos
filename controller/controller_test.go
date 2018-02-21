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

func TestIndex(t *testing.T) {
	t.Run("success index", func(t *testing.T) {

		ctx := context.Background()
		cc, _ := context.WithCancel(ctx)

		tt := []entity.Todo{
			{
				ID: 19,
				Title: "test title",
				Description: "test description",
				Completed: false,
			},
			{
				ID: 20,
				Title: "test title",
				Description: "test description",
				Completed: true,
			},
		}

		mockS := new(mockStorage)
		mockS.On("Get", cc, 10, 0).Return(tt, nil).Once()
		mockS.On("Count", cc).Return(10, nil).Once()

		ts, c, err := Index(ctx, 10, 0, mockS)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		mockS.AssertExpectations(t)

		assert.Equal(t, 10, c)
		assert.Equal(t, 2, len(ts))
	})

	t.Run("error getting list", func (t *testing.T) {
		ctx := context.Background()
		cc, _ := context.WithCancel(ctx)

		tt := []entity.Todo{
			{
				ID: 19,
				Title: "test title",
				Description: "test description",
				Completed: false,
			},
			{
				ID: 20,
				Title: "test title",
				Description: "test description",
				Completed: true,
			},
		}

		exErr := errors.New("something went wrong")

		mockS := new(mockStorage)
		mockS.On("Get", cc, 10, 0).Return(tt, nil).Maybe()
		mockS.On("Count", cc).Return(0, exErr).Once()

		ts, c, err := Index(ctx, 10, 0, mockS)
		mockS.AssertExpectations(t)

		assert.Nil(t, ts)
		assert.Equal(t, 0, c)
		assert.NotNil(t, err)
		assert.Equal(t, "could not get count of todos: something went wrong", err.Error())
	})
	t.Run("error on count", func (t *testing.T) {
		ctx := context.Background()
		cc, _ := context.WithCancel(ctx)

		exErr := errors.New("something went wrong")

		mockS := new(mockStorage)
		mockS.On("Get", cc, 10, 0).Return(nil, exErr).Once()
		mockS.On("Count", cc).Return(10, nil).Maybe()

		ts, c, err := Index(ctx, 10, 0, mockS)
		mockS.AssertExpectations(t)

		assert.Nil(t, ts)
		assert.Equal(t, 0, c)
		assert.NotNil(t, err)
		assert.Equal(t, "could not get all todos: something went wrong", err.Error())
	})
}
