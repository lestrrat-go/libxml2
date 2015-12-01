/*

Package libxml2 is an interface to libxml2, providing XML and HTML parsers
with DOM interface. The inspiration is Perl5's XML::LibXML module.

This library is still in very early stages of development. API may still change
without notice.

For the time being, the API is being written so that thye are as close as we
can get to DOM Layer 3, but some methods will, for the time being, be punted
and aliases for simpler methods that don't necessarily check for the DOM's
correctness will be used.

Also, the return values are still shaky -- I'm still debating how to handle error cases gracefully.

*/
package libxml2

/*
#cgo pkg-config: libxml-2.0
#include <stdbool.h>
#include <libxml/HTMLparser.h>
#include <libxml/HTMLtree.h>
#include <libxml/globals.h>
#include <libxml/parser.h>
#include <libxml/parserInternals.h>
#include <libxml/tree.h>
#include <libxml/xmlerror.h>
#include <libxml/xpath.h>
#include <libxml/xpathInternals.h>
#include <libxml/c14n.h>


static inline void MY_nilErrorHandler(void *ctx, const char *msg, ...) {}

static inline void MY_xmlSilenceParseErrors() {
	xmlSetGenericErrorFunc(NULL, MY_nilErrorHandler);
}

static inline void MY_xmlDefaultParseErrors() {
	// Checked in the libxml2 source code that using NULL in the second
	// argument restores the default error handler
	xmlSetGenericErrorFunc(NULL, NULL);
}

// Macro wrapper function. cgo cannot detect function-like macros,
// so this is how we avoid it
static inline void MY_xmlFree(void *p) {
	xmlFree(p);
}

// Macro wrapper function. cgo cannot detect function-like macros,
// so this is how we avoid it
static inline xmlError* MY_xmlLastError() {
	return xmlGetLastError();
}

// Change xmlIndentTreeOutput global, return old value, so caller can
// change it back to old value later
static inline int MY_setXmlIndentTreeOutput(int i) {
	int old = xmlIndentTreeOutput;
	xmlIndentTreeOutput = i;
	return old;
}

// Parse a single char out of cur
// Stolen from XML::LibXML
static inline int
MY_parseChar( xmlChar *cur, int *len )
{
	unsigned char c;
	unsigned int val;

	// We are supposed to handle UTF8, check it's valid
	// From rfc2044: encoding of the Unicode values on UTF-8:
	//
	// UCS-4 range (hex.)           UTF-8 octet sequence (binary)
	// 0000 0000-0000 007F   0xxxxxxx
	// 0000 0080-0000 07FF   110xxxxx 10xxxxxx
	// 0000 0800-0000 FFFF   1110xxxx 10xxxxxx 10xxxxxx
	//
	// Check for the 0x110000 limit too

	if ( cur == NULL || *cur == 0 ) {
		*len = 0;
		return(0);
	}

	c = *cur;
	if ( (c & 0x80) == 0 ) {
		*len = 1;
		return((int)c);
	}

	if ((c & 0xe0) == 0xe0) {
		if ((c & 0xf0) == 0xf0) {
			// 4-byte code
			*len = 4;
			val = (cur[0] & 0x7) << 18;
			val |= (cur[1] & 0x3f) << 12;
			val |= (cur[2] & 0x3f) << 6;
			val |= cur[3] & 0x3f;
		} else {
			// 3-byte code
			*len = 3;
			val = (cur[0] & 0xf) << 12;
			val |= (cur[1] & 0x3f) << 6;
			val |= cur[2] & 0x3f;
		}
	} else {
		// 2-byte code
		*len = 2;
		val = (cur[0] & 0x1f) << 6;
		val |= cur[1] & 0x3f;
	}

	if ( !IS_CHAR(val) ) {
		*len = -1;
		return 0;
	}
	return val;
}

// Checks if the given name is a valid name in XML
// stolen from XML::LibXML
static inline int
MY_test_node_name( xmlChar * name )
{
	xmlChar * cur = name;
	int tc  = 0;
	int len = 0;

	if ( cur == NULL || *cur == 0 ) {
		return 0;
	}

	tc = MY_parseChar( cur, &len );

	if ( !( IS_LETTER( tc ) || (tc == '_') || (tc == ':')) ) {
		return 0;
	}

	tc  =  0;
	cur += len;

	while (*cur != 0 ) {
		tc = MY_parseChar( cur, &len );

		if (!(IS_LETTER(tc) || IS_DIGIT(tc) || (tc == '_') ||
				(tc == '-') || (tc == ':') || (tc == '.') ||
				IS_COMBINING(tc) || IS_EXTENDER(tc)) )
		{
			return 0;
		}
		tc = 0;
		cur += len;
	}

	return(1);
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
	"strconv"
	"strings"
	"unsafe"

	"github.com/lestrrat/go-libxml2/internal/debug"
)

// ReportErrors *globally* changes the behavior of reporting errors.
// By default libxml2 spews out lots of data to stderr. When you call
// this function with a `false` value, all those messages are surpressed.
// When you call this function a `true` value, the default behavior is
// restored
func ReportErrors(b bool) {
	if b {
		C.MY_xmlDefaultParseErrors()
	} else {
		C.MY_xmlSilenceParseErrors()
	}
}

func xmlCharToString(s *C.xmlChar) string {
	return C.GoString((*C.char)(unsafe.Pointer(s)))
}

// stringToXMLChar creates a new *C.xmlChar from a Go string.
// Remember to always free this data, as C.CString creates a copy
// of the byte buffer contained in the string
func stringToXMLChar(s string) *C.xmlChar {
	return (*C.xmlChar)(unsafe.Pointer(C.CString(s)))
}

func xmlCreateMemoryParserCtxt(s string, o ParseOption) (*ParserCtxt, error) {
	ctx := C.xmlCreateMemoryParserCtxt(C.CString(s), C.int(len(s)))
	if ctx == nil {
		return nil, errors.New("error creating parser")
	}
	C.xmlCtxtUseOptions(ctx, C.int(o))

	return &ParserCtxt{
		ptr: uintptr(unsafe.Pointer(ctx)),
	}, nil
}

func validParserCtxtPtr(ctx *ParserCtxt) (*C.xmlParserCtxt, error) {
	if ptr := ctx.ptr; ptr != 0 {
		return (*C.xmlParserCtxt)(unsafe.Pointer(ptr)), nil
	}
	return nil, ErrInvalidParser
}

// Parse starts the parsing on the ParserCtxt
func (ctx ParserCtxt) Parse() error {
	ctxptr, err := validParserCtxtPtr(&ctx)
	if err != nil {
		return err
	}

	if C.xmlParseDocument(ctxptr) != C.int(0) {
		return errors.New("parse failed")
	}
	return nil
}

// Free releases the underlying C struct
func (ctx *ParserCtxt) Free() error {
	ctxptr, err := validParserCtxtPtr(ctx)
	if err != nil {
		return err
	}

	C.xmlFreeParserCtxt(ctxptr)
	ctx.ptr = 0

	return nil
}

// WellFormed returns true if the resulting document after parsing
func (ctx ParserCtxt) WellFormed() bool {
	ctxptr, err := validParserCtxtPtr(&ctx)
	if err != nil {
		return false
	}

	return ctxptr.wellFormed == C.int(0)
}

// Document returns the resulting document after parsing
func (ctx ParserCtxt) Document() (*Document, error) {
	ctxptr, err := validParserCtxtPtr(&ctx)
	if err != nil {
		return nil, err
	}

	doc := ctxptr.myDoc
	if doc != nil {
		return wrapDocument(doc), nil
	}
	return nil, errors.New("no document available")
}

func htmlReadDoc(content, url, encoding string, opts int) (*Document, error) {
	// TODO: use htmlCtxReadDoc later, so we can get the error
	ccontent := C.CString(content)
	curl := C.CString(url)
	cencoding := C.CString(encoding)
	defer C.free(unsafe.Pointer(ccontent))
	defer C.free(unsafe.Pointer(curl))
	defer C.free(unsafe.Pointer(cencoding))

	doc := C.htmlReadDoc(
		(*C.xmlChar)(unsafe.Pointer(ccontent)),
		curl,
		cencoding,
		C.int(opts),
	)

	if doc == nil {
		return nil, errors.New("failed to parse document")
	}

	return wrapDocument(doc), nil
}

func createDocument(version, encoding string) *Document {
	cver := stringToXMLChar(version)
	defer C.free(unsafe.Pointer(cver))

	doc := C.xmlNewDoc(cver)
	if encoding != "" {
		cenc := stringToXMLChar(encoding)
		defer C.free(unsafe.Pointer(cenc))

		doc.encoding = C.xmlStrdup(cenc)
	}
	return wrapDocument(doc)
}

func xmlEncodeEntitiesReentrant(doc *Document, s string) (*C.xmlChar, error) {
	docptr, err := validDocumentPtr(doc)
	if err != nil {
		return nil, err
	}

	cent := stringToXMLChar(s)
	defer C.free(unsafe.Pointer(cent))

	return C.xmlEncodeEntitiesReentrant(docptr, cent), nil
}

func myTestNodeName(n string) error {
	cn := stringToXMLChar(n)
	defer C.free(unsafe.Pointer(cn))

	if C.MY_test_node_name(cn) == 0 {
		return ErrInvalidNodeName
	}
	return nil
}

func xmlMakeSafeName(k string) (*C.xmlChar, error) {
	if err := myTestNodeName(k); err != nil {
		return nil, err
	}
	return stringToXMLChar(k), nil
}

func validNamespacePtr(ns *Namespace) (*C.xmlNs, error) {
	if ptr := ns.ptr; ptr != 0 {
		return (*C.xmlNs)(unsafe.Pointer(ptr)), nil
	}
	return nil, ErrInvalidNamespace
}

func xmlNewNode(ns *Namespace, name string) *C.xmlElement {
	var nsptr *C.xmlNs
	if ns != nil {
		nsptr = (*C.xmlNs)(unsafe.Pointer(ns.ptr))
	}

	cname := stringToXMLChar(name)
	defer C.free(unsafe.Pointer(cname))
	n := C.xmlNewNode(nsptr, cname)
	return (*C.xmlElement)(unsafe.Pointer(n))
}

func xmlNewDocProp(doc *Document, k, v string) (*C.xmlAttr, error) {
	kx, err := xmlMakeSafeName(k)
	if err != nil {
		return nil, err
	}
	defer C.free(unsafe.Pointer(kx))

	docptr, err := validDocumentPtr(doc)
	if err != nil {
		return nil, err
	}

	ent, err := xmlEncodeEntitiesReentrant(doc, v)
	if err != nil {
		return nil, err
	}
	attr := C.xmlNewDocProp(docptr, kx, ent)
	return attr, nil
}

func xmlSearchNsByHref(doc *Document, n Node, uri string) (*Namespace, error) {
	nptr, _ := validNodePtr(n)

	xcuri := stringToXMLChar(uri)
	defer C.free(unsafe.Pointer(xcuri))

	docptr, err := validDocumentPtr(doc)
	if err != nil {
		return nil, err
	}

	ns := C.xmlSearchNsByHref(docptr, nptr, xcuri)
	if ns == nil {
		return nil, ErrNamespaceNotFound{Target: uri}
	}
	return wrapNamespace(ns), nil
}

func xmlSearchNs(doc *Document, n Node, prefix string) (*Namespace, error) {
	nptr, err := validNodePtr(n)
	if err != nil {
		return nil, err
	}

	docptr, err := validDocumentPtr(doc)
	if err != nil {
		return nil, err
	}

	cprefix := stringToXMLChar(prefix)
	defer C.free(unsafe.Pointer(cprefix))

	ns := C.xmlSearchNs(docptr, nptr, cprefix)
	if ns == nil {
		return nil, ErrNamespaceNotFound{Target: prefix}
	}
	return wrapNamespace(ns), nil
}

func xmlNewDocNode(doc *Document, ns *Namespace, local, content string) (*C.xmlNode, error) {
	docptr, err := validDocumentPtr(doc)
	if err != nil {
		return nil, err
	}

	nsptr, err := validNamespacePtr(ns)
	if err != nil {
		return nil, err
	}

	var c *C.xmlChar
	if len(content) > 0 {
		c = stringToXMLChar(content)
		defer C.free(unsafe.Pointer(c))
	}

	clocal := stringToXMLChar(local)
	defer C.free(unsafe.Pointer(c))

	return C.xmlNewDocNode(docptr, nsptr, clocal, c), nil
}

func xmlNewNs(n Node, nsuri, prefix string) (*Namespace, error) {
	nptr, err := validNodePtr(n)
	if err != nil {
		return nil, err
	}

	cnsuri := stringToXMLChar(nsuri)
	cprefix := stringToXMLChar(prefix)
	defer C.free(unsafe.Pointer(cnsuri))
	defer C.free(unsafe.Pointer(cprefix))

	return wrapNamespace(C.xmlNewNs(nptr, cnsuri, cprefix)), nil
}

func xmlSetNs(n Node, ns *Namespace) error {
	nptr, err := validNodePtr(n)
	if err != nil {
		return err
	}

	nsptr, err := validNamespacePtr(ns)
	if err != nil {
		return err
	}

	C.xmlSetNs(nptr, nsptr)
	return nil
}

func xmlNewCDataBlock(doc *Document, txt string) *C.xmlNode {
	dptr, err := validDocumentPtr(doc)
	if err != nil {
		return nil
	}
	ctxt := stringToXMLChar(txt)
	defer C.free(unsafe.Pointer(ctxt))
	return C.xmlNewCDataBlock(dptr, ctxt, C.int(len(txt)))
}

func xmlNewComment(txt string) *C.xmlNode {
	ctxt := stringToXMLChar(txt)
	defer C.free(unsafe.Pointer(ctxt))
	return C.xmlNewComment(ctxt)
}

func xmlNewText(txt string) *C.xmlNode {
	ctxt := stringToXMLChar(txt)
	defer C.free(unsafe.Pointer(ctxt))
	return C.xmlNewText(ctxt)
}

func wrapNamespace(n *C.xmlNs) *Namespace {
	return &Namespace{wrapXMLNode((*C.xmlNode)(unsafe.Pointer(n)))}
}

func wrapAttribute(n *C.xmlAttr) *Attribute {
	return &Attribute{wrapXMLNode((*C.xmlNode)(unsafe.Pointer(n)))}
}

func wrapCDataSection(n *C.xmlNode) *CDataSection {
	return &CDataSection{wrapXMLNode(n)}
}

func wrapComment(n *C.xmlNode) *Comment {
	return &Comment{wrapXMLNode(n)}
}

func wrapElement(n *C.xmlElement) *Element {
	return &Element{wrapXMLNode((*C.xmlNode)(unsafe.Pointer(n)))}
}

func wrapXMLNode(n *C.xmlNode) *XMLNode {
	return &XMLNode{
		ptr: uintptr(unsafe.Pointer(n)),
	}
}

func wrapText(n *C.xmlNode) *Text {
	return &Text{wrapXMLNode(n)}
}

// WrapToNodeUnsafe is a function created with the sole purpose
// of allowing go-libxml2 consumers that can generate an xmlNodePtr
// type to create libxml2.Node types.
//
// The unsafe.Pointer variable is cast into a C.xmlNodePtr, and
// wrapped into a go-libxml2 node type. You shouldn't be using
// this function unless you know EXACTLY what you are doing
// including knowing how to allocate/free libxml2 resources
func WrapToNodeUnsafe(n unsafe.Pointer) (Node, error) {
	ptr := (*C.xmlNode)(n)
	return wrapToNode(ptr)
}

func wrapToNode(n *C.xmlNode) (Node, error) {
	switch XMLNodeType(n._type) {
	case AttributeNode:
		return wrapAttribute((*C.xmlAttr)(unsafe.Pointer(n))), nil
	case ElementNode:
		return wrapElement((*C.xmlElement)(unsafe.Pointer(n))), nil
	case TextNode:
		return wrapText(n), nil
	default:
		return nil, errors.New("unknown node: " + strconv.Itoa(int(n._type)))
	}
}

func nodeName(n Node) string {
	switch n.NodeType() {
	case XIncludeStart, XIncludeEnd, EntityRefNode, EntityNode, DTDNode, EntityDecl, DocumentTypeNode, NotationNode, NamespaceDecl:
		return xmlCharToString((*C.xmlNode)(n.Pointer()).name)
	case CommentNode:
		return "#comment"
	case CDataSectionNode:
		return "#cdata-section"
	case TextNode:
		return "#text"
	case DocumentNode, HTMLDocumentNode, DocbDocumentNode:
		return "#document"
	case DocumentFragNode:
		return "#document-fragment"
	case ElementNode, AttributeNode:
		ptr := (*C.xmlNode)(n.Pointer())
		if ns := ptr.ns; ns != nil {
			if nsstr := xmlCharToString(ns.prefix); nsstr != "" {
				return fmt.Sprintf("%s:%s", xmlCharToString(ns.prefix), xmlCharToString(ptr.name))
			}
		}
		return xmlCharToString(ptr.name)
	case ElementDecl, AttributeDecl:
		panic("unimplemented")
	default:
		panic("unknown")
	}
}

func nodeValue(n Node) string {
	switch n.NodeType() {
	case AttributeNode, ElementNode, TextNode, CommentNode, CDataSectionNode, PiNode, EntityRefNode:
		return xmlCharToString(C.xmlXPathCastNodeToString((*C.xmlNode)(n.Pointer())))
	case EntityDecl:
		np := (*C.xmlNode)(n.Pointer())
		if np.content != nil {
			return xmlCharToString(C.xmlStrdup(np.content))
		}

		panic("unimplmented")
	}

	return ""
}

// Pointer returns the pointer to the underlying C struct
func (n *XMLNode) Pointer() unsafe.Pointer {
	if n == nil || n.ptr == 0 {
		return nil
	}
	return unsafe.Pointer(n.ptr)
}

// AddChild appends the node
func (n *XMLNode) AddChild(child Node) error {
	nptr, err := validNodePtr(n)
	if err != nil {
		return err
	}

	cptr, err := validNodePtr(child)
	if err != nil {
		return err
	}

	if C.xmlAddChild(nptr, cptr) == nil {
		return errors.New("failed to add child")
	}
	return nil
}

// ChildNodes returns the child nodes
func (n *XMLNode) ChildNodes() (NodeList, error) {
	return childNodes(n)
}

func wrapDocument(n *C.xmlDoc) *Document {
	return &Document{
		ptr: uintptr(unsafe.Pointer(n)),
	}
}

// OwnerDocument returns the Document that this node belongs to
func (n *XMLNode) OwnerDocument() (*Document, error) {
	nptr, err := validNodePtr(n)
	if err != nil {
		return nil, err
	}

	if nptr.doc == nil {
		return nil, ErrInvalidDocument
	}
	return wrapDocument(nptr.doc), nil
}

// FindNodes evaluates the xpath expression and returns the matching nodes
func (n *XMLNode) FindNodes(xpath string) (NodeList, error) {
	ctx, err := NewXPathContext(n)
	if err != nil {
		return nil, err
	}
	defer ctx.Free()

	return ctx.FindNodes(xpath)
}

// FindNodesExpr evalues the pre-compiled xpath expression and returns the matching nodes
func (n *XMLNode) FindNodesExpr(expr *XPathExpression) (NodeList, error) {
	ctx, err := NewXPathContext(n)
	if err != nil {
		return nil, err
	}
	defer ctx.Free()

	return ctx.FindNodesExpr(expr)
}

// FirstChild reutrns the first child node
func (n *XMLNode) FirstChild() (Node, error) {
	if !n.HasChildNodes() {
		return nil, errors.New("no children")
	}

	nptr, err := validNodePtr(n)
	if err != nil {
		return nil, err
	}
	return wrapToNode(nptr.children)
}

// HasChildNodes returns true if the node contains children
func (n *XMLNode) HasChildNodes() bool {
	nptr, err := validNodePtr(n)
	if err != nil {
		return false
	}
	return nptr.children != nil
}

// LastChild returns the last child node
func (n *XMLNode) LastChild() (Node, error) {
	nptr, err := validNodePtr(n)
	if err != nil {
		return nil, err
	}
	return wrapToNode(nptr.last)
}

// LocalName returns the local name
func (n *XMLNode) LocalName() string {
	nptr, err := validNodePtr(n)
	if err != nil {
		return ""
	}

	switch n.NodeType() {
	case ElementNode, AttributeNode, ElementDecl, AttributeDecl:
		return xmlCharToString(nptr.name)
	}
	return ""
}

// NamespaceURI returns the namespace URI associated with this node
func (n *XMLNode) NamespaceURI() string {
	nptr, err := validNodePtr(n)
	if err != nil {
		return ""
	}

	switch n.NodeType() {
	case ElementNode, AttributeNode, PiNode:
		if ns := nptr.ns; ns != nil && ns.href != nil {
			return xmlCharToString(ns.href)
		}
	}
	return ""
}

// NextSibling returns the next sibling
func (n *XMLNode) NextSibling() (Node, error) {
	nptr, err := validNodePtr(n)
	if err != nil {
		return nil, err
	}
	return wrapToNode(nptr.next)
}

// ParentNode returns the parent node
func (n *XMLNode) ParentNode() (Node, error) {
	nptr, err := validNodePtr(n)
	if err != nil {
		return nil, err
	}
	return wrapToNode(nptr.parent)
}

// Prefix returns the prefix from the node name, if any
func (n *XMLNode) Prefix() string {
	nptr, err := validNodePtr(n)
	if err != nil {
		return ""
	}

	switch n.NodeType() {
	case ElementNode, AttributeNode, PiNode:
		if ns := nptr.ns; ns != nil && ns.prefix != nil {
			return xmlCharToString(ns.prefix)
		}
	}
	return ""
}

// PreviousSibling returns the previous sibling
func (n *XMLNode) PreviousSibling() (Node, error) {
	nptr, err := validNodePtr(n)
	if err != nil {
		return nil, err
	}
	return wrapToNode(nptr.prev)
}

// SetNodeName sets the node name
func (n *XMLNode) SetNodeName(name string) {
	nptr, err := validNodePtr(n)
	if err != nil {
		return
	}
	cname := stringToXMLChar(name)
	defer C.free(unsafe.Pointer(cname))
	C.xmlNodeSetName(nptr, cname)
}

// SetNodeValue sets the node value
func (n *XMLNode) SetNodeValue(value string) {
	nptr, err := validNodePtr(n)
	if err != nil {
		return
	}
	cvalue := stringToXMLChar(value)
	defer C.free(unsafe.Pointer(cvalue))

	// TODO: Implement this in C
	if n.NodeType() != AttributeNode {
		C.xmlNodeSetContent(nptr, cvalue)
		return
	}

	if nptr.children != nil {
		nptr.last = nil
		C.xmlFreeNodeList(nptr.children)
	}

	nptr.children = C.xmlNewText(cvalue)
	nptr.children.parent = nptr
	nptr.children.doc = nptr.doc
	nptr.last = nptr.children
}

// TextContent returns the text content
func (n *XMLNode) TextContent() string {
	nptr, err := validNodePtr(n)
	if err != nil {
		return ""
	}
	return xmlCharToString(C.xmlXPathCastNodeToString(nptr))
}

// ToString returns the string representation. (But it should probably
// be deprecated)
func (n *XMLNode) ToString(format int, docencoding bool) string {
	nptr, err := validNodePtr(n)
	if err != nil {
		return ""
	}

	// TODO: Implement htis in C
	buffer := C.xmlBufferCreate()
	defer C.xmlBufferFree(buffer)
	if format <= 0 {
		C.xmlNodeDump(buffer, nptr.doc, nptr, 0, 0)
	} else {
		oIndentTreeOutput := C.MY_setXmlIndentTreeOutput(1)
		C.xmlNodeDump(buffer, nptr.doc, nptr, 0, C.int(format))
		C.MY_setXmlIndentTreeOutput(oIndentTreeOutput)
	}
	return xmlCharToString(C.xmlBufferContent(buffer))
}

// LookupNamespacePrefix returns the prefix associated with the given URL
func (n *XMLNode) LookupNamespacePrefix(href string) (string, error) {
	nptr, err := validNodePtr(n)
	if err != nil {
		return "", err
	}

	if href == "" {
		return "", ErrNamespaceNotFound{Target: href}
	}

	chref := stringToXMLChar(href)
	defer C.free(unsafe.Pointer(chref))
	ns := C.xmlSearchNsByHref(nptr.doc, nptr, chref)
	if ns == nil {
		return "", ErrNamespaceNotFound{Target: href}
	}

	return xmlCharToString(ns.prefix), nil
}

// LookupNamespaceURI returns the URI associated with the given prefix
func (n *XMLNode) LookupNamespaceURI(prefix string) (string, error) {
	nptr, err := validNodePtr(n)
	if err != nil {
		return "", err
	}

	if prefix == "" {
		return "", ErrNamespaceNotFound{Target: prefix}
	}

	cprefix := stringToXMLChar(prefix)
	defer C.free(unsafe.Pointer(cprefix))
	ns := C.xmlSearchNs(nptr.doc, nptr, cprefix)
	if ns == nil {
		return "", ErrNamespaceNotFound{Target: prefix}
	}

	return xmlCharToString(ns.href), nil
}

// NodeType returns the XMLNodeType
func (n *XMLNode) NodeType() XMLNodeType {
	nptr, err := validNodePtr(n)
	if err != nil {
		return XMLNodeType(0)
	}
	return XMLNodeType(nptr._type)
}

// Walk traverses through all of the nodes
func (n *XMLNode) Walk(fn func(Node) error) error {
	walk(n, fn)
	return nil
}

// AutoFree allows you to free the underlying C resources. It is
// meant to be called from defer. If you don't call `MakeMortal()` or
// do call `MakePersistent()`, AutoFree is a no-op.
func (n *XMLNode) AutoFree() {
	if !n.mortal {
		return
	}
	n.Free()
}

// MakeMortal flags the node so that `AutoFree` calls Free()
// to release the underlying C resources.
func (n *XMLNode) MakeMortal() {
	n.mortal = true
}

// MakePersistent flags the node so that `AutoFree` becomes a no-op.
// Make sure to call this if you used `MakeMortal` and `AutoFree`,
// but you then decided to keep the node around.
func (n *XMLNode) MakePersistent() {
	n.mortal = false
}

// Free releases the underlying C struct
func (n *XMLNode) Free() {
	nptr, err := validNodePtr(n)
	if err != nil {
		return
	}
	C.xmlFreeNode(nptr)
	n.ptr = 0
}

func walk(n Node, fn func(Node) error) error {
	if err := fn(n); err != nil {
		return err
	}
	children, err := n.ChildNodes()
	if err != nil {
		return err
	}
	for _, c := range children {
		if err := walk(c, fn); err != nil {
			return err
		}
	}
	return nil
}

func childNodes(n Node) (NodeList, error) {
	ret := NodeList(nil)
	for chld := ((*C.xmlNode)(n.Pointer())).children; chld != nil; chld = chld.next {
		nchld, err := wrapToNode(chld)
		if err != nil {
			return nil, err
		}
		ret = append(ret, nchld)
	}
	return ret, nil
}

func splitPrefixLocal(s string) (string, string) {
	i := strings.IndexByte(s, ':')
	if i == -1 {
		return "", s
	}
	return s[:i], s[i+1:]
}

// URI returns the namespace URL
func (n *Namespace) URI() string {
	nsptr, err := validNamespacePtr(n)
	if err != nil {
		return ""
	}
	return xmlCharToString(nsptr.href)
}

// Prefix returns the prefix for this namespace
func (n *Namespace) Prefix() string {
	nsptr, err := validNamespacePtr(n)
	if err != nil {
		return ""
	}
	return xmlCharToString(nsptr.prefix)
}

// Free releases the underlying C struct
func (n *Namespace) Free() {
	nsptr, err := validNamespacePtr(n)
	if err != nil {
		return
	}
	C.MY_xmlFree(unsafe.Pointer(nsptr))
	n.ptr = 0
}

func createElement(d *Document, name string) (*Element, error) {
	dptr, err := validDocumentPtr(d)
	if err != nil {
		return nil, err
	}
	if err := myTestNodeName(name); err != nil {
		return nil, err
	}

	cname := stringToXMLChar(name)
	defer C.free(unsafe.Pointer(cname))
	newNode := C.xmlNewNode(nil, cname)
	if newNode == nil {
		return nil, errors.New("element creation failed")
	}
	// XXX hmmm...
	newNode.doc = dptr
	return wrapElement((*C.xmlElement)(unsafe.Pointer(newNode))), nil
}

func createElementNS(doc *Document, nsuri, name string) (*Element, error) {
	dptr, err := validDocumentPtr(doc)
	if err != nil {
		return nil, err
	}

	if err := myTestNodeName(name); err != nil {
		return nil, err
	}

	if nsuri == "" {
		return createElement(doc, name)
	}

	rootptr := C.xmlDocGetRootElement(dptr)

	var prefix, localname string
	var ns *C.xmlNs

	if i := strings.IndexByte(name, ':'); i > 0 {
		prefix = name[:i]
		localname = name[i+1:]
	} else {
		localname = name
	}

	xmlnsuri := stringToXMLChar(nsuri)
	xmlprefix := stringToXMLChar(prefix)
	defer C.free(unsafe.Pointer(xmlnsuri))
	defer C.free(unsafe.Pointer(xmlprefix))

	// Optimization: if rootptr is nil, then you can just
	// create the namespace
	if rootptr == nil {
		ns = C.xmlNewNs(nil, xmlnsuri, xmlprefix)
	} else if prefix != "" {
		// prefix exists, see if this is declared
		ns = C.xmlSearchNs(dptr, rootptr, xmlprefix)
		if ns == nil { // not declared. create a new one
			ns = C.xmlNewNs(nil, xmlnsuri, xmlprefix)
		} else { // declared. does uri match?
			if C.xmlStrcmp(ns.href, xmlnsuri) != C.int(0) {
				return nil, errors.New("prefix already registered to different uri")
			}
			// Namespace is already registered, we don't need to provide a
			// namespace element to xmlNewDocNode
			ns = nil

			// but the localname should be prefix:localname
			localname = name
		}
	} else {
		// If the name does not contain a prefix, check for the
		// existence of this namespace via the URI
		ns = C.xmlSearchNsByHref(dptr, rootptr, xmlnsuri)
		if ns == nil {
			ns = C.xmlNewNs(nil, xmlnsuri, nil)
		}
	}

	clocal := stringToXMLChar(localname)
	defer C.free(unsafe.Pointer(clocal))
	newNode := C.xmlNewDocNode(dptr, ns, clocal, nil)
	newNode.nsDef = ns

	return wrapElement((*C.xmlElement)(unsafe.Pointer(newNode))), nil
}

func validDocumentPtr(doc *Document) (*C.xmlDoc, error) {
	if dptr := doc.ptr; dptr != 0 {
		return (*C.xmlDoc)(unsafe.Pointer(dptr)), nil
	}
	return nil, ErrInvalidDocument
}

func documentEncoding(doc *Document) string {
	dptr, err := validDocumentPtr(doc)
	if err != nil {
		return ""
	}
	return xmlCharToString(dptr.encoding)
}

func documentStandalone(doc *Document) int {
	dptr, err := validDocumentPtr(doc)
	if err != nil {
		return 0
	}
	return int(dptr.standalone)
}

func documentURI(doc *Document) string {
	dptr, err := validDocumentPtr(doc)
	if err != nil {
		return ""
	}
	return xmlCharToString(dptr.URL)
}

func documentVersion(doc *Document) string {
	dptr, err := validDocumentPtr(doc)
	if err != nil {
		return ""
	}
	return xmlCharToString(dptr.version)
}

func documentElement(doc *Document) *C.xmlNode {
	dptr, err := validDocumentPtr(doc)
	if err != nil {
		return nil
	}

	return C.xmlDocGetRootElement(dptr)
}

func xmlFreeDoc(doc *Document) error {
	dptr, err := validDocumentPtr(doc)
	if err != nil {
		return err
	}
	C.xmlFreeDoc(dptr)
	doc.ptr = 0
	return nil
}

func documentString(doc *Document, encoding string, format bool) string {
	dptr, err := validDocumentPtr(doc)
	if err != nil {
		return ""
	}

	var xc *C.xmlChar
	var intformat C.int
	if format {
		intformat = C.int(1)
	} else {
		intformat = C.int(0)
	}

	// Ideally this shouldn't happen, but you never know.
	if encoding == "" {
		encoding = "utf-8"
	}

	i := C.int(0)

	cenc := C.CString(encoding)
	defer C.free(unsafe.Pointer(cenc))

	C.xmlDocDumpFormatMemoryEnc(dptr, &xc, &i, cenc, intformat)

	return xmlCharToString(xc)
}

func xmlNodeSetBase(doc *Document, s string) {
	dptr, err := validDocumentPtr(doc)
	if err != nil {
		return
	}

	cs := stringToXMLChar(s)
	defer C.free(unsafe.Pointer(cs))
	C.xmlNodeSetBase((*C.xmlNode)(unsafe.Pointer(dptr)), cs)
}

func validNodePtr(n Node) (*C.xmlNode, error) {
	if n == nil {
		return nil, ErrInvalidNode
	}

	nptr := (*C.xmlNode)(n.Pointer())
	if nptr == nil {
		return nil, ErrInvalidNode
	}

	return nptr, nil
}

func setDocumentElement(doc *Document, n Node) error {
	dptr, err := validDocumentPtr(doc)
	if err != nil {
		return err
	}

	nptr, err := validNodePtr(n)
	if err != nil {
		return err
	}

	C.xmlDocSetRootElement(dptr, nptr)
	return nil
}

func setDocumentEncoding(doc *Document, e string) {
	dptr, err := validDocumentPtr(doc)
	if err != nil {
		return
	}

	if dptr.encoding != nil {
		C.MY_xmlFree(unsafe.Pointer(dptr.encoding))
	}

	// note: this doesn't need to be dup'ed, as
	// C.CString is already duped/malloc'ed
	dptr.encoding = stringToXMLChar(e)
}

func setDocumentStandalone(doc *Document, v int) {
	dptr, err := validDocumentPtr(doc)
	if err != nil {
		return
	}
	dptr.standalone = C.int(v)
}

func setDocumentVersion(doc *Document, v string) {
	dptr, err := validDocumentPtr(doc)
	if err != nil {
		return
	}

	if dptr.version != nil {
		C.MY_xmlFree(unsafe.Pointer(dptr.version))
	}

	// note: this doesn't need to be dup'ed, as
	// C.CString is already duped/malloc'ed
	dptr.version = stringToXMLChar(v)
}

func xmlSetProp(n Node, name, value string) error {
	nptr, err := validNodePtr(n)
	if err != nil {
		return err
	}

	cname := stringToXMLChar(name)
	cvalue := stringToXMLChar(value)
	defer C.free(unsafe.Pointer(cname))
	defer C.free(unsafe.Pointer(cvalue))

	C.xmlSetProp(nptr, cname, cvalue)
	return nil
}

func xmlElementAttributes(n Node) ([]*Attribute, error) {
	nptr, err := validNodePtr(n)
	if err != nil {
		return nil, err
	}

	attrs := []*Attribute{}
	for attr := nptr.properties; attr != nil; attr = attr.next {
		attrs = append(attrs, wrapAttribute(attr))
	}
	return attrs, nil
}

func xmlElementNamespaces(n Node) ([]*Namespace, error) {
	nptr, err := validNodePtr(n)
	if err != nil {
		return nil, err
	}

	ret := []*Namespace{}
	for ns := nptr.nsDef; ns != nil; ns = ns.next {
		if ns.prefix == nil && ns.href == nil {
			continue
		}
		// ALERT! Allocating new C struct here
		newns := xmlCopyNamespace(ns)
		if newns == nil { // XXX this is an error, no?
			continue
		}

		ret = append(ret, wrapNamespace(newns))
	}
	return ret, nil
}

func (n *Element) getAttributeNode(name string) (*C.xmlAttr, error) {
	nptr, err := validNodePtr(n)
	if err != nil {
		return nil, err
	}

	// if this is "xmlns", look for the first namespace without
	// the prefix
	if name == "xmlns" {
		for nsdef := nptr.nsDef; nsdef != nil; nsdef = nsdef.next {
			if nsdef.prefix != nil {
				continue
			}
			debug.Printf("nsdef.href -> %s", xmlCharToString(nsdef.href))
		}
	}

	debug.Printf("n = %s", n.String())
	debug.Printf("getAttributeNode(%s)", name)

	cname := stringToXMLChar(name)
	defer C.free(unsafe.Pointer(cname))

	prop := C.xmlHasNsProp(nptr, cname, nil)
	debug.Printf("prop = %v", prop)
	if prop == nil {
		prefix, local := splitPrefixLocal(name)
		debug.Printf("prefix = %s, local = %s", prefix, local)
		if local != "" {
			cprefix := stringToXMLChar(prefix)
			defer C.free(unsafe.Pointer(cprefix))
			if ns := C.xmlSearchNs(nptr.doc, nptr, cprefix); ns != nil {
				clocal := stringToXMLChar(local)
				defer C.free(unsafe.Pointer(clocal))

				prop = C.xmlHasNsProp(nptr, clocal, ns.href)
			}

		}
	}

	if prop == nil || XMLNodeType(prop._type) != AttributeNode {
		return nil, errors.New("attribute not found")
	}

	return prop, nil
}

func xmlUnlinkNode(prop *C.xmlAttr) {
	C.xmlUnlinkNode((*C.xmlNode)(unsafe.Pointer(prop)))
}

func xmlFreeProp(attr *Attribute) {
	C.xmlFreeProp((*C.xmlAttr)(unsafe.Pointer(attr.ptr)))
}

func xmlFreeNode(n Node) {
	C.xmlFreeNode((*C.xmlNode)(unsafe.Pointer(n.Pointer())))
}

func xmlCopyNamespace(ns *C.xmlNs) *C.xmlNs {
	return C.xmlCopyNamespace(ns)
}

func xmlUnsetProp(n Node, name string) error {
	nptr := (*C.xmlNode)(unsafe.Pointer(n.Pointer()))
	if nptr == nil {
		return errors.New("invalid node")
	}

	cname := stringToXMLChar(name)
	defer C.free(unsafe.Pointer(cname))

	i := C.xmlUnsetProp(nptr, cname)
	if i == C.int(0) {
		return errors.New("failed to unset prop")
	}
	return nil
}

func xmlUnsetNsProp(n Node, ns *Namespace, name string) error {
	nptr := (*C.xmlNode)(unsafe.Pointer(n.Pointer()))
	if nptr == nil {
		return errors.New("invalid node")
	}

	cname := stringToXMLChar(name)
	defer C.free(unsafe.Pointer(cname))

	i := C.xmlUnsetNsProp(
		nptr,
		(*C.xmlNs)(unsafe.Pointer(ns.ptr)),
		cname,
	)
	if i == C.int(0) {
		return errors.New("failed to unset prop")
	}
	return nil
}

func xmlC14NDocDumpMemory(d *Document, mode C14NMode, withComments bool) (string, error) {
	dptr, err := validDocumentPtr(d)
	if err != nil {
		return "", err
	}

	var result *C.xmlChar

	var withCommentsInt C.int
	if withComments {
		withCommentsInt = 1
	}

	modeInt := C.int(mode)

	written := C.xmlC14NDocDumpMemory(
		dptr,
		nil,
		modeInt,
		nil,
		withCommentsInt,
		&result,
	)
	if written < 0 {
		e := C.MY_xmlLastError()
		return "", errors.New("c14n dump failed: " + C.GoString(e.message))
	}
	return xmlCharToString(result), nil
}

func appendText(n Node, s string) error {
	cs := stringToXMLChar(s)
	defer C.free(unsafe.Pointer(cs))

	txt := C.xmlNewText(cs)
	if txt == nil {
		return errors.New("failed to create text node")
	}

	if C.xmlAddChild((*C.xmlNode)(n.Pointer()), (*C.xmlNode)(txt)) == nil {
		return errors.New("failed to create text node")
	}
	return nil
}

func xmlDocCopyNode(n Node, d *Document, extended int) (Node, error) {
	nptr, err := validNodePtr(n)
	if err != nil {
		return nil, err
	}

	dptr, err := validDocumentPtr(d)
	if err != nil {
		return nil, err
	}

	ret := C.xmlDocCopyNode(nptr, dptr, C.int(extended))
	if ret == nil {
		return nil, errors.New("copy node failed")
	}

	return wrapToNode(ret)
}

func xmlSetTreeDoc(n Node, d *Document) error {
	nptr, err := validNodePtr(n)
	if err != nil {
		return err
	}

	dptr, err := validDocumentPtr(d)
	if err != nil {
		return err
	}

	C.xmlSetTreeDoc(nptr, dptr)
	return nil
}

func xmlParseInNodeContext(n Node, data string, o ParseOption) (Node, error) {
	nptr, err := validNodePtr(n)
	if err != nil {
		return nil, err
	}

	var ret C.xmlNodePtr
	if C.xmlParseInNodeContext(nptr, C.CString(data), C.int(len(data)), C.int(o), &ret) != 0 {
		return nil, errors.New("XXX PLACE HOLDER XXX")
	}

	return wrapToNode((*C.xmlNode)(unsafe.Pointer(ret)))
}

func validXPathContextPtr(x *XPathContext) (*C.xmlXPathContext, error) {
	if xptr := x.ptr; xptr != 0 {
		return (*C.xmlXPathContext)(unsafe.Pointer(xptr)), nil
	}
	return nil, ErrInvalidXPathContext
}

func validXPathExpressionPtr(x *XPathExpression) (*C.xmlXPathCompExpr, error) {
	if xptr := x.ptr; xptr != 0 {
		return (*C.xmlXPathCompExpr)(unsafe.Pointer(xptr)), nil
	}
	return nil, ErrInvalidXPathExpression
}

func validXPathObjectPtr(x *XPathObject) (*C.xmlXPathObject, error) {
	if xptr := x.ptr; xptr != 0 {
		return (*C.xmlXPathObject)(unsafe.Pointer(xptr)), nil
	}
	return nil, ErrInvalidXPathObject
}

func xmlXPathNewContext(n ...Node) (*XPathContext, error) {
	ctx := C.xmlXPathNewContext(nil)
	ctx.namespaces = nil

	if len(n) > 0 {
		ctx.node = (*C.xmlNode)(n[0].Pointer())
	}

	return &XPathContext{ptr: uintptr(unsafe.Pointer(ctx))}, nil
}

func xmlXPathContextSetContextNode(x *XPathContext, n Node) error {
	xptr, err := validXPathContextPtr(x)
	if err != nil {
		return err
	}

	nptr, err := validNodePtr(n)
	if err != nil {
		return err
	}

	xptr.node = nptr
	return nil
}

func xmlXPathCompile(s string) (*XPathExpression, error) {
	cs := stringToXMLChar(s)
	defer C.free(unsafe.Pointer(cs))

	if p := C.xmlXPathCompile(cs); p != nil {
		return &XPathExpression{ptr: uintptr(unsafe.Pointer(p)), expr: s}, nil
	}
	return nil, ErrXPathCompileFailure
}

func xmlXPathFreeCompExpr(x *XPathExpression) error {
	xptr, err := validXPathExpressionPtr(x)
	if err != nil {
		return err
	}
	C.xmlXPathFreeCompExpr(xptr)
	return nil
}

func xmlXPathFreeContext(x *XPathContext) error {
	xptr, err := validXPathContextPtr(x)
	if err != nil {
		return err
	}
	C.xmlXPathFreeContext(xptr)
	return nil
}

func xmlXPathNSLookup(x *XPathContext, prefix string) (string, error) {
	xptr, err := validXPathContextPtr(x)
	if err != nil {
		return "", err
	}

	cprefix := stringToXMLChar(prefix)
	defer C.free(unsafe.Pointer(cprefix))

	if s := C.xmlXPathNsLookup(xptr, cprefix); s != nil {
		return xmlCharToString(s), nil
	}

	return "", ErrNamespaceNotFound{Target: prefix}
}

func xmlXPathRegisterNS(x *XPathContext, prefix, nsuri string) error {
	xptr, err := validXPathContextPtr(x)
	if err != nil {
		return err
	}

	cprefix := stringToXMLChar(prefix)
	cnsuri := stringToXMLChar(nsuri)
	defer C.free(unsafe.Pointer(cprefix))
	defer C.free(unsafe.Pointer(cnsuri))

	if res := C.xmlXPathRegisterNs(xptr, cprefix, cnsuri); res == -1 {
		return ErrXPathNamespaceRegisterFailure
	}
	return nil
}

func evalXPath(x *XPathContext, expr *XPathExpression) (*XPathObject, error) {
	xptr, err := validXPathContextPtr(x)
	if err != nil {
		return nil, err
	}

	exprptr, err := validXPathExpressionPtr(expr)
	if err != nil {
		return nil, err
	}

	// If there is no document associated with this context,
	// then xmlXPathCompiledEval() just fails to match
	if xptr.node != nil && xptr.node.doc != nil {
		xptr.doc = xptr.node.doc
	}

	if xptr.doc == nil {
		cs := stringToXMLChar("1.0")
		defer C.free(unsafe.Pointer(cs))
		xptr.doc = C.xmlNewDoc(cs)

		defer C.xmlFreeDoc(xptr.doc)
	}

	res := C.xmlXPathCompiledEval(exprptr, xptr)
	if res == nil {
		return nil, ErrXPathEmptyResult
	}

	return &XPathObject{ptr: uintptr(unsafe.Pointer(res))}, nil
}

func xmlXPathFreeObject(x *XPathObject) {
	xptr, err := validXPathObjectPtr(x)
	if err != nil {
		return
	}
	C.xmlXPathFreeObject(xptr)
	//	if xptr.nodesetval != nil {
	//		C.xmlXPathFreeNodeSet(xptr.nodesetval)
	//	}
}

func xmlXPathObjectNodeListLen(x *XPathObject) int {
	xptr, err := validXPathObjectPtr(x)
	if err != nil {
		return 0
	}

	if xptr.nodesetval == nil {
		return 0
	}

	return int(xptr.nodesetval.nodeNr)
}

func xmlXPathObjectType(x *XPathObject) XPathObjectType {
	xptr, err := validXPathObjectPtr(x)
	if err != nil {
		return XPathUndefined
	}
	return XPathObjectType(xptr._type)
}

func xmlXPathObjectFloat64(x *XPathObject) float64 {
	xptr, err := validXPathObjectPtr(x)
	if err != nil {
		return float64(0)
	}

	return float64(xptr.floatval)
}

func xmlXPathObjectBool(x *XPathObject) bool {
	xptr, err := validXPathObjectPtr(x)
	if err != nil {
		return false
	}

	return C.int(xptr.boolval) == 1
}

func xmlXPathObjectNodeList(x *XPathObject) (NodeList, error) {
	xptr, err := validXPathObjectPtr(x)
	if err != nil {
		return nil, err
	}

	nodeset := xptr.nodesetval
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

func xmlTextData(n *Text) string {
	nptr, err := validNodePtr(n)
	if err != nil {
		return ""
	}
	return xmlCharToString(nptr.content)
}
