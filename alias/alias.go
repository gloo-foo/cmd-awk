// Package alias provides unprefixed type aliases for awk command flags.
// This allows users to import and use shorter names:
//
//	import "github.com/gloo-foo/cmd-awk/alias"
//	awk.Awk(myProgram, alias.FieldSeparator(","))
package alias

import (
	gloo "github.com/gloo-foo/framework"

	command "github.com/gloo-foo/cmd-awk"
)

// Awk re-exports the command constructor.
func Awk(program command.Program, opts ...any) gloo.Command[[]byte, []byte] {
	return command.Awk(program, opts...)
}

// -F flag: field separator
type FieldSeparator = command.AwkFieldSeparator

// output field separator
type OutputFieldSeparator = command.AwkOutputFieldSeparator

// -v flag: variable
type Variable = command.AwkVariable
