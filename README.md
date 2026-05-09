# GoFoundry

GoFoundry is a reusable Go library for building service backends with a consistent stack around authentication, HTTP transport, MongoDB data access, caching, and utility helpers.

This repository focuses on composable building blocks rather than a framework runtime, so teams can adopt only the packages they need.

## Core Composition Pattern

Most services built with GoFoundry follow this startup sequence:

1. Load environment-aware settings with `util.Configure[T]`.
2. Optionally hydrate secrets from Azure Key Vault with `util.GetVault` and `util.GetFromVault`.
3. Create and connect Mongo with `mongodb.NewMongoClient`.
4. Register collections and indexes with `IMongoClient.ConfigureCollections`.
5. Compose chi middlewares with `middlewares.AddMiddlewares`.
6. Attach OIDC guards with `oidc.ValidateToken` and `oidc.IsRoleAny`.
7. Build endpoint handlers with `generics.ApiRouter`.
8. Wire transport behavior with `server.KitHttpServerOptions` or explicit go-kit options.
9. Start HTTP server and gracefully shutdown on OS signals.

## What Is Included

- Authentication helpers for JWT/OIDC workflows.
- HTTP transport helpers for chi and go-kit patterns.
- MongoDB repository and generic CRUD service abstractions.
- Cache abstractions for Redis and MongoDB fallback.
- DTOs, config models, and utility helpers used by service modules.

## Installation

Add the module to your project:

```bash
go get github.com/mainak555/gofoundry
```

## Requirements

- Go 1.26+
- MongoDB for packages that use repository/db modules.
- Redis for optional cache primary mode.

## Quick Start Flow

For first-level setup:

1. Configure app settings with util config helpers.
2. Create a Mongo client from db/mongodb package.
3. Build repository/service layers using repository and generics packages.
4. Compose endpoints and route handlers with http/chi and generics router helpers.
5. Add auth middleware from http/oidc and auth packages where required.

Minimal startup sketch:

```go
package main

import (
	"context"

	libmongo "gofoundry/db/mongodb"
	"gofoundry/models"
	"gofoundry/util"
)

type AppConfig struct {
	MongoDb models.MongoConnection
}

func bootstrap(ctx *context.Context) (libmongo.IMongoClient, *AppConfig) {
	settings, err := util.Configure[AppConfig]("./config")
	util.PanicError(err)

	mc, err := libmongo.NewMongoClient(
		ctx,
		settings.MongoDb.ConnectionString,
		settings.MongoDb.DbName,
		nil,
		nil,
	)
	util.PanicError(err)

	return mc, settings
}
```

Detailed usage pages and example-driven guides are linked below.

## Documentation Map

- [Getting Started](docs/GETTING_STARTED.md)
- [Architecture](docs/ARCHITECTURE.md)
- [API Overview](docs/API_OVERVIEW.md)
- [HTTP Transport](docs/HTTP_TRANSPORT.md)
- [Advanced Examples](docs/ADVANCED_EXAMPLES.md)
- [Configuration](docs/CONFIGURATION.md)
- [Error Handling](docs/ERROR_HANDLING.md)
- [Code Standards](docs/CODE_STANDARDS.md)
- [Development Workflow](docs/DEVELOPMENT_WORKFLOW.md)
- [Contributing](CONTRIBUTING.md)

## Project Structure

- auth: JWT utilities and token validation.
- cache: cache client abstraction for Redis/Mongo.
- db and db/mongodb: entities and Mongo client/runtime helpers.
- dependency: generic dependency resolution utilities.
- dtos: request/response and query DTO definitions.
- generics: reusable CRUD, endpoint, and API router components.
- helpers: query, request, and validation utility helpers.
- http: chi middlewares, handlers, and transport encoding/decoding.
- repository: Mongo repository contracts and implementations.
- util: cross-cutting utility and configuration helpers.

## Notes

- Internal package imports use the module-qualified form gofoundry/<package>.
- For dependency synchronization, use the safe flow documented in [Development Workflow](docs/DEVELOPMENT_WORKFLOW.md).
