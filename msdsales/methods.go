package msdsales

import (
	"context"
	"errors"
	"github.com/amp-labs/connectors/common"
	"log/slog"
)

// get reads data from MS Sales. It handles retries and access token refreshes.
func (c *Connector) get(ctx context.Context, url string) (*common.JSONHTTPResponse, error) {
	retry := c.RetryStrategy.Start()
	for {
		rsp, err := c.Client.Get(ctx, url)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrAccessToken):
				slog.Warn("Access token invalid, retrying", "error", err)
				fallthrough
			case errors.Is(err, common.ErrRetryable):
				if retry.Completed() {
					return nil, err
				} else {
					continue
				}
			case errors.Is(err, common.ErrApiDisabled):
				fallthrough
			case errors.Is(err, common.ErrForbidden):
				fallthrough
			default:
				// Anything else is a permanent error
				return nil, err
			}
		}

		// Success
		return rsp, nil
	}

}

// post writes data to Hubspot. It handles retries and access token refreshes.
func (c *Connector) post(ctx context.Context, url string, body any) (*common.JSONHTTPResponse, error) {
	retry := c.RetryStrategy.Start()
	for {
		rsp, err := c.Client.Post(ctx, url, body)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrAccessToken):
				slog.Warn("Access token invalid, retrying", "error", err)
				fallthrough
			case errors.Is(err, common.ErrRetryable):
				if retry.Completed() {
					return nil, err
				} else {
					continue
				}
			case errors.Is(err, common.ErrApiDisabled):
				fallthrough
			case errors.Is(err, common.ErrForbidden):
				fallthrough
			default:
				// Anything else is a permanent error
				return nil, err
			}
		}

		// Success
		return rsp, nil
	}
}
