# Getting Started

This guide provides the first-level setup path for integrating GoFoundry into a service.

## 1. Add Dependency

```bash
go get github.com/mainak555/gofoundry
```

## 2. Configure Environment

Set NODE_ENV and create an environment-specific configuration file.

- NODE_ENV=local loads local.yml and local.env by convention.
- util.Configure[T] reads config from the provided path and unmarshals into your model.

## 3. Initialize Mongo Client

Use db/mongodb package to construct and connect IMongoClient with your connection string and database name.

## 4. Build Repository and Service Layers

- repository.NewMongoTRepository[T] for typed repository access.
- generics.NewCommonService[T] for common CRUD operations.

## 5. Wire HTTP Transport

- Use http/server helpers for go-kit decode/encode flows.
- Use http/chi handlers and middlewares for route-level concerns.
- Use generics.ApiRouter for reusable REST route assembly.

## 6. Add Authentication

For protected APIs:

- Use http/oidc middleware for token validation.
- Use auth package utilities for JWT parsing and signature verification when needed.

## 7. Validate Build

```bash
go test ./...
```

## Next

Example applications and runnable snippets will be added in the next documentation pass.
