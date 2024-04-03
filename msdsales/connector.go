package msdsales

import (
	"context"
	"github.com/amp-labs/connectors/common"
	"github.com/amp-labs/connectors/common/facade/interpreter"
	"github.com/amp-labs/connectors/providers"
)

type Connector struct {
	BaseURL string
	Client  *common.JSONHTTPClient
}

func (c *Connector) ListObjectMetadata(ctx context.Context, objectNames []string) (*common.ListObjectMetadataResult, error) {
	//TODO implement me
	panic("implement me")
}

func NewConnector(opts ...Option) (conn *Connector, outErr error) {
	defer common.PanicRecovery(func(cause error) {
		outErr = cause
		conn = nil
	})

	params, err := parameters{}.FromOptions(opts...)
	if err != nil {
		return nil, err
	}

	providerInfo, err := providers.ReadInfo(providers.MicrosoftDynamics365Sales, &map[string]string{
		"workspace": params.Workspace.Name,
	})
	if err != nil {
		return nil, err
	}

	baseURL := providerInfo.BaseURL
	conn = &Connector{
		BaseURL: baseURL,
		Client:  params.Client.Caller,
	}
	// connector and its client must mirror base url and provide its own error parser
	conn.Client.HTTPClient.Base = baseURL
	conn.Client.HTTPClient.ErrorHandler = interpreter.ErrorHandler{
		JSON: conn.interpretJSONError,
	}.Handle

	return conn, nil
}

func (c *Connector) Provider() providers.Provider {
	return providers.MicrosoftDynamics365Sales
}
func (c *Connector) String() string {
	return c.Provider() + ".Connector"
}
