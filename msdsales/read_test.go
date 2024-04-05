package msdsales

import (
	"context"
	"errors"
	"github.com/amp-labs/connectors/common"
	"github.com/amp-labs/connectors/common/interpreter"
	"github.com/amp-labs/connectors/common/reqrepeater"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func Test_makeQueryValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    common.ReadParams
		expected string
	}{
		{
			name:     "No params no query string",
			input:    common.ReadParams{},
			expected: "",
		},
		{
			name: "One parameter",
			input: common.ReadParams{
				Fields: []string{"cat"},
			},
			expected: "?$select=cat",
		},
		{
			name: "Many parameters",
			input: common.ReadParams{
				Fields: []string{"cat", "dog", "parrot", "hamster"},
			},
			expected: "?$select=cat,dog,parrot,hamster",
		},
		{
			name: "OData parameters with @ symbol",
			input: common.ReadParams{
				Fields: []string{"cat", "@odata.dog", "parrot"},
			},
			expected: "?$select=cat,@odata.dog,parrot",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := makeQueryValues(tt.input)
			if !reflect.DeepEqual(output, tt.expected) {
				t.Fatalf("%s: expected: (%v), got: (%v)", tt.name, tt.expected, output)
			}
		})
	}
}

func Test_Read(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		input        common.ReadParams
		server       *httptest.Server
		connector    Connector
		expected     *common.ReadResult
		expectedErrs []error
	}{
		{
			name: "Mime response header expected",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusTeapot)
			})),
			expectedErrs: []error{interpreter.MissingContentType},
		},
		{
			name: "Correct error message is understood from JSON response",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				writeBody(w, `{
					"error": {
						"code": "some-code",
						"message": "your fault"
					}
				}`)
			})),
			expectedErrs: []error{common.ErrBadRequest, errors.New("your fault")},
		},
		{
			name: "Incorrect key in payload",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				writeBody(w, `{
					"garbage": {}
				}`)
			})),
			expectedErrs: []error{errors.New("wrong request: wrong key 'value'")},
		},
		{
			name: "Incorrect data type in payload",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				writeBody(w, `{
					"value": {}
				}`)
			})),
			expectedErrs: []error{common.ErrNotArray},
		},
		// TODO there are more test to write for pagination
		//{
		//	name: "@odata.nextLink must be in payload",
		//	server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//		w.Header().Set("Content-Type", "application/json")
		//		w.WriteHeader(http.StatusOK)
		//		writeBody(w, `{
		//			"value": []
		//		}`)
		//	})),
		//	expectedErrs: []error{common.ErrNotArray},
		//},
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
			output, err := connector.Read(ctx, tt.input)
			if len(tt.expectedErrs) == 0 && err != nil {
				t.Fatalf("%s: expected no errors, got: (%v)", tt.name, err)
			}

			for _, expectedErr := range tt.expectedErrs {
				if !errors.Is(err, expectedErr) && !strings.Contains(err.Error(), expectedErr.Error()) {
					t.Fatalf("%s: expected Error: (%v), got: (%v)", tt.name, expectedErr, err)
				}
			}

			if !reflect.DeepEqual(output, tt.expected) {
				t.Fatalf("%s: expected: (%v), got: (%v)", tt.name, tt.expected, output)
			}
		})
	}
}

func writeBody(w http.ResponseWriter, body string) {
	_, _ = w.Write([]byte(body))
}
