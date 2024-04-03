package msdsales

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/amp-labs/connectors/common"
	"net/http"
)

func (c *Connector) interpretJSONError(res *http.Response, body []byte) error {
	apiError := &SalesError{}
	if err := json.Unmarshal(body, &apiError); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %w", err)
	}

	switch res.StatusCode {
	case http.StatusBadRequest:
		return errors.Join(common.ErrBadRequest, apiError)
	default:
		return common.InterpretError(res, body)
	}
}

type SalesError struct {
	// TODO ...
	SomeUsefulDetails string
}

func (s SalesError) Error() string {
	return s.SomeUsefulDetails
}
