package helpers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gofoundry/dtos"
)

// GetPagination returns pagination data from context.
func GetPagination(ctx *context.Context) *dtos.Pagination {
	return (*ctx).Value("Pagination").(*dtos.Pagination)
}

// GetFilter returns query filter data from context.
func GetFilter(ctx *context.Context) *dtos.QueryFilter {
	return (*ctx).Value("QueryFilter").(*dtos.QueryFilter)
}

// GetFilterQuery returns filter map from context when available.
func GetFilterQuery(ctx *context.Context) map[string]any {
	filter := (*ctx).Value("QueryFilter")
	if filter != nil && filter.(*dtos.QueryFilter) != nil {
		return filter.(*dtos.QueryFilter).Filter
	}
	return nil
}

// SetQueryFilterCtx stores query filter map in context.
func SetQueryFilterCtx(ctx context.Context, query map[string]any) context.Context {
	return context.WithValue(ctx, "QueryFilter", &dtos.QueryFilter{
		Filter: query,
	})
}

// GetSort returns query sort data from context.
func GetSort(ctx *context.Context) *dtos.QuerySort {
	return (*ctx).Value("QuerySort").(*dtos.QuerySort)
}

// GetSortQuery returns sort map from context when available.
func GetSortQuery(ctx *context.Context) map[string]*bool {
	sort := (*ctx).Value("QuerySort")
	if sort != nil && sort.(*dtos.QuerySort) != nil {
		return sort.(*dtos.QuerySort).Sort
	}
	return nil
}

// SetQuerySortCtx stores query sort map in context.
func SetQuerySortCtx(ctx context.Context, query map[string]*bool) context.Context {
	return context.WithValue(ctx, "QuerySort", &dtos.QuerySort{
		Sort: query,
	})
}

// GetUrlRouteParams returns route parameters from context.
func GetUrlRouteParams(ctx *context.Context) map[string]string {
	return (*ctx).Value("RouteParams").(map[string]string)
}

// GetUrlTRouteParams returns route parameters from context cast to T.
func GetUrlTRouteParams[T any](ctx *context.Context) T {
	return (*ctx).Value("RouteParams").(T)
}

// GetUrlQueryParams returns query parameters from context.
func GetUrlQueryParams(ctx *context.Context) map[string][]string {
	return (*ctx).Value("QueryParams").(map[string][]string)
}

// GetUrlTQueryParams returns query parameters from context cast to T.
func GetUrlTQueryParams[T any](ctx *context.Context) T {
	return (*ctx).Value("QueryParams").(T)
}

// ExecFilterQuery finds a filter key and optionally transforms it using exec.
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

// InvalidError creates an error prefixed with invalid for transport mapping.
func InvalidError(message string, errs ...error) error {
	msg := fmt.Sprintf("invalid: %v", message)
	if err := MergeErrors(errs...); err != nil {
		return fmt.Errorf("%v, %v", msg, err)
	}
	return errors.New(msg)
}

// NoContentError creates an error prefixed with no content err for transport mapping.
func NoContentError(message string, errs ...error) error {
	msg := fmt.Sprintf("no content err: %v", message)
	if err := MergeErrors(errs...); err != nil {
		return fmt.Errorf("%v, %v", msg, err)
	}
	return errors.New(msg)
}

// NoContent creates an error prefixed with no content for transport mapping.
func NoContent(message string) error {
	return fmt.Errorf("no content: %v", message)
}

// MergeErrors joins non-nil errors into a single error value.
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

// GetPaginated returns a page slice and total count from in-memory data.
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
