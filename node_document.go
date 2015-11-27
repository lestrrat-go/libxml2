package libxml2

import (
	"errors"
	"unsafe"
)

// CreateDocument creates a new document with version="1.0", and no encoding
func CreateDocument() *Document {
	return NewDocument("1.0", "")
}

// NewDocument creates a new document
func NewDocument(version, encoding string) *Document {
	return createDocument(version, encoding)
}

// Pointer returns the pointer to the underlying C struct
func (d *Document) Pointer() unsafe.Pointer {
	return unsafe.Pointer(d.ptr)
}

// CreateAttribute creates a new attribute
func (d *Document) CreateAttribute(k, v string) (*Attribute, error) {
	attr, err := xmlNewDocProp(d, k, v)
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

	if err := myTestNodeName(k); err != nil {
		return nil, err
	}

	root, err := d.DocumentElement()
	if err != nil {
		return nil, errors.New("attribute with namespaces require a root node")
	}

	prefix, local := splitPrefixLocal(k)

	ns := xmlSearchNsByHref(d, root, nsuri)
	if ns == nil {
		ns = xmlNewNs(root, nsuri, prefix)
		if ns == nil {
			return nil, errors.New("failed to create namespace")
		}
	}

	newAttr, err := xmlNewDocProp(d, local, v)
	if err != nil {
		return nil, err
	}
	attr := wrapAttribute(newAttr)
	xmlSetNs(attr, ns)

	return wrapAttribute(newAttr), nil
}

// CreateCDataSection creates a new CDATA section node
func (d *Document) CreateCDataSection(txt string) (*CDataSection, error) {
	cdata := xmlNewCDataBlock(d, txt)
	return wrapCDataSection(cdata), nil
}

// CreatesCommentNode creates a new comment node
func (d *Document) CreateCommentNode(txt string) (*Comment, error) {
	comment := xmlNewComment(txt)
	return wrapComment(comment), nil
}

// CreateElement creates a new element node
func (d *Document) CreateElement(name string) (*Element, error) {
	return createElement(d, name)
}

// CreateElementNS creates a new element node in the given XML namespace
func (d *Document) CreateElementNS(nsuri, name string) (*Element, error) {
	return createElementNS(d, nsuri, name)
}

// CreateTextNode creates a new text node
func (d *Document) CreateTextNode(txt string) (*Text, error) {
	t := xmlNewText(txt)
	return wrapText(t), nil
}

// DocumentElement returns the root node of the document
func (d *Document) DocumentElement() (Node, error) {
	n := documentElement(d)
	if n == nil {
		return nil, ErrNodeNotFound
	}
	return wrapToNode(n)
}

// FindNodes returns the nodes that can be selected with the
// given xpath string
func (d *Document) FindNodes(xpath string) (NodeList, error) {
	root, err := d.DocumentElement()
	if err != nil {
		return nil, err
	}
	return root.FindNodes(xpath)
}

// Encoding returns the d
func (d *Document) Encoding() string {
	return documentEncoding(d)
}

// Free releases the underlying C struct
func (d *Document) Free() {
	xmlFreeDoc(d)
}

func (d *Document) String() string {
	return documentString(d, d.Encoding(), false)
}

func (d *Document) Dump(format bool) string {
	return documentString(d, d.Encoding(), format)
}

func (d *Document) NodeType() XMLNodeType {
	return XMLNodeType(d.ptr._type)
}

func (d *Document) SetBaseURI(s string) {
	xmlNodeSetBase(d, s)
}

func (d *Document) SetDocumentElement(n Node) error {
	return setDocumentElement(d, n)
}

func (d *Document) SetEncoding(e string) {
	setDocumentEncoding(d, e)
}

func (d *Document) SetStandalone(v int) {
	setDocumentStandalone(d, v)
}

func (d *Document) SetVersion(v string) {
	setDocumentVersion(d, v)
}

func (d *Document) Standalone() int {
	return int(d.ptr.standalone)
}

func (d *Document) URI() string {
	return xmlCharToString(d.ptr.URL)
}

func (d *Document) Version() string {
	return xmlCharToString(d.ptr.version)
}

func (d *Document) Walk(fn func(Node) error) error {
	root, err := d.DocumentElement()
	if err != nil {
		return err
	}
	walk(root, fn)
	return nil
}
