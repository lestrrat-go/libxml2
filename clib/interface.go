package clib

import "errors"

// C14NMode represents the C14N mode supported by libxml2
type C14NMode int
type PtrSource interface {
	Pointer() uintptr
}

// XMLNodeType identifies the type of the underlying C struct
type XMLNodeType int

const (
	ElementNode XMLNodeType = iota + 1
	AttributeNode
	TextNode
	CDataSectionNode
	EntityRefNode
	EntityNode
	PiNode
	CommentNode
	DocumentNode
	DocumentTypeNode
	DocumentFragNode
	NotationNode
	HTMLDocumentNode
	DTDNode
	ElementDecl
	AttributeDecl
	EntityDecl
	NamespaceDecl
	XIncludeStart
	XIncludeEnd
	DocbDocumentNode
)

var (
	ErrInvalidAttribute              = errors.New("invalid attribute")
	ErrInvalidArgument               = errors.New("invalid argument")
	ErrInvalidDocument               = errors.New("invalid document")
	ErrInvalidParser                 = errors.New("invalid parser")
	ErrInvalidNamespace              = errors.New("invalid namespace")
	ErrInvalidNode                   = errors.New("invalid node")
	ErrInvalidNodeName               = errors.New("invalid node name")
	ErrInvalidXPathContext           = errors.New("invalid xpath context")
	ErrInvalidXPathExpression        = errors.New("invalid xpath expression")
	ErrInvalidXPathObject            = errors.New("invalid xpath object")
	ErrNodeNotFound                  = errors.New("node not found")
	ErrXPathEmptyResult              = errors.New("empty xpath result")
	ErrXPathCompileFailure           = errors.New("xpath compilation failed")
	ErrXPathNamespaceRegisterFailure = errors.New("cannot register namespace")
)

type ErrNamespaceNotFound struct {
	Target string
}

func (e ErrNamespaceNotFound) Error() string {
	return "namespace not found: " + e.Target
}

type XPathObjectType int

const (
	XPathUndefinedType XPathObjectType = iota
	XPathNodeSetType
	XPathBooleanType
	XPathNumberType
	XPathStringType
	XPathPointType
	XPathRangeType
	XPathLocationSetType
	XPathUsersType
	XPathXSLTTreeType
)
