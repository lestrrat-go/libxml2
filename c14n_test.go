package libxml2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestC14N(t *testing.T) {
	expected := `<?xml-stylesheet href="doc.xsl"
   type="text/xsl"   ?>
<doc>Hello, world!</doc>
<?pi-without-data?>`

	doc, err := ParseString(`<?xml version="1.0"?>
<?xml-stylesheet   href="doc.xsl"
   type="text/xsl"   ?>

<doc>Hello, world!<!-- Comment 1 --></doc>

<?pi-without-data     ?>

<!-- Comment 2 -->

<!-- Comment 3 -->


`)

	if !assert.NoError(t, err, "Parse document should succeed") {
		return
	}

	s, err := doc.C14N(true)
	if !assert.NoError(t, err, "C14N should succeed") {
		return
	}
	t.Logf("C14N -> %s", s)
	t.Logf("expected -> %s", expected)

	if !assert.Equal(t, expected, s, "C14N content matches") {
		return
	}
}

func TestC14NNonExclusive(t *testing.T) {
	p := NewParser()
	doc, err := p.ParseString(`<?xml version="1.0"?>
<Root xmlns="uri:go-libxml2:test" xmlns:test2="uri:go-libxml2:test2">
	<EmptyElement foo="bar"/>
</Root>`)

	if err != nil {
		t.Errorf("Failed to parse document: %s", err)
		return
	}

	s, err := doc.C14N(false)
	if err != nil {
		t.Errorf("Failed to format in C14N: %s", err)
		return
	}
	t.Logf("%s", s)
}
