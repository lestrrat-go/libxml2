package dom

import (
	"sync"
	"unsafe"

	"github.com/lestrrat-go/libxml2/xpath"
)

var docPool sync.Pool

func init() {
	SetupXPathCallback()
	docPool = sync.Pool{}
	docPool.New = func() interface{} {
		return &Document{}
	}
}

func SetupXPathCallback() {
	xpath.WrapNodeFunc = WrapNode
}

func WrapDocument(n unsafe.Pointer) *Document {
	//nolint:forcetypeassert
	doc := docPool.Get().(*Document)
	doc.mortal = false
	doc.ptr = n
	return doc
}
