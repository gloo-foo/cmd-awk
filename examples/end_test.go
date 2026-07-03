package awk_test

import (
	"strings"

	gloo "github.com/gloo-foo/framework/patterns"

	awk "github.com/gloo-foo/cmd-awk"
)

// endProgram demonstrates END block finalization.
// The End method runs once after processing all input lines.
// Use it for printing summaries, totals, or cleanup tasks.
type endProgram struct {
	awk.SimpleProgram
}

func (p endProgram) Action(ctx awk.Context) (awk.Context, string, bool) {
	return ctx, ctx.Field(0), true
}

func (p endProgram) End(_ awk.Context) (string, error) {
	return "Processing complete.", nil
}

func ExampleAwk_end() {
	// echo "data" | awk '{print $0} END{print "Done"}'
	if err := gloo.Run(
		awk.Awk(
			endProgram{},
			strings.NewReader("data"),
		),
	); err != nil {
		panic(err)
	}
	// Output:
	// data
	// Processing complete.
}
