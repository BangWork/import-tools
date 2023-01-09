package xml

import (
	"strings"

	"github.com/beevik/etree"
	"github.com/juju/errors"

	"github.com/bangwork/import-tools/serve/services/importer/resolve"
)

func NextElementFromReader(reader *resolve.XmlScanner) (*etree.Element, error) {
	line, err := nextLineFromReader(reader)
	if err != nil {
		return nil, err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, nil
	}

	document := etree.NewDocument()
	if err := document.ReadFromString(line); err != nil {
		return nil, errors.Errorf("parse element fail: %s", line)
	}
	return document.Root(), nil
}

func GetAttributeValue(element *etree.Element, attribute string) string {
	var e = element
	var name = attribute
	var resp string
	a := e.SelectAttr(name)
	if a == nil {
		child := e.SelectElement(name)
		if child != nil {
			resp = child.Text()
		}
	} else {
		resp = a.Value
	}
	return resp
}

func nextLineFromReader(reader *resolve.XmlScanner) (string, error) {
	e := reader.NextElement()
	if e == nil {
		return "", nil
	}
	return e.Encode(), nil
}
