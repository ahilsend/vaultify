package template

import (
	"os"

	"github.com/ahilsend/vaultify/pkg/options"
)

// Options customizes the parameters of templating.
type Options struct {
	options.CommonOptions
	// Kubernetes auth role to use
	Role string

	// Template file to be rendered (deprecated)
	TemplateFileName string

	// Template file or directory to be rendered
	TemplatePath string
	// Location of output file or directory
	OutputPath string

	// Secrets file location, where the secret leases are stored
	SecretsOutputFileName string

	// Optional, for setting variables to test the templating without vault connection.
	Variables map[string]string
}

// IsValid returns true if some values are filled into the options.
func (o *Options) IsValid() bool {
	if o == nil {
		return false
	}

	if o.TemplateFileName != "" {
		file, err := os.Stat(o.TemplateFileName)
		if err != nil {
			return false
		}

		if !file.Mode().IsRegular() {
			return false
		}
		o.TemplatePath = o.TemplateFileName
	}

	if o.TemplatePath == "" {
		return false
	}

	if len(o.Variables) > 0 {
		return true
	}

	return o.Role != "" &&
		o.OutputPath != "" &&
		o.SecretsOutputFileName != ""
}
