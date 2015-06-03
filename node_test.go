package libxml2

import "testing"

type XmlElementTypeToString struct {
	v XmlElementType
	e string
}

func TestXmlElementTypeStringer(t *testing.T) {
	values := []XmlElementTypeToString{
		XmlElementTypeToString{
			v: ElementNode,
			e: "ElementNode",
		},
		XmlElementTypeToString{
			v: AttributeNode,
			e: "AttributeNode",
		},
		XmlElementTypeToString{
			v: TextNode,
			e: "TextNode",
		},
		XmlElementTypeToString{
			v: CDataSectionNode,
			e: "CDataSectionNode",
		},
		XmlElementTypeToString{
			v: EntityRefNode,
			e: "EntityRefNode",
		},
		XmlElementTypeToString{
			v: EntityNode,
			e: "EntityNode",
		},
		XmlElementTypeToString{
			v: PiNode,
			e: "PiNode",
		},
		XmlElementTypeToString{
			v: CommentNode,
			e: "CommentNode",
		},
		XmlElementTypeToString{
			v: DocumentNode,
			e: "DocumentNode",
		},
		XmlElementTypeToString{
			v: DocumentTypeNode,
			e: "DocumentTypeNode",
		},
		XmlElementTypeToString{
			v: DocumentFragNode,
			e: "DocumentFragNode",
		},
		XmlElementTypeToString{
			v: NotationNode,
			e: "NotationNode",
		},
		XmlElementTypeToString{
			v: HTMLDocumentNode,
			e: "HTMLDocumentNode",
		},
		XmlElementTypeToString{
			v: DTDNode,
			e: "DTDNode",
		},
		XmlElementTypeToString{
			v: ElementDecl,
			e: "ElementDecl",
		},
		XmlElementTypeToString{
			v: AttributeDecl,
			e: "AttributeDecl",
		},
		XmlElementTypeToString{
			v: EntityDecl,
			e: "EntityDecl",
		},
		XmlElementTypeToString{
			v: NamespaceDecl,
			e: "NamespaceDecl",
		},
		XmlElementTypeToString{
			v: XIncludeStart,
			e: "XIncludeStart",
		},
		XmlElementTypeToString{
			v: XIncludeEnd,
			e: "XIncludeEnd",
		},
		XmlElementTypeToString{
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