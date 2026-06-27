package awk_test

import (
	. "github.com/gloo-foo/cmd-awk"
	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"
)

// This example demonstrates CSV processing from a file
func ExampleAwk_fromFile_csvProcessing() {
	// cat testdata/people.csv | awk -F, '{print $1": "$2" years old"}'
	patterns.MustRun(
		Awk(
			csvProcessingProgram{},
			AwkFieldSeparator(","),
			gloo.File("testdata/people.csv"),
		),
	)
	// Output:
	// Alice: 30 years old
	// Bob: 25 years old
	// Charlie: 35 years old
	// Diana: 28 years old
}
