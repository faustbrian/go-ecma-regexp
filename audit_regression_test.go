package ecmascript_test

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	ecmascript "github.com/faustbrian/go-ecma-regexp/v2"
)

func TestNegativeSearchStartReturnsError(t *testing.T) {
	for _, flags := range []string{"", "u", "v", "y", "uy", "vy", "g", "gu", "gv"} {
		t.Run(flags, func(t *testing.T) {
			program, err := ecmascript.Compile("a", flags, ecmascript.DefaultCompileOptions())
			if err != nil {
				t.Fatal(err)
			}
			options := ecmascript.DefaultMatchOptions()
			options.StartUTF16 = -1
			exact := ecmascript.UTF16FromString("a")
			calls := map[string]func() error{
				"Match":        func() error { _, _, err := program.Match(context.Background(), "a", options); return err },
				"MatchUTF16":   func() error { _, _, err := program.MatchUTF16(context.Background(), exact, options); return err },
				"Find":         func() error { _, _, err := program.Find(context.Background(), "a", options); return err },
				"FindUTF16":    func() error { _, _, err := program.FindUTF16(context.Background(), exact, options); return err },
				"FindAll":      func() error { _, err := program.FindAll(context.Background(), "a", options); return err },
				"FindAllUTF16": func() error { _, err := program.FindAllUTF16(context.Background(), exact, options); return err },
				"Replace":      func() error { _, err := program.Replace(context.Background(), "a", exact, options); return err },
				"ReplaceUTF16": func() error { _, err := program.ReplaceUTF16(context.Background(), exact, exact, options); return err },
			}
			for name, call := range calls {
				t.Run(name, func(t *testing.T) {
					defer func() {
						if recovered := recover(); recovered != nil {
							t.Errorf("negative start panicked: %v", recovered)
						}
					}()
					var syntax *ecmascript.SyntaxError
					if err := call(); !errors.As(err, &syntax) || syntax.Code != ecmascript.SyntaxUnexpectedToken {
						t.Errorf("negative start error = %v", err)
					}
				})
			}
			if _, err := program.Split(context.Background(), "a", options); err != nil {
				t.Errorf("Split must ignore start: %v", err)
			}
			if _, err := program.SplitUTF16(context.Background(), exact, options); err != nil {
				t.Errorf("SplitUTF16 must ignore start: %v", err)
			}
		})
	}
}

// Run hostile counts only in a killable child; a compiler regression cannot
// trap the test process in expansion work that has no context cancellation.
func TestEmptyFiniteRepetitionCompileIsBounded(t *testing.T) {
	if pattern := os.Getenv("ECMA_REGEXP_EMPTY_REPEAT_PROBE"); pattern != "" {
		options := ecmascript.DefaultCompileOptions()
		options.Limits.ProgramInstructions = 1
		program, err := ecmascript.Compile(pattern, "u", options)
		if err != nil {
			t.Fatal(err)
		}
		result, matched, err := program.Match(context.Background(), "", ecmascript.DefaultMatchOptions())
		if err != nil || !matched || result.Full().Span().End.UTF16 != 0 {
			t.Fatalf("empty repetition match = %v, %v", matched, err)
		}
		return
	}
	count := strconv.Itoa(int(^uint(0) >> 1))
	for _, pattern := range []string{"(?:){" + count + "}", "(?:(?:){" + count + "}){" + count + "}", "(?:){" + count + ",}?"} {
		t.Run(pattern, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			command := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestEmptyFiniteRepetitionCompileIsBounded$", "-test.timeout=3s")
			// The parent owns the mutation-profile lock. The child keeps its
			// hard timeout but must not reacquire that same profile lock.
			for _, variable := range os.Environ() {
				if !strings.HasPrefix(variable, "GOLIB_GREMLINS_COVERAGE_PROFILE=") {
					command.Env = append(command.Env, variable)
				}
			}
			command.Env = append(command.Env, "ECMA_REGEXP_EMPTY_REPEAT_PROBE="+pattern)
			output, err := command.CombinedOutput()
			if err != nil {
				t.Fatalf("bounded empty compilation: %v; deadline=%v; output=%s", err, ctx.Err(), output)
			}
		})
	}
}

