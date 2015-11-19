/*

libxml2 package an interface to libxml2, providing XML and HTML parsers
with DOM interface. The inspiration is Perl5's XML::LibXML module.

This library is still in very early stages of development. API may still change
without notice.

For the time being, the API is being written so that thye are as close as we
can get to DOM Layer 3, but some methods will, for the time being, be punted
and aliases for simpler methods that don't necessarily check for the DOM's
correctness will be used.

For example, `AppendChild()` must perform a lot of checks before returning
successfully, but as of this writing it's just an alias for `xmlAddChild()`
which does lots of... interesting things if you're not careful.

Also, the return values are still shaky -- I'm still debating how to handle error cases gracefully.

*/
package libxml2

/*
#cgo pkg-config: libxml-2.0
#include "libxml/xmlerror.h"

static inline void MY_nilErrorHandler(void *ctx, const char *msg, ...) {}

static inline void MY_xmlSilenceParseErrors() {
	xmlSetGenericErrorFunc(NULL, MY_nilErrorHandler);
}

static inline void MY_xmlDefaultParseErrors() {
	// Checked in the libxml2 source code that using NULL in the second
	// argument restores the default error handler
	xmlSetGenericErrorFunc(NULL, NULL);
}

*/
import "C"
import "io"

// ReportErrors *globally* changes the behavior of reporting errors.
// By default libxml2 spews out lots of data to stderr. When you call
// this function with a `false` value, all those messages are surpressed.
// When you call this function a `true` value, the default behavior is
// restored
func ReportErrors(b bool) {
	if b {
		C.MY_xmlDefaultParseErrors()
	} else {
		C.MY_xmlSilenceParseErrors()
	}
}

// Parse parses the given buffer and returns a Document.
func Parse(buf []byte, o ...ParseOption) (*Document, error) {
	p := NewParser(o...)
	return p.Parse(buf)
}

// ParseString parses the given string and returns a Document.
func ParseString(s string, o ...ParseOption) (*Document, error) {
	p := NewParser(o...)
	return p.ParseString(s)
}

func ParseReader(rdr io.Reader, o ...ParseOption) (*Document, error) {
	p := NewParser(o...)
	return p.ParseReader(rdr)
}
