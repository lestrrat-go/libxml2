package libxml2

import (
	"bytes"
	"errors"
	"strings"
)

func (n *Element) SetNamespace(uri, prefix string, activate ...bool) error {
	activateflag := false
	if len(activate) < 1 {
		activateflag = true
	} else {
		activateflag = activate[0]
	}

	if uri == "" && prefix == "" {
		// Empty namespace
		doc, err := n.OwnerDocument()
		if err != nil {
			return err
		}
		ns := xmlSearchNs(doc, n, "")
		if ns != nil && ns.URI() != "" {
			if activateflag {
				xmlSetNs(n, nil)
			}
		}
		return nil
	}

	if uri == "" {
		return errors.New("missing uri for SetNamespace")
	}
	if prefix == "" {
		return errors.New("missing prefix for SetNamespace")
	}

	ns := xmlNewNs(n, uri, prefix)
	if activateflag {
		xmlSetNs(n, ns)
	}
	return nil
}

func (n *Element) AppendText(s string) error {
	return appendText(n, s)
}

func (n *Element) SetAttribute(name, value string) error {
	return xmlSetProp(n, name, value)
}

func (n *Element) GetAttribute(name string) (*Attribute, error) {
	attrNode, err := n.getAttributeNode(name)
	if err != nil {
		return nil, err
	}
	return wrapAttribute(attrNode), nil
}

func (n *Element) Attributes() ([]*Attribute, error) {
	attrs := []*Attribute{}
	for attr := n.ptr.properties; attr != nil; attr = attr.next {
		attrs = append(attrs, wrapAttribute(attr))
	}
	return attrs, nil
}

func (n *Element) RemoveAttribute(name string) error {
	i := strings.IndexByte(name, ':')
	if i == -1 {
		return xmlUnsetProp(n, name)
	}

	// look for the prefix
	doc, err := n.OwnerDocument()
	if err != nil {
		return err
	}
	ns := xmlSearchNs(doc, n, name[:i])
	if ns == nil {
		return ErrAttributeNotFound
	}

	return xmlUnsetNsProp(n, ns, name)
}

// GetNamespaces returns Namespace objects associated with this
// element. WARNING: This method currently returns namespace
// objects which allocates C structures for each namespace.
// Therefore you MUST free the structures, or otherwise you
// WILL leak memory.
func (n *Element) GetNamespaces() []*Namespace {
	ret := []*Namespace{}
	for ns := n.ptr.nsDef; ns != nil; ns = ns.next {
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
	return ret
}

func (n Element) Literal() (string, error) {
	buf := bytes.Buffer{}
	children, err := n.ChildNodes()
	if err != nil {
		return "", err
	}
	for _, c := range children {
		l, err := c.Literal()
		if err != nil {
			return "", err
		}
		buf.WriteString(l)
	}
	return buf.String(), nil
}
