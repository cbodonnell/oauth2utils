package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cbodonnell/oauth2utils/pkg/oauth"
	"github.com/cbodonnell/oauth2utils/pkg/persistence"
	"github.com/cbodonnell/oauth2utils/pkg/term"
	"github.com/cbodonnell/oauth2utils/pkg/utils"
)

const TokenDir = ".oauth"

func main() {
	ctx := context.Background()

	oc, err := oauth.NewOIDCClient(ctx, "http://localhost:8080/realms/tunnel.farm", "tfarm-cli", []string{"profile"})
	if err != nil {
		log.Fatal(err)
	}

	token := utils.TryGetToken(ctx, oc, TokenDir)
	if !token.Valid() {
		username := term.StringPrompt("Username:")
		password := term.PasswordPrompt("Password:")
		newToken, err := oc.Password(ctx, username, password)
		if err != nil {
			log.Fatal(err)
		}
		token = newToken
	}

	// TODO: dont need to save the token if it hasn't changed
	if err := persistence.SaveToken(token, TokenDir); err != nil {
		log.Fatal(err)
	}

	userInfo, err := oc.UserInfo(ctx, oc.TokenSource(ctx, token))
	if err != nil {
		log.Fatal(err)
	}

	claims := map[string]interface{}{}
	if err := userInfo.Claims(&claims); err != nil {
		log.Fatal(err)
	}

	claimsJson, err := json.Marshal(claims)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(string(claimsJson))
}
