package main

// create a simple http server that parses the Authorization header

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cbodonnell/oauth2utils/pkg/oauth"
	"github.com/cbodonnell/oauth2utils/pkg/utils"
)

func main() {
	ctx := context.Background()

	oc, err := oauth.NewOIDCClient(ctx, "http://localhost:8080/realms/tunnel.farm", "tfarm-cli")
	if err != nil {
		log.Fatal(err)
	}

	// create an http server
	http.HandleFunc("/claims", func(w http.ResponseWriter, r *http.Request) {
		token, err := utils.ParseBearerToken(r)
		if err != nil {
			log.Printf("failed to parse token: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		idToken, err := oc.Verify(r.Context(), token)
		if err != nil {
			log.Printf("failed to verify token: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// get the claims from the token
		claims := map[string]interface{}{}
		if err := idToken.Claims(&claims); err != nil {
			log.Printf("failed to get claims: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// marshal the claims into json
		claimsJSON, err := json.Marshal(claims)
		if err != nil {
			log.Printf("failed to marshal claims: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// write the claims to the response
		fmt.Fprint(w, string(claimsJSON))
	})

	// start the server
	log.Println("starting server on :9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
