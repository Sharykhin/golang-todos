package storage

import (
	"database/sql"
	"fmt"
)

func NewSQL(driver, soruce string) (*sql.DB, error) {
	db, err := sql.Open(driver, soruce)
	if err != nil {
		return nil, fmt.Errorf("could not connect to sqlite database: %v", err)
	}

	return db, err
}
