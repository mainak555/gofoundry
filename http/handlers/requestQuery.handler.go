package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"gofoundry/dtos"
	"gofoundry/helpers"
	"gofoundry/util"

	"golang.org/x/exp/slices"
)

// RequestPaginationCtx reads pageNo/pageSize query parameters and stores pagination in request context.
func RequestPaginationCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := &dtos.Pagination{}
		var err error

		pageNoStr := r.URL.Query().Get("pageNo")
		var value int
		if value, err = strconv.Atoi(pageNoStr); err != nil {
			page.PageNo = 1
		} else {
			page.PageNo = int64(value)
		}

		pageSizeStr := r.URL.Query().Get("pageSize")
		if value, err = strconv.Atoi(pageSizeStr); err != nil {
			page.PageSize = 50
		} else {
			page.PageSize = int64(value)
		}

		if page.PageSize == -1 {
			page.PageSize = util.MAX_PAGE_SIZE
			page.PageNo = 1
		}

		ctx := context.WithValue(r.Context(), "Pagination", page)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UrlQueryFilterCtx parses the filter query parameter as JSON and stores it in context.
func UrlQueryFilterCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := make(map[string]interface{})

		jsonQuery := r.URL.Query().Get("filter")
		var byteText []byte = []byte(jsonQuery)
		_ = json.Unmarshal([]byte(byteText), &query)

		ctx := helpers.SetQueryFilterCtx(r.Context(), query)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UrlQuerySortCtx parses the sort query parameter as JSON and stores it in context.
func UrlQuerySortCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := make(map[string]*bool)
		jsonQuery := r.URL.Query().Get("sort")
		var byteText []byte = []byte(jsonQuery)
		_ = json.Unmarshal([]byte(byteText), &query)

		ctx := helpers.SetQuerySortCtx(r.Context(), query)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UrlQueryParamsCtx stores non-filter query parameters in request context.
func UrlQueryParamsCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		params := make(map[string][]string)
		for key, values := range queryParams {
			if len(values) > 0 && !slices.Contains([]string{"filter", "sort", "pageNo", "pageSize"}, key) {
				params[key] = values
			}
		}
		ctx := context.WithValue(r.Context(), "QueryParams", params)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UrlTQueryParamsCtx maps URL query parameters into a typed struct and stores it in context.
func UrlTQueryParamsCtx[T any](next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params T
		util.StructMap(func(key string) string {
			return r.URL.Query().Get(key)
		}, &params)
		ctx := context.WithValue(r.Context(), "QueryParams", params)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
