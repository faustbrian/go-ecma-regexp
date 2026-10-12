package ecmascript

// NodeKind identifies an immutable syntax-tree node.
type NodeKind uint8

const (
	// NodeEmpty represents an empty alternative or pattern.
	NodeEmpty NodeKind = iota + 1
	// NodeLiteral stores exact literal UTF-16 units.
	NodeLiteral
	// NodeDot represents dot matching under active flags.
	NodeDot
	// NodeStartAssertion represents the start assertion, including multiline mode.
	NodeStartAssertion
	// NodeEndAssertion represents the end assertion, including multiline mode.
	NodeEndAssertion
	// NodeWordBoundary represents a word-boundary assertion; Negated distinguishes its
	// complement.
	NodeWordBoundary
	// NodeCharacterClass represents stored range/builtin/property or v-mode set metadata.
	NodeCharacterClass
	// NodeBackreference refers to a numbered or resolved named capture.
	NodeBackreference
	// NodeLookaround represents a zero-width lookahead/lookbehind; Negated distinguishes
	// negative assertions.
	NodeLookaround
	// NodeConcatenation contains ordered children which match in sequence.
	NodeConcatenation
	// NodeAlternation contains ordered preferred alternatives.
	NodeAlternation
	// NodeGroup contains one capturing, noncapturing or scoped-modifier body.
	NodeGroup
	// NodeQuantifier contains one repeated child with Min, Max and Greedy metadata.
	NodeQuantifier
)

// Node is an immutable syntax-tree value. Slice accessors return copies.
type Node struct {
	kind         NodeKind
	span         Span
	text         string
	literalUnits []uint16
	children     []Node
	min          int
	max          int
	greedy       bool
	capturing    bool
	capture      int
	negated      bool
	class        []classTerm
	classStrings [][]uint16
	classOp      classOperation
	name         string
	backrefs     []int
	behind       bool
	enableFlags  uint16
	disableFlags uint16
}

type classOperation uint8

const (
	classOperationNone classOperation = iota
	classOperationUnion
	classOperationIntersection
	classOperationSubtraction
	classOperationComplement
)

// Kind reports the node variant; metadata relevance depends on this kind.
func (n Node) Kind() NodeKind { return n.kind }

// Span reports the half-open original pattern byte range.
func (n Node) Span() Span { return n.span }

// Text reports stored textual metadata, not a complete pattern serialization. Use
// Literal for exact literal units.
func (n Node) Text() string { return n.text }

// Literal returns copied exact UTF-16 literal units, including escaped lone surrogates.
func (n Node) Literal() UTF16String {
	return newUTF16String(n.literalUnits)
}

// Min reports the quantifier minimum; other node kinds leave it zero.
func (n Node) Min() int { return n.min }

// Max reports the quantifier maximum; -1 means unbounded. Other kinds leave it zero.
func (n Node) Max() int { return n.max }

// Greedy reports preferred repetition order on a quantifier.
func (n Node) Greedy() bool { return n.greedy }

// Capturing reports whether a group owns a numbered capture.
func (n Node) Capturing() bool {
	return n.capturing
}

// CaptureIndex reports a one-based group/reference index, or zero when none is
// assigned.
func (n Node) CaptureIndex() int { return n.capture }

// Negated reports complement metadata for classes or negative assertions.
func (n Node) Negated() bool { return n.negated }

// Name reports the decoded capture or backreference name, without changing source
// spans.
func (n Node) Name() string { return n.name }

// Lookbehind distinguishes backward from forward lookaround.
func (n Node) Lookbehind() bool { return n.behind }

// Children returns copied node values. Their nested storage remains private and
// immutable through public accessors.
func (n Node) Children() []Node {
	return append([]Node(nil), n.children...)
}

// CharacterRange is an inclusive code-point range in a character class.
type CharacterRange struct {
	// Start is the inclusive stored lower endpoint.
	Start rune
	// End is the inclusive stored upper endpoint; this is metadata, not an expanded
	// Unicode property table.
	End rune
}

// Ranges copies stored non-builtin range terms, not full class membership. Property
// terms are not expanded and can contribute zero-valued endpoints; v-mode child algebra
// is not flattened.
func (n Node) Ranges() []CharacterRange {
	ranges := make([]CharacterRange, 0, len(n.class))
	for _, term := range n.class {
		if term.builtin == classBuiltinNone {
			ranges = append(ranges, CharacterRange{Start: term.start, End: term.end})
		}
	}

	return ranges
}

// ClassStrings returns copied exact strings stored on this node; it does not flatten
// child set algebra.
func (n Node) ClassStrings() []UTF16String {
	values := make([]UTF16String, len(n.classStrings))
	for index, units := range n.classStrings {
		values[index] = newUTF16String(units)
	}
	return values
}

type classBuiltin uint8

const (
	classBuiltinNone classBuiltin = iota
	classBuiltinDigit
	classBuiltinSpace
	classBuiltinWord
)

type classTerm struct {
	start    rune
	end      rune
	builtin  classBuiltin
	negated  bool
	property uint16
}

// Pattern is an immutable parsed pattern.
type Pattern struct {
	source       string
	edition      Edition
	flags        Flags
	root         Node
	captureCount int
	captureNames map[string][]int
}

// Source returns retained original pattern text, which can contain sensitive data.
func (p *Pattern) Source() string { return p.source }

// Edition reports the selected closed language edition.
func (p *Pattern) Edition() Edition { return p.edition }

// Flags reports the immutable grammar flags used for parsing.
func (p *Pattern) Flags() Flags { return p.flags }

// Root returns an immutable node value; public accessors do not expose mutable tree
// storage.
func (p *Pattern) Root() Node { return p.root }

// CaptureCount excludes capture zero, the full match.
func (p *Pattern) CaptureCount() int { return p.captureCount }

// CaptureNames returns a new map selecting the first declared index for each name. Use
// CaptureNameIndices to retain disjoint duplicate groups.
func (p *Pattern) CaptureNames() map[string]int {
	names := make(map[string]int, len(p.captureNames))
	for name, indices := range p.captureNames {
		// The parser records a name only when it appends its capture index.
		names[name] = indices[0]
	}
	return names
}

// CaptureNameIndices returns a new map and copied one-based index slices for every
// name.
func (p *Pattern) CaptureNameIndices() map[string][]int {
	return cloneCaptureNames(p.captureNames)
}
