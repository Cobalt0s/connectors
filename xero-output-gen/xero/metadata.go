package xero

import (
	"context"

	"github.com/amp-labs/connectors/common"
	"github.com/amp-labs/connectors/intercom/metadata"
)

func (c *Connector) ListObjectMetadata(
	ctx context.Context, objectNames []string,
) (*common.ListObjectMetadataResult, error) {
	// Ensure that objectNames is not empty
	if len(objectNames) == 0 {
		return nil, common.ErrMissingObjects
	}

	// TODO example below is serving Metadata from static files
	// TODO if connector API has schema discoverability then make requests, process with jsonquery
	schemas, err := metadata.FileManager.LoadSchemas()
	if err != nil {
		return nil, common.ErrMetadataLoadFailure
	}

	return schemas.Select(objectNames)
}
