# Architecture

GoFoundry is organized as reusable packages that can be composed per service.

## Package Responsibilities

- auth: JWT token parsing and validation helpers.
- cache: cache contract and dual backend client (Redis and MongoDB fallback).
- db: shared data entities.
- db/mongodb: Mongo connection lifecycle and aggregation/cursor helpers.
- dependency: generic dependency graph resolution.
- dtos: shared request/response payload structures.
- generics: common CRUD services, endpoint contracts, and router composition.
- helpers: filter/sort/query and request context helpers.
- http/chi: route decoding, middleware, and renderers.
- http/server: decode/encode helpers and server options.
- repository: typed and untyped Mongo repository implementations.
- util: config loading, reflection helpers, date/crypto helpers.

## Typical Request Flow

1. Incoming HTTP request enters chi middleware stack.
2. Query/path/body decoders transform input into DTOs.
3. Endpoint layer invokes generic service/repository logic.
4. Repository executes Mongo queries through db/mongodb client.
5. Response is encoded via http/server transport helpers.

## Concrete Composition Flow

At service startup, a typical production bootstrap is:

1. `util.Configure[T]` resolves `NODE_ENV`, then loads yml/env into your config model.
2. `mongodb.NewMongoClient` creates and connects the database client.
3. `IMongoClient.ConfigureCollections` registers collections and verifies/adds indexes.
4. `middlewares.AddMiddlewares` applies CORS + request/security middleware.
5. `oidc.ValidateToken` and optional `oidc.IsRoleAny` guard protected routes.
6. `server.KitHttpServerOptions` (or explicit options) defines error/finalizer behavior.
7. `generics.ApiRouter` assembles endpoint-driven REST routes.
8. `http.Server` serves traffic and shuts down on signal.

## Request Pipeline Example

The following stack mirrors the library's intended layering:

```text
chi.Router
	-> middlewares.AddMiddlewares
	-> route-level guards (oidc.ValidateToken -> oidc.IsRoleAny)
	-> go-kit transport handler (kithttp.NewServer)
	-> endpoint
	-> generics.NewCommonService[T]
	-> repository.NewMongoTRepository[T] / mongo client
	-> server.EncodeJsonResponse or server.EncodeErrorResponse
```

## Collection and Index Lifecycle

`ConfigureCollections` supports variadic collection descriptors and handles both create-time and verify-time flows:

- Creates missing collections.
- Creates all declared indexes for new collections.
- For existing collections, compares index names and creates only missing ones.
- Registers entity type names to collection names for generic repository/service lookup.

## Design Principles

- Interface-first boundaries for substitution and testing.
- Generic abstractions where type safety improves reuse.
- Decoupled transport, business service, and persistence layers.
- Composable middleware and router extension patterns.
- Startup-first reliability: initialize config, secrets, data stores, and indexes before serving traffic.
