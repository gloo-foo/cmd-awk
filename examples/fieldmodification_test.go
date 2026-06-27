package awk_test

import (
	"strings"

	. "github.com/gloo-foo/cmd-awk"
	gloo "github.com/gloo-foo/framework/patterns"
)

// fieldModificationProgram demonstrates modifying fields.
// Use ctx.SetField() to change field values, then print with ctx.Print()
// to output with proper field separation.
type fieldModificationProgram struct {
	SimpleProgram
}

func (p fieldModificationProgram) Action(ctx *Context) (string, bool) {
	ctx.SetField(2, "MODIFIED")
	return ctx.Print(ctx.Field(1), ctx.Field(2), ctx.Field(3)), true
}

func ExampleAwk_fieldModification() {
	// echo "one two three" | awk '{$2="MODIFIED"; print}'
	gloo.MustRun(
		Awk(
			fieldModificationProgram{},
			strings.NewReader("one two three"),
		),
	)
	// Output:
	// one MODIFIED three
}
