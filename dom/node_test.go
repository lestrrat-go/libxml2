package dom_test

import (
	"fmt"
	"testing"

	"github.com/lestrrat-go/libxml2/clib"
	. "github.com/lestrrat-go/libxml2/dom"
	"github.com/lestrrat-go/libxml2/types"
	"github.com/stretchr/testify/assert"
)

func init() {
	clib.ReportErrors(false)
}

type XMLNodeTypeToString struct {
	v clib.XMLNodeType
	e string
}

func TestXMLNodeTypeStringer(t *testing.T) {
	values := []XMLNodeTypeToString{
		{
			v: ElementNode,
			e: "ElementNode",
		},
		{
			v: AttributeNode,
			e: "AttributeNode",
		},
		{
			v: TextNode,
			e: "TextNode",
		},
		{
			v: CDataSectionNode,
			e: "CDataSectionNode",
		},
		{
			v: EntityRefNode,
			e: "EntityRefNode",
		},
		{
			v: EntityNode,
			e: "EntityNode",
		},
		{
			v: PiNode,
			e: "PiNode",
		},
		{
			v: CommentNode,
			e: "CommentNode",
		},
		{
			v: DocumentNode,
			e: "DocumentNode",
		},
		{
			v: DocumentTypeNode,
			e: "DocumentTypeNode",
		},
		{
			v: DocumentFragNode,
			e: "DocumentFragNode",
		},
		{
			v: NotationNode,
			e: "NotationNode",
		},
		{
			v: HTMLDocumentNode,
			e: "HTMLDocumentNode",
		},
		{
			v: DTDNode,
			e: "DTDNode",
		},
		{
			v: ElementDecl,
			e: "ElementDecl",
		},
		{
			v: AttributeDecl,
			e: "AttributeDecl",
		},
		{
			v: EntityDecl,
			e: "EntityDecl",
		},
		{
			v: NamespaceDecl,
			e: "NamespaceDecl",
		},
		{
			v: XIncludeStart,
			e: "XIncludeStart",
		},
		{
			v: XIncludeEnd,
			e: "XIncludeEnd",
		},
		{
			v: DocbDocumentNode,
			e: "DocbDocumentNode",
		},
	}

	for _, d := range values {
		if d.v.String() != d.e {
			t.Errorf("e '%s', got '%s'", d.e, d.v.String())
		}
	}
}

func TestDOM(t *testing.T) {
	doc := CreateDocument()
	defer doc.Free()

	root, err := doc.CreateElement("root")
	if err != nil {
		t.Errorf("Failed to create root element: %s", err)
		return
	}

	_ = doc.SetDocumentElement(root)
	var toRemove types.Node
	for i := 1; i <= 3; i++ {
		child, err := doc.CreateElement(fmt.Sprintf("child%d", i))
		if !assert.NoError(t, err, "dom.CreateElement(child%d) should succeed", i) {
			return
		}
		_ = child.AppendText(fmt.Sprintf("text%d", i))
		_ = root.AddChild(child)

		if i == 2 {
			toRemove = child
		}
	}

	// Temporary test
	expected := `<?xml version="1.0" encoding="utf-8"?>
<root><child1>text1</child1><child2>text2</child2><child3>text3</child3></root>
`
	if !assert.Equal(t, expected, doc.String(), "Failed to create XML document") {
		return
	}

	if !assert.NoError(t, root.RemoveChild(toRemove), "RemoveChild should succeed") {
		return
	}
	expected = `<?xml version="1.0" encoding="utf-8"?>
<root><child1>text1</child1><child3>text3</child3></root>
`
	if !assert.Equal(t, expected, doc.String(), "XML should match") {
		return
	}
}

func TestNode_StandaloneWithNamespaces(t *testing.T) {
	uri := "http://kungfoo"
	prefix := "foo"
	name := "bar"

	doc := CreateDocument()
	elem, err := doc.CreateElementNS(uri, prefix+":"+name)
	if !assert.NoError(t, err, "CreateElementNS snould succeed") {
		return
	}

	lookedup, err := elem.LookupNamespaceURI(prefix)
	if !assert.NoError(t, err, "LookupNamespaceURI should succeed") {
		return
	}
	if !assert.Equal(t, uri, lookedup, "LookupNamespaceURI succeeds") {
		return
	}

	lookedup, err = elem.LookupNamespacePrefix(uri)
	if !assert.NoError(t, err, "LookupNamespacePrefix should succeed") {
		return
	}
	if !assert.Equal(t, prefix, lookedup, "LookupNamespacePrefix succeeds") {
		return
	}

	nslist, err := elem.GetNamespaces()
	if !assert.NoError(t, err, "GetNamespaces succeeds") {
		return
	}

	defer func() {
		for _, ns := range nslist {
			ns.Free()
		}
	}()

	if !assert.Len(t, nslist, 1, "GetNamespaces returns 1 namespace") {
		return
	}
}

func TestAttribute(t *testing.T) {
	doc := CreateDocument()
	attr, err := doc.CreateAttribute("foo", "bar")
	if !assert.NoError(t, err, "attribute created") {
		return
	}

	if !assert.NotPanics(t, func() { attr.Free() }, "free should not panic") {
		return
	}
}

func TestCreateElementNS(t *testing.T) {
	doc := CreateDocument()
	root, err := doc.CreateElementNS("http://foo.bar.baz", "foo:root")
	if !assert.NoError(t, err, "CreateElementNS should succeed") {
		return
	}
	_ = doc.SetDocumentElement(root)

	n1, err := doc.CreateElementNS("http://foo.bar.baz", "foo:n1")
	if !assert.NoError(t, err, "CreateElementNS should succeed") {
		return
	}
	_ = root.AddChild(n1)

	n2, err := doc.CreateElementNS("http://foo.bar.baz", "bar:n2")
	if !assert.NoError(t, err, "CreateElementNS should succeed") {
		return
	}
	_ = root.AddChild(n2)

	_, err = doc.CreateElementNS("http://foo.bar.baz.quux", "foo:n3")
	if !assert.Error(t, err, "CreateElementNS should fail") {
		return
	}

	t.Logf("%s", doc.Dump(false))
}
