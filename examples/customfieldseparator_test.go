package awk_test

import (
	"strings"

	gloo "github.com/gloo-foo/framework/patterns"

	awk "github.com/gloo-foo/cmd-awk"
)

// customSeparatorProgram demonstrates custom field separators.
// By default, fields are split on whitespace. Use FieldSeparator()
// to parse formats like CSV, colon-delimited files, or other structured text.
type customSeparatorProgram struct {
	awk.SimpleProgram
}

func (p customSeparatorProgram) Action(ctx awk.Context) (awk.Context, string, bool) {
	return ctx, ctx.Field(2), true
}

func ExampleAwk_customFieldSeparator() {
	// echo "one:two:three" | awk -F: '{print $2}'
	if err := gloo.Run(
		awk.Awk(
			customSeparatorProgram{},
			awk.AwkFieldSeparator(":"),
			strings.NewReader("one:two:three"),
		),
	); err != nil {
		panic(err)
	}
	// Output:
	// two
}
