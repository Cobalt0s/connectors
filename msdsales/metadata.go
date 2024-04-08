package msdsales

import (
	"context"
	"fmt"

	"github.com/amp-labs/connectors/common"
)

func (c *Connector) ListObjectMetadata(
	ctx context.Context, objectNames []string,
) (*common.ListObjectMetadataResult, error) {

	rsp, err := c.getXML(ctx, c.getURL("$metadata"))
	if err != nil {
		return nil, err
	}

	fmt.Println(rsp.Code)

	uniqueNames := make(map[string]bool)

	for _, child := range rsp.Body.Root.Children {
		if child.Name == "DataServices" {
			for _, node := range child.Children {
				if node.Name == "Schema" {
					for _, n := range node.Children {
						if n.Name == "EntityType" {
							for _, ch := range n.Children {
								uniqueNames[ch.Name] = false
							}
						} else {
							fmt.Println("") // Complex Type / EnumType / Annotation / Parameter / ReturnType
						}
					}
				}
			}
			fmt.Println("")
		}
	}
	output := ""
	for k := range uniqueNames {
		output += k + "\n"
	}
	fmt.Println(output)

	return &common.ListObjectMetadataResult{
		Result: map[string]common.ObjectMetadata{},
		Errors: nil,
	}, nil
}
