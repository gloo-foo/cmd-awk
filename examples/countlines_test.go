package awk_test

import (
	"fmt"
	"strings"

	gloo "github.com/gloo-foo/framework/patterns"

	awk "github.com/gloo-foo/cmd-awk"
)

// countLinesProgram demonstrates counting lines.
// NR tracks the line number, so by the END block it contains the total count.
// This example shows how to process silently and emit only a summary.
type countLinesProgram struct {
	awk.SimpleProgram
}

func (p countLinesProgram) Action(_ *awk.Context) (string, bool) {
	return "", false // Don't emit output per line
}

func (p countLinesProgram) End(ctx *awk.Context) (string, error) {
	return fmt.Sprintf("Total lines: %d", ctx.NR), nil
}

func ExampleAwk_countLines() {
	// echo -e "line1\nline2\nline3" | awk 'END{print NR}'
	gloo.MustRun(
		awk.Awk(
			countLinesProgram{},
			strings.NewReader("line1\nline2\nline3"),
		),
	)
	// Output:
	// Total lines: 3
}
