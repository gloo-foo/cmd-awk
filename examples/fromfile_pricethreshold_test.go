package awk_test

import (
	"fmt"
	"strconv"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"

	awk "github.com/gloo-foo/cmd-awk"
)

// priceFilterProgram filters items by price threshold
type priceFilterProgram struct {
	awk.SimpleProgram
	threshold float64
}

func (p priceFilterProgram) Condition(ctx *awk.Context) bool {
	if price, err := strconv.ParseFloat(ctx.Field(2), 64); err == nil {
		return price >= p.threshold
	}
	return false
}

func (p priceFilterProgram) Action(ctx *awk.Context) (string, bool) {
	return fmt.Sprintf("%s costs $%s", ctx.Field(1), ctx.Field(2)), true
}

// This example demonstrates filtering by numeric threshold from a file
func ExampleAwk_fromFile_priceThreshold() {
	// cat testdata/prices.txt | awk '$2 >= 2.00 {print $1" costs $"$2}'
	patterns.MustRun(
		awk.Awk(
			priceFilterProgram{threshold: 2.00},
			gloo.File("testdata/prices.txt"),
		),
	)
	// Output:
	// grape costs $2.30
	// cherry costs $3.00
}
