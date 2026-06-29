package awk_test

import (
	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"

	awk "github.com/gloo-foo/cmd-awk"
)

// This example demonstrates removing duplicate lines from a file
func ExampleAwk_fromFile_uniqueLines() {
	// cat testdata/duplicates.txt | awk '!seen[$0]++'
	patterns.MustRun(
		awk.Awk(
			uniqueLinesProgram{},
			gloo.File("testdata/duplicates.txt"),
		),
	)
	// Output:
	// apple
	// banana
	// cherry
	// grape
}
