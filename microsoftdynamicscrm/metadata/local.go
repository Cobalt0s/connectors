package metadata

import (
	"bytes"
	_ "embed"
	"errors"
	"github.com/subchen/go-xmldom"
)

//go:embed ms-crm-metadata.xml
var file []byte

func Load() (*xmldom.Node, error) {
	xmlBody, err := xmldom.Parse(bytes.NewReader(file))
	if err != nil {
		return nil, err
	}
	if xmlBody.Root == nil {
		return nil, errors.New("no root")
	}

	return xmlBody.Root, nil
}
