package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/amp-labs/connectors/common"
	msTest "github.com/amp-labs/connectors/test/msdsales"
	"github.com/amp-labs/connectors/test/utils"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// we want to compare fields returned by read and schema properties provided by metadata methods
// they must match for all such objects
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

	response, err := conn.Read(ctx, common.ReadParams{
		ObjectName: "accounts",
		PageSize:   1,
	})
	if err != nil {
		utils.Fail("error reading from microsoft sales", "error", err)
	}
	if response.Rows != 1 {
		utils.Fail("expected to read exactly one record", "error", err)
	}

	metadata, err := conn.ListObjectMetadata(ctx, []string{
		"account",
	})
	if err != nil {
		utils.Fail("error listing metadata for microsoft sales", "error", err)
	}

	fmt.Println("Compare object metadata against endpoint response:")
	mismatchErr := compareFieldsMatch(metadata, response.Data[0].Raw, "account")
	if mismatchErr != nil {
		utils.Fail("schema and payload response have mismatching fields", "error", mismatchErr)
	} else {
		fmt.Println("... success fields match.")
	}
}

func compareFieldsMatch(metadata *common.ListObjectMetadataResult, response map[string]any, objectName string) error {
	fields := make(map[string]bool)
	for field := range response {
		// ignore all fields that are OData annotations
		if !strings.Contains(field, "@") {
			fields[field] = false
		}
	}

	mismatch := make([]error, 0)
	for _, displayName := range metadata.Result[objectName].FieldsMap {
		if _, found := fields[displayName]; found {
			fields[displayName] = true
		} else {
			mismatch = append(mismatch, fmt.Errorf("read payload doesn't have %v", displayName))
		}
	}
	for name, found := range fields {
		if !found {
			mismatch = append(mismatch, fmt.Errorf("metadata schema is missing field %v", name))
		}
	}
	return errors.Join(mismatch...)
}
