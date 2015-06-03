package libxml2

/*
#cgo pkg-config: libxml-2.0
#include "libxml/parser.h"
*/
import "C"
import (
	"bytes"
	"errors"
	"io"
)

func ParseString(s string) (*XmlDoc, error) {
	doc, err := C.xmlParseDoc(stringToXmlChar(s))
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, errors.New("parse failed")
	}
	return wrapXmlDoc(doc), nil
}

func Parse(in io.Reader) (*XmlDoc, error) {
	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(in); err != nil {
		return nil, err
	}

	return ParseString(buf.String())
}

