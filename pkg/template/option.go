package template

import (
	"github.com/ahilsend/vaultify/pkg/options"
)

// Options customizes the parameters of templating.
type Options struct {
	options.CommonOptions
	// Kubernetes auth role to use
	Role string

	// Template file to be rendered
	TemplateFileName string
	// Location of the output file
	OutputFileName string
	// Secrets file location, where the secret leases are stored
	SecretsOutputFileName string

	// Optional, for setting variables to test the templating without vault connection.
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
