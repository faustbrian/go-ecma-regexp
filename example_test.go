package ecmascript_test

import (
	"context"
	"fmt"

	ecmascript "github.com/faustbrian/go-ecma-regexp"
)

func Example() {
	program, err := ecmascript.Compile(
		`(?<word>\p{Letter}+)`,
		"u",
		ecmascript.DefaultCompileOptions(),
	)
	if err != nil {
		panic(err)
	}

	result, matched, err := program.Find(
		context.Background(),
		"42 Helsinki",
		ecmascript.DefaultMatchOptions(),
	)
	if err != nil {
		panic(err)
	}
	if !matched {
		panic("expected a match")
	}

	word, ok := result.Named("word")
	if !ok {
		panic("expected the named capture")
	}

	fmt.Println(word.Value().LossyString())
	// Output: Helsinki
}
