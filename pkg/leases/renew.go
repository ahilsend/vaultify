package leases

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/ahilsend/vaultify/pkg/http"
	"github.com/ahilsend/vaultify/pkg/prometheus"
	"github.com/ahilsend/vaultify/pkg/secrets"
	"github.com/ahilsend/vaultify/pkg/vault"
)

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

	prometheus.RegisterHandler(options.MetricsPath)
	go http.Serve(options.ListenAddress)
	go vaultClient.StartAuthRenewal(ctx)
	go vaultClient.RenewLeases(ctx, secretResult.Secrets)

	return vaultClient.Wait(ctx)
}
