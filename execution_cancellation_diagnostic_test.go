package ecmascript_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	ecmascript "github.com/faustbrian/go-ecma-regexp/v2"
)

type cancellationObservation struct {
	err       error
	values    int
	matched   bool
	lastIndex int
}

// These cases observe cancellation at the public entry boundary, not during
// preparation or execution. All fixtures are ordinary one-character inputs.
func TestExecutionCancellationPrecedesInputAdmissionHosted(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("execution preparation regression runs in hosted CI")
	}
	entries := executionCancellationEntries(t)
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	expired, stop := context.WithDeadline(context.Background(), time.Unix(1, 0))
	defer stop()
	options := ecmascript.DefaultMatchOptions()
	options.Limits.InputBytes = 0
	for name, call := range entries {
		for _, test := range []struct {
			name string
			ctx  context.Context
			want error
		}{
			{"canceled", canceled, context.Canceled},
			{"deadline", expired, context.DeadlineExceeded},
			{"ordinary-limit", context.Background(), nil},
		} {
			t.Run(name+"/"+test.name, func(t *testing.T) {
				got := call(test.ctx, options)
				if test.want != nil {
					if !errors.Is(got.err, test.want) {
						t.Errorf("error = %v, want %v", got.err, test.want)
					}
				} else {
					var limit *ecmascript.LimitError
					if !errors.As(got.err, &limit) || limit.Kind != ecmascript.LimitInputBytes {
						t.Errorf("error = %v, want input-byte limit", got.err)
					}
				}
				if got.matched || got.values != 0 {
					t.Errorf("refusal returned partial output: matched=%v values=%d", got.matched, got.values)
				}
				if name == "Session.Exec" || name == "Session.ExecUTF16" {
					if got.lastIndex != 2 {
						t.Errorf("refusal changed lastIndex to %d, want 2", got.lastIndex)
					}
				}
			})
		}
	}
}

func TestCanceledSessionPreservesLastIndexHosted(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("session cancellation regression runs in hosted CI")
	}
	entries := executionCancellationEntries(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for _, name := range []string{"Session.Exec", "Session.ExecUTF16"} {
		t.Run(name, func(t *testing.T) {
			got := entries[name](ctx, ecmascript.DefaultMatchOptions())
			if !errors.Is(got.err, context.Canceled) {
				t.Errorf("error = %v, want canceled", got.err)
			}
			if got.matched || got.values != 0 || got.lastIndex != 2 {
				t.Errorf("cancellation result = %+v, want no output and lastIndex 2", got)
			}
		})
	}
}

func TestExecutionNilContextMatchesBackgroundHosted(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("execution context characterization runs in hosted CI")
	}
	for name, call := range executionCancellationEntries(t) {
		t.Run(name, func(t *testing.T) {
			options := ecmascript.DefaultMatchOptions()
			background := call(context.Background(), options)
			without := call(nil, options)
			if background.err != nil || without.err != nil {
				t.Fatalf("ordinary calls returned errors: background=%v nil=%v", background.err, without.err)
			}
			if without.values != background.values || without.matched != background.matched || without.lastIndex != background.lastIndex {
				t.Fatalf("nil context result = %+v, background result = %+v", without, background)
			}
		})
	}
}

func executionCancellationEntries(t *testing.T) map[string]func(context.Context, ecmascript.MatchOptions) cancellationObservation {
	t.Helper()
	program, err := ecmascript.Compile("a", "g", ecmascript.DefaultCompileOptions())
	if err != nil {
		t.Fatal(err)
	}
	pattern, err := ecmascript.CompileJSONSchemaPattern("a", ecmascript.DefaultCompileOptions())
	if err != nil {
		t.Fatal(err)
	}
	input := ecmascript.UTF16FromString("a")
	replacement := ecmascript.UTF16FromString("b")
	result := func(value ecmascript.Result, matched bool, err error) cancellationObservation {
		return cancellationObservation{err: err, values: len(value.Captures()), matched: matched}
	}
	return map[string]func(context.Context, ecmascript.MatchOptions) cancellationObservation{
		"Program.Match": func(ctx context.Context, options ecmascript.MatchOptions) cancellationObservation {
			return result(program.Match(ctx, "a", options))
		},
		"Program.MatchUTF16": func(ctx context.Context, options ecmascript.MatchOptions) cancellationObservation {
			return result(program.MatchUTF16(ctx, input, options))
		},
		"Program.Find": func(ctx context.Context, options ecmascript.MatchOptions) cancellationObservation {
			return result(program.Find(ctx, "a", options))
		},
		"Program.FindUTF16": func(ctx context.Context, options ecmascript.MatchOptions) cancellationObservation {
			return result(program.FindUTF16(ctx, input, options))
		},
		"Program.FindAll": func(ctx context.Context, options ecmascript.MatchOptions) cancellationObservation {
			values, err := program.FindAll(ctx, "a", options)
			return cancellationObservation{err: err, values: len(values)}
		},
		"Program.FindAllUTF16": func(ctx context.Context, options ecmascript.MatchOptions) cancellationObservation {
			values, err := program.FindAllUTF16(ctx, input, options)
			return cancellationObservation{err: err, values: len(values)}
		},
		"Program.Replace": func(ctx context.Context, options ecmascript.MatchOptions) cancellationObservation {
			value, err := program.Replace(ctx, "a", replacement, options)
			return cancellationObservation{err: err, values: len(value.Units())}
		},
		"Program.ReplaceUTF16": func(ctx context.Context, options ecmascript.MatchOptions) cancellationObservation {
			value, err := program.ReplaceUTF16(ctx, input, replacement, options)
			return cancellationObservation{err: err, values: len(value.Units())}
		},
		"Program.Split": func(ctx context.Context, options ecmascript.MatchOptions) cancellationObservation {
			values, err := program.Split(ctx, "a", options)
			return cancellationObservation{err: err, values: len(values)}
		},
		"Program.SplitUTF16": func(ctx context.Context, options ecmascript.MatchOptions) cancellationObservation {
			values, err := program.SplitUTF16(ctx, input, options)
			return cancellationObservation{err: err, values: len(values)}
		},
		"Session.Exec": func(ctx context.Context, options ecmascript.MatchOptions) cancellationObservation {
			session := ecmascript.NewSession(program)
			session.SetLastIndex(2)
			got := result(session.Exec(ctx, "a", options.Limits))
			got.lastIndex = session.LastIndex()
			return got
		},
		"Session.ExecUTF16": func(ctx context.Context, options ecmascript.MatchOptions) cancellationObservation {
			session := ecmascript.NewSession(program)
			session.SetLastIndex(2)
			got := result(session.ExecUTF16(ctx, input, options.Limits))
			got.lastIndex = session.LastIndex()
			return got
		},
		"JSONSchemaPattern.Match": func(ctx context.Context, options ecmascript.MatchOptions) cancellationObservation {
			matched, err := pattern.Match(ctx, "a", options)
			return cancellationObservation{err: err, matched: matched}
		},
	}
}
