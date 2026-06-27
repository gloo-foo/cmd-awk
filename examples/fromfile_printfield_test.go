package awk_test

import (
	. "github.com/gloo-foo/cmd-awk"
	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"
)

// This example demonstrates reading from a file instead of strings.NewReader
func ExampleAwk_fromFile_printField() {
	// cat testdata/simple_fields.txt | awk '{print $2}'
	patterns.MustRun(
		Awk(
			printFieldProgram{fieldNum: 2},
			gloo.File("testdata/simple_fields.txt"),
		),
	)
	// Output:
	// two
	// beta
	// second
}
