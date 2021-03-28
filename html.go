package libxml2

import (
	"io"

	"github.com/lestrrat-go/libxml2/parser"
	"github.com/lestrrat-go/libxml2/types"
)

// ParseHTML parses an HTML document. You can omit the options
// argument, or you can provide one bitwise-or'ed option
func ParseHTML(content []byte, options ...parser.HTMLParseOption) (types.Document, error) {
	return parser.ParseHTML(content, options...)
}

// ParseHTMLString parses an HTML document. You can omit the options
// argument, or you can provide one bitwise-or'ed option
func ParseHTMLString(content string, options ...parser.HTMLParseOption) (types.Document, error) {
	return parser.ParseHTMLString(content, options...)
}

// ParseHTMLReader parses an HTML document. You can omit the options
// argument, or you can provide one bitwise-or'ed option
func ParseHTMLReader(in io.Reader, options ...parser.HTMLParseOption) (types.Document, error) {
	return parser.ParseHTMLReader(in, options...)
}
