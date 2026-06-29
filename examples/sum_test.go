package awk_test

import (
	"fmt"
	"strconv"
	"strings"

	gloo "github.com/gloo-foo/framework/patterns"

	awk "github.com/gloo-foo/cmd-awk"
)

// sumProgram demonstrates accumulating values with BEGIN and END.
// This pattern is common for computing aggregates: initialize in Begin,
// accumulate in Action (without emitting output), then output the result in End.
type sumProgram struct {
	awk.SimpleProgram
}

func (p sumProgram) Begin(ctx *awk.Context) error {
	ctx.SetVar("sum", 0)
	return nil
}

func (p sumProgram) Action(ctx *awk.Context) (string, bool) {
	if val, err := strconv.Atoi(ctx.Field(1)); err == nil {
		sum := ctx.Var("sum").(int)
		ctx.SetVar("sum", sum+val)
	}
	return "", false // Don't emit output per line
}

func (p sumProgram) End(ctx *awk.Context) (string, error) {
	return fmt.Sprintf("Sum: %d", ctx.Var("sum")), nil
}

func ExampleAwk_sum() {
	// echo -e "10\n20\n30" | awk 'BEGIN{sum=0} {sum+=$1} END{print "Sum:",sum}'
	gloo.MustRun(
		awk.Awk(
			sumProgram{},
			strings.NewReader("10\n20\n30"),
		),
	)
	// Output:
	// Sum: 60
}
