package command

import (
	"context"
	"io"
	"strings"

	"github.com/destel/rill"
	gloo "github.com/gloo-foo/framework"
)

// awkConfig is the immutable per-construction configuration of an Awk command:
// the resolved field separators and a private copy of the initial variables.
type awkConfig struct {
	vars map[string]any
	fs   string
	ofs  string
}

// defaultSeparator is awk's default field/output separator (a single space).
const defaultSeparator = " "

// newConfig resolves the separator defaults and snapshots the variables so the
// returned config is independent of the caller's flags.
func newConfig(f flags) awkConfig {
	return awkConfig{
		fs:   orDefault(string(f.fieldSeparator)),
		ofs:  orDefault(string(f.outputFieldSeparator)),
		vars: copyVars(f.variables),
	}
}

// orDefault substitutes awk's default separator when sep is empty.
func orDefault(sep string) string {
	if sep == "" {
		return defaultSeparator
	}
	return sep
}

// copyVars returns a fresh map holding the same entries as src.
func copyVars(src map[string]any) map[string]any {
	vars := make(map[string]any, len(src))
	for k, v := range src {
		vars[k] = v
	}
	return vars
}

// newContext builds the per-execution awk context from a config.
func (cfg awkConfig) newContext() Context {
	return Context{
		Fields:    make([]string, 0, 16),
		NR:        0,
		NF:        0,
		FS:        cfg.fs,
		OFS:       cfg.ofs,
		RS:        "\n",
		Variables: copyVars(cfg.vars),
	}
}

// Awk returns a command that processes input lines through an awk-style program.
// The program defines Begin, Condition, Action, and End phases.
//
// If opts contains positional file paths or io.Reader inputs, those are used as
// the data source, overriding the upstream input stream. This lets callers do
// either Awk(prog) (transform an input stream) or Awk(prog, file...) /
// Awk(prog, reader) (self-source from those inputs).
func Awk(program Program, opts ...any) gloo.Command[[]byte, []byte] {
	params := gloo.NewParameters[gloo.File, flags](opts...)
	cfg := newConfig(params.Flags)
	selfSourced := len(params.Positional) > 0
	return gloo.FuncCommand[[]byte, []byte](func(ctx context.Context, input gloo.Stream[[]byte]) gloo.Stream[[]byte] {
		return run(ctx, program, cfg, params, selfSourced, input)
	})
}

// run executes one pipeline pass: it resolves the input source, runs the Begin
// phase, maps each line through the program, and appends any End output.
func run(
	ctx context.Context,
	program Program,
	cfg awkConfig,
	params gloo.Parameters[gloo.File, flags],
	selfSourced bool,
	input gloo.Stream[[]byte],
) gloo.Stream[[]byte] {
	input, err := resolveSource(ctx, params, selfSourced, input)
	if err != nil {
		return errStream(ctx, input, err)
	}
	state := cfg.newContext()
	if err := program.Begin(&state); err != nil {
		return errStream(ctx, input, err)
	}
	processed := rill.OrderedFilterMap(input.Chan(), 1, mapLine(program, &state))
	return gloo.WrapFrom(drain(processed, endPhase(program, &state)), input)
}

// resolveSource returns the stream the program should read from. When the
// command is self-sourced from positional inputs it replaces the upstream
// stream with a reader over those inputs; otherwise the upstream is returned
// unchanged.
func resolveSource(
	ctx context.Context,
	params gloo.Parameters[gloo.File, flags],
	selfSourced bool,
	input gloo.Stream[[]byte],
) (gloo.Stream[[]byte], error) {
	if !selfSourced {
		return input, nil
	}
	reader, err := params.Reader(nil)
	if err != nil {
		return input, err
	}
	return gloo.ByteReaderSource([]io.Reader{reader}).Stream(ctx), nil
}

// errStream yields a stream that emits a single error, tearing down upstream.
func errStream(ctx context.Context, input gloo.Stream[[]byte], err error) gloo.Stream[[]byte] {
	return gloo.GenerateFrom(ctx, input, func(_ context.Context, _ func([]byte) bool, sendErr func(error)) {
		sendErr(err)
	})
}

// mapLine returns the per-line transform passed to OrderedFilterMap: it loads
// the record into state, evaluates the program's condition, and emits the
// action's output when both the condition and the action opt to emit.
func mapLine(program Program, state *Context) func([]byte) ([]byte, bool, error) {
	return func(line []byte) ([]byte, bool, error) {
		loadRecord(state, string(line))
		if !program.Condition(state) {
			return nil, false, nil
		}
		output, emit := program.Action(state)
		if !emit {
			return nil, false, nil
		}
		return []byte(output), true, nil
	}
}

// loadRecord advances the record counter and populates state's fields for line.
func loadRecord(state *Context, line string) {
	state.NR++
	fields := splitFields(state.FS, line)
	state.Fields = append(state.Fields[:0], line) // $0
	state.Fields = append(state.Fields, fields...)
	state.NF = len(fields)
}

// splitFields splits line into fields per awk's separator rules: the default
// separator collapses runs of whitespace, any other separator splits literally.
func splitFields(fs, line string) []string {
	if fs == defaultSeparator {
		return strings.Fields(line)
	}
	if line == "" {
		return []string{}
	}
	return strings.Split(line, fs)
}

// tailFunc produces the optional trailing record drain appends after the body.
// It reports whether a record should be emitted; a false ok suppresses it.
type tailFunc func() (record rill.Try[[]byte], ok bool)

// drain forwards every mapped record and then appends tail's record, if any.
func drain(processed <-chan rill.Try[[]byte], tail tailFunc) <-chan rill.Try[[]byte] {
	out := make(chan rill.Try[[]byte])
	go func() {
		defer close(out)
		for item := range processed {
			out <- item
		}
		if record, ok := tail(); ok {
			out <- record
		}
	}()
	return out
}

// endPhase wraps the program's End phase as a tailFunc. An error is always
// emitted; otherwise output is emitted only when it is non-empty, matching
// awk's END block which appends nothing when it prints nothing.
func endPhase(program Program, state *Context) tailFunc {
	return func() (rill.Try[[]byte], bool) {
		output, err := program.End(state)
		if err != nil {
			return rill.Wrap([]byte(nil), err), true
		}
		if output == "" {
			return rill.Try[[]byte]{}, false
		}
		return rill.Wrap([]byte(output), nil), true
	}
}
