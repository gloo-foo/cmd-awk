package awk_test

import (
	"fmt"
	"strconv"
	"strings"

	. "github.com/gloo-foo/cmd-awk"
	gloo "github.com/gloo-foo/framework/patterns"
)

// averageProgram demonstrates computing average.
// This shows the full pattern for statistical analysis: initialize counters,
// accumulate values while processing, then compute and output the result.
type averageProgram struct {
	SimpleProgram
}

func (p averageProgram) Begin(ctx *Context) error {
	ctx.SetVar("sum", 0.0)
	ctx.SetVar("count", 0)
	return nil
}

func (p averageProgram) Action(ctx *Context) (string, bool) {
	if val, err := strconv.ParseFloat(ctx.Field(1), 64); err == nil {
		sum := ctx.Var("sum").(float64)
		count := ctx.Var("count").(int)
		ctx.SetVar("sum", sum+val)
		ctx.SetVar("count", count+1)
	}
	return "", false
}

func (p averageProgram) End(ctx *Context) (string, error) {
	sum := ctx.Var("sum").(float64)
	count := ctx.Var("count").(int)
	if count > 0 {
		return fmt.Sprintf("Average: %.2f", sum/float64(count)), nil
	}
	return "Average: 0.00", nil
}

func ExampleAwk_average() {
	// echo -e "10\n20\n30\n40" | awk '{sum+=$1;count++} END{print sum/count}'
	gloo.MustRun(
		Awk(
			averageProgram{},
			strings.NewReader("10\n20\n30\n40"),
		),
	)
	// Output:
	// Average: 25.00
}
