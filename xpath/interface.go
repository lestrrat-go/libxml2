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

// InvalidObject represents an invalid result as a result of the
// Evaluation -- that is, either there was a problem in the
// Context, the Expression, or the actually  was invalid.
//
// This object is returned from FindValue/FindNode so that you can
// immediatelly call StringValue/BoolValue/etc on the result of those methods
// without having to check for a second error return value
type InvalidObject struct{}

type ObjectType clib.XPathObjectType
type Object struct {
	ptr uintptr // *C.xmlObject
	// This flag controls if the StringValue should use the *contents* (literal value)
	// of the nodeset instead of stringifying the node
	ForceLiteral bool
}

type Context struct {
	ptr uintptr // *C.xmlContext
	err error
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
	// Valid returns true if the underlying Object is valid,
	// that is, the  evaluation actually succeeded. If this
	// returns false, it is most likely that there was a problem
	// with your , or somehow Context/Expression
	// was corrupted.
	Valid() bool
}


