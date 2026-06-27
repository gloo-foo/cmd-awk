// Package alias provides unprefixed type aliases for awk command flags.
// This allows users to import and use shorter names:
//
//	import "github.com/gloo-foo/cmd-awk/alias"
//	awk.Awk(myProgram, alias.FieldSeparator(","))
package alias

import command "github.com/gloo-foo/cmd-awk"

// Awk is the command constructor.
var Awk = command.Awk

// -F flag: field separator
type FieldSeparator = command.AwkFieldSeparator

// output field separator
type OutputFieldSeparator = command.AwkOutputFieldSeparator

// -v flag: variable
type Variable = command.AwkVariable
