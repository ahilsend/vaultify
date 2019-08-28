package template

import (
	"bytes"
	"strings"
	"testing"

	"github.com/hashicorp/go-hclog"

	"github.com/ahilsend/vaultify/pkg/secrets"
)

func TestRenderSimple(t *testing.T) {

	input := `
credentials:
  <{- $mySecret := vault "secret/my/key" }>
  attribute1: <{ $mySecret.Data.attribute1 }>
  attribute2: <{ $mySecret.Data.attribute2 }>
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
	renderAndCompare(t, secretReader, input, expectedOutput, []string{"secret/my/key"})
}

func TestRenderDefault(t *testing.T) {

	input := `
credentials:
  <{- $mySecret := vault "secret/my/key" }>
  attribute1: <{ $mySecret.Data.attribute1 | default "nope1" }>
  attribute2: <{ $mySecret.Data.attribute2 | default "nope2" | quote }>
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
	renderAndCompare(t, secretReader, input, expectedOutput, []string{"secret/my/key"})
}

func checkExpectedSecrets(t *testing.T, secrets *secrets.Secrets, expectedSecrets []string) {
	for _, secret := range expectedSecrets {
		if _, ok := secrets.Secrets[secret]; !ok {
			t.Errorf("didn't find expected secret %s", secret)
		}
	}

	if len(secrets.Secrets) != len(expectedSecrets) {
		t.Errorf("unexpected amount of secrets, expected %d but got %d", len(expectedSecrets), len(secrets.Secrets))
		t.Logf("secrets: %v", secrets.Secrets)
	}
}

func renderAndCompare(t *testing.T, secretReader secrets.SecretReader, input string, expectedOutput string, expectedSecrets []string) {
	template := New(hclog.Default(), secretReader)

	output := new(bytes.Buffer)
	if err := template.render(strings.NewReader(input), output); err != nil {
		t.Fatal(err)
	}
	actualResult := output.String()

	if actualResult != expectedOutput {
		t.Fatalf("expected %s but got %s", expectedOutput, actualResult)
	}

	checkExpectedSecrets(t, template.secrets, expectedSecrets)
}
