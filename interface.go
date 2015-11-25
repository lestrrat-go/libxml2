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

type XmlNodeType int

const (
	ElementNode XmlNodeType = iota + 1
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
	ErrNodeNotFound    = errors.New("node not found")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrInvalidParser   = errors.New("invalid parser")
	ErrInvalidNode     = errors.New("invalid node")
	ErrInvalidNodeName = errors.New("invalid node name")
)

type ptr interface {
	// Pointer() returns the underlying C pointer. This is an exported
	// method to allow various internal go-libxml2 packages to interoperate
	// on each other. End users are STRONGLY advised not to touch this
	// method or its return values
	Pointer() unsafe.Pointer
}

// Node defines the basic DOM interface
type Node interface {
	ptr
	AddChild(Node) error
	AppendChild(Node) error
	ChildNodes() (NodeList, error)
	OwnerDocument() *Document
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
	NodeType() XmlNodeType
	NodeValue() string
	ParetNode() (Node, error)
	PreviousSibling() (Node, error)
	SetNodeName(string)
	SetNodeValue(string)
	String() string
	TextContent() string
	ToString(int, bool) string
	Walk(func(Node) error) error
}

type NodeList []Node

type XmlNode struct {
	ptr *C.xmlNode
	mortal bool
}

type Attribute struct {
	*XmlNode
}

type CDataSection struct {
	*XmlNode
}

type Comment struct {
	*XmlNode
}

type Element struct {
	*XmlNode
}

type Document struct {
	ptr *C.xmlDoc
}

type Text struct {
	*XmlNode
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
	XmlParseRecover    ParseOption = 1 << iota /* recover on errors */
	XmlParseNoEnt                              /* substitute entities */
	XmlParseDTDLoad                            /* load the external subset */
	XmlParseDTDAttr                            /* default DTD attributes */
	XmlParseDTDValid                           /* validate with the DTD */
	XmlParseNoError                            /* suppress error reports */
	XmlParseNoWarning                          /* suppress warning reports */
	XmlParsePedantic                           /* pedantic error reporting */
	XmlParseNoBlanks                           /* remove blank nodes */
	XmlParseSAX1                               /* use the SAX1 interface internally */
	XmlParseXInclude                           /* Implement XInclude substitition  */
	XmlParseNoNet                              /* Forbid network access */
	XmlParseNoDict                             /* Do not reuse the context dictionnary */
	XmlParseNsclean                            /* remove redundant namespaces declarations */
	XmlParseNoCDATA                            /* merge CDATA as text nodes */
	XmlParseNoXIncNode                         /* do not generate XINCLUDE START/END nodes */
	XmlParseCompact                            /* compact small text nodes; no modification of the tree allowed afterwards (will possibly crash if you try to modify the tree) */
	XmlParseOld10                              /* parse using XML-1.0 before update 5 */
	XmlParseNoBaseFix                          /* do not fixup XINCLUDE xml:base uris */
	XmlParseHuge                               /* relax any hardcoded limit from the parser */
	XmlParseOldSAX                             /* parse using SAX2 interface before 2.7.0 */
	XmlParseIgnoreEnc                          /* ignore internal document encoding hint */
	XmlParseBigLines                           /* Store big lines numbers in text PSVI field */
	XmlParseMax
	XmlParseEmptyOption ParseOption = 0
)

type ParserCtxt struct {
	ptr *C.xmlParserCtxt
}

type Parser struct {
	Options ParseOption
}

type Namespace struct {
	*XmlNode
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
