package leases

import (
	"context"
	"github.com/ahilsend/vaultify/pkg/secrets"
	"github.com/ahilsend/vaultify/pkg/vault"
	"github.com/hashicorp/go-hclog"
)

func Run(logger hclog.Logger, options *Options) error {

	vaultClient, err := vault.NewVaultClient(logger, options.VaultAddress, options.Role)
	if err != nil {
		return err
	}

	ctx := context.Background()

	secretMap, err := secrets.Read(options.SecretsFileName)
	if err != nil {
		return err
	}

	go vaultClient.StartAuthRenewal(ctx)
	go vaultClient.RenewLeases(ctx, secretMap)

	return vaultClient.Wait(ctx)
}
