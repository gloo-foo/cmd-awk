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

// separator is a field or output separator value as awk resolves it.
type separator string

// record is a single input line, awk's unit of processing.
type record string

// selfSourced reports whether the command reads from its own positional
// inputs instead of the upstream stream.
type selfSourced bool

// newConfig resolves the separator defaults and snapshots the variables so the
// returned config is independent of the caller's flags.
func newConfig(f flags) awkConfig {
	return awkConfig{
		fs:   orDefault(separator(f.fieldSeparator)),
		ofs:  orDefault(separator(f.outputFieldSeparator)),
		vars: copyVars(f.variables),
	}
}

// orDefault substitutes awk's default separator when sep is empty.
func orDefault(sep separator) string {
	if sep == "" {
		return defaultSeparator
	}
	return string(sep)
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
	f, rest := fold(opts)
	params := gloo.NewParameters[gloo.File, struct{}](rest...)
	cfg := newConfig(f)
	isSelfSourced := selfSourced(len(params.Positional) > 0)
	return gloo.FuncCommand[[]byte, []byte](func(ctx context.Context, input gloo.Stream[[]byte]) gloo.Stream[[]byte] {
		return run(ctx, program, cfg, params, isSelfSourced, input)
	})
}

// run executes one pipeline pass: it resolves the input source, runs the Begin
// phase, maps each line through the program, and appends any End output. The
// context value threads through the sequential per-line steps and into the End
// phase; the channel close that ends the mapped stream orders those writes
// before drain's tail call reads the final state.
func run(
	ctx context.Context,
	program Program,
	cfg awkConfig,
	params gloo.Parameters[gloo.File, struct{}],
	isSelfSourced selfSourced,
	input gloo.Stream[[]byte],
) gloo.Stream[[]byte] {
	input, err := resolveSource(ctx, params, isSelfSourced, input)
	if err != nil {
		return errStream(ctx, input, err)
	}
	state, err := program.Begin(cfg.newContext())
	if err != nil {
		return errStream(ctx, input, err)
	}
	mapper := func(line []byte) ([]byte, bool, error) {
		var output []byte
		var isEmit bool
		state, output, isEmit = step(program, state, record(line))
		return output, isEmit, nil
	}
	processed := rill.OrderedFilterMap(input.Chan(), 1, mapper)
	tail := func() (rill.Try[[]byte], bool) { return endRecord(program, state) }
	return gloo.WrapFrom(drain(processed, tail), input)
}

// resolveSource returns the stream the program should read from. When the
// command is self-sourced from positional inputs it replaces the upstream
// stream with a reader over those inputs; otherwise the upstream is returned
// unchanged.
func resolveSource(
	ctx context.Context,
	params gloo.Parameters[gloo.File, struct{}],
	isSelfSourced selfSourced,
	input gloo.Stream[[]byte],
) (gloo.Stream[[]byte], error) {
	if !bool(isSelfSourced) {
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

// step loads line into state as the next record and applies the program's
// condition and action, returning the updated context, the action's output,
// and whether to emit it.
func step(program Program, state Context, line record) (Context, []byte, bool) {
	state = loadRecord(state, line)
	if !program.Condition(state) {
		return state, nil, false
	}
	updated, output, isEmit := program.Action(state)
	return updated, []byte(output), isEmit
}

// loadRecord advances the record counter and populates the context's fields
// for line, returning the updated context.
func loadRecord(state Context, line record) Context {
	state.NR++
	fields := splitFields(separator(state.FS), line)
	state.Fields = append(state.Fields[:0], string(line)) // $0
	state.Fields = append(state.Fields, fields...)
	state.NF = len(fields)
	return state
}

// splitFields splits line into fields per awk's separator rules: the default
// separator collapses runs of whitespace, any other separator splits literally.
func splitFields(fs separator, line record) []string {
	if fs == defaultSeparator {
		return strings.Fields(string(line))
	}
	if line == "" {
		return []string{}
	}
	return strings.Split(string(line), string(fs))
}

// tailFunc produces the optional trailing record drain appends after the body.
// It reports whether a record should be emitted; a false ok suppresses it.
type tailFunc func() (record rill.Try[[]byte], isOk bool)

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

// endRecord wraps the program's End phase as drain's tail record. An error is
// always emitted; otherwise output is emitted only when it is non-empty,
// matching awk's END block which appends nothing when it prints nothing.
func endRecord(program Program, state Context) (rill.Try[[]byte], bool) {
	output, err := program.End(state)
	if err != nil {
		return rill.Wrap([]byte(nil), err), true
	}
	if output == "" {
		return rill.Try[[]byte]{}, false
	}
	return rill.Wrap([]byte(output), nil), true
}
