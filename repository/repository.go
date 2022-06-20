package repository

import (
	"context"
)

// Repository ...
type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	Find(ctx context.Context, filter interface{}) (list []*T, err error)
	Get(ctx context.Context, id string) (entity *T, err error)
	Update(ctx context.Context, entity T) (updated *T, err error)
	Delete(ctx context.Context, id string) (deleted bool, err error)
}
