package awk_test

import (
	"fmt"
	"strings"

	gloo "github.com/gloo-foo/framework/patterns"

	awk "github.com/gloo-foo/cmd-awk"
)

// csvProcessingProgram demonstrates CSV processing.
// Combine FieldSeparator with field access and formatting to parse
// and transform structured data like CSV files.
type csvProcessingProgram struct {
	awk.SimpleProgram
}

func (p csvProcessingProgram) Action(ctx *awk.Context) (string, bool) {
	// Process CSV: name,age,city -> name: age years old
	return fmt.Sprintf("%s: %s years old", ctx.Field(1), ctx.Field(2)), true
}

func ExampleAwk_csvProcessing() {
	// echo -e "Alice,30,NYC\nBob,25,LA" | awk -F, '{print $1": "$2" years old"}'
	gloo.MustRun(
		awk.Awk(
			csvProcessingProgram{},
			awk.AwkFieldSeparator(","),
			strings.NewReader("Alice,30,NYC\nBob,25,LA"),
		),
	)
	// Output:
	// Alice: 30 years old
	// Bob: 25 years old
}
