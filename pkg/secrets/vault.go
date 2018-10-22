package secrets

import (
	"github.com/ahilsend/vaultify/pkg/vault"
)

type VaultSecretReader struct {
	vaultClient *vault.Client
}

func NewVaultReader(vaultClient *vault.Client) *VaultSecretReader {
	return &VaultSecretReader{
		vaultClient: vaultClient,
	}
}

func (reader *VaultSecretReader) Get(name string) (*Secret, error) {
	secret, err := reader.vaultClient.ApiClient.Logical().Read(name)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func (reader *VaultSecretReader) GetAuthSecret() *Secret {
	return reader.vaultClient.AuthSecret
}
