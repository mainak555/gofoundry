# Development Workflow

## Safe Dependency Sync

Avoid broad upgrade commands such as go get -u all in this repository.

Use this safe sequence:

```bash
go mod download
go mod tidy -v
go test ./...
```

This preserves declared dependency intent while regenerating go.sum safely.

## Local Validation

Run before pushing changes:

```bash
go test ./...
```

Optional targeted checks:

```bash
go list ./...
```

## Documentation Workflow

- Keep README as the high-level entry point.
- Add deep-dive package guidance under docs/.
- Add runnable examples in a follow-up pass once API narratives are stable.
