package libxml2

import "fmt"

const _XPathObjectTypeName = "XPathUndefinedXPathNodeSetXPathBooleanXPathNumberXPathStringXPathPointXPathRangeXPathLocationSetXPathUSersXPathXsltTree"

var _XPathObjectTypeIndex = [...]uint8{0, 14, 26, 38, 49, 60, 70, 80, 96, 106, 119}

// String returns the stringified version of XPathObjectType
func (i XPathObjectType) String() string {
	if i < 0 || i+1 >= XPathObjectType(len(_XPathObjectTypeIndex)) {
		return fmt.Sprintf("XPathObjectType(%d)", i)
	}
	return _XPathObjectTypeName[_XPathObjectTypeIndex[i]:_XPathObjectTypeIndex[i+1]]
}

// Bool for InvalidXPathObject always return false
func (x InvalidXPathObject) Bool() bool { return false }

// Number for InvalidXPathObject always return 0
func (x InvalidXPathObject) Number() float64 { return 0 }

// Free for InvalidXPathObject is always a no-op
func (x InvalidXPathObject) Free() {}

// NodeList for InvalidXPathObject always returns nil
func (x InvalidXPathObject) NodeList() NodeList { return nil }

// String for InvalidXPathObject always returns ""
func (x InvalidXPathObject) String() string { return "" }

// Type for InvalidXPathOBject always returns XPathUndefined
func (x InvalidXPathObject) Type() XPathObjectType { return XPathUndefined }

// Valid for InvalidXPathObject always returns false
func (x InvalidXPathObject) Valid() bool { return false }

// Valid for XPathObject always returns true
func (x XPathObject) Valid() bool {
	return true
}

// Type returns the XPathObjectType
func (x XPathObject) Type() XPathObjectType {
	return xmlXPathObjectType(&x)
}

// Number returns the floatval component of the XPathObject as float64
func (x XPathObject) Number() float64 {
	return xmlXPathObjectFloat64(&x)
}

// Bool returns the boolval component of the XPathObject
func (x XPathObject) Bool() bool {
	return xmlXPathObjectBool(&x)
}

// NodeList returns the list of nodes included in this XPathObject
func (x XPathObject) NodeList() NodeList {
	nl, err := xmlXPathObjectNodeList(&x)
	if err != nil {
		return nil
	}
	return nl
}

// String returns the stringified value of the nodes included in
// this XPathObject. If the XPathObject is anything other than a
// NodeSet, then we fallback to using fmt.Sprintf to generate
// some sort of readable output
func (x XPathObject) String() string {
	switch x.Type() {
	case XPathNodeSet:
		nl := x.NodeList()
		if nl == nil {
			return ""
		}
		if x.ForceLiteral {
			s, err := nl.Literal()
			if err == nil {
				return s
			}
			return ""
		}
		return nl.NodeValue()
	default:
		return fmt.Sprintf("%v", x)
	}
}

// Free releases the underlying C structs
func (x *XPathObject) Free() {
	xmlXPathFreeObject(x)
}

// NewXPathExpression compiles the given XPath expression string
func NewXPathExpression(s string) (*XPathExpression, error) {
	return xmlXPathCompile(s)
}

// Free releases the underlying C structs in the XPathExpression
func (x *XPathExpression) Free() {
	xmlXPathFreeCompExpr(x)
}

// NewXPathContext creates a new XPathContext, optionally providing
// with a context node.
//
// Note that although we are specifying `n... Node` for the argument,
// only the first, node is considered for the context node
func NewXPathContext(n ...Node) (*XPathContext, error) {
	return xmlXPathNewContext(n...)
}

// SetContextNode sets or resets the context node which
// XPath expressions will be evaluated against.
func (x *XPathContext) SetContextNode(n Node) error {
	return xmlXPathContextSetContextNode(x, n)
}

// Exists compiles and evaluates the xpath expression, and returns
// true if a corresponding node exists
func (x *XPathContext) Exists(xpath string) bool {
	res := x.FindValue(xpath)
	defer res.Free()

	if !res.Valid() {
		return false
	}

	obj := res.(*XPathObject)

	switch obj.Type() {
	case XPathNodeSet:
		return xmlXPathObjectNodeListLen(obj) > 0
	default:
		panic("unimplemented")
	}
	return false
}

// Free releases the underlying C structs in the XPathContext
func (x *XPathContext) Free() {
	xmlXPathFreeContext(x)
}

// FindNodes compiles a XPathExpression in string form, and then evaluates.
func (x *XPathContext) FindNodes(s string) (NodeList, error) {
	expr, err := NewXPathExpression(s)
	if err != nil {
		return nil, err
	}
	defer expr.Free()

	return x.FindNodesExpr(expr)
}

// FindNodesExpr evaluates a compiled XPathExpression.
func (x *XPathContext) FindNodesExpr(expr *XPathExpression) (NodeList, error) {
	res, err := evalXPath(x, expr)
	if err != nil {
		return nil, err
	}
	defer res.Free()

	return res.NodeList(), nil
}

// FindValue evaluates the expression s against the nodes registered
// in x. It returns the resulting data evaluated to an XPathResult.
func (x *XPathContext) FindValue(s string) XPathResult {
	expr, err := NewXPathExpression(s)
	if err != nil {
		x.err = err
		return InvalidXPathObject{}
	}
	defer expr.Free()

	return x.FindValueExpr(expr)
}

// LastError returns the error from the last operation
func (x XPathContext) LastError() error {
	return x.err
}

// FindValueExpr evaluates the given XPath expression and returns an XPathObject.
// You must call `Free()` on this returned object
func (x *XPathContext) FindValueExpr(expr *XPathExpression) XPathResult {
	res, err := evalXPath(x, expr)
	if err != nil {
		x.err = err
		return InvalidXPathObject{}
	}

	//	res.ForceLiteral = true
	return res
}

// LookupNamespaceURI looksup the namespace URI associated with prefix
func (x *XPathContext) LookupNamespaceURI(prefix string) (string, error) {
	return xmlXPathNSLookup(x, prefix)
}

// RegisterNS registers a namespace so it can be used in an XPathExpression
func (x *XPathContext) RegisterNS(name, nsuri string) error {
	return xmlXPathRegisterNS(x, name, nsuri)
}
