// Package ecmascript implements a bounded ECMA-262 regular expression engine.
//
// The package does not delegate matching to Go's regexp package. Its accepted
// syntax and behavior are closed over an explicitly selected ECMAScript
// edition. The supported profile is ECMA-262 16th edition (Edition2025),
// with generated Unicode 16.0.0 data reported by UnicodeVersion. Grammar flags
// select legacy, u or v modes; Annex B is an explicit legacy compatibility policy.
// This is not a JavaScript host-object/coercion implementation or a PCRE dialect.
//
// Parse/Compile use explicit structural budgets without a context API. Execution
// is synchronous and polls context/time cooperatively, not at every comparison
// or memory copy. Backtracking has no RE2 linear-time guarantee. Programs and
// results are immutable through public accessors; mutable Session state and any
// application cache require caller synchronization. No ambient I/O or workers
// are started by the public library.
package ecmascript

import (
	"errors"
	"fmt"
)

// Edition identifies a closed ECMA-262 language edition.
type Edition uint16

const (
	// Edition2025 is ECMA-262, 16th edition.
	Edition2025 Edition = 2025
)

// ErrUnsupportedEdition identifies an edition name outside the closed profile.
var ErrUnsupportedEdition = errors.New("unsupported ECMAScript edition")

// ParseEdition parses a public edition name without accepting aliases for
// editions whose differences have not been inventoried. The only accepted name
// is ECMAScript 2025. Errors wrap ErrUnsupportedEdition and include the supplied
// name, so they are not a sanitized public diagnostic.
func ParseEdition(name string) (Edition, error) {
	if name == Edition2025.String() {
		return Edition2025, nil
	}

	return 0, fmt.Errorf("%w: %q", ErrUnsupportedEdition, name)
}

// String returns the closed public edition name or a diagnostic for unsupported values.
func (e Edition) String() string {
	if e == Edition2025 {
		return "ECMAScript 2025"
	}

	return fmt.Sprintf("ECMAScript edition %d (unsupported)", uint16(e))
}
