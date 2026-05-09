package helpers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"dtos"
)

func GetPagination(ctx *context.Context) *dtos.Pagination {
	return (*ctx).Value("Pagination").(*dtos.Pagination)
}

func GetFilter(ctx *context.Context) *dtos.QueryFilter {
	return (*ctx).Value("QueryFilter").(*dtos.QueryFilter)
}

func GetFilterQuery(ctx *context.Context) map[string]any {
	filter := (*ctx).Value("QueryFilter")
	if filter != nil && filter.(*dtos.QueryFilter) != nil {
		return filter.(*dtos.QueryFilter).Filter
	}
	return nil
}

func SetQueryFilterCtx(ctx context.Context, query map[string]any) context.Context {
	return context.WithValue(ctx, "QueryFilter", &dtos.QueryFilter{
		Filter: query,
	})
}

func GetSort(ctx *context.Context) *dtos.QuerySort {
	return (*ctx).Value("QuerySort").(*dtos.QuerySort)
}

func GetSortQuery(ctx *context.Context) map[string]*bool {
	sort := (*ctx).Value("QuerySort")
	if sort != nil && sort.(*dtos.QuerySort) != nil {
		return sort.(*dtos.QuerySort).Sort
	}
	return nil
}

func SetQuerySortCtx(ctx context.Context, query map[string]*bool) context.Context {
	return context.WithValue(ctx, "QuerySort", &dtos.QuerySort{
		Sort: query,
	})
}

func GetUrlRouteParams(ctx *context.Context) map[string]string {
	return (*ctx).Value("RouteParams").(map[string]string)
}

func GetUrlTRouteParams[T any](ctx *context.Context) T {
	return (*ctx).Value("RouteParams").(T)
}

func GetUrlQueryParams(ctx *context.Context) map[string][]string {
	return (*ctx).Value("QueryParams").(map[string][]string)
}

func GetUrlTQueryParams[T any](ctx *context.Context) T {
	return (*ctx).Value("QueryParams").(T)
}

func ExecFilterQuery(matchKey string, filterParam *dtos.QueryFilter,
	exec func(value interface{}) (interface{}, error)) (interface{}, error) {
	if len(filterParam.Filter) > 0 {
		for k, v := range filterParam.Filter {
			if strings.ToUpper(k) == strings.ToUpper(matchKey) {
				if exec != nil {
					return exec(v)
				}
				return v, nil
			}
		}
	}
	return nil, nil
}

func InvalidError(message string, errs ...error) error {
	msg := fmt.Sprintf("invalid: %v", message)
	if err := MergeErrors(errs...); err != nil {
		return fmt.Errorf("%v, %v", msg, err)
	}
	return errors.New(msg)
}

func NoContentError(message string, errs ...error) error {
	msg := fmt.Sprintf("no content err: %v", message)
	if err := MergeErrors(errs...); err != nil {
		return fmt.Errorf("%v, %v", msg, err)
	}
	return errors.New(msg)
}

func NoContent(message string) error {
	return fmt.Errorf("no content: %v", message)
}

func MergeErrors(errs ...error) error {
	if errs != nil {
		var errStr string
		for _, e := range errs {
			if e != nil {
				errStr = fmt.Sprintf("%v; %v", errStr, e)
			}
		}
		if errStr != "" {
			return fmt.Errorf("%v", errStr)
		}
	}
	return nil
}

func GetPaginated[T any](data []T, page *dtos.Pagination) *dtos.ResponseWithCount[T] {
	start := (page.PageNo - 1) * page.PageSize
	end := start + page.PageSize
	totLen := int64(len(data))
	if start >= totLen {
		return nil
	} else if end > totLen {
		end = totLen
	}
	return &dtos.ResponseWithCount[T]{
		Data:  data[start:end],
		Count: totLen,
	}
}
