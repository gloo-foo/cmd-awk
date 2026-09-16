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

func (p beginProgram) Begin(ctx awk.Context) (awk.Context, error) {
	fmt.Println("Starting processing...")
	return ctx, nil
}

func (p beginProgram) Action(ctx awk.Context) (awk.Context, string, bool) {
	return ctx, ctx.Field(0), true
}

func ExampleAwk_begin() {
	// echo "data" | awk 'BEGIN{print "Starting..."} {print $0}'
	if err := gloo.Run(
		awk.Awk(
			beginProgram{},
			strings.NewReader("data"),
		),
	); err != nil {
		panic(err)
	}
	// Output:
	// Starting processing...
	// data
}
