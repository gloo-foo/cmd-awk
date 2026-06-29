package awk_test

import (
	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"

	awk "github.com/gloo-foo/cmd-awk"
)

// This example demonstrates CSV processing from a file
func ExampleAwk_fromFile_csvProcessing() {
	// cat testdata/people.csv | awk -F, '{print $1": "$2" years old"}'
	patterns.MustRun(
		awk.Awk(
			csvProcessingProgram{},
			awk.AwkFieldSeparator(","),
			gloo.File("testdata/people.csv"),
		),
	)
	// Output:
	// Alice: 30 years old
	// Bob: 25 years old
	// Charlie: 35 years old
	// Diana: 28 years old
}
