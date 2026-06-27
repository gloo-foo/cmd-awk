package awk_test

import (
	. "github.com/gloo-foo/cmd-awk"
	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"
)

// This example demonstrates conditional pattern matching from a file
func ExampleAwk_fromFile_condition() {
	// cat testdata/fruits.txt | awk '/ap/'
	patterns.MustRun(
		Awk(
			conditionProgram{pattern: "ap"},
			gloo.File("testdata/fruits.txt"),
		),
	)
	// Output:
	// apple
	// apricot
}
