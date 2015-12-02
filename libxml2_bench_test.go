// This file is build-tag protected because it involves loading an external
// library (xmlpath)
package libxml2_test

import (
	"bytes"
	"encoding/xml"
	"log"
	"net/http"
	"testing"

	"github.com/lestrrat/go-libxml2"
	"github.com/lestrrat/go-libxml2/dom"
	"github.com/lestrrat/go-libxml2/xpath"
	"gopkg.in/xmlpath.v1"
)

func BenchmarkXmlpathXmlpath(b *testing.B) {
	res, err := http.Get("http://mattn.kaoriya.net/index.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	root, err := xmlpath.Parse(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		p, err := xmlpath.Compile(`//loc`)
		if err != nil {
			log.Fatal(err)
		}
		it := p.Iter(root)
		for it.Next() {
			_ = it.Node()
		}
	}
}

func BenchmarkLibxml2Xmlpath(b *testing.B) {
	res, err := http.Get("http://mattn.kaoriya.net/index.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := libxml2.ParseReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		nodes := xpath.NodeList(doc.FindValue(`//loc`))
		for _, n := range nodes {
			_ = n
		}
	}
}

type Foo struct {
	XMLName xml.Name `xml:"https://github.com/lestrrat/go-libxml2/foo foo:foo"`
	Field1  string
}

func BenchmarkEncodingXMLDOM(b *testing.B) {
	var buf bytes.Buffer
	f := Foo{
		Field1: "Hello, World!",
	}
	for i := 0; i < b.N; i++ {
		buf.Reset()
		enc := xml.NewEncoder(&buf)
		enc.Encode(f)
	}
}

func BenchmarkLibxml2DOM(b *testing.B) {
	var buf bytes.Buffer
	const nsuri = `https://github.com/lestrrat/go-libxml2/foo`
	f := Foo{
		Field1: "Hello, World!",
	}
	for i := 0; i < b.N; i++ {
		d := dom.CreateDocument()
		defer d.Free()

		root, err := d.CreateElementNS(nsuri, "foo:foo")
		if err != nil {
			panic(err)
		}
		d.SetDocumentElement(root)

		f1xml, err := d.CreateElement("Field1")
		if err != nil {
			panic(err)
		}
		root.AddChild(f1xml)

		f1xml.AppendText(f.Field1)
		buf.Reset()
		buf.WriteString(d.Dump(false))
	}
}
