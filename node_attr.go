package libxml2

func (n *Attribute) Free() {
	xmlFreeProp(n)
}

func (n *Attribute) HasChildNodes() bool {
	return false
}

func (n *Attribute) Value() string {
	return nodeValue(n)
}


