package ecmascript

import "context"

// Session is caller-owned mutable RegExp execution state. Programs remain
// immutable and safe for concurrent use; a Session is not concurrency-safe.
type Session struct {
	program   *Program
	lastIndex int
}

// NewSession borrows Program and starts lastIndex at zero. It expects a non-nil
// compiled Program but does not validate the pointer. The zero Session or a nil
// Program is not a usable execution session.
func NewSession(program *Program) *Session {
	return &Session{program: program}
}

// LastIndex returns UTF-16 state, not a byte or rune offset.
func (s *Session) LastIndex() int { return s.lastIndex }

// SetLastIndex clamps negatives to zero. Values beyond input end remain stored
// until stateful Exec resets them on ordinary failure.
func (s *Session) SetLastIndex(index int) {
	s.lastIndex = max(index, 0)
}

// Exec applies RegExpBuiltinExec lastIndex, global, and sticky behavior. Limit
// failures leave lastIndex unchanged; an ordinary failed global or sticky
// match resets it to zero. Non-global/non-sticky programs ignore lastIndex.
// Empty success leaves lastIndex at that same boundary; callers looping Exec
// must apply their own progress policy. No workers or locks are supplied.
// nil context means Background.
func (s *Session) Exec(ctx context.Context, input string, limits MatchLimits) (Result, bool, error) {
	view, err := makeInputViewContext(ctx, input, limits)
	if err != nil {
		return Result{}, false, err
	}
	return s.exec(ctx, view, limits)
}

// ExecUTF16 applies stateful RegExp execution to an exact ECMAScript string.
func (s *Session) ExecUTF16(ctx context.Context, input UTF16String, limits MatchLimits) (Result, bool, error) {
	view, err := makeUTF16InputViewContext(ctx, input, limits)
	if err != nil {
		return Result{}, false, err
	}
	return s.exec(ctx, view, limits)
}

func (s *Session) exec(ctx context.Context, view *inputView, limits MatchLimits) (Result, bool, error) {
	if err := normalizeContext(ctx).Err(); err != nil {
		return Result{}, false, err
	}
	stateful := s.program.flags.Global() || s.program.flags.Sticky()
	start := 0
	if stateful {
		start = s.lastIndex
		if start > len(view.units) {
			s.lastIndex = 0
			return Result{}, false, nil
		}
	}
	executor := newExecutor(ctx, s.program, view, limits)
	result, matched, err := executor.find(start, s.program.flags.Sticky())
	if err != nil {
		return Result{}, false, err
	}
	if !matched {
		if stateful {
			s.lastIndex = 0
		}
		return Result{}, false, nil
	}
	if stateful {
		s.lastIndex = result.Full().span.End.UTF16
	}

	return result, true, nil
}
