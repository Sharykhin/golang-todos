package controller

import (
	"context"
	"testing"

	"github.com/Sharykhin/golang-todos/entity"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"sync"

	"time"

	"github.com/Sharykhin/golang-todos/utils"
	"github.com/stretchr/testify/require"
)

type mockStorage struct {
	mock.Mock
}

func (m mockStorage) Create(ctx context.Context, rt entity.CreateParams) (*entity.Todo, error) {
	args := m.Called(ctx, rt)
	t, err := args.Get(0), args.Error(1)
	if err != nil {
		return nil, err
	}
	return t.(*entity.Todo), nil
}

func (m mockStorage) Get(ctx context.Context, limit, offset int) ([]entity.Todo, error) {
	args := m.Called(ctx, limit, offset)
	t, err := args.Get(0), args.Error(1)
	if err != nil {
		return nil, err
	}
	return t.([]entity.Todo), nil
}

func (m mockStorage) Count(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	c, err := args.Get(0), args.Error(1)
	if err != nil {
		return 0, err
	}
	return c.(int), nil
}

func TestCreate(t *testing.T) {

	m := new(mockStorage)
	defer m.AssertExpectations(t)

	ctx := context.Background()
	created := utils.JSONTime(time.Now())
	expectedTodo := entity.Todo{
		Title:       "test title",
		Description: "test description",
		Completed:   false,
		Created:     created,
	}

	expectedError := errors.New("something went wrong")

	m.On("Create", ctx, entity.CreateParams{
		Title:       "test title",
		Description: "test description",
		Completed:   false,
	}).Return(&expectedTodo, nil).Once()

	m.On("Create", ctx, entity.CreateParams{
		Title:       "",
		Description: "",
		Completed:   false,
	}).Return(nil, expectedError).Once()

	tt := []struct {
		name          string
		incomeRequest entity.CreateParams
		expectedTodo  *entity.Todo
		expectedErr   error
	}{
		{
			name: "success creation",
			incomeRequest: entity.CreateParams{
				Title:       "test title",
				Description: "test description",
				Completed:   false,
			},
			expectedTodo: &expectedTodo,
			expectedErr:  nil,
		},
		{
			name: "bad creation",
			incomeRequest: entity.CreateParams{
				Title:       "",
				Description: "",
				Completed:   false,
			},
			expectedTodo: nil,
			expectedErr:  expectedError,
		},
	}

	to := &todo{
		storage: m,
	}

	var wg sync.WaitGroup

	for _, tc := range tt {
		wg.Add(1)
		t.Run(tc.name, func(t *testing.T) {
			defer wg.Done()
			actual, err := to.Create(ctx, tc.incomeRequest)
			require.Equal(t, tc.expectedTodo, actual)
			require.Equal(t, tc.expectedErr, err)

			if actual != nil {
				require.Equal(t, tc.expectedTodo.Title, "test title")
				require.Equal(t, tc.expectedTodo.Description, "test description")
				require.Equal(t, tc.expectedTodo.Completed, false)
				require.Equal(t, tc.expectedTodo.Created, created)
			}
		})
	}
	wg.Wait()

	//t.Run("success creation", func(t *testing.T) {
	//	ctx := context.Background()
	//	rt := entity.CreateParams{
	//		Title:       "test title",
	//		Description: "test desc",
	//		Completed:   false,
	//	}
	//
	//	var returnErr error = errors.New("something went wrong")
	//
	//	m := new(mockStorage)
	//	m.On("Create", ctx, rt).Return(nil, returnErr).Once()
	//
	//	to := &todo{
	//		storage: m,
	//	}
	//
	//	todo, err := to.Create(ctx, rt)
	//	if err == nil {
	//		t.Error("expected error but got nil")
	//	}
	//	m.AssertExpectations(t)
	//
	//	assert.Nil(t, todo)
	//	assert.Equal(t, returnErr.Error(), err.Error())
	//})
	//
	//t.Run("error creation", func(t *testing.T) {
	//
	//	var errExpect = errors.New("something went wrong")
	//	ctx := context.Background()
	//	rt := entity.CreateParams{
	//		Title:       "test title",
	//		Description: "test desc",
	//		Completed:   false,
	//	}
	//
	//	m := new(mockStorage)
	//	m.On("Create", ctx, rt).Return(nil, errExpect).Once()
	//
	//	to := &todo{
	//		storage: m,
	//	}
	//
	//	todo, err := to.Create(ctx, rt)
	//	if err == nil {
	//		t.Error("expected error but got nil", err)
	//	}
	//	m.AssertExpectations(t)
	//
	//	assert.Nil(t, todo)
	//	assert.Equal(t, err.Error(), errExpect.Error())
	//})
}

