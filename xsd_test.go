package libxml2_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/xsd"
	"github.com/stretchr/testify/assert"
)

func TestXSD(t *testing.T) {
	xsdfile := filepath.Join("test", "xmldsig-core-schema.xsd")
	f, err := os.Open(xsdfile)
	if !assert.NoError(t, err, "open schema") {
		return
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if !assert.NoError(t, err, "reading from schema") {
		return
	}

	s, err := xsd.Parse(buf)
	if !assert.NoError(t, err, "parsing schema") {
		return
	}
	defer s.Free()

	func() {
		d, err := libxml2.ParseString(`<foo></foo>`)
		if !assert.NoError(t, err, "parsing XML") {
			return
		}
		defer d.Free()

		err = s.Validate(d)
		if !assert.Error(t, err, "s.Validate should fail") {
			return
		}

		serr, ok := err.(xsd.SchemaValidationError)
		if !assert.True(t, ok, "error is xsd.SchemaValidationErr") {
			return
		}

		if !assert.Len(t, serr.Errors(), 1, "there's one error") {
			return
		}
		for _, e := range serr.Errors() {
			t.Logf("err (OK): '%s'", e)
		}
	}()

	func() {
		const src = `<?xml version="1.0" encoding="UTF-8"?>
  <Signature xmlns="http://www.w3.org/2000/09/xmldsig#">
    <SignedInfo>
      <CanonicalizationMethod 
           Algorithm="http://www.w3.org/TR/2001/REC-xml-c14n-
20010315#WithComments"/>
      <SignatureMethod Algorithm="http://www.w3.org/2000/09/
xmldsig#dsa-sha1"/>
      <Reference URI="">
        <Transforms>
          <Transform Algorithm="http://www.w3.org/2000/09/
xmldsig#enveloped-signature"/>
        </Transforms>
        <DigestMethod Algorithm="http://www.w3.org/2000/09/
xmldsig#sha1"/>
        <DigestValue>uooqbWYa5VCqcJCbuymBKqm17vY=</DigestValue>
      </Reference>
    </SignedInfo>
<SignatureValue>
KedJuTob5gtvYx9qM3k3gm7kbLBwVbEQRl26S2tmXjqNND7MRGtoew==
    </SignatureValue>
    <KeyInfo>
      <KeyValue>
        <DSAKeyValue>
          <P>
/KaCzo4Syrom78z3EQ5SbbB4sF7ey80etKII864WF64B81uRpH5t9jQTxe
Eu0ImbzRMqzVDZkVG9xD7nN1kuFw==
          </P>
          <Q>li7dzDacuo67Jg7mtqEm2TRuOMU=</Q>
          <G>Z4Rxsnqc9E7pGknFFH2xqaryRPBaQ01khpMdLRQnG541Awtx/
XPaF5Bpsy4pNWMOHCBiNU0NogpsQW5QvnlMpA==
          </G>
          <Y>qV38IqrWJG0V/
mZQvRVi1OHw9Zj84nDC4jO8P0axi1gb6d+475yhMjSc/
BrIVC58W3ydbkK+Ri4OKbaRZlYeRA==
         </Y>
        </DSAKeyValue>
      </KeyValue>
    </KeyInfo>
  </Signature>
`
		d, err := libxml2.ParseString(src)
		if !assert.NoError(t, err, "parsing XML") {
			return
		}
		defer d.Free()

		err = s.Validate(d)
		if !assert.NoError(t, err, "s.Validate should pass") {
			if serr, ok := err.(xsd.SchemaValidationError); ok {
				for _, e := range serr.Errors() {
					t.Logf("err: %s", e)
				}
			}
			return
		}
	}()
}

func TestXSDDefaultValue(t *testing.T) {
	const schemasrc = `<?xml version="1.0" encoding="UTF-8"?>
<xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema"
elementFormDefault="qualified">
  <xs:element name="config">
    <xs:complexType mixed="true">
      <xs:sequence>
        <xs:element ref="attribute"/>
      </xs:sequence>
    </xs:complexType>
  </xs:element>
  <xs:element name="attribute">
    <xs:complexType>
      <xs:sequence>
        <xs:element name="linguistic">
          <xs:complexType>
            <xs:attribute name="item" default="US"/>
          </xs:complexType>
        </xs:element>
      </xs:sequence>
    </xs:complexType>
  </xs:element>
</xs:schema>`
	const docsrc = `<config>
    <attribute>
        <linguistic></linguistic>
    </attribute>
</config>`

	schema, err := xsd.Parse([]byte(schemasrc))
	if !assert.NoError(t, err, `xsd.Parse should succeed`) {
		return
	}
	defer schema.Free()

	doc, err := libxml2.ParseString(docsrc)
	if !assert.NoError(t, err, "parsing XML") {
		return
	}
	defer doc.Free()
	if !assert.NoError(t, schema.Validate(doc, xsd.ValueVCCreate), `schema.Validate should succeed`) {
		return
	}

	t.Logf("%s", doc.String())

}

func TestGHIssue67(t *testing.T) {
	t.Run("Local validation", func(t *testing.T) {
		const schemafile = "test/schema/projects/go_libxml2_local.xsd"
		const docfile = "test/go_libxml2_local.xml"

		schemasrc, err := ioutil.ReadFile(schemafile)
		if !assert.NoError(t, err, `failed to read xsd file`) {
			return
		}

		docsrc, err := ioutil.ReadFile(docfile)
		if !assert.NoError(t, err, `failed to read xml file`) {
			return
		}

		schema, err := xsd.Parse(schemasrc, xsd.WithPath(schemafile))
		if !assert.NoError(t, err, `xsd.Parse should succeed`) {
			return
		}
		defer schema.Free()

		doc, err := libxml2.Parse(docsrc)
		if !assert.NoError(t, err, "parsing XML") {
			return
		}
		defer doc.Free()
		if !assert.NoError(t, schema.Validate(doc, xsd.ValueVCCreate), `schema.Validate should succeed`) {
			return
		}

		t.Logf("%s", doc.String())
	})
	t.Run("Remote validation", func(t *testing.T) {
		curdir, err := os.Getwd()
		if !assert.NoError(t, err, `os.Getwd failed`) {
			return
		}

		srv := httptest.NewServer(http.FileServer(http.Dir(curdir)))
		defer srv.Close()

		var schemafile = srv.URL + "/test/schema/projects/go_libxml2_remote.xsd"
		const docfile = "test/go_libxml2_remote.xml"

		res, err := http.Get(schemafile)
		if !assert.NoError(t, err, `failed to fetch xsd file`) {
			return
		}

		schemasrc, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		if !assert.NoError(t, err, `failed to read xsd file`) {
			return
		}

		docsrc, err := ioutil.ReadFile(docfile)
		if !assert.NoError(t, err, `failed to read xml file`) {
			return
		}

		schema, err := xsd.Parse(schemasrc, xsd.WithURI(schemafile))
		if !assert.NoError(t, err, `xsd.Parse should succeed`) {
			return
		}
		defer schema.Free()

		doc, err := libxml2.Parse(docsrc)
		if !assert.NoError(t, err, "parsing XML") {
			return
		}
		defer doc.Free()
		if !assert.NoError(t, schema.Validate(doc, xsd.ValueVCCreate), `schema.Validate should succeed`) {
			return
		}

		t.Logf("%s", doc.String())
	})
}
