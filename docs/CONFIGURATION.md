# Configuration

GoFoundry uses model-based configuration with environment-aware loading helpers.

## Environment Resolution

util.Configure[T] uses NODE_ENV to select config files:

- local.yml for NODE_ENV=local
- local.env for environment variable hydration

If NODE_ENV is missing, local is used by default.

Example:

```go
type AppConfig struct {
	MongoDb    models.MongoConnection
	Oidc       *models.OidcConfig
	Vault      models.VaultConnection
	AppInsight models.AppInsight
}

settings, err := util.Configure[AppConfig]("./config")
util.PanicError(err)
```

## Core Config Models

Defined in models/config.model.go:

- OidcConfig: OIDC issuer URL, client id, and secret.
- VaultConnection: external vault enablement and name.
- AppInsight: telemetry key.
- MongoConnection: connection string, database, and read/write behavior options.

Typical app-level config model:

```go
type AppConfig struct {
	AllowOrigins []string
	MongoDb      models.MongoConnection
	Oidc         *models.OidcConfig
	Vault        models.VaultConnection
	AppInsight   models.AppInsight
}
```

## Vault Secret Injection

When `VaultConnection.Enabled` is true, you can hydrate secrets at startup before creating clients.

```go
settings, err := util.Configure[AppConfig]("./config")
util.PanicError(err)

if settings.Vault.Enabled {
	client, err := util.GetVault(settings.Vault.Name)
	if err != nil {
		return err
	}

	password, err := util.GetFromVault(client, settings.MongoDb.PasswordVaultKey)
	if err != nil {
		return err
	}

	settings.MongoDb.ConnectionString = fmt.Sprintf(settings.MongoDb.ConnectionString, password)
}
```

Notes:

- `util.GetAzCredential` uses client secret credentials for `NODE_ENV=local` and default Azure credentials otherwise.
- Keep only template strings in yml files and inject live secrets at runtime.

## Recommended Practice

- Keep secrets in environment variables or vault providers.
- Keep non-sensitive defaults in environment-specific yml files.
- Validate required fields during service startup.
- Keep startup config loading and secret hydration in one bootstrap function.
