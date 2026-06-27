package awk_test

import (
	"strings"

	. "github.com/gloo-foo/cmd-awk"
	gloo "github.com/gloo-foo/framework/patterns"
)

// customSeparatorProgram demonstrates custom field separators.
// By default, fields are split on whitespace. Use FieldSeparator()
// to parse formats like CSV, colon-delimited files, or other structured text.
type customSeparatorProgram struct {
	SimpleProgram
}

func (p customSeparatorProgram) Action(ctx *Context) (string, bool) {
	return ctx.Field(2), true
}

func ExampleAwk_customFieldSeparator() {
	// echo "one:two:three" | awk -F: '{print $2}'
	gloo.MustRun(
		Awk(
			customSeparatorProgram{},
			AwkFieldSeparator(":"),
			strings.NewReader("one:two:three"),
		),
	)
	// Output:
	// two
}
