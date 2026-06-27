package awk_test

import (
	"fmt"
	"strings"

	. "github.com/gloo-foo/cmd-awk"
	gloo "github.com/gloo-foo/framework/patterns"
)

// fieldCountProgram demonstrates using NF for field count.
// NF is a built-in variable that contains the number of fields in the current line.
// This varies per line and is useful for validation or reporting.
type fieldCountProgram struct {
	SimpleProgram
}

func (p fieldCountProgram) Action(ctx *Context) (string, bool) {
	return fmt.Sprintf("%d fields: %s", ctx.NF, ctx.Field(0)), true
}

func ExampleAwk_fieldCount() {
	// echo -e "one two\nthree four five" | awk '{print NF" fields"}'
	gloo.MustRun(
		Awk(
			fieldCountProgram{},
			strings.NewReader("one two\nthree four five"),
		),
	)
	// Output:
	// 2 fields: one two
	// 3 fields: three four five
}
