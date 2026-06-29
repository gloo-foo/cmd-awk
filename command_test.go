package command_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/testable"

	command "github.com/gloo-foo/cmd-awk"
)

// errBoom is the sentinel the failing test programs emit, so error propagation
// can be asserted with errors.Is rather than by message.
var errBoom = errors.New("boom")

// beginFails returns errBoom from its Begin phase.
type beginFails struct{ command.SimpleProgram }

func (beginFails) Begin(*command.Context) error { return errBoom }

// endFails returns errBoom from its End phase.
type endFails struct{ command.SimpleProgram }

func (endFails) End(*command.Context) (string, error) { return "", errBoom }

// printAll prints every line unchanged.
type printAll struct{ command.SimpleProgram }

func TestAwk_PrintAllLines(t *testing.T) {
	lines, err := testable.TestLines(command.Awk(printAll{}), "alpha\nbeta\ngamma\n")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"alpha", "beta", "gamma"}
	if len(lines) != len(want) {
		t.Fatalf("got %d lines, want %d", len(lines), len(want))
	}
	for i, w := range want {
		if lines[i] != w {
			t.Errorf("line %d: got %q, want %q", i, lines[i], w)
		}
	}
}

// patternMatch only emits lines containing "err".
type patternMatch struct{ command.SimpleProgram }

func (patternMatch) Condition(ctx *command.Context) bool {
	return strings.Contains(ctx.Field(0), "err")
}

func TestAwk_PatternMatch(t *testing.T) {
	input := "info: ok\nerror: bad\nwarn: hmm\nerr: also bad\n"
	lines, err := testable.TestLines(command.Awk(patternMatch{}), input)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"error: bad", "err: also bad"}
	if len(lines) != len(want) {
		t.Fatalf("got %d lines, want %d: %q", len(lines), len(want), lines)
	}
	for i, w := range want {
		if lines[i] != w {
			t.Errorf("line %d: got %q, want %q", i, lines[i], w)
		}
	}
}

// fieldPrint prints $2 then $1 separated by OFS.
type fieldPrint struct{ command.SimpleProgram }

func (fieldPrint) Action(ctx *command.Context) (string, bool) {
	return ctx.Print(ctx.Field(2), ctx.Field(1)), true
}

func TestAwk_FieldSplitting(t *testing.T) {
	input := "alice 30\nbob 25\n"
	lines, err := testable.TestLines(command.Awk(fieldPrint{}), input)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"30 alice", "25 bob"}
	if len(lines) != len(want) {
		t.Fatalf("got %d lines, want %d: %q", len(lines), len(want), lines)
	}
	for i, w := range want {
		if lines[i] != w {
			t.Errorf("line %d: got %q, want %q", i, lines[i], w)
		}
	}
}

// lineCounter uses BEGIN to init, Action to count, END to emit the total.
type lineCounter struct{ command.SimpleProgram }

func (lineCounter) Begin(ctx *command.Context) error {
	ctx.SetVar("count", 0)
	return nil
}

func (lineCounter) Action(ctx *command.Context) (string, bool) {
	ctx.SetVar("count", ctx.Var("count").(int)+1)
	return "", false
}

func (lineCounter) End(ctx *command.Context) (string, error) {
	return fmt.Sprintf("%d", ctx.Var("count").(int)), nil
}

func TestAwk_BeginEndBlocks(t *testing.T) {
	input := "one\ntwo\nthree\n"
	lines, err := testable.TestLines(command.Awk(lineCounter{}), input)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0] != "3" {
		t.Fatalf("got %q, want [3]", lines)
	}
}

func TestAwk_EmptyInput(t *testing.T) {
	lines, err := testable.TestLines(command.Awk(printAll{}), "")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 0 {
		t.Fatalf("got %q, want empty", lines)
	}
}

func TestAwk_FieldSeparatorOption(t *testing.T) {
	input := "alice,30\nbob,25\n"
	lines, err := testable.TestLines(
		command.Awk(fieldPrint{}, command.AwkFieldSeparator(",")),
		input,
	)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"30 alice", "25 bob"}
	if len(lines) != len(want) {
		t.Fatalf("got %d lines, want %d: %q", len(lines), len(want), lines)
	}
	for i, w := range want {
		if lines[i] != w {
			t.Errorf("line %d: got %q, want %q", i, lines[i], w)
		}
	}
}

func TestAwk_BeginErrorPropagates(t *testing.T) {
	_, err := testable.TestLines(command.Awk(beginFails{}), "one\ntwo\n")
	if !errors.Is(err, errBoom) {
		t.Fatalf("got err %v, want errBoom", err)
	}
}

func TestAwk_EndErrorPropagates(t *testing.T) {
	_, err := testable.TestLines(command.Awk(endFails{}), "one\ntwo\n")
	if !errors.Is(err, errBoom) {
		t.Fatalf("got err %v, want errBoom", err)
	}
}

func TestAwk_MissingFilePropagatesError(t *testing.T) {
	// A self-sourced run over a nonexistent file surfaces the open error.
	_, err := testable.TestLines(
		command.Awk(printAll{}, gloo.File("does-not-exist.txt")),
		"",
	)
	if err == nil {
		t.Fatal("got nil error, want a file-open error")
	}
}

func TestAwk_SelfSourcedFromReader(t *testing.T) {
	// A positional io.Reader supplies the data, overriding the (empty) upstream.
	lines, err := testable.TestLines(
		command.Awk(printAll{}, strings.NewReader("x\ny\n")),
		"",
	)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"x", "y"}
	if len(lines) != len(want) {
		t.Fatalf("got %d lines, want %d: %q", len(lines), len(want), lines)
	}
	for i, w := range want {
		if lines[i] != w {
			t.Errorf("line %d: got %q, want %q", i, lines[i], w)
		}
	}
}

func TestAwk_EmptyLineWithExplicitSeparator(t *testing.T) {
	// With an explicit FS, a blank line yields zero fields ($0 is still ""),
	// distinct from the default-separator path which also yields zero fields.
	lines, err := testable.TestLines(
		command.Awk(fieldCount{}, command.AwkFieldSeparator(",")),
		"a,b\n\nc\n",
	)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"2", "0", "1"}
	if len(lines) != len(want) {
		t.Fatalf("got %d lines, want %d: %q", len(lines), len(want), lines)
	}
	for i, w := range want {
		if lines[i] != w {
			t.Errorf("line %d: got %q, want %q", i, lines[i], w)
		}
	}
}

// fieldCount emits NF (the number of fields) for each record.
type fieldCount struct{ command.SimpleProgram }

func (fieldCount) Action(ctx *command.Context) (string, bool) {
	return fmt.Sprintf("%d", ctx.NF), true
}

func ExampleAwk() {
	lines, _ := testable.TestLines(command.Awk(printAll{}), "hello\nworld\n")
	for _, line := range lines {
		fmt.Println(line)
	}
	// Output:
	// hello
	// world
}

func ExampleAwk_patternMatch() {
	lines, _ := testable.TestLines(command.Awk(patternMatch{}), "info\nerror\nwarn\n")
	for _, line := range lines {
		fmt.Println(line)
	}
	// Output:
	// error
}
