package template

// Options customizes the parameters of templating.
type Options struct {
	// TODO
	VaultAddress string
	// TODO
	Role string

	// TODO
	TemplateFileName      string
	OutputFileName        string
	SecretsOutputFileName string

	Variables map[string]string
}

// IsValid returns true if some values are filled into the options.
func (o *Options) IsValid() bool {
	if o == nil || o.TemplateFileName == "" {
		return false
	}

	if len(o.Variables) > 0 {
		return true
	}
	return o.Role != "" &&
		o.OutputFileName != "" &&
		o.SecretsOutputFileName != ""
}
