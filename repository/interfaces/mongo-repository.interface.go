package interfaces

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IBaseMongoRepository interface {
	Collection() *mongo.Collection
	UpdateMany(*context.Context, bson.M, []bson.M) (int64, error)
}

type IMongoRepository interface {
	IRepository
	IBaseMongoRepository
	UpdateOne(*context.Context, string, []bson.M) (int64, error)
	GetWithOptions(*context.Context, bson.M, *options.FindOptions) ([]any, error)
}

type TMongoRepository[T any] interface {
	TRepository[T]
	IBaseMongoRepository
	Base() IMongoRepository
	GetWithOptions(*context.Context, bson.M, *options.FindOptions) ([]T, error)
}
