package libxml2

var _XMLNodeType_index = [...]uint8{0, 11, 24, 32, 48, 61, 71, 77, 88, 100, 116, 132, 144, 160, 167, 178, 191, 201, 214, 227, 238, 254}

const _XMLNodeType_name = `ElementNodeAttributeNodeTextNodeCDataSectionNodeEntityRefNodeEntityNodePiNodeCommentNodeDocumentNodeDocumentTypeNodeDocumentFragNodeNotationNodeHTMLDocumentNodeDTDNodeElementDeclAttributeDeclEntityDeclNamespaceDeclXIncludeStartXIncludeEndDocbDocumentNode`

func (n *XMLNode) Copy() (Node, error) {
	doc, err := n.OwnerDocument()
	if err != nil {
		return nil, err
	}
	return xmlDocCopyNode(n, doc, 1)
}

func (n *XMLNode) SetDocument(d *Document) error {
	return xmlSetTreeDoc(n, d)
}

func (n *XMLNode) ParseInContext(s string, o ParseOption) (Node, error) {
	return xmlParseInNodeContext(n, s, o)
}
