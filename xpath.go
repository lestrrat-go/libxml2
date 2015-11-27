package libxml2

/*
#cgo pkg-config: libxml-2.0
#include <stdbool.h>
#include "libxml/globals.h"
#include "libxml/xpath.h"
#include <libxml/xpathInternals.h>

// Because Go can't do pointer airthmetics...
static inline xmlNodePtr MY_xmlNodeSetTabAt(xmlNodePtr *nodes, int i) {
	return nodes[i];
}

*/
import "C"
import "fmt"

const _XPathObjectType_name = "XPathUndefinedXPathNodeSetXPathBooleanXPathNumberXPathStringXPathPointXPathRangeXPathLocationSetXPathUSersXPathXsltTree"

var _XPathObjectType_index = [...]uint8{0, 14, 26, 38, 49, 60, 70, 80, 96, 106, 119}

func (i XPathObjectType) String() string {
	if i < 0 || i+1 >= XPathObjectType(len(_XPathObjectType_index)) {
		return fmt.Sprintf("XPathObjectType(%d)", i)
	}
	return _XPathObjectType_name[_XPathObjectType_index[i]:_XPathObjectType_index[i+1]]
}

func (x XPathObject) Type() XPathObjectType {
	return XPathObjectType(x.ptr._type)
}

func (x XPathObject) Float64Value() float64 {
	return float64(x.ptr.floatval)
}

func (x XPathObject) BoolValue() bool {
	return C.int(x.ptr.boolval) == 1
}

func (x XPathObject) NodeList() (NodeList, error) {
	nodeset := x.ptr.nodesetval
	if nodeset == nil {
		return nil, ErrInvalidNode
	}

	if nodeset.nodeNr == 0 {
		return nil, ErrInvalidNode
	}

	ret := make(NodeList, nodeset.nodeNr)
	for i := 0; i < int(nodeset.nodeNr); i++ {
		v, err := wrapToNode(C.MY_xmlNodeSetTabAt(nodeset.nodeTab, C.int(i)))
		if err != nil {
			return nil, err
		}
		ret[i] = v
	}

	return ret, nil
}

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

func (x *XPathObject) Free() {
	//	if x.ptr.nodesetval != nil {
	//		C.xmlXPathFreeNodeSet(x.ptr.nodesetval)
	//	}
	C.xmlXPathFreeObject(x.ptr)
}

func NewXPathExpression(s string) (*XPathExpression, error) {
	p := C.xmlXPathCompile(stringToXMLChar(s))
	if p == nil {
		return nil, ErrXPathCompileFailure
	}

	return &XPathExpression{ptr: p, expr: s}, nil
}

func (x *XPathExpression) Free() {
	if x.ptr == nil {
		return
	}
	C.xmlXPathFreeCompExpr(x.ptr)
}

// Note that although we are specifying `n... Node` for the argument,
// only the first, node is considered for the context node
func NewXPathContext(n ...Node) (*XPathContext, error) {
	ctx := C.xmlXPathNewContext(nil)
	ctx.namespaces = nil

	obj := &XPathContext{ptr: ctx}
	if len(n) > 0 {
		obj.SetContextNode(n[0])
	}
	return obj, nil
}

func (x *XPathContext) SetContextNode(n Node) {
	if n == nil {
		return
	}
	x.ptr.node = (*C.xmlNode)(n.Pointer())
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
	if x.ptr == nil {
		return
	}

	C.xmlXPathFreeContext(x.ptr)
}

func (x *XPathContext) FindNodes(s string) (NodeList, error) {
	expr, err := NewXPathExpression(s)
	if err != nil {
		return nil, err
	}
	defer expr.Free()

	return x.FindNodesExpr(expr)
}

func (x *XPathContext) evalXPath(expr *XPathExpression) (*XPathObject, error) {
	if expr == nil {
		return nil, ErrInvalidXPathExpression
	}

	// If there is no document associated with this context,
	// then xmlXPathCompiledEval() just fails to match
	ctx := x.ptr

	if ctx.node != nil && ctx.node.doc != nil {
		ctx.doc = ctx.node.doc
	}

	if ctx.doc == nil {
		ctx.doc = C.xmlNewDoc(stringToXMLChar("1.0"))
		defer C.xmlFreeDoc(ctx.doc)
	}

	res := C.xmlXPathCompiledEval(expr.ptr, ctx)
	if res == nil {
		return nil, ErrXPathEmptyResult
	}

	return &XPathObject{ptr: res}, nil
}

func (x *XPathContext) FindNodesExpr(expr *XPathExpression) (NodeList, error) {
	res, err := x.evalXPath(expr)
	if err != nil {
		return nil, err
	}
	defer res.Free()

	return res.NodeList()
}

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
	res, err := x.evalXPath(expr)
	if err != nil {
		return nil, err
	}

	res.ForceLiteral = true
	return res, nil
}

func (x *XPathContext) LookupNamespaceURI(name string) (string, error) {
	s := C.xmlXPathNsLookup(x.ptr, stringToXMLChar(name))
	if s == nil {
		return "", ErrNamespaceNotFound{Target: name}
	}
	return xmlCharToString(s), nil
}

func (x *XPathContext) RegisterNs(name, nsuri string) error {
	res := C.xmlXPathRegisterNs(x.ptr, stringToXMLChar(name), stringToXMLChar(nsuri))
	if res == -1 {
		return ErrXPathNamespaceRegisterFailure
	}
	return nil
}
