package leases

// Options customizes the parameters of templating.
type Options struct {
	// Vault api address. Can be specified via VAULT_ADDR instead
	VaultAddress string

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
