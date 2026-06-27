package awk_test

import (
	"strings"

	. "github.com/gloo-foo/cmd-awk"
	gloo "github.com/gloo-foo/framework/patterns"
)

// endProgram demonstrates END block finalization.
// The End method runs once after processing all input lines.
// Use it for printing summaries, totals, or cleanup tasks.
type endProgram struct {
	SimpleProgram
}

func (p endProgram) Action(ctx *Context) (string, bool) {
	return ctx.Field(0), true
}

func (p endProgram) End(ctx *Context) (string, error) {
	return "Processing complete.", nil
}

func ExampleAwk_end() {
	// echo "data" | awk '{print $0} END{print "Done"}'
	gloo.MustRun(
		Awk(
			endProgram{},
			strings.NewReader("data"),
		),
	)
	// Output:
	// data
	// Processing complete.
}
