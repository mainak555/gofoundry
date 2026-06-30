package handlers

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	lm "github.com/mainak555/gofoundry/db/mongodb"
	"github.com/mainak555/gofoundry/dtos"
	"github.com/mainak555/gofoundry/generics"
	"github.com/mainak555/gofoundry/http/chi"
	"github.com/mainak555/gofoundry/http/chi/renderers"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson"
)

func ValidateRouteIdCtx[T any](
	mc lm.IMongoClient,
	entityName string,
	_queryFn func(ctx *context.Context, query bson.M) bson.M,
	_setContext func(r *http.Request, Id primitive.ObjectID) context.Context,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var entity T
			if reflect.TypeOf(entity).Kind() != reflect.Struct {
				render.Render(w, r, &renderers.ErrResponse{HTTPStatusCode: 500, ReasonPhrase: "unknown data type"})
				return
			}

			ctx := r.Context()
			var routeParams any
			if routeParams, _ = chi.DecodeChiRouteParams(ctx, r); routeParams != nil && routeParams.(map[string]string)["id"] == "" {
				render.Render(w, r, &renderers.ErrResponse{HTTPStatusCode: 400, ReasonPhrase: "missing route"})
				return
			}

			Id, err := primitive.ObjectIDFromHex(routeParams.(map[string]string)["id"])
			if err != nil {
				render.Render(w, r, &renderers.ErrResponse{HTTPStatusCode: 400, ReasonPhrase: fmt.Sprintf("invalid %vId", entityName)})
				return
			}

			service := generics.NewCommonService[T](mc)
			query := bson.M{"$and": []bson.M{{"_id": Id}, {"_deleted": false}}}
			if _queryFn != nil {
				query = _queryFn(&ctx, query)
			}
			if res, err := service.Get(&ctx, query, &dtos.Pagination{
				PageNo:   1,
				PageSize: 1,
			}); err != nil {
				render.Render(w, r, &renderers.ErrResponse{HTTPStatusCode: 500, ReasonPhrase: err.Error()})
				return
			} else if res == nil || len(res) < 1 {
				render.Render(w, r, &renderers.ErrResponse{HTTPStatusCode: 404, ReasonPhrase: fmt.Sprintf("%v not found!", entityName)})
				return
			}

			if _setContext != nil {
				next.ServeHTTP(w, r.WithContext(_setContext(r, Id)))
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}
