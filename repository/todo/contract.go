package todo

import (
	"github.com/sharykhin/todoapp/entity"
	"github.com/sharykhin/todoapp/request"
)

type Repository interface {
	Get(limit, offset string) ([]entity.Todo, error)
	Create(rt request.Todo) (*entity.Todo, error)
}
