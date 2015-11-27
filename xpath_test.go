package libxml2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXPathContext(t *testing.T) {
	doc, err := ParseString(`<foo><bar a="b"></bar></foo>`)
	if err != nil {
		t.Errorf("Failed to parse string: %s", err)
	}
	defer doc.Free()

	root, err := doc.DocumentElement()
	if !assert.NoError(t, err, "DocumentElement should succeed") {
		return
	}

	ctx, err := NewXPathContext(root)
	if err != nil {
		t.Errorf("Failed to initialize XPathContext: %s", err)
		return
	}
	defer ctx.Free()

	// Use a string
	exprString := `/*`
	nodes, err := ctx.FindNodes(exprString)
	if err != nil {
		t.Errorf("Failed to execute FindNodes: %s", err)
		return
	}

	if len(nodes) != 1 {
		t.Errorf("Expected 1 nodes, got %d", len(nodes))
		return
	}

	// Use an explicitly compiled expression
	expr, err := NewXPathExpression(exprString)
	if err != nil {
		t.Errorf("Failed to compile xpath: %s", err)
		return
	}
	defer expr.Free()

	nodes, err = ctx.FindNodesExpr(expr)
	if err != nil {
		t.Errorf("Failed to execute FindNodesExpr: %s", err)
		return
	}

	if len(nodes) != 1 {
		t.Errorf("Expected 1 nodes, got %d", len(nodes))
		return
	}
}

func TestXPathContextExpression_Number(t *testing.T) {
	ctx, err := NewXPathContext()
	if err != nil {
		t.Errorf("Failed to initialize XPathContext: %s", err)
		return
	}
	defer ctx.Free()

	res, err := ctx.FindValue("1+1")
	if err != nil {
		t.Errorf("Failed to evaluate XPath expression: %s", err)
		return
	}
	defer res.Free()

	switch res.Type() {
	case XPathNumber:
		if res.Float64Value() != 2 {
			t.Errorf("Expected result number to be 2, got %f", res.Float64Value())
		}
	default:
		t.Errorf("Expected type to be XPathObjectNumber, got %s", res.Type())
	}
}

func TestXPathContextExpression_Boolean(t *testing.T) {
	ctx, err := NewXPathContext()
	if err != nil {
		t.Errorf("Failed to initialize XPathContext: %s", err)
		return
	}
	defer ctx.Free()

	res, err := ctx.FindValue("1=1")
	if err != nil {
		t.Errorf("Failed to evaluate XPath expression: %s", err)
		return
	}
	defer res.Free()

	switch res.Type() {
	case XPathBoolean:
		if !res.BoolValue() {
			t.Errorf("Expected result number to be false, got %s", res.BoolValue())
		}
	default:
		t.Errorf("Expected type to be XPathObjectBoolean, got %s", res.Type())
	}
}

func TestXPathContextExpression_NodeList(t *testing.T) {
	doc, err := ParseString(`<foo><bar a="b">baz</bar></foo>`)
	if err != nil {
		t.Errorf("Failed to parse string: %s", err)
	}
	defer doc.Free()

	root, err := doc.DocumentElement()
	if !assert.NoError(t, err, "DocumentElement should succeed") {
		return
	}

	ctx, err := NewXPathContext(root)
	if err != nil {
		t.Errorf("Failed to initialize XPathContext: %s", err)
		return
	}
	defer ctx.Free()

	res, err := ctx.FindValue("/foo/bar")
	if err != nil {
		t.Errorf("Failed to evaluate XPath expression: %s", err)
		return
	}
	defer res.Free()

	switch res.Type() {
	case XPathNodeSet:
		s, err := res.StringValue()
		if !assert.NoError(t, err, "StringValue() should succeed") {
			return
		}
		if !assert.Equal(t, "baz", s, "results match") {
			return
		}
	default:
		t.Errorf("Expected type to be XPathObjectNodeSet, got %s", res.Type())
	}
}

func TestXPathContextExpression_Namespaces(t *testing.T) {
	doc, err := ParseString(`<foo xmlns="http://example.com/foobar"><bar a="b"></bar></foo>`)
	if err != nil {
		t.Errorf("Failed to parse string: %s", err)
	}
	defer doc.Free()

	root, err := doc.DocumentElement()
	if !assert.NoError(t, err, "DocumentElement() should succeed") {
		return
	}

	ctx, err := NewXPathContext(root)
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

	nodes, err := ctx.FindNodes(`/xxx:foo`)
	if err != nil {
		t.Errorf(`Failed to evaluate "/xxx:foo": %s`, err)
		return
	}
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
