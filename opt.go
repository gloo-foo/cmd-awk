package command

type (
	AwkFieldSeparator       string
	AwkOutputFieldSeparator string
)

type AwkVariable struct {
	Name  string
	Value any
}

type flags struct {
	fieldSeparator       AwkFieldSeparator
	outputFieldSeparator AwkOutputFieldSeparator
	variables            map[string]any
}

func (f AwkFieldSeparator) Configure(flags *flags)       { flags.fieldSeparator = f }
func (o AwkOutputFieldSeparator) Configure(flags *flags) { flags.outputFieldSeparator = o }
func (v AwkVariable) Configure(flags *flags) {
	if flags.variables == nil {
		flags.variables = make(map[string]any)
	}
	flags.variables[v.Name] = v.Value
}
