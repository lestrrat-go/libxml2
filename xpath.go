package libxml2

/*
#cgo pkg-config: libxml-2.0
#include <stdbool.h>
#include "libxml/globals.h"
#include "libxml/xpath.h"
#include <libxml/xpathInternals.h>

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
import "errors"

// This compiles the expression every time. Ponder if we really need it
func findNodes(n Node, xpath string) ([]Node, error) {
	ctx := C.xmlXPathNewContext((*C.xmlNode)(n.pointer()).doc)
	defer C.xmlXPathFreeContext(ctx)

	res := C.xmlXPathEvalExpression(stringToXmlChar(xpath), ctx)
	defer C.xmlXPathFreeObject(res)
	if C.MY_xmlXPathNodeSetIsEmpty(res.nodesetval) {
		return []Node(nil), nil
	}

	ret := make([]Node, res.nodesetval.nodeNr)
	for i := 0; i < int(res.nodesetval.nodeNr); i++ {
		ret[i] = wrapToNode(C.MY_xmlNodeSetTabAt(res.nodesetval.nodeTab, C.int(i)))
	}
	return ret, nil
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

func (x *XPathContext) FindNodes(s string) ([]Node, error) {
	expr, err := NewXPathExpression(s)
	if err != nil {
		return nil, err
	}
	defer expr.Free()

	return x.FindNodesExpr(expr)
}

func (x *XPathContext) FindNodesExpr(expr *XPathExpression) ([]Node, error) {
	if expr == nil {
		return nil, errors.New("empty XPathExpression")
	}

	// If there is no document associated with this context,
	// then xmlXPathCompiledEval() just fails to match
	ctx := x.ptr
	ctx.doc = ctx.node.doc
	if ctx.doc == nil {
		ctx.doc = C.xmlNewDoc(stringToXmlChar("1.0"))
		defer C.xmlFreeDoc(ctx.doc)
	}

	res := C.xmlXPathCompiledEval(expr.ptr, ctx)
	if res == nil {
		return nil, errors.New("empty result")
	}

	if res.nodesetval.nodeNr == 0 {
		return []Node(nil), nil
	}

	ret := make([]Node, res.nodesetval.nodeNr)
	for i := 0; i < int(res.nodesetval.nodeNr); i++ {
		ret[i] = wrapToNode(C.MY_xmlNodeSetTabAt(res.nodesetval.nodeTab, C.int(i)))
	}

	C.xmlXPathFreeObject(res)

	return ret, nil
}
