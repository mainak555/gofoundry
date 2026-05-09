# Advanced Examples

This guide shows production-oriented composition patterns using concrete GoFoundry APIs.

## 1. Production Bootstrap Pattern

```go
var logger log.Logger
var appCtx context.Context
var cancel context.CancelFunc
var mongoClient mongodb.IMongoClient

type AppConfig struct {
	AllowOrigins []string
	MongoDb      models.MongoConnection
	Oidc         *models.OidcConfig
	Vault        models.VaultConnection
	AppInsight   models.AppInsight
}

var settings *AppConfig

func init() {
	appCtx, cancel = context.WithTimeout(context.Background(), 20*time.Second)

	settings, err := util.Configure[AppConfig]("./config")
	util.PanicError(err)

	if settings.Vault.Enabled {
		client, err := util.GetVault(settings.Vault.Name)
		util.PanicError(err)

		password, err := util.GetFromVault(client, settings.MongoDb.PasswordVaultKey)
		util.PanicError(err)

		settings.MongoDb.ConnectionString = fmt.Sprintf(settings.MongoDb.ConnectionString, password)
	}

	mongoClient, err = mongodb.NewMongoClient(&appCtx, settings.MongoDb.ConnectionString, settings.MongoDb.DbName, nil, nil)
	util.PanicError(err)
}
```

## 2. Configure Collections and Indexes

`ConfigureCollections` is the recommended startup checkpoint for registering data models and indexes.

```go
type MyEntity struct{}

err := mongoClient.ConfigureCollections(&appCtx,
	func() (string, interface{}, []mongo.IndexModel) {
		return "entity_collection_name", MyEntity{}, []mongo.IndexModel{
			{
				Keys: bson.D{{Key: "name", Value: 1}},
				Options: options.Index().SetName("name_1").SetUnique(true).
					SetCollation(&options.Collation{Locale: "en", Strength: 2}),
			},
			{
				Keys:    bson.D{{Key: "geoLoc", Value: mongodb.GEOSPATIAL_INDEX_2DSPHERE}},
				Options: options.Index().SetName("geoLoc_2dSphere_1"),
			},
		}
	},
	func() (string, interface{}, []mongo.IndexModel) {
		return "cache_collection_name", nil, []mongo.IndexModel{
			{Keys: bson.D{{Key: "type", Value: 1}, {Key: "createdAt", Value: 1}}, Options: options.Index().SetName("type_1_createdAt_1")},
		}
	},
)
util.PanicError(err)
```

## 3. Middleware and OIDC Stack

```go
type Claims struct {
	Sub string `json:"sub"`
}

var jwtFactory func(idToken *oidc.IDToken, clientId *string) auth.IJWT[Claims]

router := chi.NewRouter()

middlewares.AddMiddlewares(
	router,
	settings.AllowOrigins,
	middleware.Timeout(2*time.Minute),
)

router.With(
	oidc.ValidateToken(settings.Oidc, jwtFactory),
	oidc.IsRoleAny[Claims]("ROLE-A", "ROLE-B"),
).Route("/v1", func(r chi.Router) {
	// protected routes
})
```

`jwtFactory` is an app-defined adapter that returns `auth.IJWT[TClaims]`.

Order guidance:

1. Global middlewares first.
2. Authentication middleware before role checks.
3. Route/resource validation middleware near the endpoint group where context is needed.

## 4. Route Id Validation With Context Enrichment

```go
r.With(chiHandlers.ValidateRouteIdCtx[MyEntity](
	mongoClient,
	"entity",
	func(ctx *context.Context, query bson.M) bson.M {
		filter := helpers.GetFilterQuery(ctx)
		if id, ok := filter["entityId"]; ok {
			query["$and"] = append(query["$and"].([]bson.M), bson.M{"entityId": id})
		}
		return query
	},
	func(req *http.Request, id primitive.ObjectID) context.Context {
		return context.WithValue(req.Context(), "entityId", id)
	},
)).Route("/{id}", func(child chi.Router) {
	// nested handlers
})
```

## 5. Transport Error Handling + Telemetry Hook

```go
errorHandler := handlers.ErrorHandlerFunc(func(ctx context.Context, err error) {
	level.Error(logger).Log("err", err)
	// emit to telemetry backend
})

serverOptions := []kithttp.ServerOption{
	kithttp.ServerErrorHandler(errorHandler),
	kithttp.ServerErrorEncoder(server.EncodeErrorResponse),
	kithttp.ServerFinalizer(server.NewServerFinalizer(logger)),
}
```

## 6. Cache Bootstrap

```go
cache.NewMongoCache("secondary_cache", mongoClient.GetDatabase(), 60)
```

## 7. Graceful Shutdown

```go
srv := &http.Server{Addr: ":80", Handler: router}
errChan := make(chan error, 1)

go func() { errChan <- srv.ListenAndServe() }()

go func() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	errChan <- fmt.Errorf("os signal: %s | server shutdown", <-c)
	_ = srv.Shutdown(appCtx)
}()

if err := <-errChan; err != nil {
	level.Error(logger).Log("shutdown", err)
}

_ = mongoClient.Disconnect()
cancel()
```

## See Also

- `GETTING_STARTED.md` for minimal setup.
- `HTTP_TRANSPORT.md` for route and server transport details.
- `ERROR_HANDLING.md` for status mapping behavior.
