package libxml2

func (n Text) Data() string {
	return xmlCharToString(n.ptr.content)
}

func (n *Text) Walk(fn func(Node) error) {
	walk(n, fn)
}
