package libxml2_test

import (
	"log"
	"net/http"

	"github.com/lestrrat/go-libxml2"
)

func ExmapleXML() {
	res, err := http.Get("http://blog.golang.org/feed.atom")
	if err != nil {
		panic("failed to get blog.golang.org: " + err.Error())
	}

	p := &libxml2.Parser{}
	doc, err := p.Parse(res.Body)
	defer res.Body.Close()

	if err != nil {
		panic("failed to parse XML: " + err.Error())
	}
	defer doc.Free()

	doc.Walk(func(n libxml2.Node) error {
		log.Printf(n.NodeName())
		return nil
	})

	// XML with namespaces needs XPathContext, and we haven't
	// gotten around to implementing it yet
	//
	// doc.FindNodes(...)
}


func ExampleHTML() {
	res, err := http.Get("http://golang.org")
	if err != nil {
		panic("failed to get golang.org: " + err.Error())
	}

	doc, err := libxml2.ParseHTML(res.Body)
	if err != nil {
		panic("failed to parse HTML: " + err.Error())
	}
	defer doc.Free()

	doc.Walk(func(n libxml2.Node) error {
		log.Printf(n.NodeName())
		return nil
	})

	nodes, err := doc.FindNodes(`//div[@id="menu"]/a`)
	if err != nil {
		panic("failed to evaluate xpath: " + err.Error())
	}

	for i := 0; i < len(nodes); i++ {
		log.Printf("Found node: %s", nodes[i].NodeName())
	}
}
