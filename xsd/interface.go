package xsd

import "errors"

// ErrInvalidSchema is returned when the Schema struct (probably
// the pointer to the underlying C struct is not valid)
var ErrInvalidSchema = errors.New("invalid schema")

// Schema represents an XML schema.
type Schema struct {
	ptr uintptr // *C.xmlSchema
}

// SchemaValidationError is returned when the Validate() function
// finds errors. When there are multiple errors, you may access
// them using the Errors() method
type SchemaValidationError struct {
	errors []error
}
