package template

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-hclog"

	"github.com/ahilsend/vaultify/pkg/secrets"
)

func TestRenderSimple(t *testing.T) {

	input := `
credentials:
  {{- $mySecret := vault "secret/my/key" }}
  attribute1: {{ $mySecret.Data.attribute1 }}
  attribute2: {{ $mySecret.Data.attribute2 }}
`

	expectedOutput := `
credentials:
  attribute1: value1
  attribute2: value2
`
	secretReader := secrets.NewMapReader(secrets.MapSecrets{
		"secret/my/key": {
			"attribute1": "value1",
			"attribute2": "value2",
		},
	})
	if err := renderAndCompare(secretReader, input, expectedOutput); err != nil {
		t.Fatal(err)
	}
}

func TestRenderDefault(t *testing.T) {

	input := `
credentials:
  {{- $mySecret := vault "secret/my/key" }}
  attribute1: {{ $mySecret.Data.attribute1 | default "nope1" }}
  attribute2: {{ $mySecret.Data.attribute2 | default "nope2" | quote }}
`

	expectedOutput := `
credentials:
  attribute1: 1
  attribute2: "nope2"
`
	secretReader := secrets.NewMapReader(secrets.MapSecrets{
		"secret/my/key": {
			"attribute1": 1,
		},
	})
	if err := renderAndCompare(secretReader, input, expectedOutput); err != nil {
		t.Fatal(err)
	}
}

func renderAndCompare(secretReader secrets.SecretReader, input string, expectedOutput string) error {
	template := New(hclog.Default(), secretReader)

	output := new(bytes.Buffer)
	if _, err := template.render(strings.NewReader(input), output); err != nil {
		return err
	}
	actualResult := output.String()

	if actualResult != expectedOutput {
		return fmt.Errorf("expected %s but got %s", expectedOutput, actualResult)
	}
	return nil
}
