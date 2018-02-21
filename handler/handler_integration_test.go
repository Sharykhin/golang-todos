package handler

import (
	"log"
	"net/http"
	"testing"

	"database/sql"

	"github.com/gavv/httpexpect"
	_ "github.com/mattn/go-sqlite3"
)

func TestCreate(t *testing.T) {
	savedTODOIndex := todoIndex
	defer func() {
		todoIndex = savedTODOIndex
	}()

	// TODO: for integration tests this causes an issue. Since sql works inside database package and Storage uses internal link to *sql.DB
	db, err := sql.Open("sqlite3", "./foo_test.db")
	if err != nil {
		log.Fatalf("could not connect to database: %s", err)
	}

	t.Run("success creation", func(t *testing.T) {
		todoIndex = tc.todoFunc
		e := httpexpect.WithConfig(httpexpect.Config{
			Client: &http.Client{
				Transport: httpexpect.NewBinder(Handler()),
				Jar:       httpexpect.NewJar(),
			},
			Reporter: httpexpect.NewAssertReporter(t),
		})

		todoRequest := map[string]interface{}{
			"title":       "test title",
			"description": "test description",
		}

		obj := e.Request(http.MethodPost, "/create").
			WithJSON(todoRequest).
			Expect().Status(http.StatusCreated).JSON().Object()

		obj.Value("success").Equal(true)
	})
}
