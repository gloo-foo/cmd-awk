package awk_test

import (
	"fmt"
	"strings"

	. "github.com/gloo-foo/cmd-awk"
	gloo "github.com/gloo-foo/framework/patterns"
)

// beginProgram demonstrates BEGIN block initialization.
// The Begin method runs once before processing any input lines.
// Use it for printing headers, initializing variables, or setup tasks.
type beginProgram struct {
	SimpleProgram
}

func (p beginProgram) Begin(ctx *Context) error {
	fmt.Println("Starting processing...")
	return nil
}

func (p beginProgram) Action(ctx *Context) (string, bool) {
	return ctx.Field(0), true
}

func ExampleAwk_begin() {
	// echo "data" | awk 'BEGIN{print "Starting..."} {print $0}'
	gloo.MustRun(
		Awk(
			beginProgram{},
			strings.NewReader("data"),
		),
	)
	// Output:
	// Starting processing...
	// data
}
