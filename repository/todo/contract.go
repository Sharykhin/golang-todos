package todo

import (
	"github.com/sharykhin/todoapp/entity"
	"github.com/sharykhin/todoapp/request"
)

type Repositier interface {
	Get(limit, offset string) ([]entity.Todo, error)
	Create(rt request.Todo) (*entity.Todo, error)
}
