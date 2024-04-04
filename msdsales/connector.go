package msdsales

import (
	"fmt"
	"github.com/amp-labs/connectors/common"
	"github.com/amp-labs/connectors/common/facade/interpreter"
	"github.com/amp-labs/connectors/common/facade/paramsbuilder"
	"github.com/amp-labs/connectors/common/facade/repeaters"
	"github.com/amp-labs/connectors/providers"
	"strings"
	"time"
)

var DefaultModuleCRM = paramsbuilder.APIModule{ // nolint: gochecknoglobals
	Label:   "api/data",
	Version: "v9.2",
}

type Connector struct {
	BaseURL       string
	Module        string
	Client        *common.JSONHTTPClient
	RetryStrategy repeaters.Strategy
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
		Module:  params.Module.Suffix,
		Client:  params.Client.Caller,
		RetryStrategy: &repeaters.UniformRetryStrategy{ // FIXME call retry strategy could be part of options
			RetriesNum: 3,
			Interval:   time.Second,
		},
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
	return fmt.Sprintf("%s.Connector[%s]", c.Provider(), c.Module)
}

func (c *Connector) getURL(arg string) string { // FIXME should it be part of Connector interface?
	return strings.Join([]string{c.BaseURL, c.Module, arg}, "/")
}
