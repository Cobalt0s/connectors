package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/amp-labs/connectors/common"
	"github.com/amp-labs/connectors/xero"
	connTest "github.com/amp-labs/connectors/test/xero"
	"github.com/amp-labs/connectors/test/utils"
)

var (
	objectName = "contacts" // nolint: gochecknoglobals
)

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

	res, err := conn.Read(ctx, common.ReadParams{
		ObjectName: objectName, // TODO check endpoint path
		Fields: []string{
			"fullname", "emailaddress1", "fax", // TODO provide fields
		},
	})
	if err != nil {
		utils.Fail("error reading from Xero", "error", err)
	}

	fmt.Println("Reading contact..")
	utils.DumpJSON(res, os.Stdout)

	if res.Rows > xero.DefaultPageSize {
		utils.Fail(fmt.Sprintf("expected max %v rows", xero.DefaultPageSize))
	}
}
