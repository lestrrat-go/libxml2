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

	p := libxml2.NewParser()
	doc, err := p.ParseReader(res.Body)
	defer res.Body.Close()

	if err != nil {
		panic("failed to parse XML: " + err.Error())
	}
	defer doc.Free()

	doc.Walk(func(n libxml2.Node) error {
		log.Printf(n.NodeName())
		return nil
	})

	ctx, err := libxml2.NewXPathContext(doc.DocumentElement())
	if err != nil {
		log.Printf("Failed to create xpath context: %s", err)
		return
	}
	defer ctx.Free()

	ctx.RegisterNs("atom", "http://www.w3.org/2005/Atom")
	title, err := ctx.FindValue("/atom:feed/atom:title/text()")
	if err != nil {
		log.Printf("Failed to run FindValue: %s", err)
		return
	}
	log.Printf("feed title = %s", title)
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
