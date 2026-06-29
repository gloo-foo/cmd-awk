package awk_test

import (
	"fmt"
	"strings"

	gloo "github.com/gloo-foo/framework/patterns"

	awk "github.com/gloo-foo/cmd-awk"
)

// lineNumberProgram demonstrates using NR for line numbers.
// NR is a built-in variable that tracks the current line number (1-based).
// This is useful for adding line numbers to output or filtering by line position.
type lineNumberProgram struct {
	awk.SimpleProgram
}

func (p lineNumberProgram) Action(ctx *awk.Context) (string, bool) {
	return fmt.Sprintf("%d: %s", ctx.NR, ctx.Field(0)), true
}

func ExampleAwk_lineNumbers() {
	// echo -e "first\nsecond\nthird" | awk '{print NR": "$0}'
	gloo.MustRun(
		awk.Awk(
			lineNumberProgram{},
			strings.NewReader("first\nsecond\nthird"),
		),
	)
	// Output:
	// 1: first
	// 2: second
	// 3: third
}
