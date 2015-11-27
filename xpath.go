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

// Type returns the XPathObjectType
func (x XPathObject) Type() XPathObjectType {
	return xmlXPathObjectType(&x)
}

// Float64Value returns the floatval component of the XPathObject
func (x XPathObject) Float64Value() float64 {
	return xmlXPathObjectFloat64Value(&x)
}

// BoolValue returns the boolval component of the XPathObject
func (x XPathObject) BoolValue() bool {
	return xmlXPathObjectBoolValue(&x)
}

// NodeList returns the list of nodes included in this XPathObject
func (x XPathObject) NodeList() (NodeList, error) {
	return xmlXPathObjectNodeList(&x)
}

// StringValue returns the stringified value of the nodes included in
// this XPathObject. If the XPathObject is anything other than a
// NodeSet, then we fallback to using fmt.Sprintf to generate
// some sort of readable output
func (x XPathObject) StringValue() (string, error) {
	switch x.Type() {
	case XPathNodeSet:
		nl, err := x.NodeList()
		if err != nil {
			return "", err
		}
		if x.ForceLiteral {
			return nl.Literal()
		}
		return nl.String(), nil
	default:
		return fmt.Sprintf("%v", x), nil
	}
}

// Free releases the underlying C structs
func (x *XPathObject) Free() {
	xmlXPathFreeObject(x)
}

func NewXPathExpression(s string) (*XPathExpression, error) {
	return xmlXPathCompile(s)
}

func (x *XPathExpression) Free() {
	xmlXPathFreeCompExpr(x)
}

// Note that although we are specifying `n... Node` for the argument,
// only the first, node is considered for the context node
func NewXPathContext(n ...Node) (*XPathContext, error) {
	return xmlXPathNewContext(n...)
}

func (x *XPathContext) SetContextNode(n Node) error {
	return xmlXPathContextSetContextNode(x, n)
}

func (x *XPathContext) Exists(xpath string) bool {
	res, err := x.FindValue(xpath)
	if err != nil {
		return false
	}
	defer res.Free()

	switch res.Type() {
	case XPathNodeSet:
		return res.ptr.nodesetval.nodeNr > 0
	default:
		panic("unimplemented")
	}
	return false
}

func (x *XPathContext) Free() {
	xmlXPathFreeContext(x)
}

func (x *XPathContext) FindNodes(s string) (NodeList, error) {
	expr, err := NewXPathExpression(s)
	if err != nil {
		return nil, err
	}
	defer expr.Free()

	return x.FindNodesExpr(expr)
}

func (x *XPathContext) FindNodesExpr(expr *XPathExpression) (NodeList, error) {
	res, err := evalXPath(x, expr)
	if err != nil {
		return nil, err
	}
	defer res.Free()

	return res.NodeList()
}

// FindValue evaluates the expression s against the nodes registered
// in x. It returns the resulting data evaluated to an XPathObject.
func (x *XPathContext) FindValue(s string) (*XPathObject, error) {
	expr, err := NewXPathExpression(s)
	if err != nil {
		return nil, err
	}
	defer expr.Free()

	return x.FindValueExpr(expr)
}

// FindValueExpr evaluates the given XPath expression and returns an XPathObject.
// You must call `Free()` on this returned object
func (x *XPathContext) FindValueExpr(expr *XPathExpression) (*XPathObject, error) {
	res, err := evalXPath(x, expr)
	if err != nil {
		return nil, err
	}

	res.ForceLiteral = true
	return res, nil
}

// LookupNamespaceURI looksup the namespace URI associated with prefix 
func (x *XPathContext) LookupNamespaceURI(prefix string) (string, error) {
	return xmlXPathNSLookup(x, prefix)
}

// RegisterNS registers a namespace so it can be used in an XPathExpression
func (x *XPathContext) RegisterNS(name, nsuri string) error {
	return xmlXPathRegisterNS(x, name, nsuri)
}
