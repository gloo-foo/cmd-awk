package command

type (
	AwkFieldSeparator       string
	AwkOutputFieldSeparator string
)

type AwkVariable struct {
	Value any
	Name  string
}

type flags struct {
	variables            map[string]any
	fieldSeparator       AwkFieldSeparator
	outputFieldSeparator AwkOutputFieldSeparator
}

func (f AwkFieldSeparator) Configure(flags *flags)       { flags.fieldSeparator = f }
func (o AwkOutputFieldSeparator) Configure(flags *flags) { flags.outputFieldSeparator = o }
func (v AwkVariable) Configure(flags *flags) {
	if flags.variables == nil {
		flags.variables = make(map[string]any)
	}
	flags.variables[v.Name] = v.Value
}
