package dtos

import "go.mongodb.org/mongo-driver/bson"

type GetAllRequest[T any] struct {
	Header     *T
	Sort       *QuerySort
	Pagination *Pagination
	Filter     *QueryFilter
	QueryParam any
}

func (g *GetAllRequest[T]) GetHeader() *T {
	return g.Header
}

func (g *GetAllRequest[T]) GetSkipAndLimit() (int64, int64) {
	return g.Pagination.GetSkipAndLimit()
}

func (g *GetAllRequest[T]) EmptyFilter() bool {
	return g.Filter == nil || len(g.Filter.Filter) < 1
}

func (g *GetAllRequest[T]) GetFilterQuery() map[string]any {
	if !g.EmptyFilter() {
		return g.Filter.Filter
	}
	return nil
}

func (g *GetAllRequest[T]) GetSortQuery() map[string]bool {
	result := make(map[string]bool)
	if g.Sort != nil && len(g.Sort.Sort) > 0 {
		for k, v := range g.Sort.Sort {
			if v != nil {
				result[k] = *v
			}
		}
	}
	return result
}

func (g *GetAllRequest[T]) EmptySort() bool {
	return len(g.GetSortQuery()) < 1
}

func (g *GetAllRequest[T]) GetSortBsonD(sort bson.D) bson.D {
	if sq := g.GetSortQuery(); len(sq) > 0 {
		sort = bson.D{}
		for k, v := range sq {
			if v {
				sort = append(sort, bson.E{k, 1})
			} else {
				sort = append(sort, bson.E{k, -1})
			}
		}
	}
	return sort
}

func (g *GetAllRequest[T]) GetSortBsonM(sort bson.M) bson.M {
	if sq := g.GetSortQuery(); len(sq) > 0 {
		sort = bson.M{}
		for k, v := range sq {
			if v {
				sort[k] = 1
			} else {
				sort[k] = -1
			}
		}
	}
	return sort
}

func (g *GetAllRequest[T]) SetQueryParam(value any) *GetAllRequest[T] {
	g.QueryParam = value
	return g
}

func (g *GetAllRequest[T]) GetQueryParam() any {
	return g.QueryParam
}
