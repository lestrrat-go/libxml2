package dom

import (
	"errors"

	"github.com/lestrrat/go-libxml2/clib"
	"github.com/lestrrat/go-libxml2/types"
)

// CreateDocument creates a new document with version="1.0", and no encoding
func CreateDocument() *Document {
	return NewDocument("1.0", "")
}

// NewDocument creates a new document
func NewDocument(version, encoding string) *Document {
	ptr := clib.XMLCreateDocument(version, encoding)
	return &Document{ptr: ptr}
}

// Pointer returns the pointer to the underlying C struct
func (d *Document) Pointer() uintptr {
	return d.ptr
}

func (d *Document) AutoFree() {
	if !d.mortal {
		return
	}
	d.Free()
}

func (d *Document) MakeMortal() {
	d.mortal = true
}

func (d *Document) MakePersistent() {
	d.mortal = false
}

func (d *Document) IsSameNode(n types.Node) bool {
	return d.ptr == n.Pointer()
}

func (d *Document) HasChildNodes() bool {
	_, err := d.DocumentElement()
	return err != nil
}

func (d *Document) FirstChild() (types.Node, error) {
	root, err := d.DocumentElement()
	if err != nil {
		return nil, err
	}

	return root, nil
}

func (d *Document) LastChild() (types.Node, error) {
	root, err := d.DocumentElement()
	if err != nil {
		return nil, err
	}

	return root, nil
}

func (d *Document) NextSibling() (types.Node, error) {
	return nil, errors.New("document has no siblings")
}

func (d *Document) PreviousSibling() (types.Node, error) {
	return nil, errors.New("document has no siblings")
}

func (d *Document) NodeName() string {
	return ""
}

func (d *Document) SetNodeName(s string) {
//	return errors.New("cannot set node name on a document")
}

func (d *Document) NodeValue() string {
	return ""
}

func (d *Document) SetNodeValue(s string) {
//	return errors.New("cannot set node value on a document")
}

func (d *Document) OwnerDocument() (types.Document, error) {
	return d, nil
}

func (d *Document) SetDocument(n types.Document) error {
	return errors.New("cannot set document on a document")
}

func (d *Document) ParentNode() (types.Node, error) {
	return nil, errors.New("document has no parent node")
}

func (d *Document) ParseInContext(s string, n int) (types.Node, error) {
	return nil, errors.New("unimplemented")
}

func (d *Document) Literal() (string, error) {
	return d.Dump(false), nil
}

// TextContent returns the text content
func (n *Document) TextContent() string {
	return clib.XMLTextContent(n)
}

func (n *Document) ToString(x int, b bool) string {
	return n.Dump(false)
}

func (d *Document) ChildNodes() (types.NodeList, error) {
	root, err := d.DocumentElement()
	if err != nil {
		return nil, err
	}

	return []types.Node{root}, nil
}

func (d *Document) Copy() (types.Node, error) {
	// Unimplemented
	return nil, errors.New("unimplemented")
}

// AddChild is a no op for Document
func (d *Document) AddChild(n types.Node) error {
	return errors.New("method AddChild is not available for Document node")
}

// CreateAttribute creates a new attribute
func (d *Document) CreateAttribute(k, v string) (*Attribute, error) {
	attr, err := clib.XMLNewDocProp(d, k, v)
	if err != nil {
		return nil, err
	}
	return wrapAttribute(attr), nil
}

// CreateAttributeNS creates a new attribute with the given XML namespace
func (d *Document) CreateAttributeNS(nsuri, k, v string) (*Attribute, error) {
	if nsuri == "" {
		return d.CreateAttribute(k, v)
	}

	ptr, err := clib.XMLCreateAttributeNS(d, nsuri, k, v)
	if err != nil {
		return nil, err
	}
	return wrapAttribute(ptr), nil
}

// CreateCDataSection creates a new CDATA section node
func (d *Document) CreateCDataSection(txt string) (*CDataSection, error) {
	cdata, err := clib.XMLNewCDataBlock(d, txt)
	if err != nil {
		return nil, err
	}
	return wrapCDataSection(cdata), nil
}

