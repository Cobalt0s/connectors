package interpreter

import (
	"fmt"
	"github.com/amp-labs/connectors/common"
	"mime"
	"net/http"
)

// FaultyResponseHandler used to parse erroneous response
type FaultyResponseHandler func(res *http.Response, body []byte) error

// Error interpreter invokes a function that is matching a certain response media type to parse error, ex: JSON
// otherwise defaults to general error interpretation
type Error struct {
	JSON FaultyResponseHandler
}

func (h Error) Handle(res *http.Response, body []byte) error {
	mediaType, _, err := mime.ParseMediaType(res.Header.Get("Content-Type"))
	if err != nil {
		return fmt.Errorf("mime.ParseMediaType failed: %w", err)
	}

	if h.JSON != nil && mediaType == "application/json" {
		return h.JSON(res, body)
	}

	return common.InterpretError(res, body)
}
