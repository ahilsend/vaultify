package leases

import (
	"github.com/ahilsend/vaultify/pkg/options"
)

// Options customizes the parameters of templating.
type Options struct {
	options.CommonOptions

	// Secrets file location, where the secret leases are stored
	SecretsFileName string

	// Address to use to expose metrics
	MetricsAddress string
	// Path to use to expose metrics
	MetricsPath string
}

// IsValid returns true if some values are filled into the options.
func (o *Options) IsValid() bool {
	return o != nil &&
		o.SecretsFileName != "" &&
		o.MetricsAddress != "" &&
		o.MetricsPath != ""
}
