package awk_test

import (
	"strings"

	. "github.com/gloo-foo/cmd-awk"
	gloo "github.com/gloo-foo/framework/patterns"
)

// uniqueLinesProgram demonstrates deduplication using variables.
// This advanced example shows how to maintain state across lines using
// ctx.SetVar/Var and Go's native data structures like maps.
type uniqueLinesProgram struct {
	SimpleProgram
}

func (p uniqueLinesProgram) Begin(ctx *Context) error {
	ctx.SetVar("seen", make(map[string]bool))
	return nil
}

func (p uniqueLinesProgram) Condition(ctx *Context) bool {
	seen := ctx.Var("seen").(map[string]bool)
	line := ctx.Field(0)
	if seen[line] {
		return false
	}
	seen[line] = true
	return true
}

func (p uniqueLinesProgram) Action(ctx *Context) (string, bool) {
	return ctx.Field(0), true
}

func ExampleAwk_uniqueLines() {
	// echo -e "apple\nbanana\napple\ncherry\nbanana" | awk '!seen[$0]++'
	gloo.MustRun(
		Awk(
			uniqueLinesProgram{},
			strings.NewReader("apple\nbanana\napple\ncherry\nbanana"),
		),
	)
	// Output:
	// apple
	// banana
	// cherry
}
