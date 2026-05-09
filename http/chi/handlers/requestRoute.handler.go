package handlers

import (
	"context"
	"net/http"

	libChi "http/chi"
)

func UrlRouteParamsCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		routeParams, err := libChi.DecodeChiRouteParams(r.Context(), r)
		if err == nil {
			ctx := context.WithValue(r.Context(), "RouteParams", routeParams)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func UrlTRouteParamsCtx[T any](next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		routeParams, err := libChi.DecodeChiTRouteParams[T](r.Context(), r)
		if err == nil {
			ctx := context.WithValue(r.Context(), "RouteParams", routeParams)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
