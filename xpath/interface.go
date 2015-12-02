package xpath

import (
	"github.com/lestrrat/go-libxml2/clib"
	"github.com/lestrrat/go-libxml2/node"
)

const (
	UndefinedType   ObjectType = ObjectType(clib.XPathUndefinedType)
	NodeSetType                = ObjectType(clib.XPathNodeSetType)
	BooleanType                = ObjectType(clib.XPathBooleanType)
	NumberType                 = ObjectType(clib.XPathNumberType)
	StringType                 = ObjectType(clib.XPathStringType)
	PointType                  = ObjectType(clib.XPathPointType)
	RangeType                  = ObjectType(clib.XPathRangeType)
	LocationSetType            = ObjectType(clib.XPathLocationSetType)
	UsersType                  = ObjectType(clib.XPathUsersType)
	XSLTTreeType               = ObjectType(clib.XPathXSLTTreeType)
)

type ObjectType clib.XPathObjectType
type Object struct {
	ptr uintptr // *C.xmlObject
	// This flag controls if the StringValue should use the *contents* (literal value)
	// of the nodeset instead of stringifying the node
	ForceLiteral bool
}

type Context struct {
	ptr uintptr // *C.xmlContext
}

// Expression is a compiled XPath expression
type Expression struct {
	ptr uintptr // *C.xmlCompExpr
	// This exists mainly for debugging purposes
	expr string
}

type Result interface {
	Bool() bool
	Free()
	NodeList() node.List
	Number() float64
	String() string
	Type() ObjectType
}