func TestLoneSurrogateCoordinatesRemainExactOnlyAtOrigin(t *testing.T) {
	program, err := ecmascript.Compile("(?<scalar>a)", "u", ecmascript.DefaultCompileOptions())
	if err != nil {
		t.Fatal(err)
	}
	for _, units := range [][]uint16{{'a', 0xD800}, {0xD800, 'a'}, {'a', 0xDC00}} {
		subject := ecmascript.UTF16FromUnits(units)
		if _, err := subject.GoString(); !errors.Is(err, ecmascript.ErrUnpairedSurrogate) {
			t.Fatalf("subject GoString = %v", err)
		}
		result, matched, err := program.FindUTF16(context.Background(), subject, ecmascript.DefaultMatchOptions())
		if err != nil || !matched {
			t.Fatalf("FindUTF16 = %v, %v", matched, err)
		}
		capture, _ := result.Named("scalar")
		if text, err := capture.Value().GoString(); err != nil || text != "a" {
			t.Fatalf("scalar capture = %q, %v", text, err)
		}
		span := capture.Span()
		for _, index := range []ecmascript.Index{span.Start, span.End} {
			if index.UTF16 == 0 {
				if index != (ecmascript.Index{Exact: true}) {
					t.Errorf("origin = %+v", index)
				}
			} else if index.Byte != -1 || index.Rune != -1 || index.Exact {
				t.Errorf("non-origin = %+v", index)
			}
		}
	}
}

func TestFiniteRepetitionPreservesEffectfulBodies(t *testing.T) {
	for _, test := range []struct {
		pattern, input, value string
		captures              int
	}{
		{"(){3}", "", "", 1},
		{"(?:(a)?){2}", "a", "a", 1},
		{"(?:(?=a)){3}a", "a", "a", 0},
		{"(?:(?<=a)){3}b", "ab", "b", 0},
		{"(?:(?:){9}){7}a", "a", "a", 0},
	} {
		t.Run(test.pattern, func(t *testing.T) {
			program, err := ecmascript.Compile(test.pattern, "u", ecmascript.DefaultCompileOptions())
			if err != nil {
				t.Fatal(err)
			}
			result, matched, err := program.Find(context.Background(), test.input, ecmascript.DefaultMatchOptions())
			if err != nil || !matched {
				t.Fatalf("Find = %v, %v", matched, err)
			}
			if got, err := result.Full().Value().GoString(); err != nil || got != test.value {
				t.Fatalf("value = %q, %v", got, err)
			}
			if len(result.Captures()) != test.captures+1 {
				t.Fatal("capture count changed")
			}
			if test.pattern == "(){3}" && !result.Captures()[1].Participated() {
				t.Fatal("empty capture lost participation")
			}
			if test.pattern == "(?:(a)?){2}" && result.Captures()[1].Participated() {
				t.Fatal("last empty iteration failed to clear capture")
			}
		})
	}
	for _, pattern := range []string{"(){99}", "(?:(?=a)){99}", "a{99}", "(?:|){99}", "(?:){0,99}"} {
		options := ecmascript.DefaultCompileOptions()
		options.Limits.ProgramInstructions = 8
		_, err := ecmascript.Compile(pattern, "u", options)
		var limit *ecmascript.LimitError
		if !errors.As(err, &limit) || limit.Kind != ecmascript.LimitProgramInstructions {
			t.Errorf("effectful %q compile = %v", pattern, err)
		}
	}
}

func TestJSONSchemaPatternRejectsNegativeStart(t *testing.T) {
	pattern, err := ecmascript.CompileJSONSchemaPattern("a", ecmascript.DefaultCompileOptions())
	if err != nil {
		t.Fatal(err)
	}
	options := ecmascript.DefaultMatchOptions()
	options.StartUTF16 = -1
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Errorf("JSON Schema negative start panicked: %v", recovered)
		}
	}()
	matched, err := pattern.Match(context.Background(), "a", options)
	var syntax *ecmascript.SyntaxError
	if matched || !errors.As(err, &syntax) {
		t.Fatalf("JSON Schema negative start = %v, %v", matched, err)
	}
}

func TestSearchStartBoundaryContracts(t *testing.T) {
	for _, flags := range []string{"", "u", "v", "uy"} {
		program, err := ecmascript.Compile("", flags, ecmascript.DefaultCompileOptions())
		if err != nil {
			t.Fatal(err)
		}
		options := ecmascript.DefaultMatchOptions()
		options.StartUTF16 = 1
		result, matched, err := program.Find(context.Background(), "a", options)
		if err != nil || !matched || result.Full().Span().Start.UTF16 != 1 {
			t.Fatalf("end start (%s) = %v, %v", flags, matched, err)
		}
		options.StartUTF16 = 2
		_, matched, err = program.Find(context.Background(), "a", options)
		if flags == "uy" {
			var syntax *ecmascript.SyntaxError
			if matched || !errors.As(err, &syntax) {
				t.Fatalf("sticky beyond end = %v,%v", matched, err)
			}
		} else if matched || err != nil {
			t.Fatalf("search beyond end = %v,%v", matched, err)
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		options.StartUTF16 = -1
		_, _, err = program.Find(ctx, "a", options)
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("cancellation priority = %v", err)
		}
	}
}
