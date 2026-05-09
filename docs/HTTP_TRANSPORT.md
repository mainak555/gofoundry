# HTTP Transport

This guide focuses on composing chi routing with go-kit handlers using GoFoundry transport helpers.

## 1. Base Router Middleware Setup

Use `middlewares.AddMiddlewares` for baseline CORS, request id, real ip, and security headers.

```go
router := chi.NewRouter()

middlewares.AddMiddlewares(
	router,
	settings.AllowOrigins,
	middleware.Timeout(2*time.Minute),
)
```

The helper applies:

- CORS handling
- `middleware.RequestID`
- `middleware.RealIP`
- your custom middlewares
- `middleware.Recoverer`
- `middleware.NoCache`
- security headers

## 2. go-kit Server Options

Recommended defaults:

```go
serverOptions := server.KitHttpServerOptions(logger)
```

Custom stack:

```go
serverOptions := []kithttp.ServerOption{
	kithttp.ServerErrorHandler(handlers.NewLogErrorHandler(logger)),
	kithttp.ServerErrorEncoder(server.EncodeErrorResponse),
	kithttp.ServerFinalizer(server.NewServerFinalizer(logger)),
}
```

## 3. Request Decoding Helpers

Use chi decoders for route/query params:

```go
routeParams, _ := chi.DecodeChiRouteParams(ctx, req)
queryParams, _ := chi.DecodeChiQueryParams(ctx, req)

type RouteParams struct {
	Id string `json:"id"`
}
type QueryParams struct {
	Search string `json:"search"`
}

typedRoute, _ := chi.DecodeChiTRouteParams[RouteParams](ctx, req)
typedQuery, _ := chi.DecodeChiTQueryParams[QueryParams](ctx, req)

_ = routeParams
_ = queryParams
_ = typedRoute
_ = typedQuery
```

For body payloads, pair with `server.DecodeRequestBody[T]` in your decoder functions.

## 4. Endpoint Route Assembly With ApiRouter

`generics.ApiRouter` provides fluent route composition around go-kit endpoints.

```go
api := generics.NewApiRouter(
	mongoClient,
	logger,
	serverOptions,
	endpoints,
	decoders,
	nil,
)

createDecoder := func(ctx context.Context, r *http.Request) (interface{}, error) {
	return server.DecodeRequestBody[CreateRequest](ctx, r)
}

updateDecoder := func(ctx context.Context, r *http.Request) (interface{}, error) {
	return server.DecodeRequestBody[UpdateRequest](ctx, r)
}

api.
	AddGetAll().
	AddCreateOne(createDecoder).
	AddGetById().
	AddUpdateById(updateDecoder).
	AddDeleteById()
```

`CreateRequest` and `UpdateRequest` are app-defined payload types.

Nested groups:

```go
api.Route("/v1", func(r chi.Router) {
	r.Mount("/entities", api.Router)
})
```

## 5. OIDC and Role-Based Guards

Apply token validation before role checks:

```go
type Claims struct {
	Sub string `json:"sub"`
}

var jwtFactory func(idToken *oidc.IDToken, clientId *string) auth.IJWT[Claims]

router.With(
	oidc.ValidateToken(settings.Oidc, jwtFactory),
	oidc.IsRoleAny[Claims]("ROLE-A", "ROLE-B"),
).Route("/v1", func(r chi.Router) {
	// protected routes
})
```

## 6. Route Id Validation Middleware

Use `ValidateRouteIdCtx` to validate route id existence and enrich context.

```go
type MyEntity struct{}

router.With(chiHandlers.ValidateRouteIdCtx[MyEntity](
	mongoClient,
	"entity",
	nil,
	func(r *http.Request, id primitive.ObjectID) context.Context {
		return context.WithValue(r.Context(), "entityId", id)
	},
))
```

## 7. Response Encoding

Transport handlers should use:

- `server.EncodeJsonResponse` for successful responses.
- `server.EncodeErrorResponse` for prefix-based status mapping.

The error encoder maps common prefixes such as `invalid`, `unauthorize`, and `no content` to HTTP status codes.

## 8. Observability Hook

`server.NewServerFinalizer(logger)` logs status, path, and method after handler execution.

When you also need telemetry, add a custom `kithttp.ServerErrorHandler` in server options.

## See Also

- `GETTING_STARTED.md` for minimal setup.
- `ADVANCED_EXAMPLES.md` for full production bootstrap.
- `ERROR_HANDLING.md` for error mapping details.
