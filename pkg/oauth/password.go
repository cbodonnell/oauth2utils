package oauth

import (
	"context"

	"golang.org/x/oauth2"
)

func (oc *OIDCClient) Password(ctx context.Context, username, password string) (*oauth2.Token, error) {
	return oc.conf.PasswordCredentialsToken(ctx, username, password)
}
