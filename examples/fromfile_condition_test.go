package awk_test

import (
	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"

	awk "github.com/gloo-foo/cmd-awk"
)

// This example demonstrates conditional pattern matching from a file
func ExampleAwk_fromFile_condition() {
	// cat testdata/fruits.txt | awk '/ap/'
	patterns.MustRun(
		awk.Awk(
			conditionProgram{pattern: "ap"},
			gloo.File("testdata/fruits.txt"),
		),
	)
	// Output:
	// apple
	// apricot
}
