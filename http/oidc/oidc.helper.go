package oidc

import (
	"context"
	"errors"
	"net/http"

	"github.com/mainak555/gofoundry/auth"

	"github.com/coreos/go-oidc"
)

func GetClaims[TClaims any](r *http.Request) (*TClaims, error) {
	ctx := r.Context()
	return GetClaimsFromCtx[TClaims](&ctx)
}

func GetClaimsFromCtx[TClaims any](ctx *context.Context) (*TClaims, error) {
	iJwt, err := GetJWTokenFromCtx[TClaims](ctx)
	if err != nil {
		return nil, err
	}
	return GetTokenClaims[TClaims](iJwt.GetIDToken())
}

func GetTokenClaims[TClaims any](idToken *oidc.IDToken) (*TClaims, error) {
	var idTokenClaims TClaims // ID Token payload is just JSON.
	if err := idToken.Claims(&idTokenClaims); err != nil {
		return nil, errors.New("invalid claims")
	}
	return &idTokenClaims, nil
}

func GetJWToken[TClaims any](r *http.Request) (auth.IJWT[TClaims], error) {
	ctx := r.Context()
	return GetJWTokenFromCtx[TClaims](&ctx)
}

func GetJWTokenFromCtx[TClaims any](ctx *context.Context) (auth.IJWT[TClaims], error) {
	iJwt := (*ctx).Value("IJWToken")
	if iJwt == nil {
		return nil, errors.New("invalid, empty token")
	}
	return iJwt.(auth.IJWT[TClaims]), nil
}
