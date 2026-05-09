package generics

import (
	"net/http"

	libmongo "gofoundry/db/mongodb"
	"gofoundry/http/server"

	"github.com/go-chi/chi/v5"
	kithttp "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
)

// ApiRouter provides fluent route composition for endpoint-driven HTTP handlers.
type ApiRouter struct {
	Router             chi.Router
	Endpoint           IEndpoints
	Logger             kitlog.Logger
	MongoClient        libmongo.IMongoClient
	KhttpServerOptions []kithttp.ServerOption
	Decoder            IApiRouterRequestDecoder
}

// NewApiRouter creates a composable chi router wrapper for go-kit endpoints.
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

// Extend applies apiRouter to the current router and returns the same instance.
func (ar *ApiRouter) Extend(apiRouter func(ar *ApiRouter)) *ApiRouter {
	apiRouter(ar)
	return ar
}

// AddGetAll registers the default list route on the current router.
func (ar *ApiRouter) AddGetAll() *ApiRouter {
	ar.Router.Method("GET", "/", kithttp.NewServer(
		ar.Endpoint.GetAllEndpoint(),
		ar.Decoder.GetAllDeoder(),
		server.EncodeJsonResponse,
		ar.KhttpServerOptions...,
	))
	return ar
}

// AddCreateOne registers the default create route on the current router.
func (ar *ApiRouter) AddCreateOne(decoderFn kithttp.DecodeRequestFunc) *ApiRouter {
	ar.Router.Method("POST", "/", kithttp.NewServer(
		ar.Endpoint.CreateOneEndpoint(),
		decoderFn,
		server.EncodeJsonResponse,
		ar.KhttpServerOptions...,
	))
	return ar
}

// AddGetById registers the default read-by-id route on the current router.
func (ar *ApiRouter) AddGetById() *ApiRouter {
	ar.Router.Get("/{id}", kithttp.NewServer(
		ar.Endpoint.GetByIdEndpoint(),
		ar.Decoder.GetByIdDecoder(),
		server.EncodeJsonResponse,
		ar.KhttpServerOptions...,
	).ServeHTTP)
	return ar
}

// AddUpdateById registers the default update-by-id route on the current router.
func (ar *ApiRouter) AddUpdateById(decoderFn kithttp.DecodeRequestFunc) *ApiRouter {
	ar.Router.Put("/{id}", kithttp.NewServer(
		ar.Endpoint.UpdateByIdEndpoint(),
		decoderFn,
		server.EncodeJsonResponse,
		ar.KhttpServerOptions...,
	).ServeHTTP)
	return ar
}

// AddDeleteById registers the default delete-by-id route on the current router.
func (ar *ApiRouter) AddDeleteById() *ApiRouter {
	ar.Router.Delete("/{id}", kithttp.NewServer(
		ar.Endpoint.DeleteByIdEndpoint(),
		ar.Decoder.DeleteByIdDecoder(),
		server.EncodeJsonResponse,
		ar.KhttpServerOptions...,
	).ServeHTTP)
	return ar
}

// AddDeleteMany registers the default bulk-delete route on the current router.
func (ar *ApiRouter) AddDeleteMany() *ApiRouter {
	ar.Router.Put("/delete", kithttp.NewServer(
		ar.Endpoint.DeleteEndpoint(),
		ar.Decoder.DeleteManyDecoder(),
		server.EncodeJsonResponse,
		ar.KhttpServerOptions...,
	).ServeHTTP)
	return ar
}

// AddRoute creates a nested route group using endpoint and decoder overrides.
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

// With returns a child router with additional middlewares applied.
func (ar *ApiRouter) With(middlewares ...func(http.Handler) http.Handler) *ApiRouter {
	childRouter := NewApiRouter(ar.MongoClient, ar.Logger, ar.KhttpServerOptions, ar.Endpoint, ar.Decoder, ar.Router.With(middlewares...))
	return childRouter
}

// New returns a child router with endpoint and decoder overrides.
func (ar *ApiRouter) New(endpoint IEndpoints, decoder IApiRouterRequestDecoder, middlewares ...func(http.Handler) http.Handler) *ApiRouter {
	childRouter := NewApiRouter(ar.MongoClient, ar.Logger, ar.KhttpServerOptions, endpoint, decoder, ar.Router.With(middlewares...))
	return childRouter
}
