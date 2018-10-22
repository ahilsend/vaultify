package run

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/ahilsend/vaultify/pkg/secrets"
	"github.com/ahilsend/vaultify/pkg/template"
	"github.com/ahilsend/vaultify/pkg/vault"
)

func Run(logger hclog.Logger, options *Options) error {
	vaultClient, err := vault.NewClient(logger, options.VaultAddress, options.Role)
	if err != nil {
		return err
	}

	ctx := context.Background()
	go vaultClient.StartAuthRenewal(ctx)

	secretReader := secrets.NewVaultReader(vaultClient)
	vaultTemplate := template.New(logger, secretReader)

	secretMap, err := vaultTemplate.RenderToFile(options.TemplateFileName, options.OutputFileName)
	if err != nil {
		return err
	}

	go vaultClient.RenewLeases(ctx, secretMap.Secrets)

	return vaultClient.Wait(ctx)
}
