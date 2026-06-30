package mongodb

import (
	"context"

	"github.com/mainak555/gofoundry/dtos"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ExecuteCursor(ctx *context.Context, _cmd func() (*mongo.Cursor, error)) ([]bson.M, error) {
	return ExecuteTCursor[bson.M](ctx, _cmd)
}

func ExecuteTCursor[T any](ctx *context.Context, _cmd func() (*mongo.Cursor, error)) ([]T, error) {
	cur, err := _cmd()
	if err != nil {
		return nil, err
	}

	var result []T
	defer cur.Close(*ctx)
	for cur.Next(*ctx) {
		var item T
		if err = cur.Decode(&item); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}

func PipelineBottomQuery[T any](header *dtos.GetAllRequest[T], matchQuery bson.M, sortQuery bson.M, stages ...bson.M) []bson.M {
	skip, limit := header.GetSkipAndLimit()
	paginate := []bson.M{}

	if len(matchQuery) > 0 {
		paginate = append(paginate, bson.M{
			"$match": matchQuery,
		})
	}

	if len(sortQuery) > 0 {
		paginate = append(paginate, bson.M{
			"$sort": sortQuery,
		})
	}

	paginate = append(paginate,
		bson.M{"$skip": skip},
		bson.M{"$limit": limit},
	)

	if len(stages) > 0 {
		paginate = append(paginate, stages...)
	}
	return paginate
}

func PipelineFacet[T any](header *dtos.GetAllRequest[T], matchQuery bson.M, sortQuery bson.M, stages ...bson.M) []bson.M {
	pipeline := []bson.M{}
	if len(matchQuery) > 0 {
		pipeline = append(pipeline, bson.M{
			"$match": matchQuery,
		})
	}
	pipeline = append(pipeline, bson.M{
		"$facet": bson.M{
			"metadata": []bson.M{
				{
					"$group": bson.M{
						"_id": nil,
						"totalCount": bson.M{
							"$sum": 1,
						},
					},
				},
			},
			"data": PipelineBottomQuery(header, nil, sortQuery, stages...),
		},
	}, bson.M{
		"$project": bson.M{
			"totalCount": bson.M{"$arrayElemAt": []interface{}{"$metadata.totalCount", 0}},
			"data":       1,
		},
	})
	return pipeline
}
