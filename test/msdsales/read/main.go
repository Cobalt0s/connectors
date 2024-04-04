package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/amp-labs/connectors/common"
	msTest "github.com/amp-labs/connectors/test/msdsales"
	"github.com/amp-labs/connectors/test/utils"
)

type Contact struct {
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Company   string `json:"company"`
	Website   string `json:"website"`
	Lastname  string `json:"lastname"`
	Firstname string `json:"firstname"`
}

/*
	Running instructions

	Step 1: prepare "ms-dales-creds.json" file

	Create a file called "ms-sales-creds.json" in the root of the project with the following contents

		e.g. {
		"CLIENT_ID": "<client id goes here>",
		"CLIENT_SECRET": "<client secret goes here>",
		"ACCESS_TOKEN": "<access token goes here>",
		"REFRESH_TOKEN": "<refresh token goes here>"
		}

	or export to an environment variable MS_SALES_CRED_FILE by following command

	$> export MS_SALES_CRED_FILE=./ms-sales-creds.json # or the path to your ms-sales-creds.json file


	In 1password, you can find a MS Sales creds.json file in the "Shared" vault. TODO this must be in 1password
	Look for the title "MS Sales Sample OAuth Credentials".
	The 1password item has an attached file called "creds.json" that contains the JSON.

	Step 2: run the following command

		$> go run test/msdsales/read/main.go


*/

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

	res, err := conn.Read(ctx, common.ReadParams{
		ObjectName: "contacts",
		Fields:     []string{"fullname"},
		NextPage:   "",
		Since:      time.Now().Add(-5 * time.Minute),
	})
	if err != nil {
		utils.Fail("error reading from microsoft sales", "error", err)
	}

	fmt.Println("Reading contacts..")
	utils.DumpJSON(res, os.Stdout)
}
