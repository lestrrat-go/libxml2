package libxml2

import (
	"fmt"
	"os"
	"testing"

	"github.com/lestrrat-go/libxml2/dom"
	"github.com/lestrrat-go/libxml2/parser"
	"github.com/lestrrat-go/libxml2/types"
	"github.com/lestrrat-go/libxml2/xpath"
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

		p := parser.New()
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
	d := dom.CreateDocument()
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
	root.AddChild(n)

	_, err = n.GetAttribute("xmlns")
	if !assert.Error(t, err, "GetAttribute should fail with not found") ||
		!assert.Equal(t, "attribute not found", err.Error(), "error matches") {
		return
	}

	var c types.Element
	for _, name := range []string{"a", "b", "c"} {
		child, err := d.CreateElementNS("http://children", "child:"+name)
		if !assert.NoError(t, err, "CreateElementNS should succeed") {
			return
		}
		if name == "c" {
			c = child
		}
		n.AddChild(child)
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
	n.AddChild(child)

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

func TestRegressionGH7(t *testing.T) {
	doc, err := ParseHTMLString(`<!DOCTYPE html>
<html>
<body>
<div>
<style>
</style>
    1234
</div>
</body>
</html>`)

	if !assert.NoError(t, err, "ParseHTMLString should succeed") {
		return
	}

	nodes := xpath.NodeList(doc.Find(`./body/div`))
	if !assert.NotEmpty(t, nodes, "Find should succeed") {
		return
	}

	v, err := nodes.Literal()
	if !assert.NoError(t, err, "Literal() should succeed") {
		return
	}
	if !assert.NotEmpty(t, v, "Literal() should return some string") {
		return
	}
	t.Logf("v = '%s'", v)
}

func TestGHIssue43(t *testing.T) {
	d := dom.CreateDocument()
	r, _ := d.CreateElement("root")
	r.SetNamespace("http://some.uri", "pfx", true)
	d.SetDocumentElement(r)
	e, _ := d.CreateElement("elem")
	e.SetNamespace("http://other.uri", "", true)
	r.AddChild(e)
	s := d.ToString(1, true)

	if !assert.Contains(t, s, `<elem xmlns="http://other.uri"`, `default namespace works`) {
		return
	}
}
