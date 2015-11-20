package libxml2

/*
#cgo pkg-config: libxml-2.0
#include <libxml/HTMLparser.h>
#include <libxml/HTMLtree.h>
*/
import "C"
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


func htmlReadDoc(content, url, encoding string, opts int) *C.xmlDoc {
	return C.htmlReadDoc(
		C.xmlCharStrdup(C.CString(content)),
		C.CString(url),
		C.CString(encoding),
		C.int(opts),
	)
}

func ParseHTML(content []byte) (*Document, error) {
	return ParseHTMLString(string(content))
}

func ParseHTMLString(content string) (*Document, error) {
	d := htmlReadDoc(content, "", "", DefaultHtmlParseFlags)
	return &Document{ptr: d}, nil
}

func ParseHTMLReader(in io.Reader) (*Document, error) {
	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(in); err != nil {
		return nil, err
	}

	return ParseHTMLString(buf.String())
}

