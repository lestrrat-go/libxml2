package libxml2

/*
#cgo pkg-config: libxml-2.0
#include <stdbool.h>
#include "libxml/globals.h"
#include "libxml/tree.h"
#include "libxml/parser.h"
#include "libxml/parserInternals.h"
#include "libxml/xpath.h"
#include "libxml/c14n.h"

// Macro wrapper function
static inline void MY_xmlFree(void *p) {
	xmlFree(p);
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
*/
import "C"
import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"unsafe"

	"github.com/lestrrat/go-libxml2/internal/debug"
)

var _XmlNodeType_index = [...]uint8{0, 11, 24, 32, 48, 61, 71, 77, 88, 100, 116, 132, 144, 160, 167, 178, 191, 201, 214, 227, 238, 254}

const _XmlNodeType_name = `ElementNodeAttributeNodeTextNodeCDataSectionNodeEntityRefNodeEntityNodePiNodeCommentNodeDocumentNodeDocumentTypeNodeDocumentFragNodeNotationNodeHTMLDocumentNodeDTDNodeElementDeclAttributeDeclEntityDeclNamespaceDeclXIncludeStartXIncludeEndDocbDocumentNode`

func xmlNewDoc(version string) *C.xmlDoc {
	return C.xmlNewDoc(stringToXmlChar(version))
}

func xmlStrdup(s string) *C.xmlChar {
	return C.xmlStrdup(stringToXmlChar(s))
}

func xmlEncodeEntitiesReentrant(doc *Document, s string) *C.xmlChar {
	return C.xmlEncodeEntitiesReentrant(doc.ptr, stringToXmlChar(s))
}

func myTestNodeName(n string) error {
	if C.MY_test_node_name(stringToXmlChar(n)) == 0 {
		return ErrInvalidNodeName
	}
	return nil
}

func xmlMakeSafeName(k string) (*C.xmlChar, error) {
	if err := myTestNodeName(k); err != nil {
		return nil, err
	}
	return stringToXmlChar(k), nil
}

func xmlNewNode(ns *Namespace, name string) *C.xmlElement {
	var nsptr *C.xmlNs
	if ns != nil {
		nsptr = (*C.xmlNs)(unsafe.Pointer(ns.ptr))
	}

	n := C.xmlNewNode(
		nsptr,
		stringToXmlChar(name),
	)
	return (*C.xmlElement)(unsafe.Pointer(n))
}

func xmlNewDocProp(doc *Document, k, v string) (*C.xmlAttr, error) {
	kx, err := xmlMakeSafeName(k)
	if err != nil {
		return nil, err
	}

	attr := C.xmlNewDocProp(
		doc.ptr,
		kx,
		xmlEncodeEntitiesReentrant(doc, v),
	)
	return attr, nil
}

func xmlSearchNsByHref(doc *Document, n Node, uri string) *Namespace {
	var xcuri *C.xmlChar
	if len(uri) > 0 {
		xcuri = stringToXmlChar(uri)
	}

	ns := C.xmlSearchNsByHref(
		doc.ptr,
		(*C.xmlNode)(n.Pointer()),
		xcuri,
	)
	if ns == nil {
		return nil
	}
	return wrapNamespace(ns)
}

func xmlSearchNs(doc *Document, n Node, prefix string) *Namespace {
	var nptr *C.xmlNode
	if n != nil {
		nptr = (*C.xmlNode)(n.Pointer())
	}
	ns := C.xmlSearchNs(
		doc.ptr,
		nptr,
		stringToXmlChar(prefix),
	)
	if ns == nil {
		return nil
	}
	return wrapNamespace(ns)
}

func xmlNewDocNode(doc *Document, ns *Namespace, localname, content string) *C.xmlNode {
	var c *C.xmlChar
	if len(content) > 0 {
		c = stringToXmlChar(content)
	}
	return C.xmlNewDocNode(
		doc.ptr,
		(*C.xmlNs)(unsafe.Pointer(ns.ptr)),
		stringToXmlChar(localname),
		c,
	)
}

