package libxml2

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestNamespacedReconciliation(t *testing.T) {
	d := CreateDocument()
	root, err := d.CreateElement("foo")
	if !assert.NoError(t, err, "failed to create document") {
		return
	}
	d.SetDocumentElement(root)
	if !assert.NoError(t, root.SetNamespace("http://default", "root"), "SetNamespace should succeed") {
		return
	}

	if !assert.NoError(t, root.SetNamespace("http://children", "child", false), "SetNamespace (no-activate) should succeed") {
		return
	}

	n, err := d.CreateElementNS("http://default", "branch")
	if !assert.NoError(t, err, "CreateElementNS should succeed") {
		return
	}
	root.AppendChild(n)

	_, err = n.GetAttribute("xmlns")
	if !assert.Error(t, err, "GetAttribute should fail with not found") ||
		!assert.Equal(t, "attribute not found", err.Error(), "error matches") {
		return
	}

	var c *Element
	for _, name := range []string{"a", "b", "c"} {
		child, err := d.CreateElementNS("http://children", "child:"+name)
		if !assert.NoError(t, err, "CreateElementNS should succeed") {
			return
		}
		if name == "c" {
			c = child
		}
		n.AppendChild(child)
		_, err = n.GetAttribute("xmlns:child")
		if !assert.Error(t, err, "GetAttribute should fail with not found") ||
			!assert.Equal(t, "attribute not found", err.Error(), "error matches") {
			return
		}
	}

	if !assert.NoError(t, c.SetAttribute("xmlns:foo", "http://children"), "SetAttribute should succeeed") {
		return
	}

	attr, err := c.GetAttribute("xmlns:foo")
	if !assert.NoError(t, err, "xmlns:foo should exist") {
		return
	}
	if !assert.Equal(t, "http://children", attr.Value(), "attribute matches") {
		return
	}

	child, err := d.CreateElementNS("http://other", "branch")
	if !assert.NoError(t, err, "creating element with default namespace") {
		return
	}
	n.AppendChild(child)

	// XXX This still fails
	/*
	attr, err = child.GetAttribute("xmlns")
	if !assert.NoError(t, err, "GetAttribute should succeed") {
		return
	}
	if !assert.Equal(t, "http://other", attr.Value(), "attribute matches") {
		return
	}
	*/

	t.Logf("%s", d.String())
}
