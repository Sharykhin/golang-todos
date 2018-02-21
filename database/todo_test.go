package database

import (
	//"testing"
	"github.com/stretchr/testify/mock"
	"context"
	//"database/sql"
	//"github.com/stretchr/testify/assert"
)

type mockDb struct {
	mock.Mock
}

func (m mockDb) QueryRowContext(ctx context.Context, query string) *mockRow {
	ret := m.Called(ctx, query)
	return ret.Get(0).(*mockRow)
}

type mockRow struct {
	mock.Mock
}

func(m mockRow) Scan(dest ...interface{}) error {
	ret := m.Called(dest)
	err := ret.Get(0)
	if err != nil {
		return err.(error)
	}
	return nil

}

//func TestCount(t *testing.T) {
//	t.Run("success count", func(t *testing.T) {
//		var oldDb = db
//		defer func() {
//			db = oldDb
//		}()
//
//		var count int
//		query := "SELECT COUNT(id) AS `total` FROM todos"
//		ctx := context.Background()
//		m := new(mockDb)
//		mr := new(mockRow)
//		mr.On("Scan", &count).Return(nil).Once()
//		m.On("QueryRowContext", ctx, query).Return(mr).Once()
//		// TODO: how to mock?
//		db = m
//		i, err := Count(ctx)
//		if err != nil {
//			t.Errorf("unexpected error: %v", err)
//		}
//		m.AssertExpectations(t)
//		mr.AssertExpectations(t)
//		assert.Equal(t, 0, i)
//	})
//}
