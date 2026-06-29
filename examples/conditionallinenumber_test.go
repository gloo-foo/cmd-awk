package awk_test

import (
	"strings"

	gloo "github.com/gloo-foo/framework/patterns"

	awk "github.com/gloo-foo/cmd-awk"
)

// conditionalLineNumberProgram demonstrates filtering by line number.
// Combines the Condition method with NR to select specific line ranges.
// Useful for extracting headers, skipping footers, or selecting ranges.
type conditionalLineNumberProgram struct {
	awk.SimpleProgram
}

func (p conditionalLineNumberProgram) Condition(ctx *awk.Context) bool {
	return ctx.NR > 1 && ctx.NR <= 3
}

func (p conditionalLineNumberProgram) Action(ctx *awk.Context) (string, bool) {
	return ctx.Field(0), true
}

func ExampleAwk_conditionalLineNumber() {
	// echo -e "1\n2\n3\n4\n5" | awk 'NR>1 && NR<=3'
	gloo.MustRun(
		awk.Awk(
			conditionalLineNumberProgram{},
			strings.NewReader("1\n2\n3\n4\n5"),
		),
	)
	// Output:
	// 2
	// 3
}
