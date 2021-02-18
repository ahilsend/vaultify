package leases

import (
	"context"
	"time"

	"github.com/hashicorp/go-hclog"

	"github.com/ahilsend/vaultify/pkg/http"
	"github.com/ahilsend/vaultify/pkg/prometheus"
	"github.com/ahilsend/vaultify/pkg/secrets"
	"github.com/ahilsend/vaultify/pkg/vault"
)

var retries int

func Run(logger hclog.Logger, options *Options) error {
	secretResult, err := secrets.Read(options.SecretsFileName)
	if err != nil {
		return err
	}

	config := options.VaultApiConfig()
	vaultClient, err := vault.NewClientFromSecret(logger, secretResult.AuthSecret, config)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// Reset the default mux to clear handlers
	http.NewDefaultMux()
	prometheus.RegisterHandler(options.MetricsPath)
	go http.Serve(options.ListenAddress)
	go vaultClient.StartAuthRenewal(ctx)
	go vaultClient.RenewLeases(ctx, secretResult.Secrets)

	err = vaultClient.Wait(ctx)
	// We can safely retry fetching the secret as long as we get empty secret data from vault
	if err == vault.ErrRenewerNoSecretData {
		if retries <= options.MaxRetries {
			retries++
			time.Sleep(10 * time.Second)
			Run(logger, options)
		}
	}
	return err
}
