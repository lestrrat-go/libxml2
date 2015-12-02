package xpath

import (
	"github.com/lestrrat/go-libxml2/clib"
	"github.com/lestrrat/go-libxml2/node"
)

const (
	UndefinedType   = clib.XPathUndefinedType
	NodeSetType     = clib.XPathNodeSetType
	BooleanType     = clib.XPathBooleanType
	NumberType      = clib.XPathNumberType
	StringType      = clib.XPathStringType
	PointType       = clib.XPathPointType
	RangeType       = clib.XPathRangeType
	LocationSetType = clib.XPathLocationSetType
	UsersType       = clib.XPathUsersType
	XSLTTreeType    = clib.XPathXSLTTreeType
)

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

type Result node.XPathResult
