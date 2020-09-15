// Package xsd contains some of the tools available from libxml2
// that allows you to validate your XML against an XSD
//
// This is basically all you need to do:
//
//    schema, err := xsd.Parse(xsdsrc)
//    if err != nil {
//        panic(err)
//    }
//    defer schema.Free()
//    if err := schema.Validate(doc); err != nil{
//        for _, e := range err.(SchemaValidationErr).Error() {
//             println(e.Error())
//        }
//    }
//
package xsd

import (
	"github.com/lestrrat-go/libxml2/clib"
	"github.com/lestrrat-go/libxml2/types"
	"github.com/pkg/errors"
)

const ValueVCCreate = 1

// Parse is used to parse an XML Schema Document to produce a
// Schema instance. Make sure to call Free() on the instance
// when you are done with it.


func Parse(buf []byte, options ...Option) (*Schema, error) {
	// xsd.WithURI(...)
	sptr, err := clib.XMLSchemaParse(buf, options...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse input")
	}

	return &Schema{ptr: sptr}, nil
}

// ParseFromFile is used to parse an XML schema using only the file path.
// Make sure to call Free() on the instance when you are done with it.
func ParseFromFile(path string) (*Schema, error) {
	sptr, err := clib.XMLSchemaParseFromFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse input from file")
	}

	return &Schema{ptr: sptr}, nil
}

// Pointer returns the underlying C struct
func (s *Schema) Pointer() uintptr {
	return s.ptr
}

// Free frees the underlying C struct
func (s *Schema) Free() {
	if err := clib.XMLSchemaFree(s); err != nil {
		return
	}
	s.ptr = 0
}

// Validate takes in a XML document and validates it against
// the schema. If there are any problems, and error is
// returned.
func (s *Schema) Validate(d types.Document, options ...int) error {
	errs := clib.XMLSchemaValidateDocument(s, d, options...)
	if errs == nil {
		return nil
	}

	return SchemaValidationError{errors: errs}
}

// Error method fulfils the error interface
func (sve SchemaValidationError) Error() string {
	return "schema validation failed"
}

// Errors returns the list of errors found
func (sve SchemaValidationError) Errors() []error {
	return sve.errors
}
