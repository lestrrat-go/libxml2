# libxml2

Interface to libxml2, with DOM interface.

[![Build Status](https://travis-ci.org/lestrrat-go/libxml2.svg?branch=master)](https://travis-ci.org/lestrrat-go/libxml2)

[![GoDoc](https://godoc.org/github.com/lestrrat-go/libxml2?status.svg)](https://godoc.org/github.com/lestrrat-go/libxml2)

# Index

* [Why?](#why)
* [FAQ](#faq)

## Why?

I needed to write [go-xmlsec](https://github.com/lestrrat-go/xmlsec). This means we need to build trees using libxml2, and then muck with it in xmlsec: Two separate packages in Go means we cannot (safely) pass around `C.xmlFooPtr` objects (also, you pay a penalty for pointer types). This package carefully avoid references to `C.xmlFooPtr` types and uses uintptr to pass data around, so other libraries that needs to interact with libxml2 can safely interact with it.

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

  "github.com/lestrrat-go/libxml2"
  "github.com/lestrrat-go/libxml2/parser"
  "github.com/lestrrat-go/libxml2/types"
  "github.com/lestrrat-go/libxml2/xpath"
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

  nodes := xpath.NodeList(doc.Find(`//div[@id="menu"]/a`))
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

  "github.com/lestrrat-go/libxml2"
  "github.com/lestrrat-go/libxml2/xsd"
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
  defer d.Free()

  if err := s.Validate(d); err != nil {
    for _, e := range err.(xsd.SchemaValidationError).Errors() {
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
ok      github.com/lestrrat-go/libxml2  37.597s
```

## FAQ

### "It won't build"

The very first thing you need to be aware is that this is a _C binding_ to
libxml2. You should understand how to build C programs, how to debug them,
or at least be able to ask the right questions and deal with a great deal
more than Go alone.

Having said that, the most common causes for build errors are:

1. **You have not installed libxml2 / You installed it incorrectly**

The first one is obvious, but I get this a lot. You have to install libxml2.
If you are installing via some sort of package manager like apt/apk, remember
that you need to install the "development" files as well. The name of the
package differs in each environment, but it's usually something like "libxml2-dev".

The second is more subtle, and tends to happen when you install your libxml2
in a non-standard location. This causes problems for other tools such as
your C compiler or pkg-config. See more below

2. **Your header files are not in the search path**

If you don't understand what header files are or how they work, this is where
you should either look for your local C-guru, or study how these things work
before filing an issue on this repository.

Your C compiler, which is invoked via Go, needs to be able to find the libxml2
header files. If you installed them in a non-standard location, for example,
such as outside of /usr/include and /usr/local/include, you _may_ have to
configure them yourself.

How to configure them depends greatly on your environment, and again, if you
don't understand how you can fix it, you should consult your local C-guru
about it, not this repository.

3. **Your pkg-config files are not in the search path**

If you don't understand what pkg-config does, this is where you should either 
look for your local sysadmin friend, or study how these things work
before filing an issue on this repository.

pkg-config provides metadata about a installed components, such as build flags
that are required. Go uses it to figure out how to build and link Go programs
that needs to interact with things written in C.

However, pkg-config is merely a thin frontend to extract information from 
file(s) that each component provided upon installation. 
pkg-config itself needs to know where to find these files.

Make sure that the output of the following command contains `libxml-2.0`.
If not, and you don't understand how to fix this yourself, you should consult
your local sysadmin friend about it, not this repository

```
pkg-config --list-all
```

### "Fatal error: 'libxml/HTMLparser.h' file not found"

See the first FAQ entry.

### I can't build this library statically

See prior discussion: https://github.com/lestrrat-go/libxml2/issues/62

## See Also

* https://github.com/lestrrat-go/xmlsec

## Credits

* Work on this library was generously sponsored by HDE Inc (https://www.hde.co.jp)
