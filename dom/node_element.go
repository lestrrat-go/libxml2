package dom

import (
	"bytes"
	"errors"
	"strings"

	"github.com/lestrrat-go/libxml2/clib"
	"github.com/lestrrat-go/libxml2/types"
)

// SetNamespace sets up a new namespace on the given node.
// An XML namespace declaration is explicitly created, provided that
// it is not already declared in the tree.
//
// The activate flag should be considered deprecated, and will probably
// be removed in the future
func (n *Element) SetNamespace(uri, prefix string, activate ...bool) error {
	var activateflag bool
	if len(activate) < 1 {
		activateflag = true
	} else {
		activateflag = activate[0]
	}

	doc, err := n.OwnerDocument()
	if err != nil {
		return err
	}
	if uri == "" && prefix == "" {
		// Empty namespace
		nsptr, err := clib.XMLSearchNs(doc, n, "")
		if err != nil {
			return err
		}

		ns := wrapNamespaceNode(nsptr)
		if ns.URI() != "" {
			if activateflag {
				_ = clib.XMLSetNs(n, nil)
			}
		}
		return nil
	}

	if uri == "" {
		return errors.New("missing uri for SetNamespace")
	}

	root, err := doc.DocumentElement()
	if err != nil {
		return err
	}
	if _, err := clib.XMLSearchNs(doc, root, prefix); err != nil {
		// Namespace not found, create a new one
		ns, err := clib.XMLNewNs(n, uri, prefix)
		if err != nil {
			return err
		}

		if activateflag {
			if err := clib.XMLSetNs(n, wrapNamespaceNode(ns)); err != nil {
				return err
			}
		}
	}
	return nil
}

// AppendText adds a new text node
func (n *Element) AppendText(s string) error {
	return clib.XMLAppendText(n, s)
}

// SetAttribute sets an attribute
func (n *Element) SetAttribute(name, value string) error {
	return clib.XMLSetProp(n, name, value)
}

// GetAttribute retrieves the value of an attribute
func (n *Element) GetAttribute(name string) (types.Attribute, error) {
	attrNode, err := clib.XMLElementGetAttributeNode(n, name)
	if err != nil {
		return nil, err
	}
	return wrapAttributeNode(attrNode), nil
}

// Attributes returns a list of attributes on a node
func (n *Element) Attributes() ([]types.Attribute, error) {
	attrs, err := clib.XMLElementAttributes(n)
	if err != nil {
		return nil, err
	}
	ret := make([]types.Attribute, len(attrs))
	for i, attr := range attrs {
		ret[i] = wrapAttributeNode(attr)
	}
	return ret, nil
}

// RemoveAttribute completely removes an attribute from the node
func (n *Element) RemoveAttribute(name string) error {
	i := strings.IndexByte(name, ':')
	if i == -1 {
		return clib.XMLUnsetProp(n, name)
	}

	// look for the prefix
	doc, err := n.OwnerDocument()
	if err != nil {
		return err
	}
	ns, err := clib.XMLSearchNs(doc, n, name[:i])
	if err != nil {
		return ErrAttributeNotFound
	}

	return clib.XMLUnsetNsProp(n, wrapNamespaceNode(ns), name)
}

// GetNamespaces returns Namespace objects associated with this
// element. WARNING: This method currently returns namespace
// objects which allocates C structures for each namespace.
// Therefore you MUST free the structures, or otherwise you
// WILL leak memory.
func (n *Element) GetNamespaces() ([]types.Namespace, error) {
	list, err := clib.XMLElementNamespaces(n)
	if err != nil {
		return nil, err
	}
	ret := make([]types.Namespace, len(list))
	for i, nsptr := range list {
		ret[i] = wrapNamespaceNode(nsptr)
	}
	return ret, nil
}

// Literal returns a stringified version of this node and its
// children, inclusive.
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
