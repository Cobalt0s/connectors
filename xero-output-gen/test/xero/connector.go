package xero

import (
	"context"
	"net/http"

	"github.com/amp-labs/connectors"
	"github.com/amp-labs/connectors/xero"
	testUtils "github.com/amp-labs/connectors/test/utils"
	"github.com/amp-labs/connectors/utils"
)

func GetXeroConnector(ctx context.Context, filePath string) *xero.Connector {
	registry := utils.NewCredentialsRegistry()

	readers := []utils.Reader{
		&utils.JSONReader{
			FilePath: filePath,
			JSONPath: "$.CLIENT_ID",
			CredKey:  "clientId",
		},
		&utils.JSONReader{
			FilePath: filePath,
			JSONPath: "$.CLIENT_SECRET",
			CredKey:  "clientSecret",
		},
		&utils.JSONReader{
			FilePath: filePath,
			JSONPath: "$.ACCESS_TOKEN",
			CredKey:  "accessToken",
		},
		&utils.JSONReader{
			FilePath: filePath,
			JSONPath: "$.REFRESH_TOKEN",
			CredKey:  "refreshToken",
		},
		&utils.JSONReader{
			FilePath: filePath,
			JSONPath: "$.PROVIDER",
			CredKey:  "provider",
		},
	}
	_ = registry.AddReaders(readers...)

	// TODO create config and token registries
	cfg := utils.XeroConfigFromRegistry(registry)
	tok := utils.XeroTokenFromRegistry(registry)

	// TODO provide required options
	conn, err := connectors.Xero(
		xero.WithClient(ctx, http.DefaultClient, cfg, tok),
		xero.WithWorkspace(utils.XeroWorkspace),
		xero.WithModule(xero.DefaultModuleCRM),
	)
	if err != nil {
		testUtils.Fail("error creating connector", "error", err)
	}

	return conn
}
