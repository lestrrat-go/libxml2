// This file is build-tag protected because it involves loading an external
// library (xmlpath)
package libxml2_test

import (
	"bytes"
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"

	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/dom"
	"github.com/lestrrat-go/libxml2/xpath"
	"github.com/stretchr/testify/assert"
	"gopkg.in/xmlpath.v1"
)

var xmlfile = filepath.Join("test", "feed.atom")

func BenchmarkXmlpathXmlpath(b *testing.B) {
	f, err := os.Open(xmlfile)
	if err != nil {
		b.Fatalf("%s", err)
	}

	root, err := xmlpath.Parse(f)
	if err != nil {
		b.Fatalf("%s", err)
	}
	for i := 0; i < b.N; i++ {
		p, err := xmlpath.Compile(`//entry`)
		if err != nil {
			b.Fatalf("%s", err)
		}
		it := p.Iter(root)
		for it.Next() {
			n := it.Node()
			_ = n
		}
	}
}

func TestBenchmarkLibxml2Xmlpath(t *testing.T) {
	f, err := os.Open(xmlfile)
	if !assert.NoError(t, err, "os.Open succeeds") {
		return
	}

	doc, err := libxml2.ParseReader(f)
	if !assert.NoError(t, err, "ParseReader succeeds") {
		return
	}

	xpc, err := xpath.NewContext(doc)
	if !assert.NoError(t, err, "xpath.NewContext succeeds") {
		return
	}
	xpc.RegisterNS("atom", "http://www.w3.org/2005/Atom")

	res, err := xpc.Find(`//atom:entry`)
	if !assert.NoError(t, err, "xpc.Find succeeds") {
		return
	}
	defer res.Free()

	iter := res.NodeIter()
	if !assert.NotEmpty(t, iter, "res.NodeIter succeeds") {
		return
	}

	count := 0
	for iter.Next() {
		n := iter.Node()
		if !assert.NotEmpty(t, n, "iter.Node returns something") {
			return
		}
		count++
	}
	if !assert.True(t, count > 0, "there's at least 1 node") {
		return
	}
}

func BenchmarkLibxml2Xmlpath(b *testing.B) {
	f, err := os.Open(xmlfile)
	if err != nil {
		b.Fatalf("%s", err)
	}

	doc, err := libxml2.ParseReader(f)
	if err != nil {
		b.Fatalf("%s", err)
	}

	xpc, err := xpath.NewContext(doc)
	if err != nil {
		b.Fatalf("%s", err)
	}
	xpc.RegisterNS("atom", "http://www.w3.org/2005/Atom")
	for i := 0; i < b.N; i++ {
		iter := xpath.NodeIter(xpc.Find(`//atom:entry`))
		for iter.Next() {
			n := iter.Node()
			_ = n
		}
	}
}

type Foo struct {
	XMLName xml.Name `xml:"https://github.com/lestrrat-go/libxml2/foo foo:foo"`
	Field1  string
	Field2  string `xml:",attr"`
}

func BenchmarkEncodingXMLDOM(b *testing.B) {
	var buf bytes.Buffer
	f := Foo{
		Field1: "Hello, World!",
		Field2: "Hello, Attribute!",
	}
	for i := 0; i < b.N; i++ {
		buf.Reset()
		enc := xml.NewEncoder(&buf)
		enc.Encode(f)
	}
}

func BenchmarkLibxml2DOM(b *testing.B) {
	var buf bytes.Buffer
	const nsuri = `https://github.com/lestrrat-go/libxml2/foo`
	f := Foo{
		Field1: "Hello, World!",
		Field2: "Hello, Attribute!",
	}
	for i := 0; i < b.N; i++ {
		d := dom.CreateDocument()

		root, err := d.CreateElementNS(nsuri, "foo:foo")
		if err != nil {
			d.Free()
			panic(err)
		}
		d.SetDocumentElement(root)

		f1xml, err := d.CreateElement("Field1")
		if err != nil {
			d.Free()
			panic(err)
		}
		root.AddChild(f1xml)

		f1xml.SetAttribute("Field2", f.Field2)

		f1xml.AppendText(f.Field1)
		buf.Reset()
		buf.WriteString(d.Dump(false))
		d.Free()
	}
}
