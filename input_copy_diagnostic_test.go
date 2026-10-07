package ecmascript_test

import (
	"context"
	"os"
	"slices"
	"testing"
	"time"

	ecmascript "github.com/faustbrian/go-ecma-regexp/v2"
)

func TestUTF16InputSecondCopyChunkHosted(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("UTF-16 copy boundary runs in hosted CI")
	}
	program, err := ecmascript.Compile("a", "", ecmascript.DefaultCompileOptions())
	if err != nil {
		t.Fatal(err)
	}
	units := make([]uint16, 257)
	for index := range units {
		units[index] = 'x'
	}
	units[256] = 'a'
	input := ecmascript.UTF16FromUnits(units)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	options := ecmascript.DefaultMatchOptions()
	options.StartUTF16 = 256
	result, matched, err := program.MatchUTF16(ctx, input, options)
	if err != nil || !matched {
		t.Fatalf("second-chunk match = %t, %v; want success", matched, err)
	}
	full := result.Full()
	if !slices.Equal(full.Value().Units(), []uint16{'a'}) {
		t.Fatalf("second-chunk capture = %v", full.Value().Units())
	}
	span := full.Span()
	wantStart := ecmascript.Index{UTF16: 256, Rune: 256, Byte: 256, Exact: true}
	wantEnd := ecmascript.Index{UTF16: 257, Rune: 257, Byte: 257, Exact: true}
	if span.Start != wantStart || span.End != wantEnd {
		t.Fatalf("second-chunk span = %+v; want %+v to %+v", span, wantStart, wantEnd)
	}
	if !slices.Equal(input.Units(), units) {
		t.Fatal("input preparation changed caller-owned units")
	}
}
