package libxml2

import (
	"errors"
	"unsafe"
)

// HTMLParseOption represents the HTML parser options that
// can be used when parsing HTML
type HTMLParseOption int

const (
	// HTMLParseRecover enables relaxed parsing
	HTMLParseRecover HTMLParseOption = 1 << 0
	// HTMLParseNoDefDTD disables using a default doctype when absent
	HTMLParseNoDefDTD = 1 << 2
	// HTMLParseNoError suppresses error reports
	HTMLParseNoError = 1 << 5
	// HTMLParseNoWarning suppresses warning reports
	HTMLParseNoWarning = 1 << 6
	// HTMLParsePedantic enables pedantic error reporting
	HTMLParsePedantic = 1 << 7
	// HTMLParseNoBlanks removes blank nodes
	HTMLParseNoBlanks = 1 << 8
	// HTMLParseNoNet forbids network access during parsing
	HTMLParseNoNet = 1 << 11
	// HTMLParseNoImplied disables implied html/body elements
	HTMLParseNoImplied = 1 << 13
	// HTMLParseCompact enables compaction of small text nodes
	HTMLParseCompact = 1 << 16
	// HTMLParseIgnoreEnc ignores internal document encoding hints
	HTMLParseIgnoreEnc = 1 << 21
)

// DefaultHTMLParseOptions represents the default set of options
// used in the ParseHTML* functions
const DefaultHTMLParseOptions = HTMLParseCompact | HTMLParseNoBlanks | HTMLParseNoError | HTMLParseNoWarning

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
	ErrAttributeNotFound             = errors.New("attribute not found")
	ErrInvalidArgument               = errors.New("invalid argument")
	ErrInvalidDocument               = errors.New("invalid document")
	ErrInvalidParser                 = errors.New("invalid parser")
	ErrInvalidNamespace              = errors.New("invalid namespace")
	ErrInvalidNode                   = errors.New("invalid node")
	ErrInvalidNodeName               = errors.New("invalid node name")
	ErrInvalidNodeType               = errors.New("invalid node type")
	ErrInvalidXPathContext           = errors.New("invalid xpath context")
	ErrInvalidXPathExpression        = errors.New("invalid xpath expression")
	ErrInvalidXPathObject            = errors.New("invalid xpath object")
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
	ParentNode() (Node, error)
	PreviousSibling() (Node, error)
	SetDocument(d *Document) error
	SetNodeName(string)
	SetNodeValue(string)
	String() string
	TextContent() string
	ToString(int, bool) string
	Walk(func(Node) error) error

	MakeMortal()
	MakePersistent()
	AutoFree()
}

type NodeList []Node

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
	ptr uintptr // *C.xmlDoc
}

type Text struct {
	*XMLNode
}

type XPathResult interface {
	Bool() bool
	Free()
	NodeList() NodeList
	Number() float64
	String() string
	Type() XPathObjectType
	// Valid returns true if the underlying XPathObject is valid,
	// that is, the XPath evaluation actually succeeded. If this
	// returns false, it is most likely that there was a problem
	// with your XPath, or somehow XPathContext/XPathExpression
	// was corrupted.
	Valid() bool
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
	XPathUsers
	XPathXsltTree
)

// InvalidXPathObject represents an invalid result as a result of the
// XPathEvaluation -- that is, either there was a problem in the
// XPathContext, the XPathExpression, or the actually XPath was invalid.
//
// This object is returned from FindValue/FindNode so that you can
// immediatelly call StringValue/BoolValue/etc on the result of those methods
// without having to check for a second error return value
type InvalidXPathObject struct {}

type XPathObject struct {
	ptr uintptr // *C.xmlXPathObject
	// This flag controls if the StringValue should use the *contents* (literal value)
	// of the nodeset instead of stringifying the node
	ForceLiteral bool
}

type XPathContext struct {
	ptr uintptr // *C.xmlXPathContext
	err error
}

// XPathExpression is a compiled XPath.
type XPathExpression struct {
	ptr uintptr // *C.xmlXPathCompExpr
	// This exists mainly for debugging purposes
	expr string
}

// ParseOption represents the parser option bit
type ParseOption int

const (
	XMLParseRecover    ParseOption = 1 << iota /* recover on errors */
	XMLParseNoEnt                              /* substitute entities */
	XMLParseDTDLoad                            /* load the external subset */
	XMLParseDTDAttr                            /* default DTD attributes */
	XMLParseDTDValid                           /* validate with the DTD */
	XMLParseNoError                            /* suppress error reports */
	XMLParseNoWarning                          /* suppress warning reports */
	XMLParsePedantic                           /* pedantic error reporting */
	XMLParseNoBlanks                           /* remove blank nodes */
	XMLParseSAX1                               /* use the SAX1 interface internally */
	XMLParseXInclude                           /* Implement XInclude substitition  */
	XMLParseNoNet                              /* Forbid network access */
	XMLParseNoDict                             /* Do not reuse the context dictionnary */
	XMLParseNsclean                            /* remove redundant namespaces declarations */
	XMLParseNoCDATA                            /* merge CDATA as text nodes */
	XMLParseNoXIncNode                         /* do not generate XINCLUDE START/END nodes */
	XMLParseCompact                            /* compact small text nodes; no modification of the tree allowed afterwards (will possibly crash if you try to modify the tree) */
	XMLParseOld10                              /* parse using XML-1.0 before update 5 */
	XMLParseNoBaseFix                          /* do not fixup XINCLUDE xml:base uris */
	XMLParseHuge                               /* relax any hardcoded limit from the parser */
	XMLParseOldSAX                             /* parse using SAX2 interface before 2.7.0 */
	XMLParseIgnoreEnc                          /* ignore internal document encoding hint */
	XMLParseBigLines                           /* Store big lines numbers in text PSVI field */
	XMLParseMax
	XMLParseEmptyOption ParseOption = 0
)

type ParserCtxt struct {
	ptr uintptr // *C.xmlParserCtxt
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
