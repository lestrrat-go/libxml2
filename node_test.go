package libxml2

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	ReportErrors(false)
}

type XmlNodeTypeToString struct {
	v XmlNodeType
	e string
}

func TestXmlNodeTypeStringer(t *testing.T) {
	values := []XmlNodeTypeToString{
		XmlNodeTypeToString{
			v: ElementNode,
			e: "ElementNode",
		},
		XmlNodeTypeToString{
			v: AttributeNode,
			e: "AttributeNode",
		},
		XmlNodeTypeToString{
			v: TextNode,
			e: "TextNode",
		},
		XmlNodeTypeToString{
			v: CDataSectionNode,
			e: "CDataSectionNode",
		},
		XmlNodeTypeToString{
			v: EntityRefNode,
			e: "EntityRefNode",
		},
		XmlNodeTypeToString{
			v: EntityNode,
			e: "EntityNode",
		},
		XmlNodeTypeToString{
			v: PiNode,
			e: "PiNode",
		},
		XmlNodeTypeToString{
			v: CommentNode,
			e: "CommentNode",
		},
		XmlNodeTypeToString{
			v: DocumentNode,
			e: "DocumentNode",
		},
		XmlNodeTypeToString{
			v: DocumentTypeNode,
			e: "DocumentTypeNode",
		},
		XmlNodeTypeToString{
			v: DocumentFragNode,
			e: "DocumentFragNode",
		},
		XmlNodeTypeToString{
			v: NotationNode,
			e: "NotationNode",
		},
		XmlNodeTypeToString{
			v: HTMLDocumentNode,
			e: "HTMLDocumentNode",
		},
		XmlNodeTypeToString{
			v: DTDNode,
			e: "DTDNode",
		},
		XmlNodeTypeToString{
			v: ElementDecl,
			e: "ElementDecl",
		},
		XmlNodeTypeToString{
			v: AttributeDecl,
			e: "AttributeDecl",
		},
		XmlNodeTypeToString{
			v: EntityDecl,
			e: "EntityDecl",
		},
		XmlNodeTypeToString{
			v: NamespaceDecl,
			e: "NamespaceDecl",
		},
		XmlNodeTypeToString{
			v: XIncludeStart,
			e: "XIncludeStart",
		},
		XmlNodeTypeToString{
			v: XIncludeEnd,
			e: "XIncludeEnd",
		},
		XmlNodeTypeToString{
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

	doc.SetDocumentElement(root)
	for i := 1; i <= 3; i++ {
		child, err := doc.CreateElement(fmt.Sprintf("child%d", i))
		if err != nil {
			t.Errorf("Failed to create child node: %s", err)
			return
		}
		child.AppendText(fmt.Sprintf("text%d", i))
		root.AppendChild(child)
	}

	// Temporary test
	expected := `<?xml version="1.0" encoding="utf-8"?>
<root><child1>text1</child1><child2>text2</child2><child3>text3</child3></root>
`
	if doc.String() != expected {
		t.Errorf("Failed to create XML document")
		t.Logf("Expected\n%s", expected)
		t.Logf("Got\n%s", doc.String())
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

	nslist := elem.GetNamespaces()
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