func xmlNewNs(n Node, nsuri, prefix string) *Namespace {
	var nptr *C.xmlNode
	if n != nil {
		nptr = (*C.xmlNode)(n.Pointer())
	}

	return wrapNamespace(
		C.xmlNewNs(
			nptr,
			stringToXmlChar(nsuri),
			stringToXmlChar(prefix),
		),
	)
}

func xmlSetNs(n Node, ns *Namespace) {
	debug.Printf("Setting namespace for %s to %s", n.NodeName(), ns.Prefix())
	C.xmlSetNs(
		(*C.xmlNode)(n.Pointer()),
		(*C.xmlNs)(unsafe.Pointer(ns.ptr)),
	)
}

func xmlNewCDataBlock(doc *Document, txt string) *C.xmlNode {
	return C.xmlNewCDataBlock(doc.ptr, stringToXmlChar(txt), C.int(len(txt)))
}

func xmlNewComment(txt string) *C.xmlNode {
	return C.xmlNewComment(stringToXmlChar(txt))
}

func xmlNewText(txt string) *C.xmlNode {
	return C.xmlNewText(stringToXmlChar(txt))
}

func (i XmlNodeType) String() string {
	i -= 1
	if i < 0 || i+1 >= XmlNodeType(len(_XmlNodeType_index)) {
		return fmt.Sprintf("XmlNodeType(%d)", i+1)
	}
	return _XmlNodeType_name[_XmlNodeType_index[i]:_XmlNodeType_index[i+1]]
}

func (n NodeList) String() string {
	buf := bytes.Buffer{}
	for _, x := range n {
		buf.WriteString(x.String())
	}
	return buf.String()
}

func (n NodeList) Literal() (string, error) {
	buf := bytes.Buffer{}
	for _, x := range n {
		l, err := x.Literal()
		if err != nil {
			return "", err
		}
		buf.WriteString(l)
	}
	return buf.String(), nil
}

func wrapNamespace(n *C.xmlNs) *Namespace {
	return &Namespace{wrapXmlNode((*C.xmlNode)(unsafe.Pointer(n)))}
}

func wrapAttribute(n *C.xmlAttr) *Attribute {
	return &Attribute{wrapXmlNode((*C.xmlNode)(unsafe.Pointer(n)))}
}

func wrapCDataSection(n *C.xmlNode) *CDataSection {
	return &CDataSection{wrapXmlNode(n)}
}

func wrapComment(n *C.xmlNode) *Comment {
	return &Comment{wrapXmlNode(n)}
}

func wrapElement(n *C.xmlElement) *Element {
	return &Element{wrapXmlNode((*C.xmlNode)(unsafe.Pointer(n)))}
}

func wrapXmlNode(n *C.xmlNode) *XmlNode {
	return &XmlNode{
		ptr: (*C.xmlNode)(unsafe.Pointer(n)),
	}
}

func wrapText(n *C.xmlNode) *Text {
	return &Text{wrapXmlNode(n)}
}

