package node

import "github.com/lestrrat/go-libxml2/clib"

type XPathResult interface {
	Bool() bool
	Free()
	NodeList() List
	Number() float64
	String() string
	Type() clib.XPathObjectType
}


type Document interface {
	Node
	DocumentElement() (Node, error)
	Dump(bool) string
	Encoding() string
}

// Node defines the basic DOM interface
type Node interface {
	// Pointer() returns the underlying C pointer. This is an exported
	// method to allow various internal go-libxml2 packages to interoperate
	// on each other. End users are STRONGLY advised not to touch this
	// method or its return values
	Pointer() uintptr
	ParseInContext(string, int) (Node, error)

	AddChild(Node) error
	ChildNodes() (List, error)
	Copy() (Node, error)
	OwnerDocument() (Document, error)
	FindValue(string) (XPathResult, error)
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
	NodeType() clib.XMLNodeType
	NodeValue() string
	ParentNode() (Node, error)
	PreviousSibling() (Node, error)
	SetDocument(d Document) error
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

type List []Node
