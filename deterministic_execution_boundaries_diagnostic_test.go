package ecmascript

import (
	"context"
	"errors"
	"os"
	"slices"
	"testing"
)

func TestReplacementWorkBoundariesHosted(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("replacement characterization runs in hosted CI")
	}
	for _, test := range []struct {
		name, pattern, replacement, want string
		steps, refusing                  uint64
	}{
		{name: "dollar-token", pattern: "a", replacement: "$$", want: "$", steps: 4, refusing: 3},
		{name: "numeric-lookahead", pattern: "(a)", replacement: "$1x", want: "ax", steps: 8, refusing: 6},
		{name: "named-conversion", pattern: "(?<x>a)", replacement: "$<x>", want: "a", steps: 9, refusing: 8},
	} {
		t.Run(test.name, func(t *testing.T) {
			program, err := Compile(test.pattern, "", DefaultCompileOptions())
			if err != nil {
				t.Fatal(err)
			}
			for _, representation := range []string{"string", "UTF16"} {
				t.Run(representation, func(t *testing.T) {
					for _, budget := range []uint64{test.refusing, test.steps} {
						options := DefaultMatchOptions()
						options.Limits.Steps = budget
						replacement := UTF16FromString(test.replacement)
						var value UTF16String
						var err error
						if representation == "string" {
							value, err = program.Replace(context.Background(), "a", replacement, options)
						} else {
							value, err = program.ReplaceUTF16(context.Background(), UTF16FromString("a"), replacement, options)
						}
						if budget == test.steps {
							if err != nil || !slices.Equal(value.Units(), UTF16FromString(test.want).Units()) {
								t.Fatalf("exact budget %d: output = %v, error = %v; want %q", budget, value.Units(), err, test.want)
							}
							continue
						}
						var limit *LimitError
						if !errors.As(err, &limit) || limit.Kind != LimitMatchSteps || limit.Limit != budget || limit.Used != budget+1 {
							t.Fatalf("refusal budget %d: error = %v; want step limit used %d", budget, err, budget+1)
						}
						if len(value.Units()) != 0 {
							t.Fatalf("exhausted replacement returned partial output %v", value.Units())
						}
					}
				})
			}
		})
	}
}

// These are helper admission checks, not cancellation during public execution.
func TestPreparedHelperCancellationRefusalHosted(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("helper cancellation characterization runs in hosted CI")
	}
	limits := DefaultMatchOptions().Limits
	view, err := makeInputView("a", limits)
	if err != nil {
		t.Fatal(err)
	}
	program, err := Compile("a", "g", DefaultCompileOptions())
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	t.Run("replacement-admission", func(t *testing.T) {
		refusing := limits
		refusing.InputBytes = 0
		refusing.InputRunes = 0
		if err := admitReplacement(ctx, []uint16{'x'}, refusing); !errors.Is(err, context.Canceled) {
			t.Fatalf("admitReplacement() = %v, want cancellation before input refusal", err)
		}
	})
	t.Run("executor-check", func(t *testing.T) {
		executor := newExecutor(ctx, program, view, limits)
		if err := executor.check(); !errors.Is(err, context.Canceled) {
			t.Fatalf("executor.check() = %v, want cancellation", err)
		}
	})
	t.Run("prepared-session", func(t *testing.T) {
		session := NewSession(program)
		session.SetLastIndex(2)
		result, matched, err := session.exec(ctx, view, limits)
		if !errors.Is(err, context.Canceled) || matched || len(result.Captures()) != 0 || session.LastIndex() != 2 {
			t.Fatalf("session.exec() = %v, %t, %v; lastIndex = %d", result, matched, err, session.LastIndex())
		}
	})
}

func TestLoneUTF16BoundaryMappingAcrossCheckpointHosted(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("UTF16 boundary characterization runs in hosted CI")
	}
	units := make([]uint16, 256)
	for index := range units {
		units[index] = 'a'
	}
	units[0] = 0xd800
	view, err := makeUTF16InputView(UTF16FromUnits(units), DefaultMatchOptions().Limits)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(view.units, units) || len(view.boundaries) != len(units)+1 || len(view.codePointBoundary) != len(units)+1 {
		t.Fatal("UTF16 input units or boundary lengths changed")
	}
	for index, boundary := range view.boundaries {
		want := Index{UTF16: index, Rune: -1, Byte: -1}
		if index == 0 {
			want = Index{Exact: true}
		}
		if boundary != want || !view.codePointBoundary[index] {
			t.Fatalf("boundary %d = %+v, code point = %t; want %+v, true", index, boundary, view.codePointBoundary[index], want)
		}
	}
}
