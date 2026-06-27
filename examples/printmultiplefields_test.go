package awk_test

import (
	"strings"

	. "github.com/gloo-foo/cmd-awk"
	gloo "github.com/gloo-foo/framework/patterns"
)

// printMultipleFieldsProgram demonstrates printing multiple fields.
// Uses ctx.Print() to combine fields with the output field separator (OFS).
type printMultipleFieldsProgram struct {
	SimpleProgram
}

func (p printMultipleFieldsProgram) Action(ctx *Context) (string, bool) {
	return ctx.Print(ctx.Field(1), ctx.Field(3)), true
}

func ExampleAwk_printMultipleFields() {
	// echo "one two three four" | awk '{print $1, $3}'
	gloo.MustRun(
		Awk(
			printMultipleFieldsProgram{},
			strings.NewReader("one two three four"),
		),
	)
	// Output:
	// one three
}
