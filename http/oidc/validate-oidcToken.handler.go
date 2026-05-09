package oidc

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"auth"
	"http/chi/renderers"
	"models"
	"github.com/coreos/go-oidc"
	"github.com/go-chi/render"
)

func GetAuthToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		return strings.Split(authHeader, " ")[1]
	}
	return ""
}

func ValidateToken[TClaims any](config *models.OidcConfig, jwFn func(idToken *oidc.IDToken, clientId *string) auth.IJWT[TClaims]) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if os.Getenv("AUTH") != "" {
				if value, err := strconv.ParseBool(os.Getenv("AUTH")); err == nil && !value {
					next.ServeHTTP(w, r)
					return
				}
			}

			accessToken := GetAuthToken(r)
			if accessToken == "" {
				RenderError("Auth Token Missing", "", w, r)
				return
			}

			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}

			client := &http.Client{
				Timeout:   time.Duration(6000) * time.Second,
				Transport: tr,
			}

			ctx := r.Context()
			oidcCtx := oidc.ClientContext(ctx, client)
			provider, err := oidc.NewProvider(oidcCtx, config.Url)

			if err != nil {
				RenderError("authorisation failed while getting the provider", err.Error(), w, r)
				return
			}

			oidcConfig := &oidc.Config{
				ClientID: config.ClientId,
			}

			verifier := provider.Verifier(oidcConfig)
			idToken, err := verifier.Verify(ctx, accessToken)
			if err != nil {
				RenderError("authorisation failed while verifying the token", err.Error(), w, r)
				return
			}

			ctx2 := context.WithValue(ctx, "IJWToken", jwFn(idToken, &config.ClientId))
			next.ServeHTTP(w, r.WithContext(ctx2))
		})
	}
}

func RenderError(message, errTxt string, w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, &renderers.ErrResponse{
		ReasonPhrase:   message,
		HTTPStatusCode: http.StatusUnauthorized,
		ErrorText:      errTxt,
	})
}
