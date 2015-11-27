package libxml2

// Data returns the content associated with this node
func (n Text) Data() string {
	ptr := n.ptr
	if ptr == nil {
		return ""
	}
	return xmlCharToString(ptr.content)
}
