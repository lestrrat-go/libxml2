//go:generate go run internal/cmd/genwrapnode/genwrapnode.go -- dom/node_wrap.go

/*
Package libxml2 is an interface to libxml2 library, providing XML and HTML parsers
with DOM interface. The inspiration is Perl5's XML::LibXML module.

This library is still in very early stages of development. API may still change
without notice.

For the time being, the API is being written so that thye are as close as we
can get to DOM Layer 3, but some methods will, for the time being, be punted
and aliases for simpler methods that don't necessarily check for the DOM's
correctness will be used.

Also, the return values are still shaky -- I'm still debating how to handle error cases gracefully.
*/
package libxml2


