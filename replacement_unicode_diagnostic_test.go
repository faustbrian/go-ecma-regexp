package ecmascript_test

import (
	"context"
	"os"
	"slices"
	"testing"

	ecmascript "github.com/faustbrian/go-ecma-regexp/v2"
)

func TestReplacementUnicodeCaptureNamesHosted(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("Unicode replacement regression runs in hosted CI")
	}
	for _, test := range []struct {
		name, pattern, input, replacement, want string
	}{
		{"BMP", `(?<π>a)`, "a", "$<π>", "a"},
		{"supplementary", `(?<𐐀>a)`, "a", "$<𐐀>", "a"},
		{"escaped-pattern", `(?<\u03c0>a)`, "a", "$<π>", "a"},
		{"distinct-names", `(?<a>a)(?<š>b)`, "ab", "$<š>-$<a>", "b-a"},
		{"missing-name", `(?<a>a)`, "a", "$<š>", ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			program, err := ecmascript.Compile(test.pattern, "u", ecmascript.DefaultCompileOptions())
			if err != nil {
				t.Fatal(err)
			}
			for _, representation := range []string{"string", "UTF16"} {
				t.Run(representation, func(t *testing.T) {
					replacement := ecmascript.UTF16FromString(test.replacement)
					var got ecmascript.UTF16String
					var err error
					if representation == "string" {
						got, err = program.Replace(context.Background(), test.input, replacement, ecmascript.DefaultMatchOptions())
					} else {
						got, err = program.ReplaceUTF16(context.Background(), ecmascript.UTF16FromString(test.input), replacement, ecmascript.DefaultMatchOptions())
					}
					if err != nil || !slices.Equal(got.Units(), ecmascript.UTF16FromString(test.want).Units()) {
						t.Errorf("replacement = %v, %v; want %v", got.Units(), err, ecmascript.UTF16FromString(test.want).Units())
					}
				})
			}
		})
	}
}

func TestReplacementUnpairedCaptureNameHosted(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("replacement surrogate boundary runs in hosted CI")
	}
	program, err := ecmascript.Compile(`(?<x>a)`, "u", ecmascript.DefaultCompileOptions())
	if err != nil {
		t.Fatal(err)
	}
	for _, representation := range []string{"string", "UTF16"} {
		t.Run(representation, func(t *testing.T) {
			replacement := ecmascript.UTF16FromUnits([]uint16{'$', '<', 0xd800, '>'})
			var got ecmascript.UTF16String
			var err error
			if representation == "string" {
				got, err = program.Replace(context.Background(), "a", replacement, ecmascript.DefaultMatchOptions())
			} else {
				got, err = program.ReplaceUTF16(context.Background(), ecmascript.UTF16FromString("a"), replacement, ecmascript.DefaultMatchOptions())
			}
			if err != nil || len(got.Units()) != 0 {
				t.Fatalf("unpaired missing-name replacement = %v, %v; want empty output", got.Units(), err)
			}
		})
	}
}
