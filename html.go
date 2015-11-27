package libxml2

import (
	"bytes"
	"io"
)

func ParseHTML(content []byte, options ...HTMLParseOption) (*Document, error) {
	return ParseHTMLString(string(content), options...)
}

func ParseHTMLString(content string, options ...HTMLParseOption) (*Document, error) {
	var option HTMLParseOption
	if len(options) > 0 {
		option = options[0]
	} else {
		option = DefaultHTMLParseOptions
	}
	return htmlReadDoc(content, "", "", int(option))
}

func ParseHTMLReader(in io.Reader, options ...HTMLParseOption) (*Document, error) {
	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(in); err != nil {
		return nil, err
	}

	return ParseHTMLString(buf.String(), options...)
}
