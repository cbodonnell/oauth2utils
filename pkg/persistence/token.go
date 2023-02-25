package persistence

import (
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
)

const (
	TokenFile = "token.json"
)

func SaveToken(token *oauth2.Token, tokenDir string) error {
	// get the user home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// create the .tfarm directory
	dir := fmt.Sprintf("%s/%s", home, tokenDir)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	// create the token file
	tokenFile := fmt.Sprintf("%s/%s", dir, TokenFile)
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
	// get the user home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// open the token file
	tokenFile := fmt.Sprintf("%s/%s/%s", home, tokenDir, TokenFile)
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
	// get the user home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// open the token file
	tokenFile := fmt.Sprintf("%s/%s/%s", home, tokenDir, TokenFile)
	err = os.Remove(tokenFile)
	if err != nil {
		return err
	}

	return nil
}
