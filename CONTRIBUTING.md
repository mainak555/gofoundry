# Contributing

Thank you for contributing to GoFoundry.

## Development Setup

1. Install Go 1.26+.
2. Clone the repository.
3. Sync dependencies:

```bash
go mod download
go mod tidy -v
```

4. Validate:

```bash
go test ./...
```

## Pull Request Guidelines

- Keep pull requests focused and scoped.
- Include rationale and impact in PR description.
- Update docs when behavior or contracts change.
- Do not introduce unrelated formatting churn.

## Code and Comment Standards

- Follow GoDoc comments for exported symbols.
- Add intent comments for critical internal helper logic.
- Preserve module-qualified local imports: gofoundry/<package>.

## Dependency Management

- Avoid go get -u all.
- Prefer safe sync flow documented in docs/DEVELOPMENT_WORKFLOW.md.
