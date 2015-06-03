package libxml2

import "testing"

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
	doc := NewDocument("1.0")
	root := doc.CreateElement("root")

	doc.SetDocumentElement(root)

	// Temporary test
	expected := `<?xml version="1.0"?>
<root/>
`
	if doc.String() != expected {
		t.Errorf("Failed to create XML document")
		t.Logf("Expected\n%s", expected)
		t.Logf("Got\n%s", doc.String())
		return
	}
}