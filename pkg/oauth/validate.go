package oauth

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/square/go-jose"
)

func GetJWKS(ctx context.Context) (*jose.JSONWebKeySet, error) {
	// create an http request with the context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, JWKSUrl(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Retrieve the JWKS from the Keycloak server
	jwks, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve JWKS: %w", err)
	}
	defer jwks.Body.Close()

	// Read the JWKS response
	jwksData, err := io.ReadAll(jwks.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read JWKS response: %w", err)
	}

	// Unmarshal the JWKSData as JSONWebKeySet
	keys := &jose.JSONWebKeySet{}
	if err := json.Unmarshal(jwksData, keys); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JWKS: %w", err)
	}

	return keys, nil
}

func Validate(keys *jose.JSONWebKeySet, tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		header := token.Header
		kid, ok := header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("unable to parse kid")
		}

		for _, key := range keys.Keys {
			if key.KeyID == kid {
				if rsaKey, ok := key.Key.(*rsa.PublicKey); ok {
					return rsaKey, nil
				}
			}
		}

		return nil, fmt.Errorf("unable to find key")
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		issuer := claims["iss"]
		audience := claims["azp"]
		expirationTime := int64(claims["exp"].(float64))
		currentTime := time.Now().Unix()

		// Validate the issuer
		if issuer != fmt.Sprintf("%s/realms/%s", KeycloakServerUrl, RealmName) {
			return nil, fmt.Errorf("invalid iss")
		}

		// Validate the authorized party
		if audience != ClientId {
			return nil, fmt.Errorf("invalid azp")
		}

		// Validate the expiration time
		if currentTime >= expirationTime {
			return nil, fmt.Errorf("token expired")
		}

		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
