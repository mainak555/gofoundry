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

Example:

```go
type MyEntity struct{}

type AppConfig struct {
	AllowOrigins []string
	MongoDb      models.MongoConnection
	Oidc         *models.OidcConfig
	Vault        models.VaultConnection
	AppInsight   models.AppInsight
}

settings, err := util.Configure[AppConfig]("./config")
util.PanicError(err)
```

## 3. Initialize Mongo Client

Use db/mongodb package to construct and connect IMongoClient with your connection string and database name.

Example with client options:

```go
ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
defer cancel()

mongoClient, err := mongodb.NewMongoClient(
	&ctx,
	settings.MongoDb.ConnectionString,
	settings.MongoDb.DbName,
	func(co *options.ClientOptions) *options.ClientOptions {
		dialer := &net.Dialer{Timeout: 30 * time.Second, KeepAlive: 300 * time.Second}
		co.SetDialer(dialer)
		co.SetRetryReads(true)
		co.SetRetryWrites(true)
		return co
	},
	nil,
)
util.PanicError(err)
```

Register collections and indexes at startup:

```go
err = mongoClient.ConfigureCollections(&ctx,
	func() (string, interface{}, []mongo.IndexModel) {
		return "mongo_collection_name", MyEntity{}, []mongo.IndexModel{
			{Keys: bson.D{{Key: "name", Value: 1}}, Options: options.Index().SetName("name_1").SetUnique(true)},
			{Keys: bson.D{{Key: "geoLoc", Value: mongodb.GEOSPATIAL_INDEX_2DSPHERE}}, Options: options.Index().SetName("geoLoc_2dSphere_1")},
		}
	},
)
util.PanicError(err)
```

## 4. Build Repository and Service Layers

- repository.NewMongoTRepository[T] for typed repository access.
- generics.NewCommonService[T] for common CRUD operations.

Example:

```go
repo := repository.NewMongoTRepository[MyEntity](mongoClient)
service := generics.NewCommonService[MyEntity](mongoClient)

_ = repo
_ = service
```

## 5. Wire HTTP Transport

- Use http/server helpers for go-kit decode/encode flows.
- Use http/chi handlers and middlewares for route-level concerns.
- Use generics.ApiRouter for reusable REST route assembly.

Example:

```go
apiRouter := chi.NewRouter()
middlewares.AddMiddlewares(apiRouter, settings.AllowOrigins, middleware.Timeout(2*time.Minute))

serverOptions := []kithttp.ServerOption{
	kithttp.ServerErrorHandler(handlers.NewLogErrorHandler(logger)),
	kithttp.ServerErrorEncoder(server.EncodeErrorResponse),
	kithttp.ServerFinalizer(server.NewServerFinalizer(logger)),
}

_ = serverOptions
```

You can also use the packaged defaults:

```go
serverOptions := server.KitHttpServerOptions(logger)
```

## 6. Add Authentication

For protected APIs:

- Use http/oidc middleware for token validation.
- Use auth package utilities for JWT parsing and signature verification when needed.

Example:

```go
type Claims struct {
	Sub string `json:"sub"`
}

var jwtFactory func(idToken *oidc.IDToken, clientId *string) auth.IJWT[Claims]

apiRouter.With(
	oidc.ValidateToken(settings.Oidc, jwtFactory),
	oidc.IsRoleAny[Claims]("ROLE-A", "ROLE-B"),
).Route("/v1", func(r chi.Router) {
	// mount handlers
})
```

`jwtFactory` must return your application's implementation of `auth.IJWT[TClaims]`.

Middleware order matters: validate token first, then apply role checks.

## 7. Start Server and Gracefully Shutdown

```go
srv := &http.Server{Addr: ":80", Handler: apiRouter}

errChan := make(chan error, 1)
go func() { errChan <- srv.ListenAndServe() }()

go func() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	errChan <- fmt.Errorf("os signal: %s | server shutdown", <-c)
	_ = srv.Shutdown(ctx)
}()

if err := <-errChan; err != nil {
	level.Error(logger).Log("shutdown", err)
}
```

## 8. Validate Build

```bash
go test ./...
```

## Next

- See HTTP composition details in `HTTP_TRANSPORT.md`.
- See production patterns in `ADVANCED_EXAMPLES.md`.
