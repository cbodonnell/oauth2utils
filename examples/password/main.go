package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/cbodonnell/oauth2utils/pkg/oauth"
	"github.com/cbodonnell/oauth2utils/pkg/persistence"
)

func main() {
	ctx := context.Background()

	token, err := persistence.LoadToken()
	if err != nil || !token.Valid() {
		token, err = oauth.Password(ctx)
		if err != nil {
			log.Fatal(err)
		}
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
