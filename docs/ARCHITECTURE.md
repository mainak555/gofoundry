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

## Design Principles

- Interface-first boundaries for substitution and testing.
- Generic abstractions where type safety improves reuse.
- Decoupled transport, business service, and persistence layers.
- Composable middleware and router extension patterns.
