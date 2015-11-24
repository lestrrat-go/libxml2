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

const _ParseOption_name = "RecoverNoEntDTDLoadDTDAttrDTDValidNoErrorNoWarningPedanticNoBlanksSAX1XIncludeNoNetNoDictNscleanNoCDATANoXIncNodeCompactOld10NoBaseFixHugeOldSAXIgnoreEncBigLines"

var _ParseOption_map = map[int]string{
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

func (i *ParseOption) Set(options ...ParseOption) {
	v := int(*i) // current value
	for _, o := range options {
		v = v | int(o)
	}
	*i = ParseOption(v)
}

func (opts ParseOption) String() string {
	if opts == XmlParseEmptyOption {
		return "[]"
	}

	i := int(opts)
	b := bytes.Buffer{}
	b.Write([]byte{'['})
	for x := 1; x < int(XmlParseMax); x = x << 1 {
		if (i & x) == x {
			v, ok := _ParseOption_map[x]
			if !ok {
				v = "ParseOption(Unknown)"
			}
			b.WriteString(v)
			b.Write([]byte{'|'})
		}
	}
	x := b.Bytes()
	if x[len(x)-1] == '|' {
		x[len(x)-1] = ']'
	} else {
		x = append(x, ']')
	}
	return string(x)
}

func NewParser(opts ...ParseOption) *Parser {
	var o ParseOption
	if len(opts) > 0 {
		o = opts[0]
	}
	return &Parser{
		Options: o,
	}
}

func (p *Parser) Parse(buf []byte) (*Document, error) {
	return p.ParseString(string(buf))
}

func (p *Parser) ParseString(s string) (*Document, error) {
	ctx := C.xmlCreateMemoryParserCtxt(C.CString(s), C.int(len(s)))
	if ctx == nil {
		return nil, errors.New("error creating parser")
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

func (p *Parser) ParseReader(in io.Reader) (*Document, error) {
	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(in); err != nil {
		return nil, err
	}

	return p.ParseString(buf.String())
}
