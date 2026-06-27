package awk_test

import (
	"strings"

	. "github.com/gloo-foo/cmd-awk"
	gloo "github.com/gloo-foo/framework/patterns"
)

// lastFieldProgram demonstrates printing the last field.
// Use NF (number of fields) to access the last field regardless of line length.
// This is equivalent to awk's $NF syntax.
type lastFieldProgram struct {
	SimpleProgram
}

func (p lastFieldProgram) Action(ctx *Context) (string, bool) {
	return ctx.Field(ctx.NF), true
}

func ExampleAwk_lastField() {
	// echo "one two three four" | awk '{print $NF}'
	gloo.MustRun(
		Awk(
			lastFieldProgram{},
			strings.NewReader("one two three four"),
		),
	)
	// Output:
	// four
}
