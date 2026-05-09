package models

// OidcConfig holds OIDC provider and client credentials.
type OidcConfig struct {
	Url      string
	ClientId string
	Secret   string
}

// VaultConnection defines external secret-vault connection settings.
type VaultConnection struct {
	Enabled bool
	Name    string
}

// AppInsight contains telemetry integration settings.
type AppInsight struct {
	InstrumentKey string
}

// MongoConnection defines MongoDB connectivity and read/write preferences.
type MongoConnection struct {
	ConnectionString string
	PasswordVaultKey string
	DbName           string
	ReadSecondary    bool
	WriteMajority    bool
}
