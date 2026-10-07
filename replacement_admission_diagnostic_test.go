package ecmascript_test

import (
	"context"
	"errors"
	"os"
	"reflect"
	"testing"

	ecmascript "github.com/faustbrian/go-ecma-regexp/v2"
)

// Ordinary tiny templates exercise admission without a resource campaign.
func TestReplacementInputAdmissionHosted(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("replacement admission regression runs in hosted CI")
	}
	program, err := ecmascript.Compile("a", "", ecmascript.DefaultCompileOptions())
	if err != nil {
		t.Fatal(err)
	}
	for _, representation := range []string{"string", "UTF16"} {
		for _, test := range []struct {
			name        string
			units       []uint16
			bytes       uint64
			runes       uint64
			kind        ecmascript.LimitKind
			limit, used uint64
		}{
			{name: "byte-exact", units: []uint16{'b'}, bytes: 2, runes: 2},
			{name: "byte-successor", units: []uint16{'b', 'c'}, bytes: 2, runes: 2, kind: ecmascript.LimitInputBytes, limit: 2, used: 4},
			{name: "rune-exact", units: []uint16{'b'}, bytes: 4, runes: 1},
			{name: "rune-successor", units: []uint16{'b', 'c'}, bytes: 4, runes: 1, kind: ecmascript.LimitInputRunes, limit: 1, used: 2},
			{name: "supplementary-code-point", units: []uint16{0xd83d, 0xde00}, bytes: 4, runes: 1},
			{name: "lone-unit", units: []uint16{0xd800}, bytes: 2, runes: 1},
		} {
			t.Run(representation+"/"+test.name, func(t *testing.T) {
				options := ecmascript.DefaultMatchOptions()
				options.Limits.InputBytes = test.bytes
				options.Limits.InputRunes = test.runes
				replacement := ecmascript.UTF16FromUnits(test.units)
				var value ecmascript.UTF16String
				var err error
				if representation == "string" {
					value, err = program.Replace(context.Background(), "a", replacement, options)
				} else {
					value, err = program.ReplaceUTF16(context.Background(), ecmascript.UTF16FromString("a"), replacement, options)
				}
				if test.kind == 0 {
					if err != nil || !reflect.DeepEqual(value.Units(), test.units) {
						t.Fatalf("admitted replacement = %v, error %v; want %v", value.Units(), err, test.units)
					}
					return
				}
				var limit *ecmascript.LimitError
				if !errors.As(err, &limit) || limit.Kind != test.kind || limit.Limit != test.limit || limit.Used != test.used {
					t.Errorf("error = %v, want kind %v limit %d used %d", err, test.kind, test.limit, test.used)
				}
				if len(value.Units()) != 0 {
					t.Error("refused replacement returned partial output")
				}
			})
		}
	}
}

func TestReplacementConsumedNameBudgetHosted(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("replacement work regression runs in hosted CI")
	}
	program, err := ecmascript.Compile("(?<item>a)", "", ecmascript.DefaultCompileOptions())
	if err != nil {
		t.Fatal(err)
	}
	options := ecmascript.DefaultMatchOptions()
	options.Limits.Steps = 16
	for _, representation := range []string{"string", "UTF16"} {
		t.Run(representation, func(t *testing.T) {
			replacement := ecmascript.UTF16FromString("$<ordinaryMissingName>")
			var value ecmascript.UTF16String
			var err error
			if representation == "string" {
				value, err = program.Replace(context.Background(), "a", replacement, options)
			} else {
				value, err = program.ReplaceUTF16(context.Background(), ecmascript.UTF16FromString("a"), replacement, options)
			}
			var limit *ecmascript.LimitError
			if !errors.As(err, &limit) || limit.Kind != ecmascript.LimitMatchSteps || limit.Limit != 16 || limit.Used != 17 {
				t.Errorf("error = %v, want step limit 16 used 17", err)
			}
			if len(value.Units()) != 0 {
				t.Error("exhausted replacement returned partial output")
			}
		})
	}
}
