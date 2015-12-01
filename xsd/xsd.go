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

/*
#cgo pkg-config: libxml-2.0
#include <libxml/xmlschemas.h>
#include <string.h>
#include <stdio.h>

#define GO_LIBXML2_ERRWARN_ACCUMULATOR_SIZE 32
typedef struct go_libxml2_errwarn_accumulator {
	char *errors[GO_LIBXML2_ERRWARN_ACCUMULATOR_SIZE];
	char *warnings[GO_LIBXML2_ERRWARN_ACCUMULATOR_SIZE];
	int erridx;
	int warnidx;
} go_libxml2_errwarn_accumulator;

static
go_libxml2_errwarn_accumulator*
MY_createErrWarnAccumulator() {
	int i;
	go_libxml2_errwarn_accumulator *ctx;
	ctx = (go_libxml2_errwarn_accumulator *) malloc(sizeof(go_libxml2_errwarn_accumulator));
	for (i = 0; i < GO_LIBXML2_ERRWARN_ACCUMULATOR_SIZE; i++) {
		ctx->errors[i] = NULL;
		ctx->warnings[i] = NULL;
	}
	ctx->erridx = 0;
	ctx->warnidx = 0;
	return ctx;
}

static
void
MY_freeErrWarnAccumulator(go_libxml2_errwarn_accumulator* ctx) {
	int i = 0;
	for (i = 0; i < GO_LIBXML2_ERRWARN_ACCUMULATOR_SIZE; i++) {
		if (ctx->errors[i] != NULL) {
			free(ctx->errors[i]);
		}
		if (ctx->warnings[i] != NULL) {
			free(ctx->errors[i]);
		}
	}
	free(ctx);
}

static
void
MY_accumulateErr(void *ctx, const char *msg, ...) {
  char buf[1024];
  va_list args;
	go_libxml2_errwarn_accumulator *accum;
	int len;

	accum = (go_libxml2_errwarn_accumulator *) ctx;
	if (accum->erridx >= GO_LIBXML2_ERRWARN_ACCUMULATOR_SIZE) {
		return;
	}

  va_start(args, msg);
  len = vsnprintf(buf, sizeof(buf), msg, args);
  va_end(args);

	if (len == 0) {
		return;
	}

	char *out = (char *) calloc(sizeof(char), len);
	if (buf[len-1] == '\n') {
		// don't want newlines in my error values
		buf[len-1] = '\0';
		len--;
	}
	memcpy(out, buf, len);

  int i = accum->erridx++;
	accum->errors[i] = out;
}

static
void
MY_setErrWarnAccumulator(xmlSchemaValidCtxtPtr ctxt, go_libxml2_errwarn_accumulator *accum) {
	xmlSchemaSetValidErrors(ctxt, MY_accumulateErr, NULL, accum);
}

*/
import "C"
import (
	"errors"
	"unsafe"

	"github.com/lestrrat/go-libxml2/node"
)

// Parse is used to parse an XML Schema Document to produce a
// Schema instance. Make sure to call Free() on the instance
// when you are done with it.
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

	return &Schema{ptr: uintptr(unsafe.Pointer(s))}, nil
}

func validSchemaPtr(s *Schema) (*C.xmlSchema, error) {
	if s == nil || s.ptr == 0 {
		return nil, ErrInvalidSchema
	}

	return (*C.xmlSchema)(unsafe.Pointer(s.ptr)), nil
}

// Free frees the underlying C struct
func (s *Schema) Free() {
	sptr, err := validSchemaPtr(s)
	if err != nil {
		return
	}

	C.xmlSchemaFree(sptr)
	s.ptr = 0
}

// Error method fulfils the error interface
func (sve SchemaValidationError) Error() string {
	return "schema validation failed"
}

// Errors returns the list of errors found
func (sve SchemaValidationError) Errors() []error {
	return sve.errors
}

// Validate takes in a XML document and validates it against
// the schema. If there are any problems, and error is
// returned.
func (s *Schema) Validate(d node.Document) error {
	sptr, err := validSchemaPtr(s)
	if err != nil {
		return err
	}

	ctx := C.xmlSchemaNewValidCtxt(sptr)
	if ctx == nil {
		return errors.New("failed to build validator")
	}
	defer C.xmlSchemaFreeValidCtxt(ctx)

	accum := C.MY_createErrWarnAccumulator()
	defer C.MY_freeErrWarnAccumulator(accum)

	C.MY_setErrWarnAccumulator(ctx, accum)

	if C.xmlSchemaValidateDoc(ctx, (C.xmlDocPtr)(unsafe.Pointer(d.Pointer()))) != 0 {
		// Create an err
		err := SchemaValidationError{
			errors: make([]error, 0, accum.erridx),
		}
		for i := 0; i < int(accum.erridx); i++ {
			err.errors = append(err.errors, errors.New(C.GoString(accum.errors[i])))
		}
		return err
	}
	return nil
}
