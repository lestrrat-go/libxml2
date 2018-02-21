package xpath_test

import (
	"testing"

	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/xpath"
	"github.com/stretchr/testify/assert"
)

func TestXPathContext(t *testing.T) {
	doc, err := libxml2.ParseString(`<foo><bar a="b"></bar></foo>`)
	if err != nil {
		t.Errorf("Failed to parse string: %s", err)
	}
	defer doc.Free()

	root, err := doc.DocumentElement()
	if !assert.NoError(t, err, "DocumentElement should succeed") {
		return
	}

	ctx, err := xpath.NewContext(root)
	if err != nil {
		t.Errorf("Failed to initialize XPathContext: %s", err)
		return
	}
	defer ctx.Free()

	// Use a string
	exprString := `/*`
	nodes := xpath.NodeList(ctx.Find(exprString))
	if len(nodes) != 1 {
		t.Errorf("Expected 1 nodes, got %d", len(nodes))
		return
	}

	iter := xpath.NodeIter(ctx.Find(exprString))
	{
		count := 0
		for iter.Next() {
			iter.Node()
			count++
		}
		if !assert.Equal(t, 1, count, "got 1 nodes from iterator") {
			return
		}
	}

	// Use an explicitly compiled expression
	expr, err := xpath.NewExpression(exprString)
	if err != nil {
		t.Errorf("Failed to compile xpath: %s", err)
		return
	}
	defer expr.Free()

	nodes = xpath.NodeList(ctx.FindExpr(expr))
	if len(nodes) != 1 {
		t.Errorf("Expected 1 nodes, got %d", len(nodes))
		return
	}
}

func TestXPathContextExpression_Number(t *testing.T) {
	ctx, err := xpath.NewContext()
	if err != nil {
		t.Errorf("Failed to initialize XPathContext: %s", err)
		return
	}
	defer ctx.Free()

	if !assert.Equal(t, float64(2), xpath.Number(ctx.Find("1+1")), "XPath evaluates to 2") {
		return
	}
	if !assert.Equal(t, float64(0), xpath.Number(ctx.Find("1<>1")), "XPath evaluates to 0") {
		return
	}
}

func TestXPathContextExpression_Boolean(t *testing.T) {
	ctx, err := xpath.NewContext()
	if err != nil {
		t.Errorf("Failed to initialize XPathContext: %s", err)
		return
	}
	defer ctx.Free()

	if !assert.True(t, xpath.Bool(ctx.Find("1=1")), "XPath evaluates to true") {
		return
	}
	if !assert.False(t, xpath.Bool(ctx.Find("1<>1")), "XPath evaluates to false") {
		return
	}
}

func TestXPathContextExpression_NodeList(t *testing.T) {
	doc, err := libxml2.ParseString(`<foo><bar a="b">baz</bar><bar a="c">quux</bar></foo>`)
	if err != nil {
		t.Errorf("Failed to parse string: %s", err)
	}
	defer doc.Free()

	root, err := doc.DocumentElement()
	if !assert.NoError(t, err, "DocumentElement should succeed") {
		return
	}

	ctx, err := xpath.NewContext(root)
	if err != nil {
		t.Errorf("Failed to initialize XPathContext: %s", err)
		return
	}
	defer ctx.Free()

	if !assert.Len(t, xpath.NodeList(ctx.Find("/foo/bar")), 2, "XPath evaluates to 2 nodes") {
		return
	}

	if !assert.Len(t, xpath.NodeList(ctx.Find("/foo/bar[bogus")), 0, "XPath evaluates to 0 nodes") {
		return
	}

	if !assert.Equal(t, "bazquux", xpath.String(ctx.Find("/foo/bar")), "XPath evaluates to 'bazquux'") {
		return
	}

	if !assert.Equal(t, "", xpath.String(ctx.Find("/[bogus")), "XPath evaluates to ''") {
		return
	}
}

func TestXPathContextExpression_Namespaces(t *testing.T) {
	doc, err := libxml2.ParseString(`<foo xmlns="http://example.com/foobar"><bar a="b"></bar></foo>`)
	if err != nil {
		t.Errorf("Failed to parse string: %s", err)
	}
	defer doc.Free()

	root, err := doc.DocumentElement()
	if !assert.NoError(t, err, "DocumentElement() should succeed") {
		return
	}

	ctx, err := xpath.NewContext(root)
	if err != nil {
		t.Errorf("Failed to initialize XPathContext: %s", err)
		return
	}
	defer ctx.Free()

	prefix := `xxx`
	nsuri := `http://example.com/foobar`
	if err := ctx.RegisterNS(prefix, nsuri); err != nil {
		t.Errorf("Failed to register namespace: %s", err)
		return
	}

	nodes := xpath.NodeList(ctx.Find(`/xxx:foo`))
	if len(nodes) != 1 {
		t.Errorf(`Expected 1 node, got %d`, len(nodes))
		return
	}
	if nodes[0].NodeName() != "foo" {
		t.Errorf(`Expected NodeName() "foo", got "%s"`, nodes[0].NodeName())
		return
	}

	gotns, err := ctx.LookupNamespaceURI(prefix)
	if err != nil {
		t.Errorf(`LookupNamespaceURI failed: %s`, err)
		return
	}

	if gotns != nsuri {
		t.Errorf(`Expected LookupNamespaceURI("%s") "%s", got "%s"`, prefix, nsuri, gotns)
		return
	}

	if !ctx.Exists(`//xxx:bar/@a`) {
		t.Errorf(`Expected "//xxx:bar/@a" to exist`)
		return
	}
	if ctx.Exists(`//xxx:bar/@b`) {
		t.Errorf(`Expected "//xxx:bar/@b" to NOT exist`)
		return
	}
}
