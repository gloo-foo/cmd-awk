package awk_test

import (
	"strings"

	. "github.com/gloo-foo/cmd-awk"
	gloo "github.com/gloo-foo/framework/patterns"
)

// uppercaseProgram demonstrates text transformation.
// Since you have access to Go's full standard library, text transformations
// are simple and type-safe. No need to remember awk function names.
type uppercaseProgram struct {
	SimpleProgram
}

func (p uppercaseProgram) Action(ctx *Context) (string, bool) {
	return strings.ToUpper(ctx.Field(0)), true
}

func ExampleAwk_uppercase() {
	// echo "hello world" | awk '{print toupper($0)}'
	gloo.MustRun(
		Awk(
			uppercaseProgram{},
			strings.NewReader("hello world"),
		),
	)
	// Output:
	// HELLO WORLD
}