func wrapToNode(n *C.xmlNode) (Node, error) {
	switch XmlNodeType(n._type) {
	case ElementNode:
		return wrapElement((*C.xmlElement)(unsafe.Pointer(n))), nil
	case TextNode:
		return wrapText(n), nil
	default:
		return nil, errors.New("unknown node")
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
	case AttributeNode, TextNode, CommentNode, CDataSectionNode, PiNode, EntityRefNode:
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

func (n *XmlNode) Pointer() unsafe.Pointer {
	return unsafe.Pointer(n.ptr)
}

func (n *XmlNode) AddChild(child Node) {
	C.xmlAddChild(n.ptr, (*C.xmlNode)(child.Pointer()))
}

func (n *XmlNode) AppendChild(child Node) error {
	// XXX There must be lots more checks here because AddChild does things
	// under the table like merging text nodes, freeing some nodes implicitly,
	// et al
	n.AddChild(child)
	return nil
}

func (n *XmlNode) ChildNodes() (NodeList, error) {
	return childNodes(n)
}

func wrapDocument(n *C.xmlDoc) *Document {
	return &Document{ptr: n}
}

func (n *XmlNode) OwnerDocument() *Document {
	return wrapDocument(n.ptr.doc)
}

func (n *XmlNode) FindNodes(xpath string) (NodeList, error) {
	ctx, err := NewXPathContext(n)
	if err != nil {
		return nil, err
	}
	defer ctx.Free()

	return ctx.FindNodes(xpath)
}

func (n *XmlNode) FindNodesExpr(expr *XPathExpression) (NodeList, error) {
	ctx, err := NewXPathContext(n)
	if err != nil {
		return nil, err
	}
	defer ctx.Free()

	return ctx.FindNodesExpr(expr)
}

func (n *XmlNode) FirstChild() (Node, error) {
	if !n.HasChildNodes() {
		return nil, errors.New("no children")
	}

	return wrapToNode(((*C.xmlNode)(n.Pointer())).children)
}

func (n *XmlNode) HasChildNodes() bool {
	return n.ptr.children != nil
}

func (n *XmlNode) IsSameNode(other Node) bool {
	return n.Pointer() == other.Pointer()
}

func (n *XmlNode) LastChild() (Node, error) {
	return wrapToNode(n.ptr.last)
}

func (n XmlNode) Literal() (string, error) {
	return n.String(), nil
}

func (n *XmlNode) LocalName() string {
	switch n.NodeType() {
	case ElementNode, AttributeNode, ElementDecl, AttributeDecl:
		return xmlCharToString(n.ptr.name)
	}
	return ""
}

func (n *XmlNode) NamespaceURI() string {
	switch n.NodeType() {
	case ElementNode, AttributeNode, PiNode:
		if ns := n.ptr.ns; ns != nil && ns.href != nil {
			return xmlCharToString(ns.href)
		}
	}
	return ""
}

func (n *XmlNode) NodeName() string {
	return nodeName(n)
}

func (n *XmlNode) NodeValue() string {
	return nodeValue(n)
}

func (n *XmlNode) NextSibling() (Node, error) {
	return wrapToNode(n.ptr.next)
}

func (n *XmlNode) ParetNode() (Node, error) {
	return wrapToNode(n.ptr.parent)
}

func (n *XmlNode) Prefix() string {
	switch n.NodeType() {
	case ElementNode, AttributeNode, PiNode:
		if ns := n.ptr.ns; ns != nil && ns.prefix != nil {
			return xmlCharToString(ns.prefix)
		}
	}
	return ""
}

func (n *XmlNode) PreviousSibling() (Node, error) {
	return wrapToNode(n.ptr.prev)
}

func (n *XmlNode) SetNodeName(name string) {
	C.xmlNodeSetName(n.ptr, stringToXmlChar(name))
}

func (n *XmlNode) SetNodeValue(value string) {
	// TODO: Implement this in C
	if n.NodeType() != AttributeNode {
		C.xmlNodeSetContent(n.ptr, stringToXmlChar(value))
		return
	}

	ptr := n.ptr
	if ptr.children != nil {
		ptr.last = nil
		C.xmlFreeNodeList(ptr.children)
	}

	ptr.children = C.xmlNewText(stringToXmlChar(value))
	ptr.children.parent = ptr
	ptr.children.doc = ptr.doc
	ptr.last = ptr.children
}

func (n *XmlNode) String() string {
	return n.ToString(0, false)
}

func (n *XmlNode) TextContent() string {
	return xmlCharToString(C.xmlXPathCastNodeToString(n.ptr))
}

func (n *XmlNode) ToString(format int, docencoding bool) string {
	// TODO: Implement htis in C
	buffer := C.xmlBufferCreate()
	defer C.xmlBufferFree(buffer)
	if format <= 0 {
		C.xmlNodeDump(buffer, n.ptr.doc, n.ptr, 0, 0)
	} else {
		oIndentTreeOutput := C.MY_setXmlIndentTreeOutput(1)
		C.xmlNodeDump(buffer, n.ptr.doc, n.ptr, 0, C.int(format))
		C.MY_setXmlIndentTreeOutput(oIndentTreeOutput)
	}
	return xmlCharToString(C.xmlBufferContent(buffer))
}

var ErrNamespaceNotFound = errors.New("namespace not found")

func (n *XmlNode) LookupNamespacePrefix(href string) (string, error) {
	if href == "" {
		return "", ErrNamespaceNotFound
	}

	ns := C.xmlSearchNsByHref(n.ptr.doc, n.ptr, stringToXmlChar(href))
	if ns == nil {
		return "", ErrNamespaceNotFound
	}

	return xmlCharToString(ns.prefix), nil
}

func (n *XmlNode) LookupNamespaceURI(prefix string) (string, error) {
	if prefix == "" {
		return "", ErrNamespaceNotFound
	}

	ns := C.xmlSearchNs(n.ptr.doc, n.ptr, stringToXmlChar(prefix))
	if ns == nil {
		return "", ErrNamespaceNotFound
	}

	return xmlCharToString(ns.href), nil
}

func (n *XmlNode) NodeType() XmlNodeType {
	return XmlNodeType(n.ptr._type)
}

func (n *XmlNode) Walk(fn func(Node) error) error {
	walk(n, fn)
	return nil
}

func (n *XmlNode) Free() {
	C.xmlFreeNode(n.ptr)
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

func childNodes(n ptr) (NodeList, error) {
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

func (n *Namespace) URI() string {
	if ptr := n.ptr; ptr != nil {
		return xmlCharToString(((*C.xmlNs)(unsafe.Pointer(ptr))).href)
	}
	return ""
}

func (n *Namespace) Prefix() string {
	if ptr := n.ptr; ptr != nil {
		return xmlCharToString(((*C.xmlNs)(unsafe.Pointer(ptr))).prefix)
	}
	return ""
}

func (n *Namespace) Free() {
	if ptr := n.ptr; ptr != nil {
		C.MY_xmlFree(unsafe.Pointer(ptr))
	}
}

func createElementNS(doc *Document, nsuri, name string) (*Element, error) {
	if err := myTestNodeName(name); err != nil {
		return nil, err
	}

	if nsuri == "" {
		return doc.CreateElement(name)
	}

	var rootptr *C.xmlNode
	if root, err := doc.DocumentElement(); err == nil {
		rootptr = (*C.xmlNode)(root.Pointer())
	}

	var localname string
	var ns *C.xmlNs

	if i := strings.IndexByte(name, ':'); i > 0 {
		// split it into prefix and localname
		prefix := stringToXmlChar(name[:i])
		xmlnsuri := stringToXmlChar(nsuri)

		// Is this namespace prefix registered already?
		ns = C.xmlSearchNs(
			doc.ptr,
			rootptr,
			prefix,
		)
		if ns == nil {
			// nope, create a new one.
			localname = name[i+1:]
			ns = C.xmlNewNs(nil, xmlnsuri, prefix)
		} else {
			// prefix is registered. do they have matching URI?
			if xmlCharToString(ns.prefix) != name[:i] {
				return nil, errors.New("prefix registered to the wrong URI")
			}

			// Okay, so we can register this, but we won't be
			// needing to register the namespace
			return doc.CreateElement(name)
		}
	} else {
		// If the name does not contain a prefix, check for the
		// existence of this namespace via the URI
		xmlnsuri := stringToXmlChar(nsuri)
		ns = C.xmlSearchNsByHref(
			doc.ptr,
			rootptr,
			xmlnsuri,
		)
		if ns != nil {
			// Well, you can safely inherit the prefix and stuff
			return doc.CreateElement(xmlCharToString(ns.prefix) + ":" + name)
		}

		debug.Printf("Create new namespace for %s", nsuri)
		// ns doesn't exist, got to create it here
		ns = C.xmlNewNs(nil, stringToXmlChar(nsuri), nil)
		// ... and my localname shall be == name
		localname = name
	}

	newNode := C.xmlNewDocNode(doc.ptr, ns, stringToXmlChar(localname), nil)
	newNode.nsDef = ns

	return wrapElement((*C.xmlElement)(unsafe.Pointer(newNode))), nil
}

func documentElement(doc *Document) *C.xmlNode {
	if doc.ptr == nil {
		return nil
	}

	return C.xmlDocGetRootElement(doc.ptr)
}

func xmlFreeDoc(d *Document) {
	C.xmlFreeDoc(d.ptr)
	d.ptr = nil
}

func documentString(d *Document, encoding string, format bool) string {
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
	C.xmlDocDumpFormatMemoryEnc(d.ptr, &xc, &i, C.CString(encoding), intformat)

	s := xmlCharToString(xc)
	return s
}

func xmlNodeSetBase(d *Document, s string) {
	C.xmlNodeSetBase((*C.xmlNode)(unsafe.Pointer(d.ptr)), stringToXmlChar(s))
}

func setDocumentElement(d *Document, n Node) {
	C.xmlDocSetRootElement(d.ptr, (*C.xmlNode)(n.Pointer()))
}

func setDocumentEncoding(d *Document, e string) {
	if d.ptr.encoding != nil {
		C.MY_xmlFree(unsafe.Pointer(d.ptr.encoding))
	}

	d.ptr.encoding = C.xmlStrdup(stringToXmlChar(e))
}

func setDocumentStandalone(d *Document, v int) {
	d.ptr.standalone = C.int(v)
}

func setDocumentVersion(d *Document, v string) {
	if d.ptr.version != nil {
		C.MY_xmlFree(unsafe.Pointer(d.ptr.version))
	}

	d.ptr.version = C.xmlStrdup(stringToXmlChar(v))
}

func xmlSetProp(n Node, name, value string) error {
	C.xmlSetProp((*C.xmlNode)(n.Pointer()), stringToXmlChar(name), stringToXmlChar(value))
	return nil
}

func (n *Element) getAttributeNode(name string) (*C.xmlAttr, error) {
	// if this is "xmlns", look for the first namespace without
	// the prefix
	if name == "xmlns" {
		for nsdef := n.ptr.nsDef; nsdef != nil; nsdef = nsdef.next {
			if nsdef.prefix != nil {
				continue
			}
			debug.Printf("nsdef.href -> %s", xmlCharToString(nsdef.href))
		}
	}

	debug.Printf("n = %s", n.String())
	debug.Printf("getAttributeNode(%s)", name)
	prop := C.xmlHasNsProp(n.ptr, stringToXmlChar(name), nil)
	debug.Printf("prop = %v", prop)
	if prop == nil {
		prefix, local := splitPrefixLocal(name)
		debug.Printf("prefix = %s, local = %s", prefix, local)
		if local != "" {
			if ns := C.xmlSearchNs(n.ptr.doc, n.ptr, stringToXmlChar(prefix)); ns != nil {
				prop = C.xmlHasNsProp(n.ptr, stringToXmlChar(local), ns.href)
			}

		}
	}

	if prop == nil || XmlNodeType(prop._type) != AttributeNode {
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

	i := C.xmlUnsetProp(nptr, stringToXmlChar(name))
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

	i := C.xmlUnsetNsProp(
		nptr,
		(*C.xmlNs)(unsafe.Pointer(ns.ptr)),
		stringToXmlChar(name),
	)
	if i == C.int(0) {
		return errors.New("failed to unset prop")
	}
	return nil
}

func xmlC14NDocDumpMemory(d *Document, mode C14NMode, withComments bool) (string, error) {
	var result *C.xmlChar

	var withCommentsInt C.int
	if withComments {
		withCommentsInt = 1
	}

	modeInt := C.int(mode)

	ret := C.xmlC14NDocDumpMemory(
		d.ptr,
		nil,
		modeInt,
		nil,
		withCommentsInt,
		&result,
	)
	if ret == 0 {
		return "", errors.New("c14n dump failed")
	}
	return xmlCharToString(result), nil
}
