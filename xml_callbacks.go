package libxml2

import "C"
import (
	"bytes"
	"github.com/lestrrat-go/libxml2/clib"
	"io"
	"io/ioutil"
)

type inMemoryCallback struct {
	uri         string
	data        []byte
}

func (i *inMemoryCallback) Open(_ string) (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewReader(i.data)), nil
}

func (i *inMemoryCallback) Match(uri string) bool {
	return i.uri == uri
}

var _ clib.Callback = (*inMemoryCallback)(nil)

func InMemoryCallback(xmlURI string, data []byte) clib.Callback {
	return &inMemoryCallback{
		uri:  xmlURI,
		data: data,
	}
}

func RegisterInputCallback(callbacks ...clib.Callback) {
	if len(callbacks) == 0 {
		return
	}
	clib.RegisterInputCallback(callbacks...)
}

func RestoreDefaultInputCallback() {
	clib.RestoreDefaultInputCallback()
}