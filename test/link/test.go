package main

import (
	"github.com/lestrrat-go/libxml2"
)

func main() {
	doc, err := libxml2.ParseHTMLString(`<html><body><h1>Hello, World!</h1><p>Lorem Ipsum</p></body></html>`)
	if err != nil {
		panic(err)
	}
	doc.Free()
}
