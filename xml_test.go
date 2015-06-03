package libxml2

import (
	"fmt"
	"os"
	"testing"
)

// TODO: parse real XML
func TestParse(t *testing.T) {
	doc, err := ParseString(`<html><body><h1>Hello, World!</h1><p>Lorem Ipsum</p></body></html>`)
	if err != nil {
		t.Errorf("Failed to parse: %s", err)
		return
	}
	defer doc.Free()

	nodes, err := doc.FindNodes("/html/body/h1")
	if err != nil {
		t.Errorf("Failed to evaluate xpath: %s", err)
		return
	}
	if len(nodes) != 1 {
		t.Errorf("Could not find matching nodes")
		return
	}

	if nodes[0].TextContent() != "Hello, World!" {
		t.Errorf("h1 content is not 'Hello, World!', got %s", nodes[0].TextContent())
		return
	}
}

func TestEncoding(t *testing.T) {
	for _, enc := range []string{`utf-8`, `sjis`, `euc-jp`} {
		fn := fmt.Sprintf(`test/%s.xml`, enc)
		f, err := os.Open(fn)
		if err != nil {
			t.Errorf("Failed to open %s: %s", fn, err)
			return
		}
		defer f.Close()

		doc, err := Parse(f)
		if err != nil {
			t.Errorf("Failed to parse %s: %s", fn, err)
			return
		}

		if doc.Encoding() != enc {
			t.Errorf("Expected encoding %s, got %s", enc, doc.Encoding())
			return
		}
	}
}
