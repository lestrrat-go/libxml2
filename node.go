package libxml2

var _XMLNodeTypeIndex = [...]uint8{0, 11, 24, 32, 48, 61, 71, 77, 88, 100, 116, 132, 144, 160, 167, 178, 191, 201, 214, 227, 238, 254}

const _XMLNodeTypeName = `ElementNodeAttributeNodeTextNodeCDataSectionNodeEntityRefNodeEntityNodePiNodeCommentNodeDocumentNodeDocumentTypeNodeDocumentFragNodeNotationNodeHTMLDocumentNodeDTDNodeElementDeclAttributeDeclEntityDeclNamespaceDeclXIncludeStartXIncludeEndDocbDocumentNode`

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
