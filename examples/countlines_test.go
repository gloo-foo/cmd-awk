package awk_test

import (
	"fmt"
	"strings"

	. "github.com/gloo-foo/cmd-awk"
	gloo "github.com/gloo-foo/framework/patterns"
)

// countLinesProgram demonstrates counting lines.
// NR tracks the line number, so by the END block it contains the total count.
// This example shows how to process silently and emit only a summary.
type countLinesProgram struct {
	SimpleProgram
}

func (p countLinesProgram) Action(ctx *Context) (string, bool) {
	return "", false // Don't emit output per line
}

func (p countLinesProgram) End(ctx *Context) (string, error) {
	return fmt.Sprintf("Total lines: %d", ctx.NR), nil
}

func ExampleAwk_countLines() {
	// echo -e "line1\nline2\nline3" | awk 'END{print NR}'
	gloo.MustRun(
		Awk(
			countLinesProgram{},
			strings.NewReader("line1\nline2\nline3"),
		),
	)
	// Output:
	// Total lines: 3
}
