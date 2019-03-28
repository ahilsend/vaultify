package options

import (
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
