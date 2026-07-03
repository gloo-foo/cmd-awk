package command

type (
	// AwkFieldSeparator is the -F flag: the input field separator.
	AwkFieldSeparator string
	// AwkOutputFieldSeparator sets OFS: the output field separator.
	AwkOutputFieldSeparator string
)

// AwkVariable is the -v flag: a variable pre-set before the program runs.
type AwkVariable struct {
	Value any
	Name  string
}

type flags struct {
	variables            map[string]any
	fieldSeparator       AwkFieldSeparator
	outputFieldSeparator AwkOutputFieldSeparator
}

// fold partitions opts: awk's own option values are folded into the flag set,
// and every other argument is passed through unchanged for the framework's
// positional classifier.
func fold(opts []any) (flags, []any) {
	f := flags{variables: map[string]any{}}
	rest := make([]any, 0, len(opts))
	for _, o := range opts {
		switch v := o.(type) {
		case AwkFieldSeparator:
			f.fieldSeparator = v
		case AwkOutputFieldSeparator:
			f.outputFieldSeparator = v
		case AwkVariable:
			f.variables[v.Name] = v.Value
		default:
			rest = append(rest, o)
		}
	}
	return f, rest
}
