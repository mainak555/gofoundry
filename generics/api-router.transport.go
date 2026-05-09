package generics

import (
	"net/http"

	libmongo "db/mongodb"
	"http/server"

	"github.com/go-chi/chi/v5"
	kithttp "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
)

type ApiRouter struct {
	Router             chi.Router
	Endpoint           IEndpoints
	Logger             kitlog.Logger
	MongoClient        libmongo.IMongoClient
	KhttpServerOptions []kithttp.ServerOption
	Decoder            IApiRouterRequestDecoder
}

func NewApiRouter(mongoClient libmongo.IMongoClient, logger kitlog.Logger, khttpServerOptions []kithttp.ServerOption,
	endpoint IEndpoints, decoders IApiRouterRequestDecoder, router chi.Router) *ApiRouter {
	var r chi.Router
	if router != nil {
		r = router
	} else {
		r = chi.NewRouter()
	}
	return &ApiRouter{
		Router:             r,
		Endpoint:           endpoint,
		Decoder:            decoders,
		KhttpServerOptions: khttpServerOptions,
		MongoClient:        mongoClient,
		Logger:             logger,
	}
}

func (ar *ApiRouter) Extend(apiRouter func(ar *ApiRouter)) *ApiRouter {
	apiRouter(ar)
	return ar
}

func (ar *ApiRouter) AddGetAll() *ApiRouter {
	ar.Router.Method("GET", "/", kithttp.NewServer(
		ar.Endpoint.GetAllEndpoint(),
		ar.Decoder.GetAllDeoder(),
		server.EncodeJsonResponse,
		ar.KhttpServerOptions...,
	))
	return ar
}

func (ar *ApiRouter) AddCreateOne(decoderFn kithttp.DecodeRequestFunc) *ApiRouter {
	ar.Router.Method("POST", "/", kithttp.NewServer(
		ar.Endpoint.CreateOneEndpoint(),
		decoderFn,
		server.EncodeJsonResponse,
		ar.KhttpServerOptions...,
	))
	return ar
}

func (ar *ApiRouter) AddGetById() *ApiRouter {
	ar.Router.Get("/{id}", kithttp.NewServer(
		ar.Endpoint.GetByIdEndpoint(),
		ar.Decoder.GetByIdDecoder(),
		server.EncodeJsonResponse,
		ar.KhttpServerOptions...,
	).ServeHTTP)
	return ar
}

func (ar *ApiRouter) AddUpdateById(decoderFn kithttp.DecodeRequestFunc) *ApiRouter {
	ar.Router.Put("/{id}", kithttp.NewServer(
		ar.Endpoint.UpdateByIdEndpoint(),
		decoderFn,
		server.EncodeJsonResponse,
		ar.KhttpServerOptions...,
	).ServeHTTP)
	return ar
}

func (ar *ApiRouter) AddDeleteById() *ApiRouter {
	ar.Router.Delete("/{id}", kithttp.NewServer(
		ar.Endpoint.DeleteByIdEndpoint(),
		ar.Decoder.DeleteByIdDecoder(),
		server.EncodeJsonResponse,
		ar.KhttpServerOptions...,
	).ServeHTTP)
	return ar
}

func (ar *ApiRouter) AddDeleteMany() *ApiRouter {
	ar.Router.Put("/delete", kithttp.NewServer(
		ar.Endpoint.DeleteEndpoint(),
		ar.Decoder.DeleteManyDecoder(),
		server.EncodeJsonResponse,
		ar.KhttpServerOptions...,
	).ServeHTTP)
	return ar
}

func (ar *ApiRouter) AddRoute(pattern string, endpoint IEndpoints, decoder IApiRouterRequestDecoder, routeFn func(ar *ApiRouter)) *ApiRouter {
	ar.Router.Route(pattern, func(r chi.Router) {
		childRouter := NewApiRouter(ar.MongoClient, ar.Logger, ar.KhttpServerOptions, endpoint, decoder, r)
		routeFn(childRouter)
	})
	return ar
}

func (ar *ApiRouter) Get(pattern string, h http.HandlerFunc) *ApiRouter {
	ar.Router.Get(pattern, h)
	return ar
}

func (ar *ApiRouter) Post(pattern string, h http.HandlerFunc) *ApiRouter {
	ar.Router.Post(pattern, h)
	return ar
}

func (ar *ApiRouter) Put(pattern string, h http.HandlerFunc) *ApiRouter {
	ar.Router.Put(pattern, h)
	return ar
}

func (ar *ApiRouter) Patch(pattern string, h http.HandlerFunc) *ApiRouter {
	ar.Router.Patch(pattern, h)
	return ar
}

func (ar *ApiRouter) Delete(pattern string, h http.HandlerFunc) *ApiRouter {
	ar.Router.Delete(pattern, h)
	return ar
}

func (ar *ApiRouter) Mount(pattern string, h http.Handler) *ApiRouter {
	ar.Router.Mount(pattern, h)
	return ar
}

func (ar *ApiRouter) Group(fn func(r chi.Router)) *ApiRouter {
	ar.Router.Group(fn)
	return ar
}

func (ar *ApiRouter) Method(method string, pattern string, h http.Handler) *ApiRouter {
	ar.Router.Method(method, pattern, h)
	return ar
}

func (ar *ApiRouter) Route(pattern string, fn func(r chi.Router)) *ApiRouter {
	ar.Router.Route(pattern, fn)
	return ar
}

func (ar *ApiRouter) With(middlewares ...func(http.Handler) http.Handler) *ApiRouter {
	childRouter := NewApiRouter(ar.MongoClient, ar.Logger, ar.KhttpServerOptions, ar.Endpoint, ar.Decoder, ar.Router.With(middlewares...))
	return childRouter
}

func (ar *ApiRouter) New(endpoint IEndpoints, decoder IApiRouterRequestDecoder, middlewares ...func(http.Handler) http.Handler) *ApiRouter {
	childRouter := NewApiRouter(ar.MongoClient, ar.Logger, ar.KhttpServerOptions, endpoint, decoder, ar.Router.With(middlewares...))
	return childRouter
}
