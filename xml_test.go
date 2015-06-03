package libxml2

import "testing"

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
