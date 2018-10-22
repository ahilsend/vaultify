package run

// Options customizes the parameters of templating.
type Options struct {
	// TODO
	VaultAddress string
	// TODO
	Role string

	// TODO
	TemplateFileName string
	OutputFileName   string
}

// IsValid returns true if some values are filled into the options.
func (o *Options) IsValid() bool {
	return o != nil &&
		o.Role != "" &&
		o.TemplateFileName != "" &&
		o.OutputFileName != ""
}
