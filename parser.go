package libxml2

import (
	"bytes"
	"io"
)

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

// ParseReader parses XML from the given io.Reader and returns a Document.
func ParseReader(rdr io.Reader, o ...ParseOption) (*Document, error) {
	p := NewParser(o...)
	return p.ParseReader(rdr)
}

const _ParseOptionName = "RecoverNoEntDTDLoadDTDAttrDTDValidNoErrorNoWarningPedanticNoBlanksSAX1XIncludeNoNetNoDictNscleanNoCDATANoXIncNodeCompactOld10NoBaseFixHugeOldSAXIgnoreEncBigLines"

var _ParseOptionMap = map[int]string{
	1:       _ParseOptionName[0:7],
	2:       _ParseOptionName[7:12],
	4:       _ParseOptionName[12:19],
	8:       _ParseOptionName[19:26],
	16:      _ParseOptionName[26:34],
	32:      _ParseOptionName[34:41],
	64:      _ParseOptionName[41:50],
	128:     _ParseOptionName[50:58],
	256:     _ParseOptionName[58:66],
	512:     _ParseOptionName[66:70],
	1024:    _ParseOptionName[70:78],
	2048:    _ParseOptionName[78:83],
	4096:    _ParseOptionName[83:89],
	8192:    _ParseOptionName[89:96],
	16384:   _ParseOptionName[96:103],
	32768:   _ParseOptionName[103:113],
	65536:   _ParseOptionName[113:120],
	131072:  _ParseOptionName[120:125],
	262144:  _ParseOptionName[125:134],
	524288:  _ParseOptionName[134:138],
	1048576: _ParseOptionName[138:144],
	2097152: _ParseOptionName[144:153],
	4194304: _ParseOptionName[153:161],
}

// Set flips the option bit in the given ParseOption
func (o *ParseOption) Set(options ...ParseOption) {
	v := int(*o) // current value
	for _, i := range options {
		v = v | int(i)
	}
	*o = ParseOption(v)
}

// String creates a string representation of the ParseOption
func (o ParseOption) String() string {
	if o == XMLParseEmptyOption {
		return "[]"
	}

	i := int(o)
	b := bytes.Buffer{}
	b.Write([]byte{'['})
	for x := 1; x < int(XMLParseMax); x = x << 1 {
		if (i & x) == x {
			v, ok := _ParseOptionMap[x]
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

// NewParser creates a new Parser with the given options.
func NewParser(opts ...ParseOption) *Parser {
	var o ParseOption
	if len(opts) > 0 {
		o = opts[0]
	}
	return &Parser{
		Options: o,
	}
}

// Parse parses XML from the given byte buffer
func (p *Parser) Parse(buf []byte) (*Document, error) {
	return p.ParseString(string(buf))
}

// ParseString parses XML from the given string
func (p *Parser) ParseString(s string) (*Document, error) {
	ctx, err := xmlCreateMemoryParserCtxt(s, p.Options)
	if err != nil {
		return nil, err
	}
	defer ctx.Free()

	if err := ctx.Parse(); err != nil {
		return nil, err
	}

	if ctx.WellFormed() {
		return nil, ErrMalformedXML
	}

	doc, err := ctx.Document()
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// ParseReader parses XML from the given io.Reader
func (p *Parser) ParseReader(in io.Reader) (*Document, error) {
	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(in); err != nil {
		return nil, err
	}

	return p.ParseString(buf.String())
}
