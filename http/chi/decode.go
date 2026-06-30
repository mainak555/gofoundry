package chi

import (
	"context"
	"errors"
	"net/http"

	"github.com/mainak555/gofoundry/http/server"
	"github.com/mainak555/gofoundry/util"

	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slices"
)

func DecodeChiRouteParams(_ context.Context, r *http.Request) (interface{}, error) {
	if rctx := chi.RouteContext(r.Context()); rctx != nil {
		chars := []string{"*", "#"}
		params := make(map[string]string)
		for i, key := range rctx.URLParams.Keys {
			if !slices.Contains(chars, key) {
				params[key] = rctx.URLParams.Values[i]
			}
		}
		return params, nil
	}
	return nil, errors.New("route params missing")
}

func DecodeChiQueryParams(_ context.Context, r *http.Request) (interface{}, error) {
	queryParams := r.URL.Query()
	params := make(map[string]string)
	for key, values := range queryParams {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}
	return params, nil
}

func DecodeChiTRouteParams[T any](_ context.Context, r *http.Request) (interface{}, error) {
	var entity T
	err := util.StructMap[T](func(key string) string {
		return chi.URLParam(r, key)
	}, &entity)
	return entity, err
}

func DecodeChiTQueryParams[T any](_ context.Context, r *http.Request) (interface{}, error) {
	var entity T
	err := util.StructMap[T](func(key string) string {
		return r.URL.Query().Get(key)
	}, &entity)
	return entity, err
}

func DecodeChiQueryAndRouteParams[Q, R any](ctx context.Context, r *http.Request) (interface{}, error) {
	routeParams, err := DecodeChiTRouteParams[R](ctx, r)
	if err != nil {
		return nil, err
	}
	queryParams, err := DecodeChiTQueryParams[Q](ctx, r)
	if err != nil {
		return nil, err
	}

	result := make(map[string]any)
	result["RouteParams"] = routeParams
	result["QueryParams"] = queryParams
	return result, nil
}

func DecodeChiQueryParamsWithBody[Q, B any](ctx context.Context, r *http.Request) (interface{}, error) {
	queryParams, err := DecodeChiTQueryParams[Q](ctx, r)
	if err != nil {
		return nil, err
	}
	body, err := server.DecodeRequestBody[B](ctx, r)
	if err != nil {
		return nil, err
	}

	result := make(map[string]any)
	result["QueryParams"] = queryParams
	result["Body"] = body
	return result, nil
}

func DecodeChiRouteParamsWithBody[R, B any](ctx context.Context, r *http.Request) (interface{}, error) {
	routeParams, err := DecodeChiTRouteParams[R](ctx, r)
	if err != nil {
		return nil, err
	}
	body, err := server.DecodeRequestBody[B](ctx, r)
	if err != nil {
		return nil, err
	}

	result := make(map[string]any)
	result["RouteParams"] = routeParams
	result["Body"] = body
	return result, nil
}

func DecodeChiRequestParamsWithBody[Q, R, B any](ctx context.Context, r *http.Request) (interface{}, error) {
	data, err := DecodeChiQueryAndRouteParams[Q, R](ctx, r)
	if err != nil {
		return nil, err
	}
	body, err := server.DecodeRequestBody[B](ctx, r)
	if err != nil {
		return nil, err
	}

	result := data.(map[string]any)
	result["Body"] = body
	return result, nil
}
