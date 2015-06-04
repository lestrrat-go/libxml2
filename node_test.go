package libxml2

import (
	"fmt"
	"testing"
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
	expected := `<?xml version="1.0"?>
<root><child1>text1</child1><child2>text2</child2><child3>text3</child3></root>
`
	if doc.String() != expected {
		t.Errorf("Failed to create XML document")
		t.Logf("Expected\n%s", expected)
		t.Logf("Got\n%s", doc.String())
		return
	}
}