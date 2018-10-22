package secrets

import (
	"encoding/json"
	"github.com/hashicorp/vault/api"
	"os"
)

type Value map[string]interface{}
type Secret = api.Secret
type Secrets struct {
	AuthSecret *Secret
	Secrets    map[string]Secret
}

type SecretReader interface {
	Get(name string) (*Secret, error)
	GetAuthSecret() *Secret
}

func Write(filePath string, secrets *Secrets) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	file.Chmod(0600)
	return json.NewEncoder(file).Encode(secrets)
}

func Read(filePath string) (*Secrets, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	var secrets Secrets
	err = json.NewDecoder(file).Decode(&secrets)
	if err != nil {
		return nil, err
	}
	return &secrets, err
}
