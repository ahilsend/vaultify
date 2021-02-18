package run

import (
	"context"
	"time"

	"github.com/hashicorp/go-hclog"

	"github.com/ahilsend/vaultify/pkg/http"
	"github.com/ahilsend/vaultify/pkg/prometheus"
	"github.com/ahilsend/vaultify/pkg/secrets"
	"github.com/ahilsend/vaultify/pkg/template"
	"github.com/ahilsend/vaultify/pkg/vault"
)

var retries int

func Run(logger hclog.Logger, options *Options) error {
	config := options.VaultApiConfig()
	vaultClient, err := vault.NewClient(logger, options.Role, config)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// Reset the default mux to clear handlers
	http.NewDefaultMux()
	prometheus.RegisterHandler(options.MetricsPath)
	go http.Serve(options.MetricsAddress)
	go vaultClient.StartAuthRenewal(ctx)

	secretReader := secrets.NewVaultReader(vaultClient)
	vaultTemplate := template.New(logger, secretReader)
	resultSecrets, err := vaultTemplate.RenderToPath(options.CommonTemplateOptions)
	if err != nil {
		return err
	}

	go vaultClient.RenewLeases(ctx, resultSecrets.Secrets)

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
