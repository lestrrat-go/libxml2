// +build bench

// This file is build-tag protected because it involves loading an external
// library (xmlpath)
package libxml2

import (
	"log"
	"net/http"
	"testing"

	"github.com/go-xmlpath/xmlpath"
)

func BenchmarkXmlpath(b *testing.B) {
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

func BenchmarkLibxml2(b *testing.B) {
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
