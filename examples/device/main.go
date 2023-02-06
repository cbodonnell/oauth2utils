package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/cbodonnell/oauth2utils/pkg/oauth"
	"github.com/cbodonnell/oauth2utils/pkg/persistence"
	"github.com/cbodonnell/oauth2utils/pkg/utils"
)

const (
	KeycloakServerUrl = "http://localhost:8080"
	RealmName         = "tunnel-farm"
	ClientId          = "myclient"
)

func main() {
	ctx := context.Background()
	token := utils.TryGetToken(ctx)
	if !token.Valid() {
		newToken, err := oauth.DeviceCode(ctx)
		if err != nil {
			log.Fatal(err)
		}
		token = newToken
	}

	if err := persistence.SaveToken(token); err != nil {
		log.Fatal(err)
	}

	client := oauth.Client(ctx, token)
	res, err := client.Get(oauth.UserInfoURL())
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Fatal("failed to get userinfo")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(string(body))
}
