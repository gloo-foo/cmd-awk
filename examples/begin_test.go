package awk_test

import (
	"fmt"
	"strings"

	gloo "github.com/gloo-foo/framework/patterns"

	awk "github.com/gloo-foo/cmd-awk"
)

// beginProgram demonstrates BEGIN block initialization.
// The Begin method runs once before processing any input lines.
// Use it for printing headers, initializing variables, or setup tasks.
type beginProgram struct {
	awk.SimpleProgram
}

func (p beginProgram) Begin(_ *awk.Context) error {
	fmt.Println("Starting processing...")
	return nil
}

func (p beginProgram) Action(ctx *awk.Context) (string, bool) {
	return ctx.Field(0), true
}

func ExampleAwk_begin() {
	// echo "data" | awk 'BEGIN{print "Starting..."} {print $0}'
	gloo.MustRun(
		awk.Awk(
			beginProgram{},
			strings.NewReader("data"),
		),
	)
	// Output:
	// Starting processing...
	// data
}
