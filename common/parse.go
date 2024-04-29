package common

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/spyzhov/ajson"
)

type NextPageFunc func(*ajson.Node) (string, error)

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

const (
	KeysZoomSeparator = "."
)

var (
	ErrNotArray   = errors.New("JSON value is not an array")
	ErrNotObject  = errors.New("JSON value is not an object")
	ErrNotString  = errors.New("JSON value is not a string")
	ErrNotNumeric = errors.New("JSON value is not a numeric")
	ErrNotInteger = errors.New("JSON value is not an integer")

	// JSONManager is a helpful wrapper of ajson library that adds errors when querying JSON payload
	// and provides common conversion methods.
	JSONManager = jsonManager{} //nolint:gochecknoglobals
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

func (m jsonManager) GetIntegerWithDefault(node *ajson.Node, key string, defaultValue int64) (int64, error) {
	result, err := m.GetInteger(node, key, true)
	if err != nil {
		return 0, err
	}

	if result == nil {
		return defaultValue, err
	}

	return *result, nil
}

func (jsonManager) GetInteger(node *ajson.Node, key string, optional bool) (*int64, error) {
	if optional && !node.HasKey(key) {
		// null value in payload is allowed
		return nil, nil // nolint:nilnil
	}

	innerNode, err := node.GetKey(key)
	if err != nil {
		return nil, err
	}

	if innerNode.IsNull() {
		return nil, handleNullNode(key, optional)
	}

	count, err := innerNode.GetNumeric()
	if err != nil {
		return nil, ErrNotNumeric
	}

	if math.Mod(count, 1.0) != 0 {
		return nil, ErrNotInteger
	}

	result := int64(count)

	return &result, nil
}

func (jsonManager) GetArr(node *ajson.Node, key string) ([]*ajson.Node, error) {
	records, err := node.GetKey(key)
	if err != nil {
		return nil, err
	}

	arr, err := records.GetArray()
	if err != nil {
		return nil, formatProblematicKeyError(key, ErrNotArray)
	}

	return arr, nil
}

func (m jsonManager) ArrSize(node *ajson.Node, key string) (int64, error) {
	arr, err := m.GetArr(node, key)
	if err != nil {
		return 0, err
	}

	return int64(len(arr)), nil
}

func (m jsonManager) GetStringWithDefault(node *ajson.Node, key string, defaultValue string) (string, error) {
	result, err := m.GetString(node, key, true)
	if err != nil {
		return "", err
	}

	if result == nil {
		return defaultValue, err
	}

	return *result, nil
}

func (jsonManager) GetString(node *ajson.Node, key string, optional bool) (*string, error) {
	if optional && !node.HasKey(key) {
		// null value in payload is allowed
		return nil, nil // nolint:nilnil
	}

	innerNode, err := node.GetKey(key)
	if err != nil {
		return nil, err
	}

	if innerNode.IsNull() {
		return nil, handleNullNode(key, optional)
	}

	txt, err := innerNode.GetString()
	if err != nil {
		return nil, ErrNotString
	}

	return &txt, nil
}

// GetNestedNode reaches into the JSON node by zooming via dot separated keys
// Ex: keysZoom = item.shipping.address => item has object with shipping key and so on.
func (jsonManager) GetNestedNode(node *ajson.Node, keysZoom string) (inner *ajson.Node, err error) {
	inner = node
	keys := strings.Split(keysZoom, KeysZoomSeparator)

	// traverse nested JSON, use every key to zoom in
	for _, key := range keys {
		inner, err = inner.GetKey(key)
		if err != nil {
			return nil, err
		}
	}

	if !inner.IsObject() {
		return nil, ErrNotObject
	}

	return inner, nil
}

func handleNullNode(key string, optional bool) error {
	if optional {
		return nil
	}

	return formatProblematicKeyError(key, ErrNullJSON)
}

func formatProblematicKeyError(key string, baseErr error) error {
	return fmt.Errorf("problematic key: %v %w", key, baseErr)
}
