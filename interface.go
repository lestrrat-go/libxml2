package libxml2

import "errors"

var (
	ErrInvalidNodeType = errors.New("invalid node type")
)

type Serializer interface {
	Serialize(interface{}) (string, error)
}

// note: Serialize takes an interface because some serializers only allow
// Document, whereas others might allow Nodes

// C14NMode represents the C14N mode supported by libxml2
type C14NMode int

const (
	C14N1_0 C14NMode = iota
	C14NExclusive1_0
	C14N1_1
)

// C14NSerialize implements the Serializer interface, and generates
// XML in C14N format.
type C14NSerialize struct {
	Mode         C14NMode
	WithComments bool
}
