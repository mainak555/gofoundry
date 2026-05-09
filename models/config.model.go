package models

type OidcConfig struct {
	Url      string
	ClientId string
	Secret   string
}

type VaultConnection struct {
	Enabled bool
	Name    string
}

type AppInsight struct {
	InstrumentKey string
}

type MongoConnection struct {
	ConnectionString string
	PasswordVaultKey string
	DbName           string
	ReadSecondary    bool
	WriteMajority    bool
}
