package awk_test

import (
	"fmt"
	"strconv"

	. "github.com/gloo-foo/cmd-awk"
	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"
)

// studentScoreProgram calculates average score per student
type studentScoreProgram struct {
	SimpleProgram
}

func (p studentScoreProgram) Action(ctx *Context) (string, bool) {
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
		return fmt.Sprintf("%s: %.1f", name, avg), true
	}
	return "", false
}

// This example demonstrates calculating averages across multiple fields
func ExampleAwk_fromFile_studentAverage() {
	// cat testdata/scores.txt | awk '{print $1": "($2+$3+$4)/3}'
	patterns.MustRun(
		Awk(
			studentScoreProgram{},
			gloo.File("testdata/scores.txt"),
		),
	)
	// Output:
	// Alice: 85.0
	// Bob: 91.0
	// Charlie: 82.3
	// Diana: 94.3
}
