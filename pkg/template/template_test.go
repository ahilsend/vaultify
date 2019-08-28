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
	renderAndCompare(t, secretReader, input, expectedOutput)
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
	renderAndCompare(t, secretReader, input, expectedOutput)
}

func renderAndCompare(t *testing.T, secretReader secrets.SecretReader, input string, expectedOutput string) {
	template := New(hclog.Default(), secretReader)

	output := new(bytes.Buffer)
	if err := template.render(strings.NewReader(input), output); err != nil {
		t.Fatal(err)
	}
	actualResult := output.String()

	if actualResult != expectedOutput {
		t.Fatalf("expected %s but got %s", expectedOutput, actualResult)
	}
}
