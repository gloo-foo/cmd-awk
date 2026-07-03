package command_test

import (
	"testing"

	command "github.com/gloo-foo/cmd-awk"
)

// The Context accessor methods carry awk's record contract: $0 is the whole
// line, out-of-range fields read as empty, SetField grows the record, and
// SetVar seeds the variable map. The Set… methods return an updated copy and
// leave the receiver untouched. These unit tests exercise each branch of that
// contract directly, since the pipeline-level tests only reach the common path.

func TestContext_FieldOutOfRange(t *testing.T) {
	ctx := command.Context{Fields: []string{"whole", "a", "b"}}
	if got := ctx.Field(0); got != "whole" {
		t.Errorf("Field(0): got %q, want %q", got, "whole")
	}
	if got := ctx.Field(2); got != "b" {
		t.Errorf("Field(2): got %q, want %q", got, "b")
	}
	if got := ctx.Field(99); got != "" {
		t.Errorf("Field(99): got %q, want empty", got)
	}
	if got := ctx.Field(-1); got != "" {
		t.Errorf("Field(-1): got %q, want empty", got)
	}
}

func TestContext_SetFieldNegativeIsNoOp(t *testing.T) {
	ctx := command.Context{Fields: []string{"whole", "a"}}
	updated := ctx.SetField(-1, "ignored")
	if got := updated.Field(1); got != "a" {
		t.Errorf("Field(1): got %q, want %q", got, "a")
	}
	if updated.NF != 0 {
		t.Errorf("NF: got %d, want 0 (unchanged)", updated.NF)
	}
}

func TestContext_SetFieldGrowsRecord(t *testing.T) {
	ctx := command.Context{Fields: []string{"whole", "a"}}
	updated := ctx.SetField(3, "c")
	// Fields grew to [whole, a, "", c]; NF counts fields excluding $0.
	if got := updated.Field(2); got != "" {
		t.Errorf("Field(2): got %q, want empty (gap-filled)", got)
	}
	if got := updated.Field(3); got != "c" {
		t.Errorf("Field(3): got %q, want %q", got, "c")
	}
	if updated.NF != 3 {
		t.Errorf("NF: got %d, want 3", updated.NF)
	}
	if len(ctx.Fields) != 2 {
		t.Errorf("receiver Fields: got %d entries, want 2 (unchanged)", len(ctx.Fields))
	}
}

func TestContext_SetFieldReplacesInPlace(t *testing.T) {
	ctx := command.Context{Fields: []string{"whole", "a", "b"}}
	updated := ctx.SetField(1, "z")
	if got := updated.Field(1); got != "z" {
		t.Errorf("Field(1): got %q, want %q", got, "z")
	}
	if got := ctx.Field(1); got != "a" {
		t.Errorf("receiver Field(1): got %q, want %q (unchanged)", got, "a")
	}
}

func TestContext_VarNilMapReadsNil(t *testing.T) {
	ctx := command.Context{}
	if got := ctx.Var("missing"); got != nil {
		t.Errorf("Var on nil map: got %v, want nil", got)
	}
}

func TestContext_SetVarSeedsNilMapThenRead(t *testing.T) {
	ctx := command.Context{}
	updated := ctx.SetVar("k", 42)
	if got := updated.Var("k"); got != 42 {
		t.Errorf("Var(k): got %v, want 42", got)
	}
	if got := ctx.Var("k"); got != nil {
		t.Errorf("receiver Var(k): got %v, want nil (unchanged)", got)
	}
}

func TestContext_SetVarExistingMap(t *testing.T) {
	ctx := command.Context{Variables: map[string]any{"k": 1}}
	updated := ctx.SetVar("k", 2)
	if got := updated.Var("k"); got != 2 {
		t.Errorf("Var(k): got %v, want 2", got)
	}
	if got := ctx.Var("k"); got != 1 {
		t.Errorf("receiver Var(k): got %v, want 1 (unchanged)", got)
	}
}

func TestContext_PrintJoinsWithOFS(t *testing.T) {
	ctx := command.Context{OFS: "-"}
	if got := ctx.Print("a", 1, "b"); got != "a-1-b" {
		t.Errorf("Print: got %q, want %q", got, "a-1-b")
	}
}
