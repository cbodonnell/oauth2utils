package oauth

import (
	"context"
	"net/http"

	"github.com/cbodonnell/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type OIDCClient struct {
	provider *oidc.Provider
	conf     *oauth2.Config
	verifier *oidc.IDTokenVerifier
}

func NewOIDCClient(ctx context.Context, issuer string, clientID string, additionalScopes []string) (*OIDCClient, error) {
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, err
	}

	conf := &oauth2.Config{
		ClientID: clientID,
		Endpoint: provider.Endpoint(),
		Scopes:   []string{oidc.ScopeOpenID},
	}

	if additionalScopes != nil {
		conf.Scopes = append(conf.Scopes, additionalScopes...)
	}

	return &OIDCClient{
		provider: provider,
		conf:     conf,
		verifier: provider.Verifier(&oidc.Config{ClientID: clientID}),
	}, nil
}

func (oc *OIDCClient) HTTPClient(ctx context.Context, token *oauth2.Token) *http.Client {
	return oc.conf.Client(ctx, token)
}

func (oc *OIDCClient) TokenSource(ctx context.Context, token *oauth2.Token) oauth2.TokenSource {
	return oc.conf.TokenSource(ctx, token)
}

func (oc *OIDCClient) UserInfo(ctx context.Context, tokenSource oauth2.TokenSource) (*oidc.UserInfo, error) {
	return oc.provider.UserInfo(ctx, tokenSource)
}

func (oc *OIDCClient) Verify(ctx context.Context, rawIDToken string) (*oidc.IDToken, error) {
	return oc.verifier.Verify(ctx, rawIDToken)
}
