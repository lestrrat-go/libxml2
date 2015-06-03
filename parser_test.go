package libxml2

import "testing"

func TestParseEmpty(t *testing.T) {
	p := &Parser{}
	doc, err := p.ParseString(``)
	if err == nil {
		t.Errorf("Parse of empty string should fail")
		defer doc.Free()
	}
}

const stdXmlDecl = `<?xml version="1.0"?>` + "\n"

func TestParseWFStrings(t *testing.T) {
	inputs := []string{
		`<foobar/>`,
		`<foobar></foobar>`,
		`<foobar></foobar>`,
		`<?xml version="1.0" encoding="UTF-8"?>` + "\n" + `<foobar></foobar>`,
		`<?xml version="1.0" encoding="ISO-8859-1"?>` + "\n" + `<foobar></foobar>`,
		stdXmlDecl + `<foobar> </foobar>` + "\n",
		stdXmlDecl + `<foobar><foo/></foobar> `,
		stdXmlDecl + `<foobar> <foo/> </foobar> `,
		stdXmlDecl + `<foobar><![CDATA[<>&"\` + "`" + `]]></foobar>`,
		stdXmlDecl + `<foobar>&lt;&gt;&amp;&quot;&apos;</foobar>`,
		stdXmlDecl + `<foobar>&#x20;&#160;</foobar>`,
		stdXmlDecl + `<!--comment--><foobar>foo</foobar>`,
		stdXmlDecl + `<foobar>foo</foobar><!--comment-->`,
		stdXmlDecl + `<foobar>foo<!----></foobar>`,
		stdXmlDecl + `<foobar foo="bar"/>`,
		stdXmlDecl + `<foobar foo="\` + "`" + `bar>"/>`,
	}

	p := &Parser{}
	for _, s := range inputs {
		if _, err := p.ParseString(s); err != nil {
			t.Errorf("Failed to parse '%s': %s", s, err)
		}
	}
}

func TestParseBadWFStrings(t *testing.T) {
	inputs := []string{
		"",                                      // totally empty document
		stdXmlDecl,                              // only XML Declaration
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
		stdXmlDecl + `<!DOCTYPE foobar [<!ENTITY foo "bar">]><foobar &foo;="ouch"/>`,          // bad placed entity
		stdXmlDecl + `<!DOCTYPE foobar [<!ENTITY foo "bar=&quot;foo&quot;">]><foobar &foo;/>`, // even worse
		"<ouch><!---></ouch>",   // bad comment
		"<ouch><!-----></ouch>", // bad either... (is this conform with the spec????)
	}

	p := &Parser{}
	for _, s := range inputs {
		if _, err := p.ParseString(s); err == nil {
			t.Errorf("Expected failure to parse '%s'", s)
		}
	}
}