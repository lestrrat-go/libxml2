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
	go_libxml2_errwarn_accumulator *ctx;
	ctx = (go_libxml2_errwarn_accumulator *) malloc(sizeof(go_libxml2_errwarn_accumulator));
	for (int i = 0; i < GO_LIBXML2_ERRWARN_ACCUMULATOR_SIZE; i++) {
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
	for (int i = 0; i < GO_LIBXML2_ERRWARN_ACCUMULATOR_SIZE; i++) {
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

	char *out = (char *) malloc(sizeof(char) * len);
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

func (s *Schema) Free() {
	if ptr := s.ptr; ptr != nil {
		C.xmlSchemaFree(ptr)
	}
}

type SchemaValidationError struct {
	errors []error
}

func (sve SchemaValidationError) Error() string {
	return "schema validation failed"
}

func (sve SchemaValidationError) Errors() []error {
	return sve.errors
}

func (s *Schema) Validate(d *libxml2.Document) error {
	ctx := C.xmlSchemaNewValidCtxt(s.ptr)
	if ctx == nil {
		return errors.New("failed to build validator")
	}
	defer C.xmlSchemaFreeValidCtxt(ctx)

	accum := C.MY_createErrWarnAccumulator()
	defer C.MY_freeErrWarnAccumulator(accum)

	C.MY_setErrWarnAccumulator(ctx, accum)

	if C.xmlSchemaValidateDoc(ctx, (C.xmlDocPtr)(d.Pointer())) != 0 {
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
