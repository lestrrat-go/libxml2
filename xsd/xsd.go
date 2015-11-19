package xsd

/*
#cgo pkg-config: libxml-2.0
#include <libxml/xmlschemas.h>
*/
import "C"
import (
	"errors"
	"unsafe"

	"github.com/lestrrat/go-libxml2"
)

type Schema struct {
	ptr C.xmlSchemaPtr
}

func Parse(buf []byte) (*Schema, error) {
	parserCtx := C.xmlSchemaNewMemParserCtxt(
		(*C.char)(unsafe.Pointer(&buf[0])),
		C.int(len(buf)),
	)
	if parserCtx == nil {
		return nil, errors.New("failed to create parser")
	}
	defer C.xmlSchemaFreeParserCtxt(parserCtx)

	s := C.xmlSchemaParse(parserCtx)
	if s == nil {
		return nil, errors.New("failed to parse schema")
	}

	return &Schema{ptr: s}, nil
}

func (s *Schema) Close() {
	if ptr := s.ptr; ptr != nil {
		C.xmlSchemaFree(ptr)
	}
}

func (s *Schema) Validate(d *libxml2.Document) error {
	ctx := C.xmlSchemaNewValidCtxt(s.ptr)
	if ctx == nil {
		return errors.New("failed to build validator")
	}
	defer C.xmlSchemaFreeValidCtxt(ctx)

	if C.xmlSchemaValidateDoc(ctx, (C.xmlDocPtr)(d.Pointer())) != 0 {
		return errors.New("failed to validate document")
	}
	return nil
}
