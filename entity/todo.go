package entity

import (
	"time"

	"github.com/Sharykhin/golang-todos/utils"
)

// Todo represents the basis model
type Todo struct {
	ID          int64          `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Completed   bool           `json:"completed"`
	Created     utils.JSONTime `json:"created"`
}

// CreateParams used for parsing income requests
type CreateParams struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	Created     time.Time `json:"created"`
}
