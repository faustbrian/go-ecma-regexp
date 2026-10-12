package ecmascript

// Span is a half-open byte range in the original pattern source.
type Span struct {
	// Start is the inclusive source byte offset.
	Start int
	// End is the exclusive source byte offset; an empty range has equal endpoints.
	End int
}
