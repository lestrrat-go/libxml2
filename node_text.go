package libxml2

// Data returns the content associated with this node
func (n *Text) Data() string {
	return xmlTextData(n)
}
