package oauth

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

// TODO: get these from env vars and defaults
const (
	// TODO: make this less keycloak specific
	KeycloakServerUrl = "http://localhost:8080"
	RealmName         = "tunnel-farm"
	ClientId          = "myclient"
)

var conf *oauth2.Config
var deviceCodeURI string

func init() {
	conf = &oauth2.Config{
		ClientID:     ClientId,
		ClientSecret: "", // TODO: Can we remove?
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s/realms/%s/protocol/openid-connect/auth", KeycloakServerUrl, RealmName),
			TokenURL: fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", KeycloakServerUrl, RealmName),
		},
		Scopes: []string{"openid", "profile", "email"},
	}
	deviceCodeURI = fmt.Sprintf("%s/realms/%s/protocol/openid-connect/auth/device", KeycloakServerUrl, RealmName)
}

func Client(ctx context.Context, token *oauth2.Token) *http.Client {
	return conf.Client(ctx, token)
}

func UserInfoURL() string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/userinfo", KeycloakServerUrl, RealmName)
}
