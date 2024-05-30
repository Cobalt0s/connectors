package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/amp-labs/connectors/common"
	connTest "github.com/amp-labs/connectors/test/xero"
	"github.com/amp-labs/connectors/test/utils"
	"github.com/amp-labs/connectors/test/utils/mockutils"
)

var (
	objectName     = "accounts"   // nolint: gochecknoglobals
	objectNameMeta = "account" // nolint: gochecknoglobals
)

// We want to compare fields returned by read and schema properties provided by metadata methods.
// Properties from read must all be present in schema definition.
func main() {
	// Handle Ctrl-C gracefully.
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer done()

	// Set up slog logging.
	utils.SetupLogging()

	filePath := os.Getenv("XERO_CRED_FILE")
	if filePath == "" {
		filePath = "./xero-creds.json"
	}

	conn := connTest.GetXeroConnector(ctx, filePath)
	defer utils.Close(conn)

	response, err := conn.Read(ctx, common.ReadParams{
		ObjectName: objectName, // TODO check endpoint path
	})
	if err != nil {
		utils.Fail("error reading from Xero", "error", err)
	}

	if response.Rows == 0 {
		utils.Fail("expected to read at least one record", "error", err)
	}

	metadata, err := conn.ListObjectMetadata(ctx, []string{
		objectNameMeta, // TODO check endpoint path
	})
	if err != nil {
		utils.Fail("error listing metadata for Xero", "error", err)
	}

	fmt.Println("Compare object metadata against endpoint response:")

	mismatchErr := mockutils.ValidateReadConformsMetadata(objectNameMeta, response.Data[0].Raw, metadata)
	if mismatchErr != nil {
		utils.Fail("schema and payload response have mismatching fields", "error", mismatchErr)
	} else {
		fmt.Println("==> success fields match.")
	}
}
