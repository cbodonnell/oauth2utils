package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

type DeviceCodeResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`

	// Error and ErrorDescription are only set if the request failed
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type TokenResponse struct {
	AccessToken           string `json:"access_token"`
	TokenType             string `json:"token_type"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_expires_in"`
	Scope                 string `json:"scope"`
	NotBeforePolicy       int    `json:"not-before-policy"`
	SessionState          string `json:"session_state"`

	// Error and ErrorDescription are only set if the request failed
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type ErrTryLater struct {
	Err error
}

func (e ErrTryLater) Error() string {
	return e.Err.Error()
}

func (oc *OIDCClient) DeviceCode(ctx context.Context) (*oauth2.Token, error) {
	deviceCode, err := getDeviceCode(ctx, oc.conf, oc.provider.Endpoint().DeviceAuthURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get the device code: %v", err)
	}

	fmt.Println("Please go to the following URL to verify your device:")
	fmt.Println(deviceCode.VerificationURI)
	fmt.Println()
	fmt.Println("Enter the device code:", deviceCode.UserCode)
	fmt.Println()

	var token *oauth2.Token
	for i := 0; i < deviceCode.ExpiresIn; i += deviceCode.Interval {
		token, err = getAccessToken(ctx, oc.conf, deviceCode)
		if err == nil {
			break
		}

		if _, ok := err.(ErrTryLater); ok {
			time.Sleep(time.Duration(deviceCode.Interval) * time.Second)
			i += deviceCode.Interval
			continue
		}

		return nil, fmt.Errorf("failed to get the access token: %v", err)
	}

	return token, nil
}

func getDeviceCode(ctx context.Context, conf *oauth2.Config, deviceURL string) (*DeviceCodeResponse, error) {
	data := url.Values{}
	data.Set("client_id", conf.ClientID)
	data.Set("scope", strings.Join(conf.Scopes, " "))

	req, err := http.NewRequestWithContext(ctx, "POST", deviceURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create the request: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send the request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response code: %s", res.Status)
	}

	var deviceCode DeviceCodeResponse
	if err := json.Unmarshal(body, &deviceCode); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the response: %v", err)
	}

	return &deviceCode, nil
}

func getAccessToken(ctx context.Context, conf *oauth2.Config, deviceCode *DeviceCodeResponse) (*oauth2.Token, error) {
	data := url.Values{}
	data.Set("device_code", deviceCode.DeviceCode)
	data.Set("client_id", conf.ClientID)
	data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")

	req, err := http.NewRequestWithContext(ctx, "POST", conf.Endpoint.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create the request: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send the request: %v", err)
	}
	defer res.Body.Close()

	var tokenResponse TokenResponse
	if err := json.NewDecoder(res.Body).Decode(&tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to decode the response: %v", err)
	}

	if tokenResponse.Error != "" {
		if tokenResponse.Error == "authorization_pending" || tokenResponse.Error == "slow_down" {
			return nil, ErrTryLater{Err: err}
		}

		return nil, fmt.Errorf("error in the token response: %s", tokenResponse.Error)
	}

	token := &oauth2.Token{
		AccessToken:  tokenResponse.AccessToken,
		TokenType:    tokenResponse.TokenType,
		RefreshToken: tokenResponse.RefreshToken,
		Expiry:       time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second),
	}

	return token, nil
}
