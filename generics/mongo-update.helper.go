package generics

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ReplaceMongoEntityById[TEntity any](svc ICommonService[TEntity], ctx *context.Context, id string, _mapFn func(entity *TEntity) error, opts ...*options.ReplaceOptions,
) (*mongo.UpdateResult, error) {
	entity, err := svc.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := _mapFn(entity); err != nil {
		return nil, err
	}

	Id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return svc.GetRepository().Collection().ReplaceOne(*ctx, bson.M{"_id": Id}, entity, opts...)
}

func PatchMongoEntityById[TEntity, TModel any](svc ICommonService[TEntity], ctx *context.Context, id string, obj *TModel, sets ...bson.M) (int64, error) {
	Id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return -1, err
	}

	var doc bson.D
	if data, err := bson.Marshal(*obj); err != nil {
		return -1, err
	} else if err := bson.Unmarshal(data, &doc); err != nil {
		return -1, err
	}

	filter := bson.M{"$and": []bson.M{{"_id": Id}, {"_deleted": false}}} //Allow Update Only undeleted marked records
	return svc.GetRepository().UpdateMany(ctx, filter, append([]bson.M{
		{"$set": doc},
	}, sets...))
}
