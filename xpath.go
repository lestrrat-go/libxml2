package libxml2

/*
#cgo pkg-config: libxml-2.0
#include <stdbool.h>
#include "libxml/globals.h"
#include "libxml/xpath.h"
#include <libxml/xpathInternals.h>

// Macro wrapper function
static inline void MY_xmlFree(void *p) {
	xmlFree(p);
}

// Macro wrapper function
static inline bool MY_xmlXPathNodeSetIsEmpty(xmlNodeSetPtr ptr) {
	return xmlXPathNodeSetIsEmpty(ptr);
}

// Because Go can't do pointer airthmetics...
static inline xmlNodePtr MY_xmlNodeSetTabAt(xmlNodePtr *nodes, int i) {
	return nodes[i];
}

*/
import "C"
import (
	"errors"
	"fmt"
)

type XPathContext struct {
	ptr *C.xmlXPathContext
}

// XPathExpression is a compiled XPath.
type XPathExpression struct {
	ptr *C.xmlXPathCompExpr
	// This exists mainly for debugging purposes
	expr string
}

func NewXPathExpression(s string) (*XPathExpression, error) {
	p := C.xmlXPathCompile(stringToXmlChar(s))
	if p == nil {
		return nil, errors.New("xpath compilation failed")
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

	if len(n) > 0 && n[0] != nil {
		ctx.node = (*C.xmlNode)(n[0].pointer())
	}
	return &XPathContext{ptr: ctx}, nil
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
		return nil, errors.New("empty XPathExpression")
	}

	// If there is no document associated with this context,
	// then xmlXPathCompiledEval() just fails to match
	ctx := x.ptr

	if ctx.node != nil && ctx.node.doc != nil {
		ctx.doc = ctx.node.doc
	}

	if ctx.doc == nil {
		ctx.doc = C.xmlNewDoc(stringToXmlChar("1.0"))
		defer C.xmlFreeDoc(ctx.doc)
	}

	res := C.xmlXPathCompiledEval(expr.ptr, ctx)
	if res == nil {
		return nil, errors.New("empty result")
	}

	return &XPathObject{res}, nil
}

func (x *XPathContext) FindNodesExpr(expr *XPathExpression) (NodeList, error) {
	res, err := x.evalXPath(expr)
	if err != nil {
		return nil, err
	}
	defer res.Free()

	return res.NodeList(), nil
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

const _XPathObjectType_name = "XPathUndefinedXPathNodeSetXPathBooleanXPathNumberXPathStringXPathPointXPathRangeXPathLocationSetXPathUSersXPathXsltTree"

var _XPathObjectType_index = [...]uint8{0, 14, 26, 38, 49, 60, 70, 80, 96, 106, 119}

func (i XPathObjectType) String() string {
	if i < 0 || i+1 >= XPathObjectType(len(_XPathObjectType_index)) {
		return fmt.Sprintf("XPathObjectType(%d)", i)
	}
	return _XPathObjectType_name[_XPathObjectType_index[i]:_XPathObjectType_index[i+1]]
}

type XPathObject struct {
	ptr *C.xmlXPathObject
}

func (x XPathObject) Type() XPathObjectType {
	return XPathObjectType(x.ptr._type)
}

func (x XPathObject) Float64() float64 {
	return float64(x.ptr.floatval)
}

func (x XPathObject) Bool() bool {
	return C.int(x.ptr.boolval) == 1
}

func (x XPathObject) NodeList() NodeList {
	if x.ptr.nodesetval.nodeNr == 0 {
		return NodeList(nil)
	}

	ret := make(NodeList, x.ptr.nodesetval.nodeNr)
	for i := 0; i < int(x.ptr.nodesetval.nodeNr); i++ {
		ret[i] = wrapToNode(C.MY_xmlNodeSetTabAt(x.ptr.nodesetval.nodeTab, C.int(i)))
	}

	return ret
}

func (x XPathObject) String() string {
	switch x.Type() {
	case XPathNodeSet:
		return x.NodeList().String()
	default:
		return fmt.Sprintf("%v", x)
	}
}

func (x *XPathObject) Free() {
	//	if x.ptr.nodesetval != nil {
	//		C.xmlXPathFreeNodeSet(x.ptr.nodesetval)
	//	}
	C.xmlXPathFreeObject(x.ptr)
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

	return res, nil
}

func (x *XPathContext) LookupNamespaceURI(name string) (string, error) {
	s := C.xmlXPathNsLookup(x.ptr, stringToXmlChar(name))
	if s == nil {
		return "", errors.New("not found")
	}
	return xmlCharToString(s), nil
}

func (x *XPathContext) RegisterNs(name, nsuri string) error {
	res := C.xmlXPathRegisterNs(x.ptr, stringToXmlChar(name), stringToXmlChar(nsuri))
	if res == -1 {
		return errors.New("cannot register namespace")
	}
	return nil
}
