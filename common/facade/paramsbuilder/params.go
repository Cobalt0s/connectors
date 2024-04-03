package paramsbuilder

import (
	"context"
	"errors"
	"github.com/amp-labs/connectors/common"
	"golang.org/x/oauth2"
	"net/http"
)

// ParamAssurance checks that param data is valid
// Every param instance must implement it
//
type ParamAssurance interface {
	ValidateParams() error
}

var (
	ErrMissingClient    = errors.New("JSON http client not set")
	ErrMissingWorkspace   = errors.New("missing workspace name")
)

// Client params sets up authenticated proxy HTTP client
// This can be reused among other param builders by composition
type Client struct {
	Caller *common.JSONHTTPClient
}

func (p *Client) ValidateParams() error {
	if p.Caller == nil {
		return ErrMissingClient
	}
	return nil
}

func (p *Client) WithClient(
	ctx context.Context, client *http.Client,
	config *oauth2.Config, token *oauth2.Token,
	opts ...common.OAuthOption) {

	options := []common.OAuthOption{
		common.WithClient(client),
		common.WithOAuthConfig(config),
		common.WithOAuthToken(token),
	}

	oauthClient, err := common.NewOAuthHTTPClient(ctx, append(options, opts...)...)
	if err != nil {
		panic(err) // caught in NewConnector
	}

	p.WithAuthenticatedClient(oauthClient)
}

func (p *Client) WithAuthenticatedClient(client common.AuthenticatedHTTPClient) {
	p.Caller = &common.JSONHTTPClient{
		HTTPClient: &common.HTTPClient{
			Client:       client,
			ErrorHandler: common.InterpretError,
		},
	}
}

// Workspace params sets up varying workspace name
type Workspace struct {
	Name string
}

func (p *Workspace) ValidateParams() error {
	if len(p.Name) == 0 {
		return ErrMissingWorkspace
	}
	return nil
}

func (p *Workspace) WithWorkspace(workspaceRef string) {
	p.Name = workspaceRef
}
