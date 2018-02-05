package entity

import "github.com/sharykhin/todoapp/utils"

type Todo struct {
	Id          int            `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Completed   bool           `json:"completed"`
	Created     utils.JSONTime `json:"created"`
}
