package libxml2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestC14N(t *testing.T) {
	expected := `<?xml-stylesheet href="doc.xsl"
   type="text/xsl"   ?>
<doc>Hello, world!<!-- Comment 1 --></doc>
<?pi-without-data?>
<!-- Comment 2 -->
<!-- Comment 3 -->`

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

	s, err := C14NSerialize{Mode: C14NExclusive1_0, WithComments: true}.Serialize(doc)
	if !assert.NoError(t, err, "C14N should succeed") {
		return
	}
	t.Logf("C14N -> %s", s)
	t.Logf("expected -> %s", expected)

	if !assert.Equal(t, expected, s, "C14N content matches") {
		return
	}
}
