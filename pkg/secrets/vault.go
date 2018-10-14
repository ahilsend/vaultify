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

/*

auth secret:
{
  "request_id": "e3bbf043-95d0-46d5-0db5-92120de452f2",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": null,
  "wrap_info": null,
  "warnings": null,
  "auth": {
    "client_token": "36f6ba13-1af4-a35c-579d-a44f01c08eab",
    "accessor": "afac37e5-b3a2-fd77-35f1-6c0dbb223488",
    "policies": [
      "default",
      "maindb-admin"
    ],
    "token_policies": [
      "default",
      "maindb-admin"
    ],
    "metadata": {
      "role": "maindb-admin",
      "service_account_name": "vault",
      "service_account_namespace": "vault",
      "service_account_secret_name": "vault-token-nm2bl",
      "service_account_uid": "7ef19604-ce1a-11e8-a247-0800270b34b8"
    },
    "lease_duration": 2764800,
    "renewable": true,
    "entity_id": "02e9b4f5-4b8c-b52d-2a6b-f7e27e06d736"
  }
}

db creds:
{
  "request_id": "83cde93c-1b52-44c9-ecae-3b8dd008e913",
  "lease_id": "database/creds/maindb-admin/48mtIpnwiWOQeBtDnScQ2Myb",
  "renewable": true,
  "lease_duration": 900,
  "data": {
    "password": "A1a-2NpdNbcLm3MoXbbT",
    "username": "v-token-maindb-adm-53v7oAIdwgd8S"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}


*/
