package dependency

import (
	"context"

	"gofoundry/util"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type ResolveResult[TId comparable] map[TId]Result[TId]
type ResolveFn[TId comparable] func(ctx *context.Context, items []TId) (ResolveResult[TId], error)

type Result[TId comparable] struct {
	EntityName string                 `json:"entityName"`
	Dependency map[TId][]Result[TId]  `json:"dependency"`
	Parameters map[TId]map[string]any `json:"-"`
}

type EntityRelation[TId comparable] struct {
	resolver   map[string][]ResolveFn[TId]
	dependency map[string][]string
}

func NewEntityRelation[TId comparable]() *EntityRelation[TId] {
	return &EntityRelation[TId]{
		resolver:   make(map[string][]ResolveFn[TId]),
		dependency: make(map[string][]string),
	}
}

func (er *EntityRelation[TId]) register(parentEntityName, entityName string, resolverFn ResolveFn[TId]) {
	if _, ok := er.resolver[parentEntityName]; !ok {
		er.resolver[parentEntityName] = make([]ResolveFn[TId], 0)
		er.dependency[parentEntityName] = make([]string, 0)
	}

	if !slices.Contains(er.dependency[parentEntityName], entityName) {
		er.dependency[parentEntityName] = append(er.dependency[parentEntityName], entityName)
		er.resolver[parentEntityName] = append(er.resolver[parentEntityName], resolverFn)
	}
}

func Register[TParent, TEntity any, TId comparable](dm *EntityRelation[TId], resolverFn ResolveFn[TId]) {
	parentEntityName, err := util.GetEntityName[TParent]()
	if err != nil {
		panic(err)
	}

	entityName, err := util.GetEntityName[TEntity]()
	if err != nil {
		panic(err)
	}

	dm.register(parentEntityName, entityName, resolverFn)
}

func Resolve[T any, TId comparable](dm *EntityRelation[TId], ctx *context.Context, cascade bool, ids ...TId) (*Result[TId], error) {
	entityName, err := util.GetEntityName[T]()
	if err != nil {
		return nil, err
	}
	res, err := dm.dependencyLoop(dm.resolver[entityName], ctx, cascade, ids)
	if err != nil {
		return nil, err
	}
	return &Result[TId]{
		EntityName: entityName,
		Dependency: res,
	}, nil
}

func (er *EntityRelation[TId]) dependencyLoop(resolvers []ResolveFn[TId], ctx *context.Context, cascade bool, ids []TId) (map[TId][]Result[TId], error) {
	if len(resolvers) < 1 || len(ids) < 1 {
		return nil, nil
	}

	result := make(map[TId][]Result[TId])
	for _, resolve := range resolvers {
		if res, err := resolve(ctx, ids); err != nil {
			return nil, err
		} else if len(res) > 0 {
			for k, v := range res {
				if _, ok := result[k]; !ok {
					result[k] = make([]Result[TId], 0)
				}

				var nxtIds []TId
				for _, nxt := range maps.Keys(v.Dependency) {
					if len(v.Parameters) < 1 || v.Parameters[nxt]["skip"] == nil || !v.Parameters[nxt]["skip"].(bool) {
						nxtIds = append(nxtIds, nxt)
					}
				}

				if cascade && len(er.dependency[v.EntityName]) > 0 && len(nxtIds) > 0 {
					res, err := er.dependencyLoop(er.resolver[v.EntityName], ctx, cascade, nxtIds)
					if err != nil {
						return nil, err
					} else if len(res) > 0 {
						for i, j := range res {
							v.Dependency[i] = j
						}
					}
				}

				result[k] = append(result[k], v)
			}
		}
	}
	return result, nil
}

// fn return: (*parentId, *entityId)
func MapResult[TEntity any, TId comparable](items []TEntity, fn func(item *TEntity) (parentId, entityId *TId), params ...func(parentId, entityId *TId, item *TEntity) map[string]any) ResolveResult[TId] {
	tmp := make(map[TId]map[TId][]Result[TId])
	ext := make(map[TId]map[TId]map[string]any)

	for _, item := range items {
		k, v := fn(&item)
		if k == nil || v == nil {
			continue
		}
		if _, ok := tmp[*k]; !ok {
			tmp[*k] = make(map[TId][]Result[TId])
		}
		if _, ok := tmp[*k][*v]; !ok {
			tmp[*k][*v] = nil
		}
		if len(params) > 0 {
			if _, ok := ext[*k]; !ok {
				ext[*k] = make(map[TId]map[string]any)
			}
			ext[*k][*v] = params[0](k, v, &item)
		}
	}

	entityName, _ := util.GetEntityName[TEntity]()
	result := make(map[TId]Result[TId])
	for k, v := range tmp {
		res := Result[TId]{
			EntityName: entityName,
			Dependency: v,
		}
		if p, ok := ext[k]; ok && len(p) > 0 {
			res.Parameters = p
		}
		result[k] = res
	}
	return result
}
