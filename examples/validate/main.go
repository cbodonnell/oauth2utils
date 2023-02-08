package main

// create a simple http server that parses the Authorization header

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/cbodonnell/oauth2utils/pkg/oauth"
)

func main() {
	ctx := context.Background()

	oc, err := oauth.NewOIDCClient(ctx, "http://localhost:8080/realms/tunnel.farm", "tfarm-cli")
	if err != nil {
		log.Fatal(err)
	}

	// create an http server
	http.HandleFunc("/claims", func(w http.ResponseWriter, r *http.Request) {
		// get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("no auth header")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// split the header into its parts
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			log.Println("invalid auth header, wrong number of parts")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// make sure the header is a bearer token
		if parts[0] != "Bearer" {
			log.Println("invalid auth header, not a bearer token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// get the token from the header
		token := parts[1]

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
