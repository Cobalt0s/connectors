package main

import (
	"context"
	"fmt"
	msTest "github.com/amp-labs/connectors/test/msdsales"
	"github.com/amp-labs/connectors/test/utils"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Handle Ctrl-C gracefully.
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer done()

	// Set up slog logging.
	utils.SetupLogging()

	filePath := os.Getenv("MS_SALES_CRED_FILE")
	if filePath == "" {
		filePath = "./ms-sales-creds.json"
	}

	conn := msTest.GetMSDynamics365SalesConnector(ctx, filePath)
	defer utils.Close(conn)

	res, err := conn.ListObjectMetadata(ctx, []string{
		// TODO
	})
	if err != nil {
		utils.Fail("error reading from microsoft sales", "error", err)
	}

	fmt.Println("List object metadata..")
	utils.DumpJSON(res, os.Stdout)

}
