package msdsales

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"strings"

	"github.com/amp-labs/connectors/common"
)

// TODO MS Sales allows fine control of reading
// TODO basic read
// $select = ReadParams{fields} [COMPLETED]
// $expand = nested response
// $orderby = list of fields with asc/desc keyword
// TODO batch
// $apply = batch
// TODO search
// $filter = query functions, comparisons
// TODO pagination
// $top = <int> of entries to return (ignored if header <Prefer: odata.maxpagesize>)
// $count = counts all existing rows (@odata.count)

var (
	AnnotationsHeader = common.Header{
		Key:   "Prefer",
		Value: "odata.include-annotations=\"*\"", // TODO we can specify which annotations to include
	}
)

func newPaginationHeader(pageSize int) common.Header {
	return common.Header{
		Key:   "Prefer",
		Value: fmt.Sprintf("odata.maxpagesize=%v", pageSize),
	}
}

func (c *Connector) Read(ctx context.Context, config common.ReadParams) (*common.ReadResult, error) {
	var fullURL string

	if len(config.NextPage) == 0 {
		// First page
		relativeURL := config.ObjectName + makeQueryValues(config)
		fullURL = c.getURL(relativeURL)
	} else {
		// Next page
		fullURL = config.NextPage
	}
	// TODO given that one of the fields is annotation we can automatically add annotation header (how the hell the end user gonna know about the names of those fields)
	rsp, err := c.get(ctx, fullURL, newPaginationHeader(resolvePageSize(config)), AnnotationsHeader)
	if err != nil {
		return nil, err
	}

	fmt.Println(rsp.Body.String())

	return common.ParseResult(
		rsp,
		getTotalSize,
		getRecords,
		getNextRecordsURL,
		getMarshaledData,
		config.Fields,
	)
}

// TODO this must be tested very well, must follow MS query syntax
func makeQueryValues(config common.ReadParams) string {
	queryValues := url.Values{}

	if len(config.Fields) != 0 {
		queryValues.Add("$select", strings.Join(config.Fields, ","))
	}

	result := queryValues.Encode()
	if len(result) != 0 {
		// FIXME this is a hack. net/url encodes $. But we rely heavily on it
		// same problem with "@" ex: @Microsoft.Dynamics.CRM.totalrecordcountlimitexceeded
		// @ symbol is removed
		for before, after := range map[string]string{
			"%24select": "$select",
		} {
			result = strings.Replace(result, before, after, -1)
		}

		result = strings.Replace(result, "%40", "@", -1)
		result = strings.Replace(result, "%2C", ",", -1)

		return "?" + result
	}

	return result
}

func resolvePageSize(config common.ReadParams) int {
	return int(math.Min(MaxPageSize, float64(config.PageSize)))
}
