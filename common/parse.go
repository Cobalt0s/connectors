package common

import (
	"errors"
	"math"
	"strings"

	"github.com/spyzhov/ajson"
)

// ParseResult parses the response from a provider into a ReadResult. A 2xx return type is assumed.
// The sizeFunc, recordsFunc, nextPageFunc, and marshalFunc are used to extract the relevant data from the response.
// The sizeFunc returns the total number of records in the response.
// The recordsFunc returns the records in the response.
// The nextPageFunc returns the URL for the next page of results.
// The marshalFunc is used to structure the data into an array of ReadResultRows.
// The fields are used to populate ReadResultRow.Fields.
func ParseResult(
	resp *JSONHTTPResponse,
	sizeFunc func(*ajson.Node) (int64, error),
	recordsFunc func(*ajson.Node) ([]map[string]any, error),
	nextPageFunc func(*ajson.Node) (string, error),
	marshalFunc func([]map[string]any, []string) ([]ReadResultRow, error),
	fields []string,
) (*ReadResult, error) {
	totalSize, err := sizeFunc(resp.Body)
	if err != nil {
		return nil, err
	}

	records, err := recordsFunc(resp.Body)
	if err != nil {
		return nil, err
	}

	nextPage, err := nextPageFunc(resp.Body)
	if err != nil {
		return nil, err
	}

	done := nextPage == ""

	marshaledData, err := marshalFunc(records, fields)
	if err != nil {
		return nil, err
	}

	return &ReadResult{
		Rows:     totalSize,
		Data:     marshaledData,
		NextPage: NextPageToken(nextPage),
		Done:     done,
	}, nil
}

// ExtractLowercaseFieldsFromRaw returns a map of fields from a record.
// The fields are all returned in lowercase.
func ExtractLowercaseFieldsFromRaw(fields []string, record map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(fields))

	// Modify all record keys to lowercase
	lowercaseRecord := make(map[string]interface{}, len(record))
	for key, value := range record {
		lowercaseRecord[strings.ToLower(key)] = value
	}

	for _, field := range fields {
		// Lowercase the field name to make lookup case-insensitive.
		lowercaseField := strings.ToLower(field)

		if value, ok := lowercaseRecord[lowercaseField]; ok {
			out[lowercaseField] = value
		}
	}

	return out
}

var (
	ErrNotArray   = errors.New("results is not an array")
	ErrNotObject  = errors.New("result is not an object")
	ErrNotString  = errors.New("link is not a string")
	ErrNotNumeric = errors.New("value is not numeric")
	ErrNotInteger = errors.New("value is not integer")

	// JsonManager is a helpful wrapper of ajson library that adds errors when querying JSON payload
	// and provides common conversion methods.
	JsonManager = jsonManager{}
)

type jsonManager struct{}

func (jsonManager) ArrToMap(arr []*ajson.Node) ([]map[string]any, error) {
	output := make([]map[string]any, 0, len(arr))

	for _, v := range arr {
		if !v.IsObject() {
			return nil, ErrNotObject
		}

		data, err := v.Unpack()
		if err != nil {
			return nil, err
		}

		m, ok := data.(map[string]interface{})
		if !ok {
			return nil, ErrNotObject
		}

		output = append(output, m)
	}

	return output, nil
}

func (jsonManager) GetInteger(node *ajson.Node, key string) (int64, error) {
	innerNode, err := node.GetKey(key)
	if err != nil {
		return 0, err
	}

	count, err := innerNode.GetNumeric()
	if err != nil {
		return 0, ErrNotNumeric
	}

	if math.Mod(count, 1.0) != 0 {
		return 0, ErrNotInteger
	}

	return int64(count), nil
}

func (jsonManager) GetArr(node *ajson.Node, key string) ([]*ajson.Node, error) {
	records, err := node.GetKey(key)
	if err != nil {
		return nil, err
	}

	arr, err := records.GetArray()
	if err != nil {
		return nil, ErrNotArray
	}

	return arr, nil
}

func (jsonManager) ArrSize(node *ajson.Node, key string) (int64, error) {
	innerNode, err := node.GetKey(key)
	if err != nil {
		return 0, err
	}

	if !innerNode.IsArray() {
		return 0, ErrNotArray
	}

	return int64(innerNode.Size()), nil
}

func (jsonManager) GetString(node *ajson.Node, key string) (string, error) {
	innerNode, err := node.GetKey(key)
	if err != nil {
		return "", nil
	}

	txt, err := innerNode.GetString()
	if err != nil {
		return "", ErrNotString
	}

	return txt, nil
}
