package awk_test

import (
	"fmt"
	"strings"

	. "github.com/gloo-foo/cmd-awk"
	gloo "github.com/gloo-foo/framework/patterns"
)

// csvProcessingProgram demonstrates CSV processing.
// Combine FieldSeparator with field access and formatting to parse
// and transform structured data like CSV files.
type csvProcessingProgram struct {
	SimpleProgram
}

func (p csvProcessingProgram) Action(ctx *Context) (string, bool) {
	// Process CSV: name,age,city -> name: age years old
	return fmt.Sprintf("%s: %s years old", ctx.Field(1), ctx.Field(2)), true
}

func ExampleAwk_csvProcessing() {
	// echo -e "Alice,30,NYC\nBob,25,LA" | awk -F, '{print $1": "$2" years old"}'
	gloo.MustRun(
		Awk(
			csvProcessingProgram{},
			AwkFieldSeparator(","),
			strings.NewReader("Alice,30,NYC\nBob,25,LA"),
		),
	)
	// Output:
	// Alice: 30 years old
	// Bob: 25 years old
}
