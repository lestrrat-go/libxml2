package libxml2

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

const stdXMLDecl = `<?xml version="1.0"?>` + "\n"

var (
	goodWFNSStrings = []string{
		stdXMLDecl + `<foobar xmlns:bar="xml://foo" bar:foo="bar"/>` + "\n",
		stdXMLDecl + `<foobar xmlns="xml://foo" foo="bar"><foo/></foobar>` + "\n",
		stdXMLDecl + `<bar:foobar xmlns:bar="xml://foo" foo="bar"><foo/></bar:foobar>` + "\n",
		stdXMLDecl + `<bar:foobar xmlns:bar="xml://foo" foo="bar"><bar:foo/></bar:foobar>` + "\n",
		stdXMLDecl + `<bar:foobar xmlns:bar="xml://foo" bar:foo="bar"><bar:foo/></bar:foobar>` + "\n",
	}
	goodWFStrings = []string{
		`<foobar/>`,
		`<foobar></foobar>`,
		`<foobar></foobar>`,
		`<?xml version="1.0" encoding="UTF-8"?>` + "\n" + `<foobar></foobar>`,
		`<?xml version="1.0" encoding="ISO-8859-1"?>` + "\n" + `<foobar></foobar>`,
		stdXMLDecl + `<foobar> </foobar>` + "\n",
		stdXMLDecl + `<foobar><foo/></foobar> `,
		stdXMLDecl + `<foobar> <foo/> </foobar> `,
		stdXMLDecl + `<foobar><![CDATA[<>&"\` + "`" + `]]></foobar>`,
		stdXMLDecl + `<foobar>&lt;&gt;&amp;&quot;&apos;</foobar>`,
		stdXMLDecl + `<foobar>&#x20;&#160;</foobar>`,
		stdXMLDecl + `<!--comment--><foobar>foo</foobar>`,
		stdXMLDecl + `<foobar>foo</foobar><!--comment-->`,
		stdXMLDecl + `<foobar>foo<!----></foobar>`,
		stdXMLDecl + `<foobar foo="bar"/>`,
		stdXMLDecl + `<foobar foo="\` + "`" + `bar>"/>`,
	}
	goodWFDTDStrings = []string{
		stdXMLDecl + `<!DOCTYPE foobar [` + "\n" + `<!ENTITY foo " test ">` + "\n" + `]>` + "\n" + `<foobar>&foo;</foobar>`,
		stdXMLDecl + `<!DOCTYPE foobar [<!ENTITY foo "bar">]><foobar>&foo;</foobar>`,
		stdXMLDecl + `<!DOCTYPE foobar [<!ENTITY foo "bar">]><foobar>&foo;&gt;</foobar>`,
		stdXMLDecl + `<!DOCTYPE foobar [<!ENTITY foo "bar=&quot;foo&quot;">]><foobar>&foo;&gt;</foobar>`,
		stdXMLDecl + `<!DOCTYPE foobar [<!ENTITY foo "bar">]><foobar>&foo;&gt;</foobar>`,
		stdXMLDecl + `<!DOCTYPE foobar [<!ENTITY foo "bar">]><foobar foo="&foo;"/>`,
		stdXMLDecl + `<!DOCTYPE foobar [<!ENTITY foo "bar">]><foobar foo="&gt;&foo;"/>`,
	}
	badWFStrings = []string{
		"",                                      // totally empty document
		stdXMLDecl,                              // only XML Declaration
		"<!--ouch-->",                           // comment only is like an empty document
		`<!DOCTYPE ouch [<!ENTITY foo "bar">]>`, // no good either ...
		"<ouch>",                // single tag (tag mismatch)
		"<ouch/>foo",            // trailing junk
		"foo<ouch/>",            // leading junk
		"<ouch foo=bar/>",       // bad attribute
		`<ouch foo="bar/>`,      // bad attribute
		"<ouch>&</ouch>",        // bad char
		`<ouch>&//0x20;</ouch>`, // bad chart
		"<foob<e4>r/>",          // bad encoding
		"<ouch>&foo;</ouch>",    // undefind entity
		"<ouch>&gt</ouch>",      // unterminated entity
		stdXMLDecl + `<!DOCTYPE foobar [<!ENTITY foo "bar">]><foobar &foo;="ouch"/>`,          // bad placed entity
		stdXMLDecl + `<!DOCTYPE foobar [<!ENTITY foo "bar=&quot;foo&quot;">]><foobar &foo;/>`, // even worse
		"<ouch><!---></ouch>",   // bad comment
		"<ouch><!-----></ouch>", // bad either... (is this conform with the spec????)
	}
)

