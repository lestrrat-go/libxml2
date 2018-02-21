package libxml2

import (
	"regexp"
	"testing"

	"github.com/lestrrat-go/libxml2/dom"
	"github.com/lestrrat-go/libxml2/types"

	"github.com/lestrrat-go/libxml2/clib"
	"github.com/lestrrat-go/libxml2/parser"
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

func parseShouldSucceed(t *testing.T, opts parser.Option, inputs []string) {
	t.Logf("Test parsing with parser %v", opts)
	for _, s := range inputs {
		d, err := ParseString(s, opts)
		if !assert.NoError(t, err, "Parse should succeed") {
			return
		}
		d.Free()
	}
}

func parseShouldFail(t *testing.T, opts parser.Option, inputs []string) {
	for _, s := range inputs {
		d, err := ParseString(s, opts)
		if err == nil {
			d.Free()
			t.Errorf("Expected failure to parse '%s'", s)
		}
	}
}

type ParseOptionToString struct {
	v parser.Option
	e string
}

func TestParseOptionStringer(t *testing.T) {
	values := []ParseOptionToString{
		ParseOptionToString{
			v: parser.XMLParseRecover,
			e: "Recover",
		},
		ParseOptionToString{
			v: parser.XMLParseNoEnt,
			e: "NoEnt",
		},
		ParseOptionToString{
			v: parser.XMLParseDTDLoad,
			e: "DTDLoad",
		},
		ParseOptionToString{
			v: parser.XMLParseDTDAttr,
			e: "DTDAttr",
		},
		ParseOptionToString{
			v: parser.XMLParseDTDValid,
			e: "DTDValid",
		},
		ParseOptionToString{
			v: parser.XMLParseNoError,
			e: "NoError",
		},
		ParseOptionToString{
			v: parser.XMLParseNoWarning,
			e: "NoWarning",
		},
		ParseOptionToString{
			v: parser.XMLParsePedantic,
			e: "Pedantic",
		},
		ParseOptionToString{
			v: parser.XMLParseNoBlanks,
			e: "NoBlanks",
		},
		ParseOptionToString{
			v: parser.XMLParseSAX1,
			e: "SAX1",
		},
		ParseOptionToString{
			v: parser.XMLParseXInclude,
			e: "XInclude",
		},
		ParseOptionToString{
			v: parser.XMLParseNoNet,
			e: "NoNet",
		},
		ParseOptionToString{
			v: parser.XMLParseNoDict,
			e: "NoDict",
		},
		ParseOptionToString{
			v: parser.XMLParseNsclean,
			e: "Nsclean",
		},
		ParseOptionToString{
			v: parser.XMLParseNoCDATA,
			e: "NoCDATA",
		},
		ParseOptionToString{
			v: parser.XMLParseNoXIncNode,
			e: "NoXIncNode",
		},
		ParseOptionToString{
			v: parser.XMLParseCompact,
			e: "Compact",
		},
		ParseOptionToString{
			v: parser.XMLParseOld10,
			e: "Old10",
		},
		ParseOptionToString{
			v: parser.XMLParseNoBaseFix,
			e: "NoBaseFix",
		},
		ParseOptionToString{
			v: parser.XMLParseHuge,
			e: "Huge",
		},
		ParseOptionToString{
			v: parser.XMLParseOldSAX,
			e: "OldSAX",
		},
		ParseOptionToString{
			v: parser.XMLParseIgnoreEnc,
			e: "IgnoreEnc",
		},
		ParseOptionToString{
			v: parser.XMLParseBigLines,
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
	clib.ReportErrors(false)
	defer clib.ReportErrors(true)

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
		parseShouldSucceed(t, parser.XMLParseNoBlanks, input)
	}
}

func TestRoundtripNoBlanks(t *testing.T) {
	doc, err := ParseString(`<a>    <b/> </a>`, parser.XMLParseNoBlanks)
	if err != nil {
		t.Errorf("failed to parse string: %s", err)
		return
	}

	if !assert.Regexp(t, regexp.MustCompile(`<a><b/></a>`), doc.Dump(false), "stringified xml should have no blanks") {
		return
	}
}

func TestOptionsShouldCombine(t *testing.T) {
	opts := map[parser.Option][]parser.Option{
		parser.Option(64): []parser.Option{parser.XMLParseNoWarning},
		parser.Option(96): []parser.Option{parser.XMLParseNoWarning, parser.XMLParseNoError},
	}

	for expected, options := range opts {
		p := parser.New(options...)
		assert.Equal(t, expected, p.Options)
	}
}

func TestGHIssue23(t *testing.T) {
	const src = `<?xml version=1.0?>
<rootnode>
    <greeting>Hello</greeting>
    <goodbye>Goodbye!</goodbye>
</rootnode>`

	doc, err := ParseString(src, parser.XMLParseRecover, parser.XMLParseNoWarning, parser.XMLParseNoError)
	if !assert.NoError(t, err, "should pass") {
		return
	}
	doc.Free()
}

func TestCommentWrapNodeIssue(t *testing.T) {

	// should wrap comment node
	const testHTML = "<p><!-- test --></p><!-- test --><p><!-- test --></p>"

	doc, err := ParseHTMLString(testHTML, parser.HTMLParseRecover)
	if err != nil {
		t.Fatalf("Got error when parsing HTML: %v", err)
	}

	bodyRes, err := doc.Find("//body")
	if err != nil {
		t.Fatalf("Got error when grabbing body: %v", err)
	}

	bodyChildren, err := bodyRes.NodeList().First().ChildNodes()
	if err != nil {
		t.Fatalf("Got error when grabbing body's children: %v", err)
	}

	if str := bodyChildren.String(); str != testHTML {
		t.Fatalf("HTML did not convert back correctly, expected: %v, got: %v.", testHTML, str)
	}
}

func TestPiWrapNodeIssue(t *testing.T) {

	// should wrap Pi node
	const textXML = "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n<a>test <?test?></a>\n"
	doc, err := ParseString(textXML)
	if err != nil {
		t.Fatalf("Got error when parsing xml: %v", err)
	}

	nodes, err := doc.ChildNodes()
	if err != nil {
		t.Fatalf("Got error when getting childnodes: %v", err)
	}

	for _, node := range nodes {
		if node.HasChildNodes() {
			if _, err := node.ChildNodes(); err != nil {
				t.Fatalf("Got error when getting childnodes of childnodes: %v", err)
			}
		}
	}

	if str := doc.String(); str != textXML {
		t.Fatalf("XML did not convert back correctly, expected: %v, got: %v", textXML, str)
	}
}

func TestGetNonexistentAttributeReturnsRecoverableError(t *testing.T) {
	const src = `<?xml version="1.0"?><rootnode/>`
	doc, err := ParseString(src)
	if !assert.NoError(t, err, "Should parse") {
		return
	}
	defer doc.Free()

	rootNode, err := doc.DocumentElement()
	if !assert.NoError(t, err, "Should find root element") {
		return
	}

	el, ok := rootNode.(types.Element)
	if !ok {
		t.Fatalf("Root node was not an element")
	}

	_, err = el.GetAttribute("non-existant")
	if err != dom.ErrAttributeNotFound {
		t.Fatalf("GetAttribute() error not comparable to existing library")
	}
}
