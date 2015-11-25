package libxml2

func (n Text) Data() string {
	ptr := n.ptr
	if ptr == nil {
		return ""
	}
	return xmlCharToString(ptr.content)
}
