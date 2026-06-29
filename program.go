package command

import (
	"fmt"
	"strings"
)

// Context provides access to awk's execution context for each line
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
func (c *Context) Field(index int) string {
	if index < 0 || index >= len(c.Fields) {
		return ""
	}
	return c.Fields[index]
}

// SetField sets the value of a field
func (c *Context) SetField(index int, value string) {
	if index < 0 {
		return
	}
	// Expand fields if necessary
	for len(c.Fields) <= index {
		c.Fields = append(c.Fields, "")
	}
	c.Fields[index] = value
	c.NF = len(c.Fields) - 1 // Don't count $0
}

// Var returns a variable value
func (c *Context) Var(name string) any {
	if c.Variables == nil {
		return nil
	}
	return c.Variables[name]
}

// SetVar sets a variable value
func (c *Context) SetVar(name string, value any) {
	if c.Variables == nil {
		c.Variables = make(map[string]any)
	}
	c.Variables[name] = value
}

// Print formats and returns a string with fields separated by OFS
func (c *Context) Print(values ...any) string {
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
	// Use this for initialization
	Begin(ctx *Context) error

	// Condition is called for each line to determine if Action should run
	// Return true to run the action, false to skip
	Condition(ctx *Context) bool

	// Action is called for each line where Condition returns true
	// Return the output string and whether to emit it
	Action(ctx *Context) (output string, emit bool)

	// End is called once after processing all lines
	// Return any final output and an error if needed
	End(ctx *Context) (output string, err error)
}

// SimpleProgram provides default implementations for all Program methods
// Embed this in your program struct and override only what you need
type SimpleProgram struct{}

func (SimpleProgram) Begin(*Context) error               { return nil }
func (SimpleProgram) Condition(*Context) bool            { return true }
func (SimpleProgram) Action(ctx *Context) (string, bool) { return ctx.Field(0), true }
func (SimpleProgram) End(*Context) (string, error)       { return "", nil }
