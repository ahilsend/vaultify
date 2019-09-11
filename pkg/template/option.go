package template

import (
	"github.com/ahilsend/vaultify/pkg/options"
)

// Options customizes the parameters of templating.
type Options struct {
	options.CommonOptions
	options.CommonTemplateOptions

	// Secrets file location, where the secret leases are stored
	SecretsOutputFileName string
}

// IsValid returns true if some values are filled into the options.
func (o *Options) IsValid() bool {
	if o == nil {
		return false
	}

	if !o.CommonTemplateOptions.IsValid() {
		return false
	}

	return len(o.Variables) == 0 || o.SecretsOutputFileName != ""
}
