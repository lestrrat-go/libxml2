package dom

import (
	"errors"

	"github.com/lestrrat/go-libxml2/clib"
	"github.com/lestrrat/go-libxml2/node"
	"github.com/lestrrat/go-libxml2/xpath"
)

func init () {
	xpath.WrapNodeFunc = WrapNode
}


func WrapDocument(n uintptr) *Document {
	return &Document{
		ptr: n,
	}
}

func wrapNamespace(n uintptr) *Namespace {
	return &Namespace{
		XMLNode: wrapXMLNode(n),
	}
}

func wrapAttribute(n uintptr) *Attribute {
	return &Attribute{
		XMLNode: wrapXMLNode(n),
	}
}

func wrapCDataSection(n uintptr) *CDataSection {
	return &CDataSection{
		XMLNode: wrapXMLNode(n),
	}
}

func wrapComment(n uintptr) *Comment {
	return &Comment{
		XMLNode: wrapXMLNode(n),
	}
}

func wrapElement(n uintptr) *Element {
	return &Element{
		XMLNode: wrapXMLNode(n),
	}
}

func wrapText(n uintptr) *Text {
	return &Text{wrapXMLNode(n)}
}

func wrapXMLNode(n uintptr) *XMLNode {
	return &XMLNode{ptr: n}
}

// WrapNode is a function created with the sole purpose of allowing
// go-libxml2 consumers that can generate a C.xmlNode pointer to
// create libxml2.Node types, e.g. go-xmlsec.
func WrapNode(n uintptr) (node.Node, error) {
	switch clib.XMLGetNodeTypeRaw(n) {
	case clib.AttributeNode:
		return wrapAttribute(n), nil
	case clib.ElementNode:
		return wrapElement(n), nil
	case clib.TextNode:
		return wrapText(n), nil
	default:
		return nil, errors.New("unknown node")
	}
}

