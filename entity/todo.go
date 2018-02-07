package entity

import (
	"time"

	"github.com/sharykhin/todoapp/utils"
)

type Todo struct {
	Id          int            `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Completed   bool           `json:"completed"`
	Created     utils.JSONTime `json:"created"`
}

type CreateParams struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	Created     time.Time `json:"created"`
}
