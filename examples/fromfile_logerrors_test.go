package awk_test

import (
	"strings"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"

	awk "github.com/gloo-foo/cmd-awk"
)

// logErrorFilterProgram filters log entries by level
type logErrorFilterProgram struct {
	awk.SimpleProgram
}

func (p logErrorFilterProgram) Condition(ctx *awk.Context) bool {
	return strings.Contains(ctx.Field(0), "ERROR")
}

func (p logErrorFilterProgram) Action(ctx *awk.Context) (string, bool) {
	return ctx.Field(0), true
}

// This example demonstrates filtering log entries from a file
func ExampleAwk_fromFile_logErrors() {
	// cat testdata/log_entries.txt | awk '/ERROR/'
	patterns.MustRun(
		awk.Awk(
			logErrorFilterProgram{},
			gloo.File("testdata/log_entries.txt"),
		),
	)
	// Output:
	// 2024-01-15 10:24:12 ERROR Connection failed
	// 2024-01-15 10:25:45 ERROR Database timeout
}
