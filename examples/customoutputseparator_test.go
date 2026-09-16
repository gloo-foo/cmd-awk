package awk_test

import (
	"strings"

	gloo "github.com/gloo-foo/framework/patterns"

	awk "github.com/gloo-foo/cmd-awk"
)

// outputSeparatorProgram demonstrates custom output field separator.
// By default, fields are joined with a space. Use OutputFieldSeparator()
// to format output as CSV, TSV, or other delimited formats.
type outputSeparatorProgram struct {
	awk.SimpleProgram
}

func (p outputSeparatorProgram) Action(ctx awk.Context) (awk.Context, string, bool) {
	return ctx, ctx.Print(ctx.Field(1), ctx.Field(2), ctx.Field(3)), true
}

func ExampleAwk_customOutputSeparator() {
	// echo "one two three" | awk 'BEGIN{OFS=","} {print $1,$2,$3}'
	if err := gloo.Run(
		awk.Awk(
			outputSeparatorProgram{},
			awk.AwkOutputFieldSeparator(","),
			strings.NewReader("one two three"),
		),
	); err != nil {
		panic(err)
	}
	// Output:
	// one,two,three
}
