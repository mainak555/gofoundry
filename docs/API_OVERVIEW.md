# API Overview

This page summarizes key extension points and public contracts.

## Authentication

- auth.GetJWT
- auth.GetClaims
- auth.ValidateJWTIssuer
- auth.ValidateJWTSignature
- http/oidc.ValidateToken middleware
- http/oidc.IsRoleAny middleware

Key middleware signatures:

```go
func ValidateToken[TClaims any](
	config *models.OidcConfig,
	jwFn func(idToken *oidc.IDToken, clientId *string) auth.IJWT[TClaims],
) func(http.Handler) http.Handler

func IsRoleAny[TClaims any](roles ...string) func(http.Handler) http.Handler
```

## Caching

- cache.ICache, cache.IMongoCache, cache.IRedisCache
- cache.NewRedisCache, cache.NewMongoCache, cache.NewCache
- cache.Get, cache.GetMany, cache.Set, cache.SetJson, cache.Delete

Example bootstrap:

```go
cache.NewMongoCache("secondary_cache", mongoClient.GetDatabase(), 60)
```

## MongoDB and Repository

- db/mongodb.IMongoClient and MongoClient lifecycle methods.
- repository/interfaces.IMongoRepository and TMongoRepository[T].
- repository.NewMongoRepository and repository.NewMongoTRepository[T].
- generics.ICommonService[T] and generics.NewCommonService[T].

Key signature:

```go
func NewMongoClient(
	ctx *context.Context,
	conStr string,
	dbName string,
	_clientOpts func(*options.ClientOptions) *options.ClientOptions,
	_dbOps func(*options.DatabaseOptions) *options.DatabaseOptions,
) (IMongoClient, error)
```

## HTTP Transport

- http/server decode and encode helpers.
- http/chi decode helpers and middleware package.
- generics.ApiRouter fluent route composition methods.

Useful transport helpers:

- `server.KitHttpServerOptions(logger)` for default go-kit server option stack.
- `server.EncodeErrorResponse` for prefix-based error-to-status mapping.
- `server.NewServerFinalizer(logger)` for request status logging.
- `middlewares.AddMiddlewares(router, origins, custom...)` for standardized chi middleware setup.
- `handlers.ValidateRouteIdCtx[T]` for route id validation and context enrichment.

ApiRouter methods used most often:

- `AddGetAll`, `AddCreateOne`, `AddGetById`, `AddUpdateById`, `AddDeleteById`, `AddDeleteMany`
- `With`, `Route`, `Mount`, `Group`, `Method`

## DTO and Helper Utilities

- dtos pagination/filter/sort request structures.
- helpers query and context accessor helpers.
- util shared conversion and reflection utilities.

Common context helpers:

- `helpers.GetPagination(ctx)`
- `helpers.GetFilter(ctx)`
- `helpers.GetSort(ctx)`
- `helpers.SetQueryFilterCtx(ctx, filter)`
- `helpers.SetQuerySortCtx(ctx, sort)`

See also:

- `HTTP_TRANSPORT.md` for route and middleware composition examples.
- `ADVANCED_EXAMPLES.md` for production startup wiring.
