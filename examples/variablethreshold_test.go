package awk_test

import (
	"strconv"
	"strings"

	. "github.com/gloo-foo/cmd-awk"
	gloo "github.com/gloo-foo/framework/patterns"
)

// variableThresholdProgram demonstrates using initialized variables.
// Pass variables at initialization with Variable{Name, Value} to make
// your programs reusable with different parameters.
type variableThresholdProgram struct {
	SimpleProgram
}

func (p variableThresholdProgram) Condition(ctx *Context) bool {
	threshold := ctx.Var("threshold").(int)
	if val, err := strconv.Atoi(ctx.Field(1)); err == nil {
		return val > threshold
	}
	return false
}

func (p variableThresholdProgram) Action(ctx *Context) (string, bool) {
	return ctx.Field(0), true
}

func ExampleAwk_variableThreshold() {
	// echo -e "10\n25\n30\n15" | awk -v threshold=20 '$1>threshold'
	gloo.MustRun(
		Awk(
			variableThresholdProgram{},
			AwkVariable{Name: "threshold", Value: 20},
			strings.NewReader("10\n25\n30\n15"),
		),
	)
	// Output:
	// 25
	// 30
}
