# Error Handling

The library follows transport-layer error mapping for HTTP responses and helper-level error composition.

## Transport Error Mapping

http/server.EncodeErrorResponse maps known error prefixes to status codes:

- invalid -> 400
- unauthorize -> 401
- mongo: no documents in result -> 404
- no content err -> 404
- no content -> 204
- default -> 500

Use helper constructors to keep prefixes consistent:

```go
return helpers.InvalidError("missing name")
return helpers.NoContentError("entity not found")
return helpers.NoContent("nothing to return")
```

## Helper Error Constructors

helpers package provides helpers for consistent error prefixes:

- helpers.InvalidError
- helpers.NoContentError
- helpers.NoContent
- helpers.MergeErrors

## go-kit Server Error Wiring

Default option stack:

```go
serverOptions := server.KitHttpServerOptions(logger)
```

Equivalent explicit stack:

```go
serverOptions := []kithttp.ServerOption{
	kithttp.ServerErrorHandler(handlers.NewLogErrorHandler(logger)),
	kithttp.ServerErrorEncoder(server.EncodeErrorResponse),
	kithttp.ServerFinalizer(server.NewServerFinalizer(logger)),
}
```

## Custom Error Handler Hook

If you need telemetry, provide your own `kithttp.ServerErrorHandler`:

```go
errorHandler := handlers.ErrorHandlerFunc(func(ctx context.Context, err error) {
	level.Error(logger).Log("err", err)
	// send exception to your telemetry backend
})

serverOptions := []kithttp.ServerOption{
	kithttp.ServerErrorHandler(errorHandler),
	kithttp.ServerErrorEncoder(server.EncodeErrorResponse),
	kithttp.ServerFinalizer(server.NewServerFinalizer(logger)),
}
```

## Best Practice

- Return typed/structured errors where possible at service layer.
- Preserve root cause in wrapped errors.
- Avoid leaking sensitive backend details in HTTP messages.
- Keep service-layer error messages meaningful but prefix-compatible for transport mapping.
