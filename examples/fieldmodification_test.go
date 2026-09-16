package awk_test

import (
	"strings"

	gloo "github.com/gloo-foo/framework/patterns"

	awk "github.com/gloo-foo/cmd-awk"
)

// fieldModificationProgram demonstrates modifying fields.
// Use ctx.SetField() to change field values, then print with ctx.Print()
// to output with proper field separation.
type fieldModificationProgram struct {
	awk.SimpleProgram
}

func (p fieldModificationProgram) Action(ctx awk.Context) (awk.Context, string, bool) {
	ctx = ctx.SetField(2, "MODIFIED")
	return ctx, ctx.Print(ctx.Field(1), ctx.Field(2), ctx.Field(3)), true
}

func ExampleAwk_fieldModification() {
	// echo "one two three" | awk '{$2="MODIFIED"; print}'
	if err := gloo.Run(
		awk.Awk(
			fieldModificationProgram{},
			strings.NewReader("one two three"),
		),
	); err != nil {
		panic(err)
	}
	// Output:
	// one MODIFIED three
}
