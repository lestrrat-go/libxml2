package parser

import "github.com/lestrrat-go/option"

// Non-standard (i.e. non-libxml2 native options) go here

type optkeyWithEncoding struct {}

type nonNativeXMLParseOption struct {
	option.Interface
	enabled bool
}

func (*nonNativeXMLParseOption) xmlParseOption() {}
func (o *nonNativeXMLParseOption) Enabled() bool {
	return o.enabled
}

type nonNativeHTMLParseOption struct {
	option.Interface
	enabled bool
}

func (*nonNativeHTMLParseOption) htmlParseOption() {}
func (o *nonNativeHTMLParseOption) Enabled() bool {
	return o.enabled
}

// Specifies the encoding when parsing documents.
func WithXMLEncoding(s string) XMLParseOption {
	return &nonNativeXMLParseOption{
		Interface: option.New(optkeyWithEncoding{}, s),
	}
}

// Specifies the encoding when parsing documents.
func WithHTMLEncoding(s string) HTMLParseOption {
	return &nonNativeHTMLParseOption{
		Interface: option.New(optkeyWithEncoding{}, s),
	}
}

