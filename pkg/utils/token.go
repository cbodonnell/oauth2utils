package utils

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cbodonnell/oauth2utils/pkg/oauth"
	"github.com/cbodonnell/oauth2utils/pkg/persistence"
	"golang.org/x/oauth2"
)

// TryGetToken attempts to load a token from the persistence layer. If the token
// is expired, it will attempt to refresh it. The function returns nil if the
// token is invalid or the refresh fails.
func TryGetToken(ctx context.Context, oc *oauth.OIDCClient) *oauth2.Token {
	token, err := persistence.LoadToken()
	if err == nil {
		token, _ = oc.TokenSource(ctx, token).Token()
	}
	return token
}

func ParseBearerToken(r *http.Request) (string, error) {
	// get the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("no auth header")
	}

	// split the header into its parts
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid auth header, wrong number of parts")
	}

	// make sure the header is a bearer token
	if parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid auth header, not a bearer token")
	}

	return parts[1], nil
}
