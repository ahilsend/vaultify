package leases

// Options customizes the parameters of templating.
type Options struct {
	// TODO
	VaultAddress string
	// TODO
	Role string

	// TODO
	SecretsFileName string
}

// IsValid returns true if some values are filled into the options.
func (o *Options) IsValid() bool {
	return o != nil &&
		o.Role != "" &&
		o.SecretsFileName != ""
}
