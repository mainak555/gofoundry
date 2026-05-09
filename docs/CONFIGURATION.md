# Configuration

GoFoundry uses model-based configuration with environment-aware loading helpers.

## Environment Resolution

util.Configure[T] uses NODE_ENV to select config files:

- local.yml for NODE_ENV=local
- local.env for environment variable hydration

If NODE_ENV is missing, local is used by default.

## Core Config Models

Defined in models/config.model.go:

- OidcConfig: OIDC issuer URL, client id, and secret.
- VaultConnection: external vault enablement and name.
- AppInsight: telemetry key.
- MongoConnection: connection string, database, and read/write behavior options.

## Recommended Practice

- Keep secrets in environment variables or vault providers.
- Keep non-sensitive defaults in environment-specific yml files.
- Validate required fields during service startup.
