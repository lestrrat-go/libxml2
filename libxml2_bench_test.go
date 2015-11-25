// +build bench

// This file is build-tag protected because it involves loading an external
// library (xmlpath)
package libxml2

import (
	"bytes"
	"encoding/xml"
	"log"
	"net/http"
	"testing"

	"github.com/go-xmlpath/xmlpath"
)

func BenchmarkXmlpath_Xmlpath(b *testing.B) {
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

func BenchmarkXmlpath_Libxml2(b *testing.B) {
	res, err := http.Get("http://mattn.kaoriya.net/index.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := ParseReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		nodes, _ := doc.FindNodes(`//loc`)
		for _, n := range nodes {
			_ = n
		}
	}
}

type Foo struct {
	XMLName xml.Name `xml:"https://github.com/lestrrat/go-libxml2/foo foo:foo"`
	Field1  string
}

func BenchmarkDOM_EncodingXml(b *testing.B) {
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

func BenchmarkDOM_Libxml2(b *testing.B) {
	f := Foo{
		Field1: "Hello, World!",
	}
	for i := 0; i < b.N; i++ {
		d := CreateDocument()
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
		root.AppendChild(f1xml)

		f1xml.AppendText(f.Field1)
		buf.Reset()
		buf.WriteString(d.Dump(false))
	}
}

func Benchmark_stringToXmlChar(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xmlchar := stringToXmlChar("Hello, World")
		_ = xmlchar
	}
}

func Benchmark_xmlCharToString(b *testing.B) {
	xmlchar := stringToXmlChar("Hello, World")
	for i := 0; i < b.N; i++ {
		_ = xmlCharToString(xmlchar)
	}
}
