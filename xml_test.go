package libxml2

import (
	"fmt"
	"os"
	"testing"
)

func TestEncoding(t *testing.T) {
	for _, enc := range []string{`utf-8`, `sjis`, `euc-jp`} {
		fn := fmt.Sprintf(`test/%s.xml`, enc)
		f, err := os.Open(fn)
		if err != nil {
			t.Errorf("Failed to open %s: %s", fn, err)
			return
		}
		defer f.Close()

		p := NewParser()
		doc, err := p.ParseReader(f)
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
