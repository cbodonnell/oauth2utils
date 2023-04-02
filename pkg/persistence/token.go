package persistence

import (
	"encoding/json"
	"os"
	"path"

	"golang.org/x/oauth2"
)

const (
	TokenFile = "token.json"
)

func SaveToken(token *oauth2.Token, tokenDir string) error {
	err := os.MkdirAll(tokenDir, 0755)
	if err != nil {
		return err
	}

	// create the token file
	tokenFile := path.Join(tokenDir, TokenFile)
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

func LoadToken(tokenDir string) (*oauth2.Token, error) {
	// open the token file
	tokenFile := path.Join(tokenDir, TokenFile)
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

func DeleteToken(tokenDir string) error {
	// open the token file
	tokenFile := path.Join(tokenDir, TokenFile)
	err := os.Remove(tokenFile)
	if err != nil {
		return err
	}

	return nil
}
