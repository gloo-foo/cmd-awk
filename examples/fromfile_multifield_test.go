package awk_test

import (
	"fmt"
	"strconv"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"

	awk "github.com/gloo-foo/cmd-awk"
)

// studentScoreProgram calculates average score per student
type studentScoreProgram struct {
	awk.SimpleProgram
}

func (p studentScoreProgram) Action(ctx awk.Context) (awk.Context, string, bool) {
	name := ctx.Field(1)
	var sum float64
	var count int

	// Sum scores from fields 2, 3, 4
	for i := 2; i <= 4; i++ {
		if score, err := strconv.ParseFloat(ctx.Field(i), 64); err == nil {
			sum += score
			count++
		}
	}

	if count > 0 {
		avg := sum / float64(count)
		return ctx, fmt.Sprintf("%s: %.1f", name, avg), true
	}
	return ctx, "", false
}

// This example demonstrates calculating averages across multiple fields
func ExampleAwk_fromFile_studentAverage() {
	// cat testdata/scores.txt | awk '{print $1": "($2+$3+$4)/3}'
	if err := patterns.Run(
		awk.Awk(
			studentScoreProgram{},
			gloo.File("testdata/scores.txt"),
		),
	); err != nil {
		panic(err)
	}
	// Output:
	// Alice: 85.0
	// Bob: 91.0
	// Charlie: 82.3
	// Diana: 94.3
}
