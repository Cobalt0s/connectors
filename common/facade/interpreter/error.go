package interpreter

import (
	"errors"
	"fmt"
	"github.com/amp-labs/connectors/common"
	"mime"
	"net/http"
)

var (
	ErrUnmarshal = errors.New("unmarshal failed")
)

// FaultyResponseHandler used to parse erroneous response
type FaultyResponseHandler func(res *http.Response, body []byte) error

// ErrorHandler invokes a function that is matching a certain response media type to parse error, ex: JSON
// otherwise defaults to general error interpretation
type ErrorHandler struct {
	JSON FaultyResponseHandler
}

func (h ErrorHandler) Handle(res *http.Response, body []byte) error {
	mediaType, _, err := mime.ParseMediaType(res.Header.Get("Content-Type"))
	if err != nil {
		return fmt.Errorf("mime.ParseMediaType failed: %w", err)
	}

	if h.JSON != nil && mediaType == "application/json" {
		return h.JSON(res, body)
	}

	return common.InterpretError(res, body)
}
