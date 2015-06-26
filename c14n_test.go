package libxml2

import "testing"

func TestC14N(t *testing.T) {
	p := &Parser{}
	doc, err := p.ParseString(`<?xml version="1.0"?>
<Root>
	<EmptyElement foo="bar"/>
</Root>`)

	if err != nil {
		t.Errorf("Failed to parse document: %s", err)
		return
	}

	s, err := doc.ToStringC14N(true)
	if err != nil {
		t.Errorf("Failed to format in C14N: %s", err)
		return
	}
	t.Logf("%s", s)
}