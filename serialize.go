package libxml2

import "errors"

/* Serialize produces serialization of the document, canonicalized.
 * Below document is taken from libxml2 directly. Pay special attention
 * to the required settings when parsing the document to be canonicalized.
 *
 * ---
 * Canonical form of an XML document could be created if and only if
 *  a) default attributes (if any) are added to all nodes
 *  b) all character and parsed entity references are resolved
 * In order to achive this in libxml2 the document MUST be loaded with
 * following global setings:
 *
 *    xmlLoadExtDtdDefaultValue = XML_DETECT_IDS | XML_COMPLETE_ATTRS;
 *    xmlSubstituteEntitiesDefault(1);
 *
 * or corresponding parser context setting:
 *    xmlParserCtxtPtr ctxt;
 *
 *    ...
 *    ctxt->loadsubset = XML_DETECT_IDS | XML_COMPLETE_ATTRS;
 *    ctxt->replaceEntities = 1;
 *    ...
 * ---
 *
 * In go-libxml2, this translates to:
 *
 *    options = ParseDtdLoad | ParseDtdAttrs | ParseNoEnt
 *
 */
func (s C14NSerialize) Serialize(n interface{}) (string, error) {
	switch n.(type) {
	case *Document:
	default:
		return "", errors.New("invalid node type: *Document is required")
	}

	return xmlC14NDocDumpMemory(n.(*Document), s.Mode, s.WithComments)
}