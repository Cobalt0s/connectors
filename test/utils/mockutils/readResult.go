package mockutils

import (
	"reflect"

	"github.com/amp-labs/connectors/common"
)

var (
	ReadResultComparator = readResultComparator{}
)

type readResultComparator struct{}

// SubsetRaw checks that expected ReadResult.Raw is a subset of actual ReadResult.Raw
func (readResultComparator) SubsetRaw(actual, expected *common.ReadResult) bool {
	if len(actual.Data) < len(expected.Data) {
		return false
	}

	for i := range expected.Data {
		if len(expected.Data[i].Raw) == 0 {
			panic("invalid test, there is no point to check if empty set belongs to any set; " +
				"please specify expected Raw response")
		}

		for field := range expected.Data[i].Raw {
			if !reflect.DeepEqual(actual.Data[i].Raw[field], expected.Data[i].Raw[field]) {
				return false
			}
		}
	}

	return true
}

// SubsetFields checks that expected ReadResult.Fields is a subset of actual ReadResult.Fields
func (readResultComparator) SubsetFields(actual, expected *common.ReadResult) bool {
	if len(actual.Data) < len(expected.Data) {
		return false
	}

	for i := range expected.Data {
		if len(expected.Data[i].Fields) == 0 {
			panic("invalid test, there is no point to check if empty set belongs to any set; " +
				"please specify expected Fields response")
		}

		for field := range expected.Data[i].Fields {
			if !reflect.DeepEqual(actual.Data[i].Fields[field], expected.Data[i].Fields[field]) {
				return false
			}
		}
	}

	return true
}
