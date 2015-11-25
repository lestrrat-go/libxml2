package libxml2

import (
	"bytes"
	"io"
)

const (
	HtmlParseRecover = 1 << 0
	HtmlParseNoError = 1<<iota + 5
	HtmlParseNoWarning
	HtmlParsePedantic
	HtmlParseNoBlanks
	HtmlParseNoNet
	HtmlParseCompact
)

const DefaultHtmlParseFlags = HtmlParseCompact | HtmlParseNoBlanks | HtmlParseNoError | HtmlParseNoWarning

func ParseHTML(content []byte) (*Document, error) {
	return ParseHTMLString(string(content))
}

func ParseHTMLString(content string) (*Document, error) {
	return htmlReadDoc(content, "", "", DefaultHtmlParseFlags)
}

func ParseHTMLReader(in io.Reader) (*Document, error) {
	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(in); err != nil {
		return nil, err
	}

	return ParseHTMLString(buf.String())
}
