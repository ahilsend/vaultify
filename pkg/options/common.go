package options

import (
	"os"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/vault/api"
	"golang.org/x/time/rate"
)

type CommonOptions struct {
	VaultAddress string        // VAULT_ADDR
	Timeout      time.Duration // VAULT_CLIENT_TIMEOUT
	MaxRetries   int           // VAULT_MAX_RETRIES

	RateLimit      time.Duration
	RateLimitBurst int
}

type CommonTemplateOptions struct {
	Role string

	// Template file to be rendered (deprecated)
	TemplateFileName string

	// Template file or directory to be rendered
	TemplatePath string
	// Location of output file or directory
	OutputPath string

	// Optional, for setting variables to test the templating without vault connection.
	Variables map[string]string
}

// IsValid returns true if some values are filled into the options.
func (o *CommonTemplateOptions) IsValid() bool {
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

	if o.TemplatePath == "" || o.OutputPath == "" {
		return false
	}

	return len(o.Variables) > 0 || o.Role != ""
}

func (o *CommonOptions) VaultApiConfig() *api.Config {
	var limiter *rate.Limiter
	if o.RateLimit != 0 && o.RateLimitBurst != 0 {
		limiter = rate.NewLimiter(rate.Every(o.RateLimit), o.RateLimitBurst)
	}

	return &api.Config{
		Address:    o.VaultAddress,
		Timeout:    o.Timeout,
		MaxRetries: o.MaxRetries,
		Limiter:    limiter,
		Backoff:    retryablehttp.DefaultBackoff,
	}
}
