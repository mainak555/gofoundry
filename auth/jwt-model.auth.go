package auth

import (
	"github.com/coreos/go-oidc"
)

type IJWT[TClaims any] interface {
	GetClient() *string
	GetIDToken() *oidc.IDToken
	Roles(clientId *string) ([]string, error)
	Scopes(clientId *string) ([]string, error)
}

type JWKS struct {
	Keys []struct {
		E   string   `json:"e"`
		N   string   `json:"n"`
		Kid string   `json:"kid"`
		Kty string   `json:"kty"`
		Alg string   `json:"alg"`
		Use string   `json:"use"`
		X5c []string `json:"x5c"`
	} `json:"keys"`
}

type OpenIDWellKnown struct {
	IssUrl   string `json:"issuer"`
	JwksUrl  string `json:"jwks_uri"`
	AuthUrl  string `json:"authorization_endpoint"`
	TokenUrl string `json:"token_endpoint"`
}
