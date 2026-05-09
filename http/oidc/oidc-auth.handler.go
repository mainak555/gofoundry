package oidc

import (
	"net/http"
	"os"
	"strconv"

	"golang.org/x/exp/slices"
)

func IsRoleAny[TClaims any](roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if os.Getenv("AUTH") != "" {
				if value, err := strconv.ParseBool(os.Getenv("AUTH")); err == nil && !value {
					next.ServeHTTP(w, r)
					return
				}
			}

			iJwt, err := GetJWToken[TClaims](r)
			if err != nil {
				RenderError("unauthorize", err.Error(), w, r)
				return
			}

			jwtRoles, _ := iJwt.Roles(iJwt.GetClient())
			if len(jwtRoles) > 0 && len(roles) > 0 {
				for _, b := range roles {
					if slices.Contains(jwtRoles, b) {
						next.ServeHTTP(w, r)
						return
					}
				}
			} else if len(jwtRoles) < 1 && len(roles) < 1 {
				next.ServeHTTP(w, r)
				return
			}
			RenderError("unauthorized", "role missing", w, r)

		})
	}
}

func IsClientRoleAny[TClaims any](clientId string, roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if os.Getenv("AUTH") != "" {
				if value, err := strconv.ParseBool(os.Getenv("AUTH")); err == nil && !value {
					next.ServeHTTP(w, r)
					return
				}
			}

			iJwt, err := GetJWToken[TClaims](r)
			if err != nil {
				RenderError("unauthorize", err.Error(), w, r)
				return
			}

			jwtRoles, _ := iJwt.Roles(&clientId)
			if len(jwtRoles) > 0 && len(roles) > 0 {
				for _, b := range roles {
					if slices.Contains(jwtRoles, b) {
						next.ServeHTTP(w, r)
						return
					}
				}
			} else if len(jwtRoles) < 1 && len(roles) < 1 {
				next.ServeHTTP(w, r)
				return
			}
			RenderError("unauthorized", "role missing", w, r)
		})
	}
}

func IsScopeAny[TClaims any](scopes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if os.Getenv("AUTH") != "" {
				if value, err := strconv.ParseBool(os.Getenv("AUTH")); err == nil && !value {
					next.ServeHTTP(w, r)
					return
				}
			}

			iJwt, err := GetJWToken[TClaims](r)
			if err != nil {
				RenderError("unauthorize", err.Error(), w, r)
				return
			}

			jwtScopes, _ := iJwt.Scopes(iJwt.GetClient())
			if len(jwtScopes) > 0 && len(scopes) > 0 {
				for _, b := range scopes {
					if slices.Contains(jwtScopes, b) {
						next.ServeHTTP(w, r)
						return
					}
				}
			} else if len(jwtScopes) < 1 && len(scopes) < 1 {
				next.ServeHTTP(w, r)
				return
			}
			RenderError("unauthorized", "scope missing", w, r)
		})
	}
}

func IsClientScopeAny[TClaims any](clientId string, scopes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if os.Getenv("AUTH") != "" {
				if value, err := strconv.ParseBool(os.Getenv("AUTH")); err == nil && !value {
					next.ServeHTTP(w, r)
					return
				}
			}

			iJwt, err := GetJWToken[TClaims](r)
			if err != nil {
				RenderError("unauthorize", err.Error(), w, r)
				return
			}

			jwtScopes, _ := iJwt.Scopes(&clientId)
			if len(jwtScopes) > 0 && len(scopes) > 0 {
				for _, b := range scopes {
					if slices.Contains(jwtScopes, b) {
						next.ServeHTTP(w, r)
						return
					}
				}
			} else if len(jwtScopes) < 1 && len(scopes) < 1 {
				next.ServeHTTP(w, r)
				return
			}
			RenderError("unauthorized", "scope missing", w, r)
		})
	}
}

