package libxml2

/*
#cgo pkg-config: libxml-2.0
#include "libxml/parser.h"
#include "libxml/parserInternals.h"
*/
import "C"
import (
	"bytes"
	"errors"
	"io"
)

const (
	XmlParseRecover    = 1 << iota /* recover on errors */
	XmlParseNoEnt                  /* substitute entities */
	XmlParseDTDLoad                /* load the external subset */
	XmlParseDTDAttr                /* default DTD attributes */
	XmlParseDTDValid               /* validate with the DTD */
	XmlParseNoError                /* suppress error reports */
	XmlParseNoWarning              /* suppress warning reports */
	XmlParsePedantic               /* pedantic error reporting */
	XmlParseNoBlanks               /* remove blank nodes */
	XmlParseSAX1                   /* use the SAX1 interface internally */
	XmlParseXInclude               /* Implement XInclude substitition  */
	XmlParseNoNet                  /* Forbid network access */
	XmlParseNoDict                 /* Do not reuse the context dictionnary */
	XmlParseNsclean                /* remove redundant namespaces declarations */
	XmlParseNoCDATA                /* merge CDATA as text nodes */
	XmlParseNoXIncNode             /* do not generate XINCLUDE START/END nodes */
	XmlParseCompact                /* compact small text nodes; no modification of the tree allowed afterwards (will possibly crash if you try to modify the tree) */
	XmlParseOld10                  /* parse using XML-1.0 before update 5 */
	XmlParseNoBaseFix              /* do not fixup XINCLUDE xml:base uris */
	XmlParseHuge                   /* relax any hardcoded limit from the parser */
	XmlParseOldSAX                 /* parse using SAX2 interface before 2.7.0 */
	XmlParseIgnoreEnc              /* ignore internal document encoding hint */
	XmlParseBigLines               /* Store big lines numbers in text PSVI field */
)

type Parser struct {
	Options int
}

func (p *Parser) ParseString(s string) (*Document, error) {
	ctx := C.xmlCreateMemoryParserCtxt(C.CString(s), C.int(len(s)))
	if ctx == nil {
		return nil, errors.New("error createing parser")
	}
	defer C.xmlFreeParserCtxt(ctx)

	C.xmlCtxtUseOptions(ctx, C.int(p.Options))
	C.xmlParseDocument(ctx)

	if ctx.wellFormed == C.int(0) {
		return nil, errors.New("malformed XML")
	}

	doc := ctx.myDoc
	if doc == nil {
		return nil, errors.New("parse failed")
	}
	return wrapDocument(doc), nil
}

func (p *Parser) Parse(in io.Reader) (*Document, error) {
	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(in); err != nil {
		return nil, err
	}

	return p.ParseString(buf.String())
}
