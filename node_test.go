package libxml2

import "testing"

type XmlElementTypeToString struct {
	v XmlElementType
	e string
}

func TestXmlElementTypeStringer(t *testing.T) {
	values := []XmlElementTypeToString{
		XmlElementTypeToString{
			v: XmlElementNode,
			e: "ElementNode",
		},
		XmlElementTypeToString{
			v: XmlAttributeNode,
			e: "AttributeNode",
		},
		XmlElementTypeToString{
			v: XmlTextNode,
			e: "TextNode",
		},
		XmlElementTypeToString{
			v: XmlCDataSectionNode,
			e: "CDataSectionNode",
		},
		XmlElementTypeToString{
			v: XmlEntityRefNode,
			e: "EntityRefNode",
		},
		XmlElementTypeToString{
			v: XmlEntityNode,
			e: "EntityNode",
		},
		XmlElementTypeToString{
			v: XmlPiNode,
			e: "PiNode",
		},
		XmlElementTypeToString{
			v: XmlCommentNode,
			e: "CommentNode",
		},
		XmlElementTypeToString{
			v: XmlDocumentNode,
			e: "DocumentNode",
		},
		XmlElementTypeToString{
			v: XmlDocumentTypeNode,
			e: "DocumentTypeNode",
		},
		XmlElementTypeToString{
			v: XmlDocumentFragNode,
			e: "DocumentFragNode",
		},
		XmlElementTypeToString{
			v: XmlNotationNode,
			e: "NotationNode",
		},
		XmlElementTypeToString{
			v: XmlHTMLDocumentNode,
			e: "HTMLDocumentNode",
		},
		XmlElementTypeToString{
			v: XmlDTDNode,
			e: "DTDNode",
		},
		XmlElementTypeToString{
			v: XmlElementDecl,
			e: "ElementDecl",
		},
		XmlElementTypeToString{
			v: XmlAttributeDecl,
			e: "AttributeDecl",
		},
		XmlElementTypeToString{
			v: XmlEntityDecl,
			e: "EntityDecl",
		},
		XmlElementTypeToString{
			v: XmlNamespaceDecl,
			e: "NamespaceDecl",
		},
		XmlElementTypeToString{
			v: XmlXIncludeStart,
			e: "XIncludeStart",
		},
		XmlElementTypeToString{
			v: XmlXIncludeEnd,
			e: "XIncludeEnd",
		},
		XmlElementTypeToString{
			v: XmlDocbDocumentNode,
			e: "DocbDocumentNode",
		},
	}

	for _, d := range values {
		if d.v.String() != d.e {
			t.Errorf("e '%s', got '%s'", d.e, d.v.String())
		}
	}
}