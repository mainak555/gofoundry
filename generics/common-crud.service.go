package generics

import (
	"context"

	libmongo "gofoundry/db/mongodb"
	"gofoundry/dtos"
	"gofoundry/helpers"
	"gofoundry/repository"
	"gofoundry/repository/interfaces"

	"go.mongodb.org/mongo-driver/bson"
)

type ICommonService[TEntity any] interface {
	GetRepository() interfaces.TMongoRepository[TEntity]
	GetById(ctx *context.Context, objectId string) (*TEntity, error)
	HardDelete(ctx *context.Context, objectIds ...string) (int64, error)
	Get(ctx *context.Context, filter map[string]interface{}, pagination *dtos.Pagination) ([]TEntity, error)
	SoftDelete(ctx *context.Context, deletedBy string, objectIds ...string) (int64, error)
}

type CommonService[TEntity any] struct {
	_repository interfaces.TMongoRepository[TEntity]
}

func NewCommonService[TEntity any](mongoClient libmongo.IMongoClient) ICommonService[TEntity] {
	return &CommonService[TEntity]{
		_repository: repository.NewMongoTRepository[TEntity](mongoClient),
	}
}

func (svc *CommonService[TEntity]) GetRepository() interfaces.TMongoRepository[TEntity] {
	return svc._repository
}

func (svc *CommonService[TEntity]) GetById(ctx *context.Context, objectId string) (*TEntity, error) {
	res, err := svc._repository.GetOne(ctx, objectId)
	return res, err
}

func (svc *CommonService[TEntity]) Get(ctx *context.Context, filter map[string]interface{}, pagination *dtos.Pagination) ([]TEntity, error) {
	query := bson.M{}
	if filter != nil {
		query = filter
	}
	res, err := svc._repository.GetByPage(ctx, query, pagination.PageNo, pagination.PageSize)
	return res, err
}

// HardDelete implements ICommonService.
func (svc *CommonService[TEntity]) HardDelete(ctx *context.Context, objectIds ...string) (int64, error) {
	res, err := svc._repository.Erase(ctx, helpers.MongoQueryInByIds(objectIds...))
	return res, err
}

// SoftDelete implements ICommonService.
func (svc *CommonService[TEntity]) SoftDelete(ctx *context.Context, deletedBy string, objectIds ...string) (int64, error) {
	res, err := svc._repository.Delete(ctx, helpers.MongoQueryInByIds(objectIds...), helpers.MongoQueryDeletedBy(deletedBy))
	return res, err
}
