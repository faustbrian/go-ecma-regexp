package ecmascript

import (
	"fmt"
	"time"
)

// SyntaxCode classifies a pattern error without requiring message matching.
type SyntaxCode uint8

const (
	// SyntaxUnexpectedToken reports a token invalid in its grammar context.
	SyntaxUnexpectedToken SyntaxCode = iota + 1
	// SyntaxUnexpectedEOF reports prematurely ended pattern syntax.
	SyntaxUnexpectedEOF
	// SyntaxUnclosedGroup reports a group missing its closing delimiter.
	SyntaxUnclosedGroup
	// SyntaxInvalidQuantifier reports an invalid repetition declaration.
	SyntaxInvalidQuantifier
	// SyntaxInvalidEscape reports an invalid mode-specific escape.
	SyntaxInvalidEscape
	// SyntaxInvalidBackreference reports an unresolved or disallowed capture reference.
	SyntaxInvalidBackreference
	// SyntaxUnsupported reports syntax or an edition outside the selected engine support.
	SyntaxUnsupported
)

// SyntaxError reports a source-byte diagnostic. Library parser messages are
// bounded; the exported fields do not sanitize caller-constructed diagnostics.
type SyntaxError struct {
	// Code classifies syntax failure without message matching.
	Code SyntaxCode
	// Span locates the pattern failure in source bytes.
	Span Span
	// Message contains diagnostic text; library parser messages are capped at 160 bytes,
	// but caller-built errors are not sanitized.
	Message string
}

// Error formats code-independent diagnostic text and source span; use Code for
// classification.
func (e *SyntaxError) Error() string {
	return fmt.Sprintf("invalid pattern at bytes %d..%d: %s", e.Span.Start, e.Span.End, e.Message)
}

// LimitKind identifies a finite parser or execution resource.
type LimitKind uint8

const (
	// LimitPatternBytes counts pattern source bytes.
	LimitPatternBytes LimitKind = iota + 1
	// LimitASTDepth counts nested group or Unicode Sets depth.
	LimitASTDepth
	// LimitCaptures counts numbered capturing groups.
	LimitCaptures
	// LimitCharacterClasses counts charged bracket classes and property escapes.
	LimitCharacterClasses
	// LimitASTNodes counts charged parser nodes.
	LimitASTNodes
	// LimitProgramInstructions counts emitted instructions.
	LimitProgramInstructions
	// LimitInputBytes counts admitted subject or replacement storage bytes.
	LimitInputBytes
	// LimitInputRunes counts admitted subject or replacement code points.
	LimitInputRunes
	// LimitMatchSteps counts charged dispatch/substitution/name work.
	LimitMatchSteps
	// LimitBacktracks counts resumed alternative frames.
	LimitBacktracks
	// LimitStackDepth counts pending alternatives in one VM invocation.
	LimitStackDepth
	// LimitRecursionDepth counts nested lookaround invocations.
	LimitRecursionDepth
	// LimitAllocations counts selected logical VM/collection reservations, not all Go
	// allocations.
	LimitAllocations
	// LimitResults counts collected matches, replacements or split values.
	LimitResults
	// LimitOutputUTF16 counts replacement or split output code units.
	LimitOutputUTF16
)

// LimitError distinguishes resource exhaustion from invalid syntax and a
// successful no-match result.
type LimitError struct {
	// Kind identifies the exhausted allowance.
	Kind LimitKind
	// Limit is the configured ceiling in that resource kind.
	Limit uint64
	// Used is the charged or observed count causing refusal, not a complete heap/CPU
	// measurement.
	Used uint64
}

// Error formats the kind and charged counts; callers should use typed fields, not
// message matching.
func (e *LimitError) Error() string {
	return fmt.Sprintf("regular expression resource limit exceeded: kind=%d limit=%d used=%d", e.Kind, e.Limit, e.Used)
}

// TimeoutError reports synchronous wall-time budget exhaustion. The matcher
// never creates a worker goroutine to enforce time.
type TimeoutError struct {
	// Limit is the configured execution wall-time allowance.
	Limit time.Duration
	// Elapsed is the sampled elapsed duration, not a guaranteed whole-operation deadline.
	Elapsed time.Duration
}

// Error formats sampled execution timing; callers should inspect the typed error.
func (e *TimeoutError) Error() string {
	return fmt.Sprintf("regular expression wall-time limit exceeded: limit=%s elapsed=%s", e.Limit, e.Elapsed)
}
