package request

import "time"

type Todo struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	Created     time.Time `json:"created"`
}
