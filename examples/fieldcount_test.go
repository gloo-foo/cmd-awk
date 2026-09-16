package awk_test

import (
	"fmt"
	"strings"

	gloo "github.com/gloo-foo/framework/patterns"

	awk "github.com/gloo-foo/cmd-awk"
)

// fieldCountProgram demonstrates using NF for field count.
// NF is a built-in variable that contains the number of fields in the current line.
// This varies per line and is useful for validation or reporting.
type fieldCountProgram struct {
	awk.SimpleProgram
}

func (p fieldCountProgram) Action(ctx awk.Context) (awk.Context, string, bool) {
	return ctx, fmt.Sprintf("%d fields: %s", ctx.NF, ctx.Field(0)), true
}

func ExampleAwk_fieldCount() {
	// echo -e "one two\nthree four five" | awk '{print NF" fields"}'
	if err := gloo.Run(
		awk.Awk(
			fieldCountProgram{},
			strings.NewReader("one two\nthree four five"),
		),
	); err != nil {
		panic(err)
	}
	// Output:
	// 2 fields: one two
	// 3 fields: three four five
}
