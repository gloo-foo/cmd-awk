package command

import (
	"fmt"
	"strings"
)

// Context provides access to awk's execution context for each line. It is an
// immutable value: the Set… methods return an updated copy rather than
// mutating in place, and programs hand the updated context back to the
// pipeline (Begin and Action return it).
type Context struct {
	Variables map[string]any
	FS        string
	OFS       string
	RS        string
	Fields    []string
	NR        int64
	NF        int
}

// Field returns the field at the given index (0 = whole line, 1 = first field, etc.)
func (c Context) Field(index int) string {
	if index < 0 || index >= len(c.Fields) {
		return ""
	}
	return c.Fields[index]
}

// SetField returns a copy of the context with the field at index set, growing
// the record as needed. A negative index returns the context unchanged.
func (c Context) SetField(index int, value string) Context {
	if index < 0 {
		return c
	}
	fields := make([]string, max(len(c.Fields), index+1))
	copy(fields, c.Fields)
	fields[index] = value
	c.Fields = fields
	c.NF = len(fields) - 1 // Don't count $0
	return c
}

// Var returns a variable value
func (c Context) Var(name string) any {
	if c.Variables == nil {
		return nil
	}
	return c.Variables[name]
}

// SetVar returns a copy of the context with the variable set.
func (c Context) SetVar(name string, value any) Context {
	vars := make(map[string]any, len(c.Variables)+1)
	for k, v := range c.Variables {
		vars[k] = v
	}
	vars[name] = value
	c.Variables = vars
	return c
}

// Print formats and returns a string with fields separated by OFS
func (c Context) Print(values ...any) string {
	parts := make([]string, len(values))
	for i, v := range values {
		parts[i] = fmt.Sprint(v)
	}
	return strings.Join(parts, c.OFS)
}

// Program defines the interface for awk-style programs
// all methods are optional - implement only what you need
type Program interface {
	// Begin is called once before processing any lines
	// Use this for initialization; return the (possibly updated) context
	Begin(ctx Context) (Context, error)

	// Condition is called for each line to determine if Action should run
	// Return true to run the action, false to skip
	Condition(ctx Context) bool

	// Action is called for each line where Condition returns true
	// Return the updated context, the output string, and whether to emit it
	Action(ctx Context) (updated Context, output string, isEmit bool)

	// End is called once after processing all lines
	// Return any final output and an error if needed
	End(ctx Context) (output string, err error)
}

// SimpleProgram provides default implementations for all Program methods
// Embed this in your program struct and override only what you need
type SimpleProgram struct{}

func (SimpleProgram) Begin(ctx Context) (Context, error) { return ctx, nil }
func (SimpleProgram) Condition(Context) bool             { return true }
func (SimpleProgram) Action(ctx Context) (Context, string, bool) {
	return ctx, ctx.Field(0), true
}
func (SimpleProgram) End(Context) (string, error) { return "", nil }
