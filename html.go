package libxml2

import (
	"bytes"
	"io"
)

const (
	HTMLParserRecover = 1 << 0
	HTMLParserNoError = 1<<iota + 5
	HTMLParserNoWarning
	HTMLParserPedantic
	HTMLParserNoBlanks
	HTMLParserNoNet
	HTMLParserCompact
)

const DefaultHTMLParserFlags = HTMLParserCompact | HTMLParserNoBlanks | HTMLParserNoError | HTMLParserNoWarning

func ParseHTML(content []byte) (*Document, error) {
	return ParseHTMLString(string(content))
}

func ParseHTMLString(content string) (*Document, error) {
	return htmlReadDoc(content, "", "", DefaultHTMLParserFlags)
}

func ParseHTMLReader(in io.Reader) (*Document, error) {
	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(in); err != nil {
		return nil, err
	}

	return ParseHTMLString(buf.String())
}
