# go-libxml2

[![Build Status](https://travis-ci.org/lestrrat/go-libxml2.svg?branch=master)](https://travis-ci.org/lestrrat/go-libxml2)

[![GoDoc](https://godoc.org/github.com/lestrrat/go-libxml2?status.svg)](https://godoc.org/github.com/lestrrat/go-libxml2)

Interface to libxml2, with DOM interface.

## Status

* This library should be considered alpha grade. API may still change.
* Much of commonly used functionalities from libxml2 that *I* use are there already, and are known to be functional

## Package Layout:

| Name    | Description                                                 |
|---------|-------------------------------------------------------------|
| libxml2 | Globally available utility functions, such as `ParseString` |
| types   | Common data types, such as `types.Node`                     |
| parser  | Parser routines                                             |
| dom     | DOM-like manipulation of XML document/nodes                 |
| xpath   | XPath related tools                                         |
| xsd     | XML Schema related tools                                    |
| clib    | Wrapper around C libxml2 library - DO NOT TOUCH IF UNSURE   |

## Features

Create XML documents using DOM-like interface:

```go
  d := dom.CreateDocument()
  e, err := d.CreateElement("foo")
  if err != nil {
    println(err)
    return
  }
  d.SetDocumentElement(e)
  ...
```

Parse documents:

```go
  d, err := libxml2.ParseString(xmlstring)
  if err != nil {
    println(err)
    return
  }
```

Use XPath to extract node values:

```go
  text := xpath.String(node.Find("//xpath/expression"))
```

## Examples

### Basic XML Example

```go
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

  doc.Walk(func(n types.Node) error {
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
BenchmarkXmlpathXmlpath-4     500000         11737 ns/op         721 B/op          6 allocs/op
BenchmarkLibxml2Xmlpath-4    1000000          7627 ns/op         368 B/op         15 allocs/op
BenchmarkEncodingXMLDOM-4    2000000          4079 ns/op        4560 B/op          9 allocs/op
BenchmarkLibxml2DOM-4        1000000         11454 ns/op         264 B/op          7 allocs/op
ok      github.com/lestrrat/go-libxml2  37.597s
```

## See Also

* https://github.com/lestrrat/go-xmlsec

## Credits

* Work on this library was generously sponsored by HDE Inc (https://www.hde.co.jp)
