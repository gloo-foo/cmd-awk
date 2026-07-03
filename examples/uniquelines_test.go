package awk_test

import (
	"strings"

	gloo "github.com/gloo-foo/framework/patterns"

	awk "github.com/gloo-foo/cmd-awk"
)

// uniqueLinesProgram demonstrates deduplication using variables.
// This advanced example shows how to maintain state across lines using
// ctx.SetVar/Var and Go's native data structures like maps.
type uniqueLinesProgram struct {
	awk.SimpleProgram
}

func (p uniqueLinesProgram) Begin(ctx awk.Context) (awk.Context, error) {
	ctx = ctx.SetVar("seen", make(map[string]bool))
	return ctx, nil
}

func (p uniqueLinesProgram) Condition(ctx awk.Context) bool {
	seen := ctx.Var("seen").(map[string]bool)
	line := ctx.Field(0)
	if seen[line] {
		return false
	}
	seen[line] = true
	return true
}

func (p uniqueLinesProgram) Action(ctx awk.Context) (awk.Context, string, bool) {
	return ctx, ctx.Field(0), true
}

func ExampleAwk_uniqueLines() {
	// echo -e "apple\nbanana\napple\ncherry\nbanana" | awk '!seen[$0]++'
	if err := gloo.Run(
		awk.Awk(
			uniqueLinesProgram{},
			strings.NewReader("apple\nbanana\napple\ncherry\nbanana"),
		),
	); err != nil {
		panic(err)
	}
	// Output:
	// apple
	// banana
	// cherry
}
