package awk_test

import (
	. "github.com/gloo-foo/cmd-awk"
	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"
)

// This example demonstrates calculating average from a file
func ExampleAwk_fromFile_average() {
	// cat testdata/numbers.txt | awk '{sum+=$1;count++} END{print sum/count}'
	patterns.MustRun(
		Awk(
			averageProgram{},
			gloo.File("testdata/numbers.txt"),
		),
	)
	// Output:
	// Average: 25.00
}
