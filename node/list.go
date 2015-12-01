package node

import "bytes"

// String returns the string representation of the List
func (n List) String() string {
	buf := bytes.Buffer{}
	for _, x := range n {
		buf.WriteString(x.String())
	}
	return buf.String()
}

// NodeValue returns the concatenation of NodeValue within the nodes in List
func (n List) NodeValue() string {
	buf := bytes.Buffer{}
	for _, x := range n {
		buf.WriteString(x.NodeValue())
	}
	return buf.String()
}

// Literal returns the string representation of the List (using Literal())
func (n List) Literal() (string, error) {
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
