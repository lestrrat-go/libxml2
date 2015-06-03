package libxml2

/*
#include <stdlib.h>
#include "libxml/xmlstring.h"

static inline xmlChar* to_xmlcharptr(const char *s) {
	return (xmlChar *) s;
}
static inline char * to_charptr(const xmlChar *s) {
	return (char *) s;
}

*/
import "C"

var emptyStringBytes = []byte{0}

func toCString(s string) *C.char {
	return C.CString(s)
}

func xmlCharToString(s *C.xmlChar) string {
	return C.GoString(C.to_charptr(s))
}

func stringToXmlChar(s string) *C.xmlChar {
	return C.to_xmlcharptr(C.CString(s))
}

func AppendCStringTerminator(b []byte) []byte {
	if n := len(b); n > 0 {
		if b[n-1] != 0 {
			return append(b, 0)
		}
	}
	return b
}

func GetCString(b []byte) []byte {
	b = AppendCStringTerminator(b)
	if len(b) == 0 {
		return emptyStringBytes
	}
	return b
}