package secrets

import (
	"fmt"
)

type MapSecrets map[string]Value

type MapSecretReader struct {
	values MapSecrets
}

func NewMapReader(values MapSecrets) *MapSecretReader {
	return &MapSecretReader{
		values: values,
	}
}

func (reader *MapSecretReader) Get(name string) (*Secret, error) {
	if value, ok := reader.values[name]; ok {
		return &Secret{
			Renewable: false,
			Data:      value,
		}, nil
	}

	return nil, fmt.Errorf("unknown key '%s'", name)
}

func (reader *MapSecretReader) GetAuthSecret() *Secret {
	return nil
}
