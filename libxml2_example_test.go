package libxml2_test

import (
	"log"
	"net/http"

	"github.com/lestrrat/go-libxml2"
	"github.com/lestrrat/go-libxml2/parser"
	"github.com/lestrrat/go-libxml2/types"
	"github.com/lestrrat/go-libxml2/xpath"
)

func ExampleXML() {
	res, err := http.Get("http://blog.golang.org/feed.atom")
	if err != nil {
		panic("failed to get blog.golang.org: " + err.Error())
	}

	p := parser.New()
	doc, err := p.ParseReader(res.Body)
	defer res.Body.Close()

	if err != nil {
		panic("failed to parse XML: " + err.Error())
	}
	defer doc.Free()

	doc.Walk(func(n types.Node) error {
		log.Printf(n.NodeName())
		return nil
	})

	root, err := doc.DocumentElement()
	if err != nil {
		log.Printf("Failed to fetch document element: %s", err)
		return
	}

	ctx, err := xpath.NewContext(root)
	if err != nil {
		log.Printf("Failed to create xpath context: %s", err)
		return
	}
	defer ctx.Free()

	ctx.RegisterNS("atom", "http://www.w3.org/2005/Atom")
	title := xpath.String(ctx.Find("/atom:feed/atom:title/text()"))
	log.Printf("feed title = %s", title)
}

func ExampleHTML() {
	res, err := http.Get("http://golang.org")
	if err != nil {
		panic("failed to get golang.org: " + err.Error())
	}

	doc, err := libxml2.ParseHTMLReader(res.Body)
	if err != nil {
		panic("failed to parse HTML: " + err.Error())
	}
	defer doc.Free()

	doc.Walk(func(n types.Node) error {
		log.Printf(n.NodeName())
		return nil
	})

	nodes := xpath.NodeList(doc.Find(`//div[@id="menu"]/a`))
	for i := 0; i < len(nodes); i++ {
		log.Printf("Found node: %s", nodes[i].NodeName())
	}
}
