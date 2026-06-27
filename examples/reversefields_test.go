package awk_test

import (
	"strings"

	. "github.com/gloo-foo/cmd-awk"
	gloo "github.com/gloo-foo/framework/patterns"
)

// reverseFieldsProgram demonstrates field reordering.
// Iterate through fields in reverse order using NF and a standard Go loop.
// This pattern works for any field manipulation or reordering.
type reverseFieldsProgram struct {
	SimpleProgram
}

func (p reverseFieldsProgram) Action(ctx *Context) (string, bool) {
	var reversed []string
	for i := ctx.NF; i >= 1; i-- {
		reversed = append(reversed, ctx.Field(i))
	}
	return strings.Join(reversed, " "), true
}

func ExampleAwk_reverseFields() {
	// echo "one two three" | awk '{for(i=NF;i>=1;i--)print $i}'
	gloo.MustRun(
		Awk(
			reverseFieldsProgram{},
			strings.NewReader("one two three"),
		),
	)
	// Output:
	// three two one
}