func TestTodo_Index3(t *testing.T) {

	m := new(mockStorage)
	defer m.AssertExpectations(t)

	to := &todo{
		storage: m,
	}

	ctx := context.Background()
	cc, cancel := context.WithCancel(ctx)
	defer cancel()

	expectedList := []entity.Todo{
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

	expectedError := errors.New("something went wring")

	m.On("Get", cc, 10, 0).Return(expectedList, nil).Once()
	m.On("Count", cc).Return(10, nil).Once()
	m.On("Get", cc, 0, 0).Return(nil, expectedError).Once()

	tt := []struct {
		name          string
		limit         int
		offset        int
		expectedList  []entity.Todo
		expectedCount int
		expectedErr   error
	}{
		{
			name:          "success index",
			limit:         10,
			offset:        0,
			expectedList:  expectedList,
			expectedCount: 10,
			expectedErr:   nil,
		},
		{
			name:          "bad list",
			limit:         0,
			offset:        0,
			expectedList:  nil,
			expectedCount: 0,
			expectedErr:   expectedError,
		},
	}

	var wg sync.WaitGroup
	for _, tc := range tt {
		wg.Add(1)
		go t.Run(tc.name, func(t *testing.T) {
			defer wg.Done()
			actual, count, err := to.Index3(ctx, tc.limit, tc.offset)

			require.Equal(t, tc.expectedList, actual)
			require.Equal(t, tc.expectedCount, count)
			require.Equal(t, tc.expectedErr, err)

			if actual != nil {
				require.Equal(t, len(actual), 2)
				require.Equal(t, actual[0].Title, expectedList[0].Title)
				require.Equal(t, actual[0].Description, expectedList[0].Description)
				require.Equal(t, actual[0].Completed, expectedList[0].Completed)
			}
		})
	}
	wg.Wait()

	//t.Run("success index", func(t *testing.T) {
	//
	//	ctx := context.Background()
	//	cc, cancel := context.WithCancel(ctx)
	//	defer cancel()
	//
	//	tt := []entity.Todo{
	//		{
	//			ID:          19,
	//			Title:       "test title",
	//			Description: "test description",
	//			Completed:   false,
	//		},
	//		{
	//			ID:          20,
	//			Title:       "test title",
	//			Description: "test description",
	//			Completed:   true,
	//		},
	//	}
	//
	//	m := new(mockStorage)
	//	m.On("Get", cc, 10, 0).Return(tt, nil).Once()
	//	m.On("Count", cc).Return(10, nil).Once()
	//
	//	to := &todo{
	//		storage: m,
	//	}
	//
	//	ts, c, err := to.Index(ctx, 10, 0)
	//	if err != nil {
	//		t.Errorf("unexpected error: %v", err)
	//	}
	//	m.AssertExpectations(t)
	//
	//	assert.Equal(t, 10, c)
	//	assert.Equal(t, 2, len(ts))
	//})

	//t.Run("error getting list", func(t *testing.T) {
	//	ctx := context.Background()
	//	cc, cancel := context.WithCancel(ctx)
	//	defer cancel()
	//
	//	tt := []entity.Todo{
	//		{
	//			ID:          19,
	//			Title:       "test title",
	//			Description: "test description",
	//			Completed:   false,
	//		},
	//		{
	//			ID:          20,
	//			Title:       "test title",
	//			Description: "test description",
	//			Completed:   true,
	//		},
	//	}
	//
	//	exErr := errors.New("something went wrong")
	//
	//	m := new(mockStorage)
	//	m.On("Get", cc, 10, 0).Return(tt, nil).Maybe()
	//	m.On("Count", cc).Return(0, exErr).Once()
	//
	//	to := &todo{
	//		storage: m,
	//	}
	//	ts, c, err := to.Index(ctx, 10, 0)
	//	m.AssertExpectations(t)
	//
	//	assert.Nil(t, ts)
	//	assert.Equal(t, 0, c)
	//	assert.NotNil(t, err)
	//	assert.Equal(t, "could not get count of todos: something went wrong", err.Error())
	//})
	//t.Run("error on count", func(t *testing.T) {
	//	ctx := context.Background()
	//	cc, cancel := context.WithCancel(ctx)
	//	defer cancel()
	//
	//	exErr := errors.New("something went wrong")
	//
	//	m := new(mockStorage)
	//	m.On("Get", cc, 10, 0).Return(nil, exErr).Once()
	//	m.On("Count", cc).Return(10, nil).Maybe()
	//
	//	to := &todo{
	//		storage: m,
	//	}
	//
	//	ts, c, err := to.Index(ctx, 10, 0)
	//	m.AssertExpectations(t)
	//
	//	assert.Nil(t, ts)
	//	assert.Equal(t, 0, c)
	//	assert.NotNil(t, err)
	//	assert.Equal(t, "could not get all todos: something went wrong", err.Error())
	//})
}
