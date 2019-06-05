package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ahilsend/vaultify/pkg/prometheus"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent/auth"
	"github.com/hashicorp/vault/command/agent/auth/kubernetes"
)

type Client struct {
	ApiClient   *api.Client
	AuthSecret  *api.Secret
	authRenewer *api.Renewer
	role        string
	doneCh      chan error
	logger      hclog.Logger
}

func NewClient(logger hclog.Logger, role string, config *api.Config) (*Client, error) {
	return createClient(logger, func(client *api.Client) (*api.Secret, string, error) {
		authSecret, err := kubernetesAuthentication(client, logger, role)
		return authSecret, role, err
	}, config)
}

func NewClientFromSecret(logger hclog.Logger, authSecret *api.Secret, config *api.Config) (*Client, error) {
	return createClient(logger, func(client *api.Client) (*api.Secret, string, error) {
		metadata, err := authSecret.TokenMetadata()
		return authSecret, metadata["role"], err
	}, config)
}

func createClient(logger hclog.Logger, auth func(*api.Client) (*api.Secret, string, error), config *api.Config) (*Client, error) {
	vaultConfig := mergeConfig(api.DefaultConfig(), config)

	client, err := api.NewClient(vaultConfig)
	if err != nil {
		return nil, err
	}

	authSecret, role, err := auth(client)
	if err != nil {
		return nil, err
	}

	client.SetToken(authSecret.Auth.ClientToken)
	renewer, err := client.NewRenewer(&api.RenewerInput{
		Secret: authSecret,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		ApiClient:   client,
		AuthSecret:  authSecret,
		authRenewer: renewer,
		role:        role,
		doneCh:      make(chan error, 1),
		logger:      logger,
	}, err
}

func mergeConfig(vaultConfig *api.Config, config *api.Config) *api.Config {
	if config.Address != "" {
		vaultConfig.Address = config.Address
	}
	if config.HttpClient != nil {
		vaultConfig.HttpClient = config.HttpClient
	}
	if config.MaxRetries != 0 {
		vaultConfig.MaxRetries = config.MaxRetries
	}
	if config.Timeout != 0 {
		vaultConfig.Timeout = config.Timeout
	}
	if config.Error != nil {
		vaultConfig.Error = config.Error
	}
	if config.Backoff != nil {
		vaultConfig.Backoff = config.Backoff
	}
	if config.Limiter != nil {
		vaultConfig.Limiter = config.Limiter
	}

	return vaultConfig
}

func kubernetesAuthentication(v *api.Client, logger hclog.Logger, role string) (*api.Secret, error) {
	authMethod, err := kubernetes.NewKubernetesAuthMethod(&auth.AuthConfig{
		MountPath: "auth/kubernetes",
		Logger:    logger,
		Config: map[string]interface{}{
			"role": role,
		},
	})
	if err != nil {
		return nil, err
	}
	path, data, err := authMethod.Authenticate(context.Background(), v)
	if err != nil {
		return nil, err
	}

	return v.Logical().Write(path, data)
}

func (v *Client) Wait(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			v.logger.Info("shutdown triggered, stopping auth renewal loop")
			return nil

		case err := <-v.DoneCh():
			v.logger.Info("vault client done channel triggered")
			if err != nil {
				v.logger.Error("error renewing secret", "error", err)
			}
			return err
		}
	}
}

func (v *Client) DoneCh() chan error {
	return v.doneCh
}

func (v *Client) StartAuthRenewal(ctx context.Context) {
	v.logger.Info("starting auth lease renewal")
	go v.authRenewer.Renew()

	for {
		select {
		case <-ctx.Done():
			v.logger.Info("shutdown triggered, stopping auth lease renewal")
			v.authRenewer.Stop()
			return

		case err := <-v.authRenewer.DoneCh():
			prometheus.IncAuthLeaseFailed(v.role)
			v.logger.Warn("auth leaese renewer done channel triggered")
			v.doneCh <- fmt.Errorf("auth lease renewer done: %v", err)
			return

		case renewed := <-v.authRenewer.RenewCh():
			hasWarnings := len(renewed.Secret.Warnings) > 0
			prometheus.IncAuthLeaseRenewed(v.role, hasWarnings)
			if v.logger.IsTrace() {
				bytes, _ := json.MarshalIndent(renewed.Secret, "", "  ")
				v.logger.Trace("renewed lease for auth token", "secret", string(bytes))
			} else {
				v.logger.Info("renewed lease for auth token")
			}

			if hasWarnings {
				v.logger.Warn("Lease warning", "lease_warning", renewed.Secret.Warnings)
			}
			break
		}
	}
}

func (v *Client) RenewLeases(ctx context.Context, secretMap map[string]api.Secret) {
	for name, secret := range secretMap {
		if !secret.Renewable {
			continue
		}

		// local copy of secret or the reference to it will be incorrect
		renewerSecret := secret
		renewer, err := v.ApiClient.NewRenewer(&api.RenewerInput{
			Secret: &renewerSecret,
		})
		if err != nil {
			v.doneCh <- err
			return
		}

		go v.startRenewal(ctx, name, renewer)
	}

	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (v *Client) startRenewal(ctx context.Context, name string, renewer *api.Renewer) {
	v.logger.Info("starting lease renewal for secret", "name", name)
	go renewer.Renew()

	for {
		select {
		case <-ctx.Done():
			v.logger.Info("shutdown triggered, stopping lease renewer", "name", name)
			renewer.Stop()
			return

		case err := <-renewer.DoneCh():
			prometheus.IncSecretLeaseFailed(v.role, name)
			v.logger.Warn("lease renewer done channel triggered", "name", name)
			v.doneCh <- fmt.Errorf("lease renewer done: %v", err)
			return

		case renewed := <-renewer.RenewCh():
			hasWarnings := len(renewed.Secret.Warnings) > 0
			prometheus.IncSecretLeaseRenewed(v.role, name, hasWarnings)
			if v.logger.IsTrace() {
				bytes, _ := json.MarshalIndent(renewed.Secret, "", "  ")
				v.logger.Trace("renewed lease for secret",
					"name", name,
					"secret", string(bytes))
			} else if v.logger.IsDebug() {
				v.logger.Debug("renewed lease for secret",
					"name", name,
					"leaseId", renewed.Secret.LeaseID)
			} else {
				v.logger.Info("renewed lease for secret",
					"name", name)
			}

			if hasWarnings {
				v.logger.Warn("Lease warning", "lease_warning", renewed.Secret.Warnings)
			}
			break
		}
	}
}
