package oauth

import (
	"context"
	"log"

	"github.com/cbodonnell/oauth2utils/pkg/term"
	"golang.org/x/oauth2"
)

func (oc *OIDCClient) Password(ctx context.Context) (*oauth2.Token, error) {
	username := term.StringPrompt("Username:")
	password := term.PasswordPrompt("Password:")

	token, err := oc.conf.PasswordCredentialsToken(ctx, username, password)
	if err != nil {
		log.Fatal(err)
	}

	return token, nil
}
