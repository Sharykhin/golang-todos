package command

import (
	"log"

	"github.com/Sharykhin/golang-todos/database"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/sqlite3"
	_ "github.com/mattes/migrate/source/file"
)

var m *migrate.Migrate

// Migrate runs all migration from the scratch
func Migrate() {
	driver, err := sqlite3.WithInstance(database.DB(), &sqlite3.Config{})
	if err != nil {
		log.Fatalf("could not get driveer: %v", err)
	}

	// TODO: I don't like passes with relative journey
	m, err = migrate.NewWithDatabaseInstance(
		"file://migration",
		"sqlite3", driver,
	)

	if err != nil {
		log.Fatalf("could not get migrate instance: %v", err)
	}
	m.Down()
	m.Up()
}
