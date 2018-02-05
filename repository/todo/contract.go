package todo

import "github.com/sharykhin/todoapp/entity"

type Repository interface {
	Get(limit int, offset int) ([]entity.Todo, error)
}
