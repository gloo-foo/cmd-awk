package alias_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/gloo-foo/testable"

	awk "github.com/gloo-foo/cmd-awk"
	alias "github.com/gloo-foo/cmd-awk/alias"
)

// The alias package re-exports the constructor and flag types under unprefixed
// names. A mis-wired re-export (say, OutputFieldSeparator aliased to the input
// separator type, or Awk bound to the wrong function) compiles cleanly, so only
// behavior can prove the wiring. Each test drives one re-export through the
// constructor and asserts the exact output it must produce.

func assertLines(t *testing.T, got, want []string) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

// swapFields prints $2 then $1, joined by OFS — observing both FS (the split)
// and OFS (the join).
type swapFields struct{ awk.SimpleProgram }

func (swapFields) Action(ctx awk.Context) (awk.Context, string, bool) {
	return ctx, ctx.Print(ctx.Field(2), ctx.Field(1)), true
}

// printVar emits the value of the user variable named "tag" for each line,
// observing the -v re-export.
type printVar struct{ awk.SimpleProgram }

func (printVar) Action(ctx awk.Context) (awk.Context, string, bool) {
	return ctx, fmt.Sprint(ctx.Var("tag")), true
}

func TestAlias_AwkConstructorPassesThrough(t *testing.T) {
	// alias.Awk must be the real constructor: a default program echoes input.
	lines, err := testable.TestLines(alias.Awk(awk.SimpleProgram{}), "alpha\nbeta\n")
	if err != nil {
		t.Fatal(err)
	}
	assertLines(t, lines, []string{"alpha", "beta"})
}

func TestAlias_FieldSeparatorSplitsOnComma(t *testing.T) {
	// alias.FieldSeparator(",") must split records on commas before $1/$2 read.
	lines, err := testable.TestLines(
		alias.Awk(swapFields{}, alias.FieldSeparator(",")),
		"alice,30\nbob,25\n",
	)
	if err != nil {
		t.Fatal(err)
	}
	assertLines(t, lines, []string{"30 alice", "25 bob"})
}

func TestAlias_OutputFieldSeparatorJoinsWithPipe(t *testing.T) {
	// alias.OutputFieldSeparator("|") must set OFS, changing how Print joins.
	lines, err := testable.TestLines(
		alias.Awk(swapFields{}, alias.OutputFieldSeparator("|")),
		"alice 30\nbob 25\n",
	)
	if err != nil {
		t.Fatal(err)
	}
	assertLines(t, lines, []string{"30|alice", "25|bob"})
}

func TestAlias_VariableIsReadableByProgram(t *testing.T) {
	// alias.Variable must seed ctx.Var with the given name/value.
	lines, err := testable.TestLines(
		alias.Awk(printVar{}, alias.Variable{Name: "tag", Value: "X"}),
		"one\ntwo\n",
	)
	if err != nil {
		t.Fatal(err)
	}
	assertLines(t, lines, []string{"X", "X"})
}
