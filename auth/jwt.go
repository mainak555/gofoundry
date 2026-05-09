package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// GetJWT parses tokenStr without signature verification and returns the token object.
func GetJWT(tokenStr string) (*jwt.Token, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
	return token, err
}

// GetClaims returns map claims from tokenStr without signature verification.
func GetClaims(tokenStr string) (*jwt.MapClaims, error) {
	if token, err := GetJWT(tokenStr); err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(jwt.MapClaims); !ok {
		return nil, errors.New("invalid claims")
	} else {
		return &claims, nil
	}
}

// ValidateJWTIssuer validates issuer claims against issuerUrl and returns token claims.
func ValidateJWTIssuer(tokenStr, issuerUrl string, audiences ...string) (*jwt.MapClaims, error) {
	if token, err := GetJWT(tokenStr); err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(jwt.MapClaims); !ok {
		return nil, errors.New("invalid claims")
	} else if iss, ok := claims["iss"].(string); !ok || !strings.Contains(iss, issuerUrl) {
		return nil, errors.New("issuer not matched")
	} else {
		return &claims, nil
	}
}

// ValidateJWTSignature validates tokenStr against a JWKS endpoint and returns claims.
func ValidateJWTSignature(tokenStr, jwksUrl string) (*jwt.MapClaims, error) {
	if token, err := GetJWT(tokenStr); err != nil {
		return nil, err
	} else if pubKey, err := func(token *jwt.Token) (*rsa.PublicKey, error) {
		var jwks JWKS
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("missing kid in token header")
		}

		if resp, err := http.Get(jwksUrl); err != nil {
			return nil, err
		} else if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
			return nil, err
		} else {
			resp.Body.Close()
		}

		parseBase64 := func(input string) ([]byte, error) {
			decoded := make([]byte, base64.RawURLEncoding.DecodedLen(len(input)))
			n, err := base64.RawURLEncoding.Decode(decoded, []byte(input))
			if err != nil {
				return nil, err
			}
			return decoded[:n], nil
		}

		for _, key := range jwks.Keys {
			if key.Kid == kid {
				n, err := parseBase64(key.N)
				if err != nil {
					return nil, err
				}
				e, err := parseBase64(key.E)
				if err != nil {
					return nil, err
				}
				return &rsa.PublicKey{
					N: new(big.Int).SetBytes(n),
					E: int(new(big.Int).SetBytes(e).Int64()),
				}, nil
			}
		}
		return nil, errors.New("no public key found")
	}(token); err != nil {
		return nil, err
	} else if token, err = jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); ok {
			return pubKey, nil
		}
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}); err != nil {
		return nil, err
	} else if !token.Valid {
		return nil, errors.New("invalid token")
	} else {
		c := token.Claims.(jwt.MapClaims)
		return &c, nil
	}
}
