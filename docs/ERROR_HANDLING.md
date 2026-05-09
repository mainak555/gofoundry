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

## Helper Error Constructors

helpers package provides helpers for consistent error prefixes:

- helpers.InvalidError
- helpers.NoContentError
- helpers.NoContent
- helpers.MergeErrors

## Best Practice

- Return typed/structured errors where possible at service layer.
- Preserve root cause in wrapped errors.
- Avoid leaking sensitive backend details in HTTP messages.
