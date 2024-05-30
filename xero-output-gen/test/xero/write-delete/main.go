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
	"github.com/amp-labs/connectors/test/utils/mockutils"
)

type leadPayload struct { // TODO fill in properties
	Name       string `json:"name,omitempty"`
	View       string `json:"view,omitempty"`
	ViewParams string `json:"view_params,omitempty"` // JSON object of list view parameters
	IsDefault  *bool  `json:"is_default,omitempty"`
	Shared     bool   `json:"shared,omitempty"`
}

var (
	objectName = "leads" // nolint: gochecknoglobals
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

	fmt.Println("> TEST Create/Update/Delete lead")
	fmt.Println("Creating lead")

	// NOTE: list view must have unique `Name`
	view := createLead(ctx, conn, &leadPayload{
		Name:       "Tom's Prospects",
		View:       "companies",
		ViewParams: "",
		IsDefault:  mockutils.Pointers.Bool(true),
		Shared:     false,
	})

	fmt.Println("Updating some lead properties")
	updateLead(ctx, conn, view.RecordId, &leadPayload{
		Name:      "Jerry's Prospects",
		View:      "companies",
		IsDefault: mockutils.Pointers.Bool(false),
	})

	fmt.Println("View that lead has changed accordingly")

	res := readLead(ctx, conn)

	updatedView := searchLead(res, "id", view.RecordId)
	for k, v := range map[string]string{
		"name":       "Jerry's Prospects",
		"view":       "companies",
		"is_default": "false",
		"shared":     "false",
	} {
		if !mockutils.DoesObjectCorrespondToString(updatedView[k], v) {
			utils.Fail("error updated properties do not match", k, v, updatedView[k])
		}
	}

	fmt.Println("Removing this lead")
	removeLead(ctx, conn, view.RecordId)
	fmt.Println("> Successful test completion")
}

func searchLead(res *common.ReadResult, key, value string) map[string]any {
	for _, data := range res.Data {
		if mockutils.DoesObjectCorrespondToString(data.Fields[key], value) {
			return data.Raw
		}
	}

	utils.Fail("error finding lead")

	return nil
}

func readLead(ctx context.Context, conn *xero.Connector) *common.ReadResult {
	res, err := conn.Read(ctx, common.ReadParams{
		ObjectName: objectName,
		Fields: []string{
			"id", "view", "name",
		},
	})
	if err != nil {
		utils.Fail("error reading from Xero", "error", err)
	}

	return res
}

func createLead(ctx context.Context, conn *xero.Connector, payload *leadPayload) *common.WriteResult {
	res, err := conn.Write(ctx, common.WriteParams{
		ObjectName: objectName,
		RecordId:   "",
		RecordData: payload,
	})
	if err != nil {
		utils.Fail("error writing to Xero", "error", err)
	}

	if !res.Success {
		utils.Fail("failed to create a lead")
	}

	return res
}

func updateLead(ctx context.Context, conn *xero.Connector, viewID string, payload *leadPayload) *common.WriteResult {
	res, err := conn.Write(ctx, common.WriteParams{
		ObjectName: objectName,
		RecordId:   viewID,
		RecordData: payload,
	})
	if err != nil {
		utils.Fail("error writing to Xero", "error", err)
	}

	if !res.Success {
		utils.Fail("failed to update a lead")
	}

	return res
}

func removeLead(ctx context.Context, conn *xero.Connector, viewID string) {
	res, err := conn.Delete(ctx, common.DeleteParams{
		ObjectName: objectName,
		RecordId:   viewID,
	})
	if err != nil {
		utils.Fail("error deleting for Xero", "error", err)
	}

	if !res.Success {
		utils.Fail("failed to remove a lead")
	}
}
