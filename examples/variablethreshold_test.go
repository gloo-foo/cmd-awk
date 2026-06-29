package awk_test

import (
	"strconv"
	"strings"

	gloo "github.com/gloo-foo/framework/patterns"

	awk "github.com/gloo-foo/cmd-awk"
)

// variableThresholdProgram demonstrates using initialized variables.
// Pass variables at initialization with Variable{Name, Value} to make
// your programs reusable with different parameters.
type variableThresholdProgram struct {
	awk.SimpleProgram
}

func (p variableThresholdProgram) Condition(ctx *awk.Context) bool {
	threshold := ctx.Var("threshold").(int)
	if val, err := strconv.Atoi(ctx.Field(1)); err == nil {
		return val > threshold
	}
	return false
}

func (p variableThresholdProgram) Action(ctx *awk.Context) (string, bool) {
	return ctx.Field(0), true
}

func ExampleAwk_variableThreshold() {
	// echo -e "10\n25\n30\n15" | awk -v threshold=20 '$1>threshold'
	gloo.MustRun(
		awk.Awk(
			variableThresholdProgram{},
			awk.AwkVariable{Name: "threshold", Value: 20},
			strings.NewReader("10\n25\n30\n15"),
		),
	)
	// Output:
	// 25
	// 30
}