// CreateCommentNode creates a new comment node
func (d *Document) CreateCommentNode(txt string) (*Comment, error) {
	ptr, err := clib.XMLNewComment(txt)
	if err != nil {
		return nil, err
	}
	return wrapComment(ptr), nil
}

// CreateElement creates a new element node
func (d *Document) CreateElement(name string) (types.Element, error) {
	ptr, err := clib.XMLCreateElement(d, name)
	if err != nil {
		return nil, err
	}
	return wrapElement(ptr), nil
}

// CreateElementNS creates a new element node in the given XML namespace
func (d *Document) CreateElementNS(nsuri, name string) (types.Element, error) {
	ptr, err := clib.XMLCreateElementNS(d, nsuri, name)
	if err != nil {
		return nil, err
	}
	return wrapElement(ptr), nil
}

// CreateTextNode creates a new text node
func (d *Document) CreateTextNode(txt string) (*Text, error) {
	ptr, err := clib.XMLNewText(txt)
	if err != nil {
		return nil, err
	}
	return wrapText(ptr), nil
}

// DocumentElement returns the root node of the document
func (d *Document) DocumentElement() (types.Node, error) {
	n, err := clib.XMLDocumentElement(d)
	if err != nil {
		return nil, err
	}
	return WrapNode(n)
}

// Find returns the nodes that can be selected with the
// given xpath string
func (d *Document) Find(xpath string) (types.XPathResult, error) {
	root, err := d.DocumentElement()
	if err != nil {
		return nil, err
	}
	return root.Find(xpath)
}

// Encoding returns the d
func (d *Document) Encoding() string {
	return clib.XMLDocumentEncoding(d)
}

// Free releases the underlying C struct
func (d *Document) Free() {
	clib.XMLFreeDoc(d)
	d.ptr = 0
}

// String formats the document, always without formatting.
func (d *Document) String() string {
	return clib.XMLDocumentString(d, d.Encoding(), false)
}

// Dump formats the document with or withour formatting.
func (d *Document) Dump(format bool) string {
	return clib.XMLDocumentString(d, d.Encoding(), format)
}

// NodeType returns the XMLNodeType
func (d *Document) NodeType() clib.XMLNodeType {
	return DocumentNode
}

// SetBaseURI sets the base URI
func (d *Document) SetBaseURI(s string) {
	clib.XMLNodeSetBase(d, s)
}

// SetDocumentElement sets the document element
func (d *Document) SetDocumentElement(n types.Node) error {
	return clib.XMLSetDocumentElement(d, n)
}

// SetEncoding sets the encoding of the document
func (d *Document) SetEncoding(e string) {
	clib.XMLSetDocumentEncoding(d, e)
}

// SetStandalone sets the standalone flag
func (d *Document) SetStandalone(v int) {
	clib.XMLSetDocumentStandalone(d, v)
}

// SetVersion sets the version of the document
func (d *Document) SetVersion(v string) {
	clib.XMLSetDocumentVersion(d, v)
}

// Standalone returns the value of the standalone flag
func (d *Document) Standalone() int {
	return clib.XMLDocumentStandalone(d)
}

// URI returns the document URI
func (d *Document) URI() string {
	return clib.XMLDocumentURI(d)
}

// Version returns the version of the document
func (d *Document) Version() string {
	return clib.XMLDocumentVersion(d)
}

// Walk traverses the nodes in the document
func (d *Document) Walk(fn func(types.Node) error) error {
	root, err := d.DocumentElement()
	if err != nil {
		return err
	}
	walk(root, fn)
	return nil
}

func (d *Document) LookupNamespacePrefix(href string) (string, error) {
	root, err := d.DocumentElement()
	if err != nil {
		return "", err
	}

	return root.LookupNamespacePrefix(href)
}

func (d *Document) LookupNamespaceURI(prefix string) (string, error) {
	root, err := d.DocumentElement()
	if err != nil {
		return "", err
	}

	return root.LookupNamespaceURI(prefix)
}

