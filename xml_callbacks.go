package libxml2

import (
	"strings"

	"github.com/lestrrat-go/libxml2/clib"
)

type XMLCallback struct {
	uri         string
	matchingFnc func(receivedURI, expectedURI string) bool
	data        []byte
}

func (i XMLCallback) CanHandle(uri string) bool {
	return i.matchingFnc(uri, i.uri)
}

func (i XMLCallback) GetData(_ string) []byte {
	return i.data
}

var _ clib.XMLCallback = (*XMLCallback)(nil)

type Option func(*XMLCallback)

func WithURIEquals() Option {
	return func(callback *XMLCallback) {
		callback.matchingFnc = func(receivedURI, expectedURI string) bool {
			return receivedURI == expectedURI
		}
	}
}

func WithURIContains() Option {
	return func(callback *XMLCallback) {
		callback.matchingFnc = strings.Contains
	}
}

func defaultOptions() []Option {
	return []Option{
		WithURIEquals(),
	}
}

func InMemoryCallback(xmlURI string, data []byte, options ...Option) clib.XMLCallback {
	callback := &XMLCallback{
		uri:  xmlURI,
		data: data,
	}
	for _, opt := range append(defaultOptions(), options...) {
		opt(callback)
	}
	return callback
}
