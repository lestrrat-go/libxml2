# go-libxml2

[![Build Status](https://travis-ci.org/lestrrat/go-libxml2.svg?branch=master)](https://travis-ci.org/lestrrat/go-libxml2)

[![GoDoc](https://godoc.org/github.com/lestrrat/go-libxml2?status.svg)](https://godoc.org/github.com/lestrrat/go-libxml2)

Interface to libxml2, with DOM interface.

## Status

* This library should be considered alpha grade. API may still change.
* Much of commonly used functionalities from libxml2 that *I* use are there already, and are known to be functional

Having said the above, AFAIK this is the most complete Go wrapper for libxml2 as of this writing
(Nov, 2015).

## Examples

### Basic XML Example

```go
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

  ctx.RegisterNS("atom", "http://www.w3.org/2005/Atom")
  title := ctx.FindValue("/atom:feed/atom:title/text()").String()
  log.Printf("feed title = %s", title)
}
```

### Basic HTML Example

```go
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

### XSD Validation

```go
import (
  "io/ioutil"
  "log"
  "os"
  "path/filepath"

  "github.com/lestrrat/go-libxml2"
  "github.com/lestrrat/go-libxml2/xsd"
)

func ExampleXSD() {
  xsdfile := filepath.Join("test", "xmldsig-core-schema.xsd")
  f, err := os.Open(xsdfile)
  if err != nil {
    log.Printf("failed to open file: %s", err)
    return
  }
  defer f.Close()

  buf, err := ioutil.ReadAll(f)
  if err != nil {
    log.Printf("failed to read file: %s", err)
    return
  }

  s, err := xsd.Parse(buf)
  if err != nil {
    log.Printf("failed to parse XSD: %s", err)
    return
  }
  defer s.Free()

  d, err := libxml2.ParseString(`<foo></foo>`)
  if err != nil {
    log.Printf("failed to parse XML: %s", err)
    return
  }

  if err := s.Validate(d); err != nil {
    for _, e := range err.(SchemaValidationErr).Errors() {
      log.Printf("error: %s", e.Error())
    }
    return
  }

  log.Printf("validation successful!")
}
```

## Caveats

### Other libraries

There exists many similar libraries. I want speed, I want DOM, and I want XPath.When all of these are met, I'd be happy to switch to another library.

For now my closest contender was [xmlpath](https://github.com/go-xmlpath/xmlpath), but as of this writing it suffers in the speed (for xpath) area a bit:

```
shoebill% go test -v -run=none -benchmem -benchtime=5s -bench .
PASS
BenchmarkXmlpath_Xmlpath-4    200000         43330 ns/op        1394 B/op          6 allocs/op
BenchmarkXmlpath_Libxml2-4    200000         25488 ns/op          80 B/op          5 allocs/op
BenchmarkDOM_EncodingXml-4   2000000          3418 ns/op        4496 B/op          8 allocs/op
BenchmarkDOM_Libxml2-4        500000         21086 ns/op         287 B/op          8 allocs/op
ok      github.com/lestrrat/go-libxml2  42.319s
```

## See Also

* https://github.com/lestrrat/go-xmlsec

## Credits

* Work on this library was generously sponsored by HDE Inc (https://www.hde.co.jp)
