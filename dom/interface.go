package dom

import (
	"errors"

	"github.com/lestrrat/go-libxml2/clib"
)

var (
	ErrAttributeNotFound = errors.New("attribute not found")
)

// XMLNodeType identifies the type of the underlying C struct
type XMLNodeType clib.XMLNodeType

const (
	ElementNode      = clib.ElementNode
	AttributeNode    = clib.AttributeNode
	TextNode         = clib.TextNode
	CDataSectionNode = clib.CDataSectionNode
	EntityRefNode    = clib.EntityRefNode
	EntityNode       = clib.EntityNode
	PiNode           = clib.PiNode
	CommentNode      = clib.CommentNode
	DocumentNode     = clib.DocumentNode
	DocumentTypeNode = clib.DocumentTypeNode
	DocumentFragNode = clib.DocumentFragNode
	NotationNode     = clib.NotationNode
	HTMLDocumentNode = clib.HTMLDocumentNode
	DTDNode          = clib.DTDNode
	ElementDecl      = clib.ElementDecl
	AttributeDecl    = clib.AttributeDecl
	EntityDecl       = clib.EntityDecl
	NamespaceDecl    = clib.NamespaceDecl
	XIncludeStart    = clib.XIncludeStart
	XIncludeEnd      = clib.XIncludeEnd
	DocbDocumentNode = clib.DocbDocumentNode
)

type XMLNode struct {
	ptr    uintptr // *C.xmlNode
	mortal bool
}

type Attribute struct {
	*XMLNode
}

type CDataSection struct {
	*XMLNode
}

type Comment struct {
	*XMLNode
}

type Element struct {
	*XMLNode
}

type Document struct {
	ptr    uintptr // *C.xmlDoc
	mortal bool
}

type Text struct {
	*XMLNode
}

type Namespace struct {
	*XMLNode
}
