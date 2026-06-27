package awk_test

import (
	. "github.com/gloo-foo/cmd-awk"
	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"
)

// This example demonstrates calculating sum from a file
func ExampleAwk_fromFile_sum() {
	// cat testdata/numbers.txt | awk 'BEGIN{sum=0} {sum+=$1} END{print "Sum:",sum}'
	patterns.MustRun(
		Awk(
			sumProgram{},
			gloo.File("testdata/numbers.txt"),
		),
	)
	// Output:
	// Sum: 100
}
