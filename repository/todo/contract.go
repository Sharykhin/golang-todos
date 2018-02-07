package todo

import (
	"context"

	"github.com/sharykhin/todoapp/entity"
	"github.com/sharykhin/todoapp/request"
)

type Repositier interface {
	Get(ctx context.Context, limit, offset int) (*[]entity.Todo, error)
	Create(ctx context.Context, rt request.Todo) (*entity.Todo, error)
	Count(ctx context.Context) (int, error)
}
