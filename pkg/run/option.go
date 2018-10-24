package run

// Options customizes the parameters of templating.
type Options struct {
	// Vault api address. Can be specified via VAULT_ADDR instead
	VaultAddress string
	// Kubernetes auth role to use
	Role string

	// Template file to be rendered
	TemplateFileName string
	// Location of the output file
	OutputFileName string

	// Address to use to expose metrics
	MetricsAddress string
	// Path to use to expose metrics
	MetricsPath string
}

// IsValid returns true if some values are filled into the options.
func (o *Options) IsValid() bool {
	return o != nil &&
		o.Role != "" &&
		o.TemplateFileName != "" &&
		o.OutputFileName != "" &&
		o.MetricsAddress != "" &&
		o.MetricsPath != ""
}
