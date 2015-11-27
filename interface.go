package libxml2

/*
#cgo pkg-config: libxml-2.0
#include "libxml/tree.h"
#include "libxml/xpath.h"
#include <libxml/xpathInternals.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

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
	ErrAttributeNotFound             = errors.New("attribute not found")
	ErrInvalidArgument               = errors.New("invalid argument")
	ErrInvalidDocument               = errors.New("invalid document")
	ErrInvalidParser                 = errors.New("invalid parser")
	ErrInvalidNode                   = errors.New("invalid node")
	ErrInvalidNodeName               = errors.New("invalid node name")
	ErrInvalidNodeType               = errors.New("invalid node type")
	ErrInvalidXPathExpression        = errors.New("empty xpath expression")
	ErrMalformedXML                  = errors.New("malformed XML")
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

type Libxml2Node interface {
	// Pointer() returns the underlying C pointer. This is an exported
	// method to allow various internal go-libxml2 packages to interoperate
	// on each other. End users are STRONGLY advised not to touch this
	// method or its return values
	Pointer() unsafe.Pointer
	ParseInContext(string, ParseOption) (Node, error)
}

// Node defines the basic DOM interface
type Node interface {
	Libxml2Node
	AddChild(Node) error
	ChildNodes() (NodeList, error)
	Copy() (Node, error)
	OwnerDocument() (*Document, error)
	FindNodes(string) (NodeList, error)
	FirstChild() (Node, error)
	Free()
	HasChildNodes() bool
	IsSameNode(Node) bool
	LastChild() (Node, error)
	// Literal is almost the same as String(), except for things like Element
	// and Attribute nodes. String() will return the XML stringification of
	// these, but Literal() will return the "value" associated with them.
	Literal() (string, error)
	NextSibling() (Node, error)
	NodeName() string
	NodeType() XMLNodeType
	NodeValue() string
	ParetNode() (Node, error)
	PreviousSibling() (Node, error)
	SetDocument(d *Document) error
	SetNodeName(string)
	SetNodeValue(string)
	String() string
	TextContent() string
	ToString(int, bool) string
	Walk(func(Node) error) error
}

type NodeList []Node

type XMLNode struct {
	ptr    *C.xmlNode
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
	ptr *C.xmlDoc
}

type Text struct {
	*XMLNode
}

type XPathObjectType int

const (
	XPathUndefined XPathObjectType = iota
	XPathNodeSet
	XPathBoolean
	XPathNumber
	XPathString
	XPathPoint
	XPathRange
	XPathLocationSet
	XPathUSers
	XPathXsltTree
)

type XPathObject struct {
	ptr *C.xmlXPathObject
	// This flag controls if the StringValue should use the *contents* (literal value)
	// of the nodeset instead of stringifying the node
	ForceLiteral bool
}

type XPathContext struct {
	ptr *C.xmlXPathContext
}

// XPathExpression is a compiled XPath.
type XPathExpression struct {
	ptr *C.xmlXPathCompExpr
	// This exists mainly for debugging purposes
	expr string
}

// ParseOption represents the parser option bit
type ParseOption int

const (
	XMLParserRecover    ParseOption = 1 << iota /* recover on errors */
	XMLParserNoEnt                              /* substitute entities */
	XMLParserDTDLoad                            /* load the external subset */
	XMLParserDTDAttr                            /* default DTD attributes */
	XMLParserDTDValid                           /* validate with the DTD */
	XMLParserNoError                            /* suppress error reports */
	XMLParserNoWarning                          /* suppress warning reports */
	XMLParserPedantic                           /* pedantic error reporting */
	XMLParserNoBlanks                           /* remove blank nodes */
	XMLParserSAX1                               /* use the SAX1 interface internally */
	XMLParserXInclude                           /* Implement XInclude substitition  */
	XMLParserNoNet                              /* Forbid network access */
	XMLParserNoDict                             /* Do not reuse the context dictionnary */
	XMLParserNsclean                            /* remove redundant namespaces declarations */
	XMLParserNoCDATA                            /* merge CDATA as text nodes */
	XMLParserNoXIncNode                         /* do not generate XINCLUDE START/END nodes */
	XMLParserCompact                            /* compact small text nodes; no modification of the tree allowed afterwards (will possibly crash if you try to modify the tree) */
	XMLParserOld10                              /* parse using XML-1.0 before update 5 */
	XMLParserNoBaseFix                          /* do not fixup XINCLUDE xml:base uris */
	XMLParserHuge                               /* relax any hardcoded limit from the parser */
	XMLParserOldSAX                             /* parse using SAX2 interface before 2.7.0 */
	XMLParserIgnoreEnc                          /* ignore internal document encoding hint */
	XMLParserBigLines                           /* Store big lines numbers in text PSVI field */
	XMLParserMax
	XMLParserEmptyOption ParseOption = 0
)

type ParserCtxt struct {
	ptr *C.xmlParserCtxt
}

type Parser struct {
	Options ParseOption
}

type Namespace struct {
	*XMLNode
}

type Serializer interface {
	Serialize(interface{}) (string, error)
}

// note: Serialize takes an interface because some serializers only allow
// Document, whereas others might allow Nodes

// C14NMode represents the C14N mode supported by libxml2
type C14NMode int

const (
	C14N1_0 C14NMode = iota
	C14NExclusive1_0
	C14N1_1
)

// C14NSerialize implements the Serializer interface, and generates
// XML in C14N format.
type C14NSerialize struct {
	Mode         C14NMode
	WithComments bool
}
