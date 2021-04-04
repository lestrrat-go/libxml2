package xsd

import (
	"net/url"
	"os"
	"path/filepath"

	"github.com/lestrrat-go/libxml2/internal/option"
)

// WithPath provides a hint to the XSD parser as to where the
// document being parsed is located at.
//
// This is useful when you must resolve relative paths inside a
// document, because to use relative paths the parser needs to
// know the reference location (i.e. location of the document
// being parsed). In case where you are parsing using `ParseFromFile()`
// this is handled automatically by the `ParseFromFile` method,
// but if you are using `Parse` method this is required
//
// If the path is provided as a relative path, the current directory
// should be obtainable via `os.Getwd` when this call is made, otherwise
// path resolution may fail in weird ways.
func WithPath(path string) Option {
	if !filepath.IsAbs(path) {
		if curdir, err := os.Getwd(); err == nil {
			path = filepath.Join(curdir, path)
		}
	}

	return WithURI(
		(&url.URL{
			Scheme: `file`,
			Path:   path,
		}).String(),
	)
}

func WithURI(v string) Option {
	return option.New(option.OptKeyWithURI, v)
}
