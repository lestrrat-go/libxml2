package libxml2

import (
	"io"

	"github.com/lestrrat-go/libxml2/parser"
	"github.com/lestrrat-go/libxml2/types"
)

// Parse parses the given buffer and returns a Document.
func Parse(buf []byte, o ...parser.Option) (types.Document, error) {
	p := parser.New(o...)
	return p.Parse(buf)
}

// ParseString parses the given string and returns a Document.
func ParseString(s string, o ...parser.Option) (types.Document, error) {
	p := parser.New(o...)
	return p.ParseString(s)
}

// ParseReader parses XML from the given io.Reader and returns a Document.
func ParseReader(rdr io.Reader, o ...parser.Option) (types.Document, error) {
	p := parser.New(o...)
	return p.ParseReader(rdr)
}
