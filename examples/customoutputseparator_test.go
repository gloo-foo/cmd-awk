package awk_test

import (
	"strings"

	. "github.com/gloo-foo/cmd-awk"
	gloo "github.com/gloo-foo/framework/patterns"
)

// outputSeparatorProgram demonstrates custom output field separator.
// By default, fields are joined with a space. Use OutputFieldSeparator()
// to format output as CSV, TSV, or other delimited formats.
type outputSeparatorProgram struct {
	SimpleProgram
}

func (p outputSeparatorProgram) Action(ctx *Context) (string, bool) {
	return ctx.Print(ctx.Field(1), ctx.Field(2), ctx.Field(3)), true
}

func ExampleAwk_customOutputSeparator() {
	// echo "one two three" | awk 'BEGIN{OFS=","} {print $1,$2,$3}'
	gloo.MustRun(
		Awk(
			outputSeparatorProgram{},
			AwkOutputFieldSeparator(","),
			strings.NewReader("one two three"),
		),
	)
	// Output:
	// one,two,three
}
