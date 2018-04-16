package contract

import (
	"context"

	"github.com/Sharykhin/golang-todos/entity"
)

type (
	// TodoCreator interface describes creation method
	TodoProvider interface {
		Create(ctx context.Context, rt entity.CreateParams) (*entity.Todo, error)
		Get(ctx context.Context, limit, offset int) ([]entity.Todo, error)
		Count(ctx context.Context) (int, error)
	}
)
