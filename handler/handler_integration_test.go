package handler

import (
	"net/http"
	"testing"

	"fmt"

	"log"

	"os"

	"github.com/Sharykhin/golang-todos/database"
	"github.com/gavv/httpexpect"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/sqlite3"
	_ "github.com/mattes/migrate/source/file"
	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite
}

var m *migrate.Migrate

func (suite *HandlerTestSuite) SetupTest() {
	fmt.Println("run migrations here")
	driver, err := sqlite3.WithInstance(database.DB(), &sqlite3.Config{})
	if err != nil {
		log.Fatalf("could not get driveer: %v", err)
	}

	m, err = migrate.NewWithDatabaseInstance(
		"file://../migration",
		"sqlite3", driver,
	)

	if err != nil {
		log.Fatalf("could not get migrate instance: %v", err)
	}
	m.Up()
}

func (suite *HandlerTestSuite) TearDownTest() {
	m.Down()
	os.Remove(os.Getenv("DB_SOURCE"))
}

func (suite *HandlerTestSuite) TestCreate() {
	fmt.Println("Run integration test for creation")
	suite.T().Run("success creation", func(t *testing.T) {
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
		obj.Value("error").Equal("")
		obj.Value("meta").Equal(nil)
		obj.Value("data").Object().Value("title").Equal("test title")
	})
}

func TestCreate(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
