package libxml2

/*
#cgo pkg-config: libxml-2.0
#include <stdbool.h>
#include "libxml/tree.h"
#include "libxml/parser.h"
#include "libxml/xpath.h"

static inline bool MY_xmlXPathNodeSetIsEmpty(xmlNodeSetPtr ptr) {
	return ptr == NULL ||
		ptr->nodeNr == 0 ||
		ptr->nodeTab == NULL;
}

static inline xmlNodePtr MY_xmlNodeSetTabAt(xmlNodePtr *nodes, int i) {
	return nodes[i];
}

*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

type XmlElementType int

const (
	XmlElementNode XmlElementType = iota + 1
	XmlAttributeNode
	XmlTextNode
	XmlCDataSectionNode
	XmlEntityRefNode
	XmlEntityNode
	XmlPiNode
	XmlCommentNode
	XmlDocumentNode
	XmlDocumentTypeNode
	XmlDocumentFragNode
	XmlNotationNode
	XmlHTMLDocumentNode
	XmlDTDNode
	XmlElementDecl
	XmlAttributeDecl
	XmlEntityDecl
	XmlNamespaceDecl
	XmlXIncludeStart
	XmlXIncludeEnd
	XmlDocbDocumentNode
)

var _XmlElementType_index = [...]uint8{0, 11, 24, 32, 48, 61, 71, 77, 88, 100, 116, 132, 144, 160, 167, 178, 191, 201, 214, 227, 238, 254}

const _XmlElementType_name = `ElementNodeAttributeNodeTextNodeCDataSectionNodeEntityRefNodeEntityNodePiNodeCommentNodeDocumentNodeDocumentTypeNodeDocumentFragNodeNotationNodeHTMLDocumentNodeDTDNodeElementDeclAttributeDeclEntityDeclNamespaceDeclXIncludeStartXIncludeEndDocbDocumentNode`

func (i XmlElementType) String() string {
	i -= 1
	if i < 0 || i+1 >= XmlElementType(len(_XmlElementType_index)) {
		return fmt.Sprintf("XmlElementType(%d)", i+1)
	}
	return _XmlElementType_name[_XmlElementType_index[i]:_XmlElementType_index[i+1]]
}

var ErrNodeNotFound = errors.New("node not found")
var ErrInvalidArgument = errors.New("invalid argument")

// Node defines the basic DOM interface
type Node interface {
	// pointer() returns the underlying C pointer. Only we are allowed to
	// slice it, dice it, do whatever the heck with it.
	pointer() unsafe.Pointer

	ChildNodes() []Node
	OwnerDocument() *XmlDoc
	FindNodes(string) ([]Node, error)
	IsSameNode(Node) bool
	LastChild() Node
	NodeName() string
	NextSibling() Node
	ParetNode() Node
	PreviousSibling() Node
	SetNodeName(string)
	TextContent() string
	Type() XmlElementType
	Walk(func(Node) error)
}

type xmlNode struct {
	ptr *C.xmlNode
}

type XmlNode struct {
	*xmlNode
}

type XmlElement struct {
	*XmlNode
}

type XmlDoc struct {
	ptr  *C.xmlDoc
	root *C.xmlNode
}

type XmlText struct {
	*XmlNode
}

func wrapXmlElement(n *C.xmlElement) *XmlElement {
	return &XmlElement{wrapXmlNode((*C.xmlNode)(unsafe.Pointer(n)))}
}

func wrapXmlNode(n *C.xmlNode) *XmlNode {
	return &XmlNode{
		&xmlNode{
			ptr: (*C.xmlNode)(unsafe.Pointer(n)),
		},
	}
}

func wrapToNode(n *C.xmlNode) Node {
	switch XmlElementType(n._type) {
	case XmlElementNode:
		return wrapXmlElement((*C.xmlElement)(unsafe.Pointer(n)))
	case XmlTextNode:
		return &XmlText{&XmlNode{&xmlNode{ptr: n}}}
	default:
		return &XmlNode{&xmlNode{ptr: n}}
	}
}

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

func (n *xmlNode) pointer() unsafe.Pointer {
	return unsafe.Pointer(n.ptr)
}

func (n *xmlNode) ChildNodes() []Node {
	return childNodes(n)
}

func wrapXmlDoc(n *C.xmlDoc) *XmlDoc {
	r := C.xmlDocGetRootElement(n) // XXX Should check for n == nil
	return &XmlDoc{ptr: n, root: r}
}

func (n *xmlNode) OwnerDocument() *XmlDoc {
	return wrapXmlDoc(n.ptr.doc)
}

func (n *xmlNode) FindNodes(xpath string) ([]Node, error) {
	return findNodes(n, xpath)
}

func (n *xmlNode) IsSameNode(other Node) bool {
	return n.pointer() == other.pointer()
}

func (n *xmlNode) LastChild() Node {
	return wrapToNode(n.ptr.last)
}

func (n *xmlNode) NodeName() string {
	return xmlCharToString(n.ptr.name)
}

func (n *xmlNode) NextSibling() Node {
	return wrapToNode(n.ptr.next)
}

func (n *xmlNode) ParetNode() Node {
	return wrapToNode(n.ptr.parent)
}

func (n *xmlNode) PreviousSibling() Node {
	return wrapToNode(n.ptr.prev)
}

func (n *xmlNode) SetNodeName(name string) {
	C.xmlNodeSetName(n.ptr, stringToXmlChar(name))
}

func (n *xmlNode) TextContent() string {
	return xmlCharToString(C.xmlXPathCastNodeToString(n.ptr))
}

func (n *xmlNode) Type() XmlElementType {
	return XmlElementType(n.ptr._type)
}

func (n *xmlNode) Walk(fn func(Node) error) {
	panic("should not call walk on internal struct")
}

func (n *XmlNode) Walk(fn func(Node) error) {
	walk(n, fn)
}

func walk(n Node, fn func(Node) error) {
	if err := fn(n); err != nil {
		return
	}
	for _, c := range n.ChildNodes() {
		walk(c, fn)
	}
}

func childNodes(n Node) []Node {
	ret := []Node(nil)
	for chld := ((*C.xmlNode)(n.pointer())).children; chld != nil; chld = chld.next {
		ret = append(ret, wrapToNode(chld))
	}
	return ret
}

func (d *XmlDoc) pointer() unsafe.Pointer {
	return unsafe.Pointer(d.ptr)
}

func (d *XmlDoc) DocumentElement() Node {
	if d.ptr == nil || d.root == nil {
		return nil
	}

	return wrapToNode(d.root)
}

func (d *XmlDoc) FindNodes(xpath string) ([]Node, error) {
	root := d.DocumentElement()
	if root == nil {
		return nil, ErrNodeNotFound
	}
	return root.FindNodes(xpath)
}

func (d *XmlDoc) Encoding() string {
	return xmlCharToString(d.ptr.encoding)
}

func (d *XmlDoc) Free() {
	C.xmlFreeDoc(d.ptr)
	d.ptr = nil
	d.root = nil
}

func (d *XmlDoc) String() string {
	var xc *C.xmlChar
	i := C.int(0)
	C.xmlDocDumpMemory(d.ptr, &xc, &i)
	return xmlCharToString(xc)
}

func (d *XmlDoc) Type() XmlElementType {
	return XmlElementType(d.ptr._type)
}

func (n *XmlDoc) Walk(fn func(Node) error) {
	walk(wrapXmlNode(n.root), fn)
}

func (n *XmlText) Data() string {
	return xmlCharToString(n.ptr.content)
}

func (n *XmlText) Walk(fn func(Node) error) {
	walk(n, fn)
}