func parseShouldSucceed(t *testing.T, opts ParseOption, inputs []string) {
	t.Logf("Test parsing with parser %v", opts)
	for _, s := range inputs {
		d, err := ParseString(s, opts)
		if !assert.NoError(t, err, "Parse should succeed") {
			return
		}
		d.Free()
	}
}

func parseShouldFail(t *testing.T, opts ParseOption, inputs []string) {
	for _, s := range inputs {
		d, err := ParseString(s, opts)
		if err == nil {
			d.Free()
			t.Errorf("Expected failure to parse '%s'", s)
		}
	}
}

type ParseOptionToString struct {
	v ParseOption
	e string
}

func TestParseOptionStringer(t *testing.T) {
	values := []ParseOptionToString{
		ParseOptionToString{
			v: XMLParseRecover,
			e: "Recover",
		},
		ParseOptionToString{
			v: XMLParseNoEnt,
			e: "NoEnt",
		},
		ParseOptionToString{
			v: XMLParseDTDLoad,
			e: "DTDLoad",
		},
		ParseOptionToString{
			v: XMLParseDTDAttr,
			e: "DTDAttr",
		},
		ParseOptionToString{
			v: XMLParseDTDValid,
			e: "DTDValid",
		},
		ParseOptionToString{
			v: XMLParseNoError,
			e: "NoError",
		},
		ParseOptionToString{
			v: XMLParseNoWarning,
			e: "NoWarning",
		},
		ParseOptionToString{
			v: XMLParsePedantic,
			e: "Pedantic",
		},
		ParseOptionToString{
			v: XMLParseNoBlanks,
			e: "NoBlanks",
		},
		ParseOptionToString{
			v: XMLParseSAX1,
			e: "SAX1",
		},
		ParseOptionToString{
			v: XMLParseXInclude,
			e: "XInclude",
		},
		ParseOptionToString{
			v: XMLParseNoNet,
			e: "NoNet",
		},
		ParseOptionToString{
			v: XMLParseNoDict,
			e: "NoDict",
		},
		ParseOptionToString{
			v: XMLParseNsclean,
			e: "Nsclean",
		},
		ParseOptionToString{
			v: XMLParseNoCDATA,
			e: "NoCDATA",
		},
		ParseOptionToString{
			v: XMLParseNoXIncNode,
			e: "NoXIncNode",
		},
		ParseOptionToString{
			v: XMLParseCompact,
			e: "Compact",
		},
		ParseOptionToString{
			v: XMLParseOld10,
			e: "Old10",
		},
		ParseOptionToString{
			v: XMLParseNoBaseFix,
			e: "NoBaseFix",
		},
		ParseOptionToString{
			v: XMLParseHuge,
			e: "Huge",
		},
		ParseOptionToString{
			v: XMLParseOldSAX,
			e: "OldSAX",
		},
		ParseOptionToString{
			v: XMLParseIgnoreEnc,
			e: "IgnoreEnc",
		},
		ParseOptionToString{
			v: XMLParseBigLines,
			e: "BigLines",
		},
	}

	for _, d := range values {
		if d.v.String() != "["+d.e+"]" {
			t.Errorf("e '%s', got '%s'", d.e, d.v.String())
		}
	}
}

func TestParseEmpty(t *testing.T) {
	doc, err := ParseString(``)
	if err == nil {
		t.Errorf("Parse of empty string should fail")
		defer doc.Free()
	}
}

func TestParse(t *testing.T) {
	inputs := [][]string{
		goodWFStrings,
		goodWFNSStrings,
		goodWFDTDStrings,
	}

	for _, input := range inputs {
		parseShouldSucceed(t, 0, input)
	}
}

func TestParseBad(t *testing.T) {
	inputs := [][]string{
		badWFStrings,
	}

	for _, input := range inputs {
		parseShouldFail(t, 0, input)
	}
}

func TestParseNoBlanks(t *testing.T) {
	inputs := [][]string{
		goodWFStrings,
		goodWFNSStrings,
		goodWFDTDStrings,
	}
	for _, input := range inputs {
		parseShouldSucceed(t, XMLParseNoBlanks, input)
	}
}

func TestRoundtripNoBlanks(t *testing.T) {
	doc, err := ParseString(`<a>    <b/> </a>`, XMLParseNoBlanks)
	if err != nil {
		t.Errorf("failed to parse string: %s", err)
		return
	}

	if !assert.Regexp(t, regexp.MustCompile(`<a><b/></a>`), doc.Dump(false), "stringified xml should have no blanks") {
		return
	}
}
