package msdsales

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/amp-labs/connectors/common"
	"github.com/amp-labs/connectors/common/interpreter"
	"github.com/amp-labs/connectors/common/reqrepeater"
	"github.com/go-test/deep"
)

var (
	metadataTestInputFile = "metadata.xml"
)

func Test_ListObjectMetadata(t *testing.T) {
	t.Parallel()

	fakeServerResp, err := readTestFile(metadataTestInputFile)
	if err != nil {
		t.Fatalf("failed to start test, input file missing, %v", err)
	}

	tests := []struct {
		name                string
		input               []string
		server              *httptest.Server
		connector           Connector
		expected            *common.ListObjectMetadataResult
		expectedFieldsCount map[string]int // used instead of `expected` when response result is too big
		expectedErrs        []error
	}{
		{
			name:  "At least one object name must be queried",
			input: nil,
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusTeapot)
			})),
			expectedErrs: []error{common.ErrMissingObjects},
		},
		{
			name:  "Mime response header expected",
			input: []string{"msfp_surveyinvite"},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusTeapot)
			})),
			expectedErrs: []error{interpreter.MissingContentType},
		},
		{
			name:  "Missing XML response on status OK",
			input: []string{"msfp_surveyinvite"},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/xml")
				w.WriteHeader(http.StatusOK)
				writeBody(w, "")
			})),
			expectedErrs: []error{common.ErrNotXML},
		},
		{
			name:  "Missing XML root",
			input: []string{"msfp_surveyinvite"},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/xml")
				w.WriteHeader(http.StatusOK)
				writeBody(w, `<?xml version="1.0" encoding="utf-8"?>`)
			})),
			expectedErrs: []error{common.ErrNoXMLRoot},
		},
		{
			name:  "Server response without Sales Schema",
			input: []string{"msfp_surveyinvite"},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/xml")
				w.WriteHeader(http.StatusOK)
				writeBody(w, `
<?xml version="1.0" encoding="utf-8"?>
<edmx:Edmx Version="4.0" xmlns:edmx="http://docs.oasis-open.org/odata/ns/edmx"></edmx:Edmx>`)
			})),
			expectedErrs: []error{ErrMissingSchema},
		},
		{
			name:  "Object name cannot be found from server response",
			input: []string{"msfp_surveyinvite"},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/xml")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(fakeServerResp)
			})),
			expectedErrs: []error{ErrObjectNotFound, errors.New("unknown name msfp_surveyinvite")},
		},
		{
			name:  "Correctly list metadata for account leads and invite contact",
			input: []string{"accountleads", "adx_invitation_invitecontacts"},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/xml")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(fakeServerResp)
			})),
			expected: &common.ListObjectMetadataResult{
				Result: map[string]common.ObjectMetadata{
					"accountleads": {
						DisplayName: "accountleads",
						FieldsMap: map[string]string{
							"accountid":                 "accountid",
							"accountleadid":             "accountleadid",
							"importsequencenumber":      "importsequencenumber",
							"leadid":                    "leadid",
							"name":                      "name",
							"overriddencreatedon":       "overriddencreatedon",
							"timezoneruleversionnumber": "timezoneruleversionnumber",
							"utcconversiontimezonecode": "utcconversiontimezonecode",
							"versionnumber":             "versionnumber",
						},
					},
					"adx_invitation_invitecontacts": {
						DisplayName: "adx_invitation_invitecontacts",
						FieldsMap: map[string]string{
							"adx_invitation_invitecontactsid": "adx_invitation_invitecontactsid",
							"adx_invitationid":                "adx_invitationid",
							"contactid":                       "contactid",
							"versionnumber":                   "versionnumber",
						},
					},
				},
				Errors: nil,
			},
			expectedErrs: nil,
		},
		{
			// In total phonecall will have 65 fields, where
			// phonecall 		(has 7 fields) and inherits from
			// activitypointer 	(has 58 fields), which in turn inherits from
			// crmbaseentity 	(has 0 fields)
			name:  "Correctly list metadata for phone calls including inherited fields",
			input: []string{"phonecall", "activitypointer", "crmbaseentity"},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/xml")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(fakeServerResp)
			})),
			expectedFieldsCount: map[string]int{
				"phonecall":       65,
				"activitypointer": 58,
				"crmbaseentity":   0,
			},
			expectedErrs: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()

			ctx := context.Background()

			connector, err := NewConnector(
				WithAuthenticatedClient(http.DefaultClient),
				WithWorkspace("test-workspace"),
			)
			if err != nil {
				t.Fatalf("%s: error in test while constructin connector %v", tt.name, err)
			}

			// for testing we want to redirect calls to our server
			connector.BaseURL = tt.server.URL
			connector.Client.HTTPClient.Base = tt.server.URL
			// failed requests will be not retried
			connector.RetryStrategy = &reqrepeater.NullStrategy{}

			// start of tests
			output, err := connector.ListObjectMetadata(ctx, tt.input)
			if len(tt.expectedErrs) == 0 && err != nil {
				t.Fatalf("%s: expected no errors, got: (%v)", tt.name, err)
			}

			if len(tt.expectedErrs) != 0 && err == nil {
				t.Fatalf("%s: expected errors (%v), but got nothing", tt.name, tt.expectedErrs)
			}

			for _, expectedErr := range tt.expectedErrs {
				if !errors.Is(err, expectedErr) && !strings.Contains(err.Error(), expectedErr.Error()) {
					t.Fatalf("%s: expected Error: (%v), got: (%v)", tt.name, expectedErr, err)
				}
			}
			if tt.expectedFieldsCount != nil {
				// we are comparing if number of fields match under ListObjectMetadataResult.Result
				for entityName, count := range tt.expectedFieldsCount {
					entity, ok := output.Result[entityName]
					if !ok {
						t.Fatalf("%s: expected entity was missing: (%v)", tt.name, entityName)
					}
					got := len(entity.FieldsMap)
					if got != count {
						t.Fatalf("%s: expected entity '%v' to have (%v) fields got: (%v)",
							tt.name, entityName, count, got)
					}
				}

			} else {
				// usual comparison of ListObjectMetadataResult
				if !reflect.DeepEqual(output, tt.expected) {
					diff := deep.Equal(output, tt.expected)
					t.Fatalf("%s:, \nexpected: (%v), \ngot: (%v), \ndiff: (%v)",
						tt.name, tt.expected, output, diff)
				}
			}
		})
	}
}

func readTestFile(testFileName string) ([]byte, error) {
	_, runnerLocation, _, _ := runtime.Caller(0)
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	relativePath, _ := strings.CutPrefix(runnerLocation, workingDir)
	testDir := path.Join(".", relativePath, "../test")
	return os.ReadFile(testDir + "/" + testFileName)
}
