package run

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/ahilsend/vaultify/pkg/prometheus"
	"github.com/ahilsend/vaultify/pkg/secrets"
	"github.com/ahilsend/vaultify/pkg/template"
	"github.com/ahilsend/vaultify/pkg/vault"
)

func Run(logger hclog.Logger, options *Options) error {
	config := options.VaultApiConfig()
	vaultClient, err := vault.NewClient(logger, options.Role, config)
	if err != nil {
		return err
	}

	ctx := context.Background()
	go prometheus.StartServer(options.MetricsAddress, options.MetricsPath)
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
