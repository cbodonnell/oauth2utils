package persistence

import (
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
)

// TODO: Get these from env vars and defaults
const (
	// TokenFile is the name of the file where the token is stored
	TokenDir  = ".tfarm"
	TokenFile = "token.json"
)

func SaveToken(token *oauth2.Token) error {
	// get the user home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// create the .tfarm directory
	tfarmDir := fmt.Sprintf("%s/%s", home, TokenDir)
	err = os.MkdirAll(tfarmDir, 0755)
	if err != nil {
		return err
	}

	// create the token file
	tokenFile := fmt.Sprintf("%s/%s", tfarmDir, TokenFile)
	f, err := os.Create(tokenFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// write the token to the file
	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		return err
	}

	return nil
}

func LoadToken() (*oauth2.Token, error) {
	// get the user home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// open the token file
	tokenFile := fmt.Sprintf("%s/%s/%s", home, TokenDir, TokenFile)
	f, err := os.Open(tokenFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// read the token from the file
	var token oauth2.Token
	err = json.NewDecoder(f).Decode(&token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func DeleteToken() error {
	// get the user home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// open the token file
	tokenFile := fmt.Sprintf("%s/%s/%s", home, TokenDir, TokenFile)
	err = os.Remove(tokenFile)
	if err != nil {
		return err
	}

	return nil
}
