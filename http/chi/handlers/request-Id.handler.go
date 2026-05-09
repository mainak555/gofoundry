package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"gofoundry/http/chi/renderers"
)

func ValidateIdCtx(idValidator func(ctx *context.Context, id string) error) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			strId := chi.URLParam(r, "id")
			if strId == "" {
				render.Render(w, r, renderers.ErrBadRequest)
				return
			} else if idValidator != nil {
				if err := idValidator(&ctx, strId); err != nil {
					render.Render(w, r, renderers.ErrInvalidRequest(err))
					return
				}
			}
			ctx = context.WithValue(ctx, "Id", strId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
