package libxml2

import (
	"errors"
	"unsafe"
)

func CreateDocument() *Document {
	return NewDocument("1.0", "")
}

func NewDocument(version, encoding string) *Document {
	return createDocument(version, encoding)
}

func (d *Document) Pointer() unsafe.Pointer {
	return unsafe.Pointer(d.ptr)
}

func (d *Document) CreateAttribute(k, v string) (*Attribute, error) {
	attr, err := xmlNewDocProp(d, k, v)
	if err != nil {
		return nil, err
	}
	return wrapAttribute(attr), nil
}

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

func (d *Document) CreateCDataSection(txt string) (*CDataSection, error) {
	cdata := xmlNewCDataBlock(d, txt)
	return wrapCDataSection(cdata), nil
}

func (d *Document) CreateCommentNode(txt string) (*Comment, error) {
	comment := xmlNewComment(txt)
	return wrapComment(comment), nil
}

func (d *Document) CreateElement(name string) (*Element, error) {
	return createElement(d, name)
}

func (d *Document) CreateElementNS(nsuri, name string) (*Element, error) {
	return createElementNS(d, nsuri, name)
}

func (d *Document) CreateTextNode(txt string) (*Text, error) {
	t := xmlNewText(txt)
	return wrapText(t), nil
}

func (d *Document) DocumentElement() (Node, error) {
	n := documentElement(d)
	if n == nil {
		return nil, ErrNodeNotFound
	}
	return wrapToNode(n)
}

func (d *Document) FindNodes(xpath string) (NodeList, error) {
	root, err := d.DocumentElement()
	if err != nil {
		return nil, err
	}
	return root.FindNodes(xpath)
}

func (d *Document) Encoding() string {
	return xmlCharToString(d.ptr.encoding)
}

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
