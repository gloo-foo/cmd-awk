package awk_test

import (
	"strings"

	gloo "github.com/gloo-foo/framework/patterns"

	awk "github.com/gloo-foo/cmd-awk"
)

// printMultipleFieldsProgram demonstrates printing multiple fields.
// Uses ctx.Print() to combine fields with the output field separator (OFS).
type printMultipleFieldsProgram struct {
	awk.SimpleProgram
}

func (p printMultipleFieldsProgram) Action(ctx awk.Context) (awk.Context, string, bool) {
	return ctx, ctx.Print(ctx.Field(1), ctx.Field(3)), true
}

func ExampleAwk_printMultipleFields() {
	// echo "one two three four" | awk '{print $1, $3}'
	if err := gloo.Run(
		awk.Awk(
			printMultipleFieldsProgram{},
			strings.NewReader("one two three four"),
		),
	); err != nil {
		panic(err)
	}
	// Output:
	// one three
}
