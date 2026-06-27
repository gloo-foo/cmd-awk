package awk_test

import (
	"strings"

	. "github.com/gloo-foo/cmd-awk"
	gloo "github.com/gloo-foo/framework/patterns"
)

// conditionProgram demonstrates conditional processing.
// The Condition method filters which lines are processed by Action.
// Return true to process the line, false to skip it.
type conditionProgram struct {
	SimpleProgram
	pattern string
}

func (p conditionProgram) Condition(ctx *Context) bool {
	return strings.Contains(ctx.Field(0), p.pattern)
}

func (p conditionProgram) Action(ctx *Context) (string, bool) {
	return ctx.Field(0), true
}

func ExampleAwk_condition() {
	// echo -e "apple\nbanana\napricot" | awk '/ap/'
	gloo.MustRun(
		Awk(
			conditionProgram{pattern: "ap"},
			strings.NewReader("apple\nbanana\napricot"),
		),
	)
	// Output:
	// apple
	// apricot
}
