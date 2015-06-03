# go-libxml2

[![Build Status](https://travis-ci.org/lestrrat/go-libxml2.svg?branch=master)](https://travis-ci.org/lestrrat/go-libxml2)

[![GoDoc](https://godoc.org/github.com/lestrrat/go-libxml2?status.svg)](https://godoc.org/github.com/lestrrat/go-libxml2)

Interface to libxml2, with DOM interface

This library is still in very early stages of development. API may still change

```go
import (
	"log"
	"net/http"

	"github.com/lestrrat/go-libxml2"
)

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
```

# Caveats

## Other libraries

There exists many similar libraries. I want speed, I want DOM, and I want XPath.When all of these are met, I'd be happy to switch to another library.

For now my closest contender was [xmlpath](https://github.com/go-xmlpath/xmlpath), but as of this writing it suffers in the speed (for xpath) area a bit:

```
go test -tags bench -bench=. -benchtime=5s
PASS
BenchmarkXmlpath      100000         88764 ns/op
BenchmarkLibxml2      300000         22509 ns/op
ok      github.com/lestrrat/go-libxml2  24.926s
```
