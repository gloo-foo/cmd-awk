package awk_test

import (
	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"

	awk "github.com/gloo-foo/cmd-awk"
)

// This example demonstrates reading from a file instead of strings.NewReader
func ExampleAwk_fromFile_printField() {
	// cat testdata/simple_fields.txt | awk '{print $2}'
	if err := patterns.Run(
		awk.Awk(
			printFieldProgram{fieldNum: 2},
			gloo.File("testdata/simple_fields.txt"),
		),
	); err != nil {
		panic(err)
	}
	// Output:
	// two
	// beta
	// second
}
