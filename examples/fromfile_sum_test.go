package awk_test

import (
	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"

	awk "github.com/gloo-foo/cmd-awk"
)

// This example demonstrates calculating sum from a file
func ExampleAwk_fromFile_sum() {
	// cat testdata/numbers.txt | awk 'BEGIN{sum=0} {sum+=$1} END{print "Sum:",sum}'
	patterns.MustRun(
		awk.Awk(
			sumProgram{},
			gloo.File("testdata/numbers.txt"),
		),
	)
	// Output:
	// Sum: 100
}
