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
	"fmt"
	"io"
	"strings"
)

// ParseOption represents each of the parser option bit
type ParseOption int

// ParseOption represents the parser option bit set
type ParseOptions int

const (
	XmlParseRecover    ParseOption = 1 << iota /* recover on errors */
	XmlParseNoEnt                              /* substitute entities */
	XmlParseDTDLoad                            /* load the external subset */
	XmlParseDTDAttr                            /* default DTD attributes */
	XmlParseDTDValid                           /* validate with the DTD */
	XmlParseNoError                            /* suppress error reports */
	XmlParseNoWarning                          /* suppress warning reports */
	XmlParsePedantic                           /* pedantic error reporting */
	XmlParseNoBlanks                           /* remove blank nodes */
	XmlParseSAX1                               /* use the SAX1 interface internally */
	XmlParseXInclude                           /* Implement XInclude substitition  */
	XmlParseNoNet                              /* Forbid network access */
	XmlParseNoDict                             /* Do not reuse the context dictionnary */
	XmlParseNsclean                            /* remove redundant namespaces declarations */
	XmlParseNoCDATA                            /* merge CDATA as text nodes */
	XmlParseNoXIncNode                         /* do not generate XINCLUDE START/END nodes */
	XmlParseCompact                            /* compact small text nodes; no modification of the tree allowed afterwards (will possibly crash if you try to modify the tree) */
	XmlParseOld10                              /* parse using XML-1.0 before update 5 */
	XmlParseNoBaseFix                          /* do not fixup XINCLUDE xml:base uris */
	XmlParseHuge                               /* relax any hardcoded limit from the parser */
	XmlParseOldSAX                             /* parse using SAX2 interface before 2.7.0 */
	XmlParseIgnoreEnc                          /* ignore internal document encoding hint */
	XmlParseBigLines                           /* Store big lines numbers in text PSVI field */
	XmlParseMax
)

const _ParseOption_name = "RecoverNoEntDTDLoadDTDAttrDTDValidNoErrorNoWarningPedanticNoBlanksSAX1XIncludeNoNetNoDictNscleanNoCDATANoXIncNodeCompactOld10NoBaseFixHugeOldSAXIgnoreEncBigLines"

var _ParseOption_map = map[ParseOption]string{
	1:       _ParseOption_name[0:7],
	2:       _ParseOption_name[7:12],
	4:       _ParseOption_name[12:19],
	8:       _ParseOption_name[19:26],
	16:      _ParseOption_name[26:34],
	32:      _ParseOption_name[34:41],
	64:      _ParseOption_name[41:50],
	128:     _ParseOption_name[50:58],
	256:     _ParseOption_name[58:66],
	512:     _ParseOption_name[66:70],
	1024:    _ParseOption_name[70:78],
	2048:    _ParseOption_name[78:83],
	4096:    _ParseOption_name[83:89],
	8192:    _ParseOption_name[89:96],
	16384:   _ParseOption_name[96:103],
	32768:   _ParseOption_name[103:113],
	65536:   _ParseOption_name[113:120],
	131072:  _ParseOption_name[120:125],
	262144:  _ParseOption_name[125:134],
	524288:  _ParseOption_name[134:138],
	1048576: _ParseOption_name[138:144],
	2097152: _ParseOption_name[144:153],
	4194304: _ParseOption_name[153:161],
}

func (i ParseOption) String() string {
	if str, ok := _ParseOption_map[i]; ok {
		return str
	}
	return fmt.Sprintf("ParseOption(%d)", i)
}

func (i *ParseOptions) Set(options ...ParseOption) {
	v := int(*i) // current value
	for _, o := range options {
		v = v | int(o)
	}
	*i = ParseOptions(v)
}

func (i ParseOptions) String() string {
	if int(i) == 0 {
		return "[]"
	}

	list := make([]string, 0, 24)
	for x := 1; x < int(XmlParseMax); x = x << 1 {
		if (int(i) & x) == x {
			list = append(list, ParseOption(x).String())
		}
	}

	return "[ " + strings.Join(list, " | ") + " ]"
}

type Parser struct {
	Options ParseOptions
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