func IsRolesAll[TClaims any](roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if os.Getenv("AUTH") != "" {
				if value, err := strconv.ParseBool(os.Getenv("AUTH")); err == nil && !value {
					next.ServeHTTP(w, r)
					return
				}
			}

			iJwt, err := GetJWToken[TClaims](r)
			if err != nil {
				RenderError("unauthorize", err.Error(), w, r)
				return
			}

			jwtRoles, _ := iJwt.Roles(iJwt.GetClient())
			if len(jwtRoles) > 0 && len(roles) > 0 {
				for _, b := range roles {
					if !slices.Contains(jwtRoles, b) {
						RenderError("unauthorized", "role missing", w, r)
						return
					}
				}
			} else if (len(jwtRoles) < 1 && len(roles) > 0) || (len(jwtRoles) > 0 && len(roles) < 1) {
				RenderError("unauthorized", "role missing", w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func IsClientRolesAll[TClaims any](clientId string, roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if os.Getenv("AUTH") != "" {
				if value, err := strconv.ParseBool(os.Getenv("AUTH")); err == nil && !value {
					next.ServeHTTP(w, r)
					return
				}
			}

			iJwt, err := GetJWToken[TClaims](r)
			if err != nil {
				RenderError("unauthorize", err.Error(), w, r)
				return
			}

			jwtRoles, _ := iJwt.Roles(&clientId)
			if len(jwtRoles) > 0 && len(roles) > 0 {
				for _, b := range roles {
					if !slices.Contains(jwtRoles, b) {
						RenderError("unauthorized", "role missing", w, r)
						return
					}
				}
			} else if (len(jwtRoles) < 1 && len(roles) > 0) || (len(jwtRoles) > 0 && len(roles) < 1) {
				RenderError("unauthorized", "role missing", w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func IsScopesAll[TClaims any](scopes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if os.Getenv("AUTH") != "" {
				if value, err := strconv.ParseBool(os.Getenv("AUTH")); err == nil && !value {
					next.ServeHTTP(w, r)
					return
				}
			}

			iJwt, err := GetJWToken[TClaims](r)
			if err != nil {
				RenderError("unauthorize", err.Error(), w, r)
				return
			}

			jwtScopes, _ := iJwt.Scopes(iJwt.GetClient())
			if len(jwtScopes) > 0 && len(scopes) > 0 {
				for _, b := range scopes {
					if !slices.Contains(jwtScopes, b) {
						RenderError("unauthorized", "scope missing", w, r)
						return
					}
				}
			} else if (len(jwtScopes) < 1 && len(scopes) > 0) || (len(jwtScopes) > 0 && len(scopes) < 1) {
				RenderError("unauthorized", "scope missing", w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func IsClientScopesAll[TClaims any](clientId string, scopes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if os.Getenv("AUTH") != "" {
				if value, err := strconv.ParseBool(os.Getenv("AUTH")); err == nil && !value {
					next.ServeHTTP(w, r)
					return
				}
			}

			iJwt, err := GetJWToken[TClaims](r)
			if err != nil {
				RenderError("unauthorize", err.Error(), w, r)
				return
			}

			jwtScopes, _ := iJwt.Scopes(&clientId)
			if len(jwtScopes) > 0 && len(scopes) > 0 {
				for _, b := range scopes {
					if !slices.Contains(jwtScopes, b) {
						RenderError("unauthorized", "scope missing", w, r)
						return
					}
				}
			} else if (len(jwtScopes) < 1 && len(scopes) > 0) || (len(jwtScopes) > 0 && len(scopes) < 1) {
				RenderError("unauthorized", "scope missing", w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func IsScopesOrRoles[TClaims any](scopes, roles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if os.Getenv("AUTH") != "" {
				if value, err := strconv.ParseBool(os.Getenv("AUTH")); err == nil && !value {
					next.ServeHTTP(w, r)
					return
				}
			}

			iJwt, err := GetJWToken[TClaims](r)
			if err != nil {
				RenderError("unauthorize", err.Error(), w, r)
				return
			}

			scopeFn := func() bool {
				jwtScopes, _ := iJwt.Scopes(iJwt.GetClient())
				if len(jwtScopes) > 0 && len(scopes) > 0 {
					for _, b := range scopes {
						if !slices.Contains(jwtScopes, b) {
							return false
						}
					}
					return true
				}
				return false
			}

			roleFn := func() bool {
				jwtRoles, _ := iJwt.Roles(iJwt.GetClient())
				if len(jwtRoles) > 0 && len(roles) > 0 {
					for _, b := range roles {
						if !slices.Contains(jwtRoles, b) {
							return false
						}
					}
					return true
				}
				return false
			}

			if scopeFn() || roleFn() {
				next.ServeHTTP(w, r)
				return
			}
			RenderError("unauthorized", "scope, role missing", w, r)
		})
	}
}

func IsClientScopesOrRoles[TClaims any](clientId string, scopes, roles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if os.Getenv("AUTH") != "" {
				if value, err := strconv.ParseBool(os.Getenv("AUTH")); err == nil && !value {
					next.ServeHTTP(w, r)
					return
				}
			}

			iJwt, err := GetJWToken[TClaims](r)
			if err != nil {
				RenderError("unauthorize", err.Error(), w, r)
				return
			}

			scopeFn := func() bool {
				jwtScopes, _ := iJwt.Scopes(&clientId)
				if len(jwtScopes) > 0 && len(scopes) > 0 {
					for _, b := range scopes {
						if !slices.Contains(jwtScopes, b) {
							return false
						}
					}
					return true
				}
				return false
			}

			roleFn := func() bool {
				jwtRoles, _ := iJwt.Roles(&clientId)
				if len(jwtRoles) > 0 && len(roles) > 0 {
					for _, b := range roles {
						if !slices.Contains(jwtRoles, b) {
							return false
						}
					}
					return true
				}
				return false
			}

			if scopeFn() || roleFn() {
				next.ServeHTTP(w, r)
				return
			}
			RenderError("unauthorized", "scope, roles missing", w, r)
		})
	}
}
