package awk_test

import (
	"strings"

	gloo "github.com/gloo-foo/framework/patterns"

	awk "github.com/gloo-foo/cmd-awk"
)

// uppercaseProgram demonstrates text transformation.
// Since you have access to Go's full standard library, text transformations
// are simple and type-safe. No need to remember awk function names.
type uppercaseProgram struct {
	awk.SimpleProgram
}

func (p uppercaseProgram) Action(ctx awk.Context) (awk.Context, string, bool) {
	return ctx, strings.ToUpper(ctx.Field(0)), true
}

func ExampleAwk_uppercase() {
	// echo "hello world" | awk '{print toupper($0)}'
	if err := gloo.Run(
		awk.Awk(
			uppercaseProgram{},
			strings.NewReader("hello world"),
		),
	); err != nil {
		panic(err)
	}
	// Output:
	// HELLO WORLD
}
