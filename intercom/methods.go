package intercom

import (
	"context"
	"errors"
	"log/slog"

	"github.com/amp-labs/connectors/common"
)

// TODO these methods look repetitive.
func (c *Connector) get(ctx context.Context,
	url string, headers ...common.Header,
) (*common.JSONHTTPResponse, error) {
	rsp, err := c.Client.Get(ctx, url, headers...)
	if err = handleError(err); err != nil {
		return nil, err
	}

	return rsp, nil
}

func handleError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, common.ErrAccessToken):
		slog.Warn("Access token invalid, retrying", "error", err)

		fallthrough
	case errors.Is(err, common.ErrRetryable):
		fallthrough
	case errors.Is(err, common.ErrApiDisabled):
		fallthrough
	case errors.Is(err, common.ErrForbidden):
		fallthrough
	default:
		// Anything else is a permanent error
		return err
	}
}