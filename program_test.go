package command_test

import (
	"testing"

	command "github.com/gloo-foo/cmd-awk"
)

// The Context accessor methods carry awk's record contract: $0 is the whole
// line, out-of-range fields read as empty, SetField grows the record, and the
// variable map is lazily created. These unit tests exercise each branch of that
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
	ctx.SetField(-1, "ignored")
	if got := ctx.Field(1); got != "a" {
		t.Errorf("Field(1): got %q, want %q", got, "a")
	}
	if ctx.NF != 0 {
		t.Errorf("NF: got %d, want 0 (unchanged)", ctx.NF)
	}
}

func TestContext_SetFieldGrowsRecord(t *testing.T) {
	ctx := command.Context{Fields: []string{"whole", "a"}}
	ctx.SetField(3, "c")
	// Fields grew to [whole, a, "", c]; NF counts fields excluding $0.
	if got := ctx.Field(2); got != "" {
		t.Errorf("Field(2): got %q, want empty (gap-filled)", got)
	}
	if got := ctx.Field(3); got != "c" {
		t.Errorf("Field(3): got %q, want %q", got, "c")
	}
	if ctx.NF != 3 {
		t.Errorf("NF: got %d, want 3", ctx.NF)
	}
}

func TestContext_VarNilMapReadsNil(t *testing.T) {
	ctx := command.Context{}
	if got := ctx.Var("missing"); got != nil {
		t.Errorf("Var on nil map: got %v, want nil", got)
	}
}

func TestContext_SetVarLazyInitThenRead(t *testing.T) {
	ctx := command.Context{}
	ctx.SetVar("k", 42)
	if got := ctx.Var("k"); got != 42 {
		t.Errorf("Var(k): got %v, want 42", got)
	}
}

func TestContext_SetVarExistingMap(t *testing.T) {
	ctx := command.Context{Variables: map[string]any{"k": 1}}
	ctx.SetVar("k", 2)
	if got := ctx.Var("k"); got != 2 {
		t.Errorf("Var(k): got %v, want 2", got)
	}
}

func TestContext_PrintJoinsWithOFS(t *testing.T) {
	ctx := command.Context{OFS: "-"}
	if got := ctx.Print("a", 1, "b"); got != "a-1-b" {
		t.Errorf("Print: got %q, want %q", got, "a-1-b")
	}
}
