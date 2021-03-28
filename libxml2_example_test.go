package libxml2_test

import (
	"context"
	"log"
	"net/http"

	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/parser"
	"github.com/lestrrat-go/libxml2/types"
	"github.com/lestrrat-go/libxml2/xpath"
)

//nolint:govet
func ExampleXML() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://blog.golang.org/feed.atom", nil)
	if err != nil {
		panic("failed to create request: " + err.Error())
	}

	res, err := http.DefaultClient.Do(req)
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
		log.Println(n.NodeName())
		return nil
	})

	root, err := doc.DocumentElement()
	if err != nil {
		log.Printf("Failed to fetch document element: %s", err)
		return
	}

	xctx, err := xpath.NewContext(root)
	if err != nil {
		log.Printf("Failed to create xpath context: %s", err)
		return
	}
	defer xctx.Free()

	xctx.RegisterNS("atom", "http://www.w3.org/2005/Atom")
	title := xpath.String(xctx.Find("/atom:feed/atom:title/text()"))
	log.Printf("feed title = %s", title)
}

//nolint:govet
func ExampleHTML() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://golang.org", nil)
	if err != nil {
		panic("failed to create request: " + err.Error())
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic("failed to get golang.org: " + err.Error())
	}

	doc, err := libxml2.ParseHTMLReader(res.Body)
	if err != nil {
		panic("failed to parse HTML: " + err.Error())
	}
	defer res.Body.Close()

	defer doc.Free()

	doc.Walk(func(n types.Node) error {
		log.Println(n.NodeName())
		return nil
	})

	nodes := xpath.NodeList(doc.Find(`//div[@id="menu"]/a`))
	for i := 0; i < len(nodes); i++ {
		log.Printf("Found node: %s", nodes[i].NodeName())
	}
}
