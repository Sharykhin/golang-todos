package provider

import (
	"database/sql"
	"log"

	"github.com/sharykhin/todoapp/config"
	"github.com/sharykhin/todoapp/service/storage"
)

type Provider struct {
	Initialized bool
	Storage     *sql.DB
}

func Register() *Provider {
	var provider Provider
	db, err := storage.NewSQL(config.DB_DRIVER, config.DB_SOURCE)
	if err != nil {
		log.Fatal(err)
	}

	provider.Storage = db
	// TODO: experimental field
	provider.Initialized = true
	return &provider
}
