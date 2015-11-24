package libxml2

func (n Text) Data() string {
	return xmlCharToString(n.ptr.content)
}
