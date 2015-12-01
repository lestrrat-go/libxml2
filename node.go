package libxml2

import (
	"bytes"
	"fmt"
)

var _XMLNodeTypeIndex = [...]uint8{0, 11, 24, 32, 48, 61, 71, 77, 88, 100, 116, 132, 144, 160, 167, 178, 191, 201, 214, 227, 238, 254}

const _XMLNodeTypeName = `ElementNodeAttributeNodeTextNodeCDataSectionNodeEntityRefNodeEntityNodePiNodeCommentNodeDocumentNodeDocumentTypeNodeDocumentFragNodeNotationNodeHTMLDocumentNodeDTDNodeElementDeclAttributeDeclEntityDeclNamespaceDeclXIncludeStartXIncludeEndDocbDocumentNode`

// String returns the string representation of this XMLNodeType
func (i XMLNodeType) String() string {
	x := i - 1
	if x < 0 || x+1 >= XMLNodeType(len(_XMLNodeTypeIndex)) {
		return fmt.Sprintf("XMLNodeType(%d)", x+1)
	}
	return _XMLNodeTypeName[_XMLNodeTypeIndex[x]:_XMLNodeTypeIndex[x+1]]
}

// String returns the string representation
func (n *XMLNode) String() string {
	return n.ToString(0, false)
}

// NodeName returns the node name
func (n *XMLNode) NodeName() string {
	return nodeName(n)
}

// NodeValue returns the node value
func (n *XMLNode) NodeValue() string {
	return nodeValue(n)
}

// Literal returns the literal string value
func (n XMLNode) Literal() (string, error) {
	return n.String(), nil
}

// IsSameNode returns true if two nodes point to the same node
func (n *XMLNode) IsSameNode(other Node) bool {
	return n.Pointer() == other.Pointer()
}

// Copy creates a copy of the node
func (n *XMLNode) Copy() (Node, error) {
	doc, err := n.OwnerDocument()
	if err != nil {
		return nil, err
	}
	return xmlDocCopyNode(n, doc, 1)
}

// SetDocument sets the document of this node and its descendants
func (n *XMLNode) SetDocument(d *Document) error {
	return xmlSetTreeDoc(n, d)
}

// ParseInContext parses a chunk of XML in the context of the current
// node. This makes it safe to append the resulting node to the current
// node or other nodes in the same document.
func (n *XMLNode) ParseInContext(s string, o ParseOption) (Node, error) {
	return xmlParseInNodeContext(n, s, o)
}

// String returns the string representation of the NodeList
func (n NodeList) String() string {
	buf := bytes.Buffer{}
	for _, x := range n {
		buf.WriteString(x.String())
	}
	return buf.String()
}

func (n NodeList) NodeValue() string {
	buf := bytes.Buffer{}
	for _, x := range n {
		buf.WriteString(x.NodeValue())
	}
	return buf.String()
}

// Literal returns the string representation of the NodeList (using Literal())
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
