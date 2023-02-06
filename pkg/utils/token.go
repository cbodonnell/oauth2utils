package utils

import (
	"context"

	"github.com/cbodonnell/oauth2utils/pkg/oauth"
	"github.com/cbodonnell/oauth2utils/pkg/persistence"
	"golang.org/x/oauth2"
)

// TryGetToken attempts to load a token from the persistence layer. If the token
// is expired, it will attempt to refresh it. The function returns nil if the
// token is invalid or the refresh fails.
func TryGetToken(ctx context.Context) *oauth2.Token {
	token, err := persistence.LoadToken()
	if err == nil {
		token, _ = oauth.RefreshToken(ctx, token)
	}
	return token
}
