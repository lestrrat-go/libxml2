package xpath

import (
	"github.com/lestrrat/go-libxml2/node"
)

// String returns the String component of the result, and as a side effect
// releases the Result by calling Free() on it
func String(r Result, err error) string {
	if err != nil {
		return ""
	}

	defer r.Free()
	return r.String()
}

func Bool(r Result, err error) bool {
	if err != nil {
		return false
	}

	defer r.Free()
	return r.Bool()
}

func Number(r Result, err error) float64 {
	if err != nil {
		return 0
	}

	defer r.Free()
	return r.Number()
}

func NodeList(r Result, err error) node.List {
	if err != nil {
		return nil
	}

	defer r.Free()
	return r.NodeList()
}
