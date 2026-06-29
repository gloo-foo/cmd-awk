package awk_test

import (
	"strings"

	gloo "github.com/gloo-foo/framework/patterns"

	awk "github.com/gloo-foo/cmd-awk"
)

// lastFieldProgram demonstrates printing the last field.
// Use NF (number of fields) to access the last field regardless of line length.
// This is equivalent to awk's $NF syntax.
type lastFieldProgram struct {
	awk.SimpleProgram
}

func (p lastFieldProgram) Action(ctx *awk.Context) (string, bool) {
	return ctx.Field(ctx.NF), true
}

func ExampleAwk_lastField() {
	// echo "one two three four" | awk '{print $NF}'
	gloo.MustRun(
		awk.Awk(
			lastFieldProgram{},
			strings.NewReader("one two three four"),
		),
	)
	// Output:
	// four
}
