package oauth

import (
	"context"
	"fmt"
	"log"

	"github.com/cbodonnell/oauth2utils/pkg/term"
	"golang.org/x/oauth2"
)

func Password(ctx context.Context) (*oauth2.Token, error) {
	username := term.StringPrompt("Username:")
	password := term.PasswordPrompt("Password:")

	conf := &oauth2.Config{
		ClientID:     ClientId,
		ClientSecret: "",
		Endpoint: oauth2.Endpoint{
			TokenURL: fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", KeycloakServerUrl, RealmName),
		},
		Scopes: []string{"openid", "profile", "email"},
	}

	token, err := conf.PasswordCredentialsToken(ctx, username, password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Success! You are now logged in.")

	return token, nil
}