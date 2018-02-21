package libxml2_test

import (
	"testing"

	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/xpath"
	"github.com/stretchr/testify/assert"
)

func TestParseHTML(t *testing.T) {
	doc, err := libxml2.ParseHTMLString(`<html><body><h1>Hello, World!</h1><p>Lorem Ipsum</p></body></html>`)
	if err != nil {
		t.Errorf("Failed to parse: %s", err)
		return
	}
	defer doc.Free()

	root, err := doc.DocumentElement()
	if !assert.NoError(t, err, "DocumentElement() should succeed") {
		return
	}
	if !assert.True(t, root.IsSameNode(root), "root == root") {
		return
	}

	nodes := xpath.NodeList(doc.Find("/html/body/h1"))
	if len(nodes) != 1 {
		t.Errorf("Could not find matching nodes")
		return
	}

	if nodes[0].TextContent() != "Hello, World!" {
		t.Errorf("h1 content is not 'Hello, World!', got %s", nodes[0].TextContent())
		return
	}
}
