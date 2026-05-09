package interfaces

import (
	"context"
)

type IBaseRepository interface {
	Get(*context.Context, map[string]any) ([]any, error)
	GetOne(*context.Context, string) (any, error)
	CreateOne(*context.Context, any) error
}

type IExtendedRepository interface {
	//Create
	//Update
}

type ISoftDeleteRepository interface {
	DeleteOne(*context.Context, string, func() map[string]any) (int64, error)
	Delete(*context.Context, map[string]any, func() map[string]any) (int64, error)
}

type IHardDeleteRepository interface {
	Erase(*context.Context, map[string]any) (int64, error)
	EraseOne(*context.Context, string) (int64, error)
}

type IRepository interface {
	IBaseRepository
	ISoftDeleteRepository
	IHardDeleteRepository
	GetByPage(*context.Context, map[string]any, int64, int64) ([]any, error)
}

type TBaseRepository[T any] interface {
	Get(*context.Context, map[string]any) ([]T, error)
	GetOne(*context.Context, string) (*T, error)
	CreateOne(*context.Context, *T) error
	UpdateOne(*context.Context, string, *T) (int64, error)
}

type TRepository[T any] interface {
	TBaseRepository[T]
	ISoftDeleteRepository
	IHardDeleteRepository
	GetByPage(*context.Context, map[string]any, int64, int64) ([]T, error)
}
