package awk_test

import (
	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"

	awk "github.com/gloo-foo/cmd-awk"
)

// This example demonstrates adding line numbers to file content
// It reuses the lineNumberProgram defined in linenumbers_test.go
func ExampleAwk_fromFile_lineNumbers() {
	// cat testdata/fruits.txt | awk '{print NR": "$0}'
	if err := patterns.Run(
		awk.Awk(
			lineNumberProgram{},
			gloo.File("testdata/fruits.txt"),
		),
	); err != nil {
		panic(err)
	}
	// Output:
	// 1: apple
	// 2: banana
	// 3: apricot
	// 4: cherry
	// 5: orange
}
