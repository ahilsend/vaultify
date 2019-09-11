package run

import (
	"github.com/ahilsend/vaultify/pkg/options"
)

// Options customizes the parameters of templating.
type Options struct {
	options.CommonOptions
	options.CommonTemplateOptions

	// Address to use to expose metrics
	MetricsAddress string
	// Path to use to expose metrics
	MetricsPath string
}

// IsValid returns true if some values are filled into the options.
func (o *Options) IsValid() bool {
	if o == nil {
		return false
	}

	return o.CommonTemplateOptions.IsValid() &&
		o.MetricsAddress != "" &&
		o.MetricsPath != ""
}
