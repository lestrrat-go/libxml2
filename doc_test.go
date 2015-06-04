package libxml2

import "testing"

// Tests for DOM Layer 3

func TestDocumentAttributes(t *testing.T) {
	doc := CreateDocument()
	if doc.Encoding() != "" {
		t.Errorf("Encoding should be empty string at first, got '%s'", doc.Encoding())
	}

	if doc.Version() != "1.0" {
		t.Errorf("Version should be 1.0 by default, got '%s'", doc.Version())
	}

	for _, enc := range []string{"utf-8", "euc-jp", "sjis", "iso-8859-1"} {
		doc.SetEncoding(enc)
		if doc.Encoding() != enc {
			t.Errorf("Expected encoding '%s', got '%s'", enc, doc.Encoding())
		}
	}

	for _, v := range []string{"1.5", "4.12", "12.5"} {
		doc.SetVersion(v)
		if doc.Version() != v {
			t.Errorf("Expected version '%s', got '%s'", v, doc.Version())
		}
	}
}