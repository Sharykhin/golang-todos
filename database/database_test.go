package database

import (
	//"testing"
	"github.com/stretchr/testify/mock"
	"context"
	//"database/sql"
	//"github.com/stretchr/testify/assert"
	"testing"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"github.com/Sharykhin/golang-todos/entity"
	"github.com/stretchr/testify/assert"
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

func TestCreate(t *testing.T) {
	var err error
	var mockS sqlmock.Sqlmock
	db, mockS, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctx := context.Background()
	cc, _ := context.WithCancel(ctx)
	rt := entity.CreateParams{
		Title: "test title",
		Description: "test description",
		Completed: false,
	}

	mockS.ExpectExec("INSERT INTO todos").
		WithArgs("test title", "test description", false).
		WillReturnResult(sqlmock.NewResult(1,1))


	nt, err := Create(cc, rt)
	if err != nil {
		t.Errorf("error was not expected while inserting a new row: %v", err)
	}
	if err := mockS.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
	assert.Equal(t, "test title", nt.Title)
	assert.Equal(t, "test description", nt.Description)

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
