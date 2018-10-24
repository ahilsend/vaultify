package leases

import (
	"context"
	"github.com/ahilsend/vaultify/pkg/prometheus"
	"github.com/ahilsend/vaultify/pkg/secrets"
	"github.com/ahilsend/vaultify/pkg/vault"
	"github.com/hashicorp/go-hclog"
)

func Run(logger hclog.Logger, options *Options) error {
	secretResult, err := secrets.Read(options.SecretsFileName)
	if err != nil {
		return err
	}

	vaultClient, err := vault.NewClientFromSecret(logger, options.VaultAddress, secretResult.AuthSecret)
	if err != nil {
		return err
	}

	ctx := context.Background()
	go prometheus.StartServer(options.MetricsAddress, options.MetricsPath)
	go vaultClient.StartAuthRenewal(ctx)
	go vaultClient.RenewLeases(ctx, secretResult.Secrets)

	return vaultClient.Wait(ctx)
}
