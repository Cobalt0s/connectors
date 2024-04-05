package msdsales

import (
	"reflect"
	"testing"

	"github.com/amp-labs/connectors/common"
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
