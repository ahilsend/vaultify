package vault

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent/auth"
	"github.com/hashicorp/vault/command/agent/auth/kubernetes"
)

type Client struct {
	ApiClient   *api.Client
	AuthSecret  *api.Secret
	authRenewer *api.Renewer
	doneCh      chan error
	logger      hclog.Logger
}

func NewClient(logger hclog.Logger, vaultAddr string, role string) (*Client, error) {
	client, err := createClient(vaultAddr)
	if err != nil {
		return nil, err
	}

	authSecret, err := kubernetesAuthentication(client, logger, role)
	if err != nil {
		return nil, err
	}
	renewer, err := authRenewer(client, authSecret)
	if err != nil {
		return nil, err
	}

	return &Client{
		ApiClient:   client,
		AuthSecret:  authSecret,
		authRenewer: renewer,
		doneCh:      make(chan error, 1),
		logger:      logger,
	}, err
}

func NewClientFromSecret(logger hclog.Logger, vaultAddr string, authSecret *api.Secret) (*Client, error) {
	client, err := createClient(vaultAddr)
	if err != nil {
		return nil, err
	}

	renewer, err := authRenewer(client, authSecret)
	if err != nil {
		return nil, err
	}

	return &Client{
		ApiClient:   client,
		AuthSecret:  authSecret,
		authRenewer: renewer,
		doneCh:      make(chan error, 1),
		logger:      logger,
	}, err
}

func createClient(vaultAddr string) (*api.Client, error) {
	vaultConfig := api.DefaultConfig()
	if vaultAddr != "" {
		vaultConfig.Address = vaultAddr
	}

	return api.NewClient(vaultConfig)
}

func authRenewer(v *api.Client, authSecret *api.Secret) (*api.Renewer, error) {
	v.SetToken(authSecret.Auth.ClientToken)
	return v.NewRenewer(&api.RenewerInput{
		Secret: authSecret,
	})
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
			v.logger.Info("auth leaese renewer done channel triggered")
			v.doneCh <- fmt.Errorf("auth lease renewer done: %v", err)
			return

		case <-v.authRenewer.RenewCh():
			v.logger.Info("renewed lease for auth token")
			break
		}
	}
}

func (v *Client) RenewLeases(ctx context.Context, secretMap map[string]api.Secret) {
	for name, secret := range secretMap {
		if !secret.Renewable {
			continue
		}

		renewer, err := v.ApiClient.NewRenewer(&api.RenewerInput{
			Secret: &secret,
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
			v.logger.Info("lease renewer done channel triggered", "name", name)
			v.doneCh <- fmt.Errorf("lease renewer done: %v", err)
			return

		case <-renewer.RenewCh():
			v.logger.Info("renewed lease for secret", "name", name)
			break
		}
	}
}
