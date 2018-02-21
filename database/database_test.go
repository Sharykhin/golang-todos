package database

import (
	"context"
	"testing"

	"github.com/Sharykhin/golang-todos/entity"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestCreate(t *testing.T) {
	t.Run("success creation", func(t *testing.T) {
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
			Title:       "test title",
			Description: "test description",
			Completed:   false,
		}

		mockS.ExpectExec("INSERT INTO todos").
			WithArgs(rt.Title, rt.Description, rt.Completed).
			WillReturnResult(sqlmock.NewResult(1, 1))

		nt, err := Storage.Create(cc, rt)
		if err != nil {
			t.Errorf("error was not expected while inserting a new row: %v", err)
		}
		if err := mockS.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %v", err)
		}
		assert.Equal(t, "test title", nt.Title)
		assert.Equal(t, "test description", nt.Description)
	})

	t.Run("error inserting", func(t *testing.T) {
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
			Title:       "test title",
			Description: "test description",
			Completed:   false,
		}

		mockS.ExpectExec("INSERT INTO todos").
			WithArgs(rt.Title, rt.Description, rt.Completed).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("could not insert row")))

		nt, err := Storage.Create(cc, rt)
		if err := mockS.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %v", err)
		}
		assert.Nil(t, nt)
		assert.Equal(t, "could not get last insert id: could not insert row", err.Error())
	})
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
