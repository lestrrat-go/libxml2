package libxml2

/*
#cgo pkg-config: libxml-2.0
#include <stdbool.h>
#include "libxml/globals.h"
#include "libxml/tree.h"
#include "libxml/parser.h"
#include "libxml/parserInternals.h"
#include "libxml/xpath.h"

// Macro wrapper function
static inline bool MY_xmlXPathNodeSetIsEmpty(xmlNodeSetPtr ptr) {
	return xmlXPathNodeSetIsEmpty(ptr);
}

// Macro wrapper function
static inline void MY_xmlFree(void *p) {
	xmlFree(p);
}

// Because Go can't do pointer airthmetics...
static inline xmlNodePtr MY_xmlNodeSetTabAt(xmlNodePtr *nodes, int i) {
	return nodes[i];
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
)

type XmlNodeType int

const (
	ElementNode XmlNodeType = iota + 1
	AttributeNode
	TextNode
	CDataSectionNode
	EntityRefNode
	EntityNode
	PiNode
	CommentNode
	DocumentNode
	DocumentTypeNode
	DocumentFragNode
	NotationNode
	HTMLDocumentNode
	DTDNode
	ElementDecl
	AttributeDecl
	EntityDecl
	NamespaceDecl
	XIncludeStart
	XIncludeEnd
	DocbDocumentNode
)

var _XmlNodeType_index = [...]uint8{0, 11, 24, 32, 48, 61, 71, 77, 88, 100, 116, 132, 144, 160, 167, 178, 191, 201, 214, 227, 238, 254}

const _XmlNodeType_name = `ElementNodeAttributeNodeTextNodeCDataSectionNodeEntityRefNodeEntityNodePiNodeCommentNodeDocumentNodeDocumentTypeNodeDocumentFragNodeNotationNodeHTMLDocumentNodeDTDNodeElementDeclAttributeDeclEntityDeclNamespaceDeclXIncludeStartXIncludeEndDocbDocumentNode`

func (i XmlNodeType) String() string {
	i -= 1
	if i < 0 || i+1 >= XmlNodeType(len(_XmlNodeType_index)) {
		return fmt.Sprintf("XmlNodeType(%d)", i+1)
	}
	return _XmlNodeType_name[_XmlNodeType_index[i]:_XmlNodeType_index[i+1]]
}

var (
	ErrNodeNotFound = errors.New("node not found")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrInvalidNodeName = errors.New("invalid node name")
)

// Node defines the basic DOM interface
type Node interface {
	// pointer() returns the underlying C pointer. Only we are allowed to
	// slice it, dice it, do whatever the heck with it.
	pointer() unsafe.Pointer

	AddChild(Node)
	AppendChild(Node) error
	ChildNodes() []Node
	OwnerDocument() *Document
	FindNodes(string) ([]Node, error)
	FirstChild() Node
	HasChildNodes() bool
	IsSameNode(Node) bool
	LastChild() Node
	NextSibling() Node
	NodeName() string
	NodeType() XmlNodeType
	NodeValue() string
	ParetNode() Node
	PreviousSibling() Node
	SetNodeName(string)
	String() string
	TextContent() string
	ToString(int, bool) string
	Walk(func(Node) error)
}

type xmlNode struct {
	ptr *C.xmlNode
}

type XmlNode struct {
	*xmlNode
}

type Attribute struct {
	*XmlNode
}

type CDataSection struct {
	*XmlNode
}

type Comment struct {
	*XmlNode
}

type Element struct {
	*XmlNode
}

type Document struct {
	ptr  *C.xmlDoc
	root *C.xmlNode
}

type Text struct {
	*XmlNode
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
		&xmlNode{
			ptr: (*C.xmlNode)(unsafe.Pointer(n)),
		},
	}
}

func wrapText(n *C.xmlNode) *Text {
	return &Text{wrapXmlNode(n)}
}

func wrapToNode(n *C.xmlNode) Node {
	switch XmlNodeType(n._type) {
	case ElementNode:
		return wrapElement((*C.xmlElement)(unsafe.Pointer(n)))
	case TextNode:
		return &Text{&XmlNode{&xmlNode{ptr: n}}}
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

func nodeName(n Node) string {
	switch n.NodeType() {
	case XIncludeStart, XIncludeEnd, EntityRefNode, EntityNode, DTDNode, EntityDecl, DocumentTypeNode, NotationNode, NamespaceDecl:
		return xmlCharToString((*C.xmlNode)(n.pointer()).name)
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
		ptr := (*C.xmlNode)(n.pointer())
		if ns := ptr.ns; ns != nil {
			return fmt.Sprintf("%s:%s", xmlCharToString(ns.prefix), xmlCharToString(ptr.name))
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
		return xmlCharToString(C.xmlXPathCastNodeToString((*C.xmlNode)(n.pointer())))
	case EntityDecl:
		np := (*C.xmlNode)(n.pointer())
		if np.content != nil {
			return xmlCharToString(C.xmlStrdup(np.content))
		}

		panic("unimplmented")
	}

	return ""
}

func (n *xmlNode) pointer() unsafe.Pointer {
	return unsafe.Pointer(n.ptr)
}

func (n *xmlNode) AddChild(child Node) {
	C.xmlAddChild(n.ptr, (*C.xmlNode)(child.pointer()))
}

func (n *xmlNode) AppendChild(child Node) error {
	// XXX There must be lots more checks here because AddChild does things
	// under the table like merging text nodes, freeing some nodes implicitly,
	// et al
	n.AddChild(child)
	return nil
}

func (n *xmlNode) ChildNodes() []Node {
	return childNodes(n)
}

func wrapDocument(n *C.xmlDoc) *Document {
	r := C.xmlDocGetRootElement(n) // XXX Should check for n == nil
	return &Document{ptr: n, root: r}
}

func (n *xmlNode) OwnerDocument() *Document {
	return wrapDocument(n.ptr.doc)
}

func (n *xmlNode) FindNodes(xpath string) ([]Node, error) {
	return findNodes(n, xpath)
}

func (n *xmlNode) FirstChild() Node {
	if !n.HasChildNodes() {
		return nil
	}

	return wrapToNode(((*C.xmlNode)(n.pointer())).children)
}

func (n *xmlNode) HasChildNodes() bool {
	return n.ptr.children != nil
}

func (n *xmlNode) IsSameNode(other Node) bool {
	return n.pointer() == other.pointer()
}

func (n *xmlNode) LastChild() Node {
	return wrapToNode(n.ptr.last)
}

func (n *xmlNode) LocalName() string {
	switch n.NodeType() {
	case ElementNode, AttributeNode, ElementDecl, AttributeDecl:
		return xmlCharToString(n.ptr.name)
	}
	return ""
}

func (n *xmlNode) NamespaceURI() string {
	switch n.NodeType() {
	case ElementNode, AttributeNode, PiNode:
		if ns := n.ptr.ns; ns != nil && ns.href != nil {
			return xmlCharToString(ns.href)
		}
	}
	return ""
}

func (n *xmlNode) NodeName() string {
	return nodeName(n)
}

func (n *xmlNode) NodeValue() string {
	return nodeValue(n)
}

func (n *xmlNode) NextSibling() Node {
	return wrapToNode(n.ptr.next)
}

func (n *xmlNode) ParetNode() Node {
	return wrapToNode(n.ptr.parent)
}

func (n *xmlNode) Prefix() string {
	switch n.NodeType() {
	case ElementNode, AttributeNode, PiNode:
		if ns := n.ptr.ns; ns != nil && ns.prefix != nil {
			return xmlCharToString(ns.prefix)
		}
	}
	return ""
}

func (n *xmlNode) PreviousSibling() Node {
	return wrapToNode(n.ptr.prev)
}

func (n *xmlNode) SetNodeName(name string) {
	C.xmlNodeSetName(n.ptr, stringToXmlChar(name))
}

func (n *xmlNode) String() string {
	return n.ToString(0, false)
}

func (n *xmlNode) TextContent() string {
	return xmlCharToString(C.xmlXPathCastNodeToString(n.ptr))
}

func (n *xmlNode) ToString(format int, docencoding bool) string {
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

func (n *xmlNode) NodeType() XmlNodeType {
	return XmlNodeType(n.ptr._type)
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

func CreateDocument() *Document {
	return NewDocument("1.0", "")
}

func NewDocument(version, encoding string) *Document {
	doc := C.xmlNewDoc(stringToXmlChar(version))
	if encoding != "" {
		doc.encoding = C.xmlStrdup(stringToXmlChar(encoding))
	}
	return wrapDocument(doc)
}

func (d *Document) pointer() unsafe.Pointer {
	return unsafe.Pointer(d.ptr)
}

func (d *Document) CreateAttribute(k, v string) (*Attribute, error) {
	kx := stringToXmlChar(k)
	vx := stringToXmlChar(v)
	if C.MY_test_node_name(kx) == 0 {
		return nil, ErrInvalidNodeName
	}

	buf := C.xmlEncodeEntitiesReentrant(d.ptr, vx)
	newAttr := C.xmlNewDocProp(d.ptr, kx, buf)

	return wrapAttribute((*C.xmlAttr)(unsafe.Pointer(newAttr))), nil
}

func (d *Document) CreateAttributeNS(nsuri, k, v string) (*Attribute, error) {
	if nsuri == "" {
		return d.CreateAttribute(k, v)
	}

	kx := stringToXmlChar(k)
	if C.MY_test_node_name(kx) == 0 {
		return nil, ErrInvalidNodeName
	}

	root := d.DocumentElement()
	if root == nil {
		return nil, errors.New("attribute with namespaces require a root node")
	}

	prefix, local := splitPrefixLocal(k)

	ns := C.xmlSearchNsByHref(d.ptr, (*C.xmlNode)(root.pointer()), stringToXmlChar(nsuri))
	if ns == nil {
		ns = C.xmlNewNs((*C.xmlNode)(root.pointer()), stringToXmlChar(nsuri), stringToXmlChar(prefix))
		if ns == nil {
			return nil, errors.New("failed to create namespace")
		}
	}

	vx := stringToXmlChar(v)
	buf := C.xmlEncodeEntitiesReentrant(d.ptr, vx)
	newAttr := C.xmlNewDocProp(d.ptr, stringToXmlChar(local), buf)
	C.xmlSetNs((*C.xmlNode)(unsafe.Pointer(newAttr)), ns)

	return wrapAttribute((*C.xmlAttr)(unsafe.Pointer(newAttr))), nil
}

func (d *Document) CreateCDataSection(txt string) (*CDataSection, error) {
	return wrapCDataSection(C.xmlNewCDataBlock(d.ptr, stringToXmlChar(txt), C.int(len(txt)))), nil
}

func (d *Document) CreateCommentNode(txt string) (*Comment, error) {
	return wrapComment(C.xmlNewComment(stringToXmlChar(txt))), nil
}

func (d *Document) CreateElement(name string) (*Element, error) {
	if C.MY_test_node_name(stringToXmlChar(name)) == 0 {
		return nil, ErrInvalidNodeName
	}

	newNode := C.xmlNewNode(nil, stringToXmlChar(name))
	if newNode == nil {
		return nil, errors.New("element creation failed")
	}
	// XXX hmmm...
	newNode.doc = d.ptr
	return wrapElement((*C.xmlElement)(unsafe.Pointer(newNode))), nil
}

func (d *Document) CreateElementNS(nsuri, name string) (*Element, error) {
	if C.MY_test_node_name(stringToXmlChar(name)) == 0 {
		return nil, ErrInvalidNodeName
	}

	i := strings.IndexByte(name, ':')
	nsuriDup := stringToXmlChar(nsuri)
	prefix := stringToXmlChar(name[:i])
	localname := stringToXmlChar(name[i+1:])

	ns := C.xmlNewNs(nil, nsuriDup, prefix)
	newNode := C.xmlNewDocNode(d.ptr, ns, localname, nil)
	newNode.nsDef = ns

	return wrapElement((*C.xmlElement)(unsafe.Pointer(newNode))), nil
}

func (d *Document) CreateTextNode(txt string) (*Text, error) {
	return wrapText(C.xmlNewText(stringToXmlChar(txt))), nil
}

func (d *Document) DocumentElement() Node {
	if d.ptr == nil {
		return nil
	}

	if d.root == nil {
		n := C.xmlDocGetRootElement(d.ptr)
		if n == nil {
			return nil
		}
		d.root = n
	}

	return wrapToNode(d.root)
}

func (d *Document) FindNodes(xpath string) ([]Node, error) {
	root := d.DocumentElement()
	if root == nil {
		return nil, ErrNodeNotFound
	}
	return root.FindNodes(xpath)
}

func (d *Document) Encoding() string {
	return xmlCharToString(d.ptr.encoding)
}

func (d *Document) Free() {
	C.xmlFreeDoc(d.ptr)
	d.ptr = nil
	d.root = nil
}

func (d *Document) String() string {
	var xc *C.xmlChar
	i := C.int(0)
	C.xmlDocDumpMemory(d.ptr, &xc, &i)
	return xmlCharToString(xc)
}

func (d *Document) NodeType() XmlNodeType {
	return XmlNodeType(d.ptr._type)
}

func (d *Document) SetBaseURI(s string) {
	C.xmlNodeSetBase((*C.xmlNode)(unsafe.Pointer(d.ptr)), stringToXmlChar(s))
}

func (d *Document) SetDocumentElement(n Node) {
	C.xmlDocSetRootElement(d.ptr, (*C.xmlNode)(n.pointer()))
	d.root = (*C.xmlNode)(n.pointer())
}

func (d *Document) SetEncoding(e string) {
	if d.ptr.encoding != nil {
		C.MY_xmlFree(unsafe.Pointer(d.ptr.encoding))
	}

	d.ptr.encoding = C.xmlStrdup(stringToXmlChar(e))
}

func (d *Document) SetStandalone(v int) {
	d.ptr.standalone = C.int(v)
}

func (d *Document) SetVersion(e string) {
	if d.ptr.version != nil {
		C.MY_xmlFree(unsafe.Pointer(d.ptr.version))
	}

	d.ptr.version = C.xmlStrdup(stringToXmlChar(e))
}

func (d *Document) Standalone() int {
	return int(d.ptr.standalone)
}

func (d *Document) ToString(skipXmlDecl bool) string {
	buf := &bytes.Buffer{}
	for _, n := range childNodes(wrapXmlNode((*C.xmlNode)(d.pointer()))) {
		if n.NodeType() == DTDNode {
			continue
		}
		buf.WriteString(n.String())
	}

	return buf.String()
}

func (d *Document) URI() string {
	return xmlCharToString(C.xmlStrdup(d.ptr.URL))
}

func (d *Document) Version() string {
	return xmlCharToString(d.ptr.version)
}

func (d *Document) Walk(fn func(Node) error) {
	walk(wrapXmlNode(d.root), fn)
}

func (n *Element) AppendText(s string) error {
	txt, err := n.OwnerDocument().CreateTextNode(s)
	if err != nil {
		return err
	}
	return n.AppendChild(txt)
}

func splitPrefixLocal(s string) (string, string) {
	i := strings.IndexByte(s, ':')
	if i == -1 {
		return "", s
	}
	return s[:i], s[i+1:]
}

func (n *Element) getAttributeNode(name string) (*C.xmlAttr, error) {
	prop := C.xmlHasNsProp(n.ptr, stringToXmlChar(name), nil)
	if prop != nil {
		prefix, local := splitPrefixLocal(name)
		if local != "" {
			ns := C.xmlSearchNs(n.ptr.doc, n.ptr, stringToXmlChar(prefix))
			if ns != nil {
				prop = C.xmlHasNsProp(n.ptr, stringToXmlChar(local), ns.href)
			}
		}
	}

	if prop == nil || XmlNodeType(prop._type) != AttributeNode {
		return nil, errors.New("attribute not found")
	}

	return prop, nil
}

func (n *Element) RemoveAttribute(name string) error {
	prop, err := n.getAttributeNode(name)
	if err != nil {
		return err
	}

	C.xmlUnlinkNode((*C.xmlNode)(unsafe.Pointer(prop)))
	C.xmlFreeProp(prop)

	return nil
}

func (n *Text) Data() string {
	return xmlCharToString(n.ptr.content)
}

func (n *Text) Walk(fn func(Node) error) {
	walk(n, fn)
}

func (n *Attribute) HasChildNodes() bool {
	return false
}
