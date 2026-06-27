package awk_test

import (
	. "github.com/gloo-foo/cmd-awk"
	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"
)

// This example demonstrates removing duplicate lines from a file
func ExampleAwk_fromFile_uniqueLines() {
	// cat testdata/duplicates.txt | awk '!seen[$0]++'
	patterns.MustRun(
		Awk(
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
