package awk_test

import (
	"strings"

	gloo "github.com/gloo-foo/framework/patterns"

	awk "github.com/gloo-foo/cmd-awk"
)

// printFieldProgram demonstrates printing a specific field.
// This is the most basic awk operation - selecting and printing
// a single field by its position (e.g., $1, $2, $3).
type printFieldProgram struct {
	awk.SimpleProgram
	fieldNum int
}

func (p printFieldProgram) Action(ctx awk.Context) (awk.Context, string, bool) {
	return ctx, ctx.Field(p.fieldNum), true
}

func ExampleAwk_printField() {
	// echo "one two three" | awk '{print $2}'
	if err := gloo.Run(
		awk.Awk(
			printFieldProgram{fieldNum: 2},
			strings.NewReader("one two three"),
		),
	); err != nil {
		panic(err)
	}
	// Output:
	// two
}
