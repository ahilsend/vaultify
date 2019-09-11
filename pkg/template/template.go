package template

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/hashicorp/go-hclog"

	"github.com/ahilsend/vaultify/pkg/options"
	"github.com/ahilsend/vaultify/pkg/secrets"
	"github.com/ahilsend/vaultify/pkg/vault"
)

const (
	templateName = "vaultify"
)

type VaultifyTemplate struct {
	secretReader secrets.SecretReader
	logger       hclog.Logger
	funcMap      map[string]interface{}
	secrets      *secrets.Secrets
}

func Run(logger hclog.Logger, options *Options) error {
	secretReader, err := createSecretReader(logger, options)
	if err != nil {
		return err
	}

	vaultTemplate := New(logger, secretReader)
	resultSecrets, err := vaultTemplate.RenderToPath(options.CommonTemplateOptions)
	if err != nil {
		return err
	}

	if options.SecretsOutputFileName == "" {
		return nil
	}
	return secrets.Write(options.SecretsOutputFileName, resultSecrets)
}

func createSecretReader(logger hclog.Logger, options *Options) (secrets.SecretReader, error) {
	if len(options.Variables) > 0 {
		values := secrets.MapSecrets{}
		for name, jsonString := range options.Variables {
			var secret secrets.Value
			err := json.Unmarshal([]byte(jsonString), &secret)
			if err != nil {
				return nil, err
			}
			values[name] = secret
		}
		return secrets.NewMapReader(values), nil
	}

	config := options.VaultApiConfig()
	vaultClient, err := vault.NewClient(logger, options.Role, config)
	if err != nil {
		return nil, err
	}

	return secrets.NewVaultReader(vaultClient), nil
}

func New(logger hclog.Logger, secretReader secrets.SecretReader) *VaultifyTemplate {
	t := &VaultifyTemplate{
		secretReader: secretReader,
		logger:       logger,
		funcMap:      sprig.GenericFuncMap(),
		secrets: &secrets.Secrets{
			AuthSecret: secretReader.GetAuthSecret(),
			Secrets:    map[string]secrets.Secret{},
		},
	}

	t.funcMap["vault"] = t.getVaultSecret
	return t
}

func (t *VaultifyTemplate) getVaultSecret(name string) (*secrets.Secret, error) {
	if name == "" {
		return nil, errors.New("you need to pass a name to the 'vault' function")
	}

	secret, err := t.secretReader.Get(name)
	if err != nil {
		return nil, err
	}
	t.secrets.Secrets[name] = *secret
	return secret, err
}

func (t *VaultifyTemplate) RenderToPath(options options.CommonTemplateOptions) (*secrets.Secrets, error) {
	file, err := os.Stat(options.TemplatePath)
	if err != nil {
		return nil, err
	}

	if file.Mode().IsRegular() {
		return t.RenderToFile(options.TemplatePath, options.OutputPath)
	} else if file.Mode().IsDir() {
		return t.RenderToDirectory(options.TemplatePath, options.OutputPath)
	}
	return nil, errors.New("Path is not a file or a directory")
}

func (t *VaultifyTemplate) RenderToFile(templateFile string, outputFile string) (*secrets.Secrets, error) {
	t.logger.Info("Rendering template", "template", templateFile)
	templateBytes, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return nil, err
	}

	var output io.Writer
	if outputFile == "" {
		output = os.Stdout
	} else {
		file, err := os.Create(outputFile)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		file.Chmod(0600)
		output = file
	}

	err = t.render(bytes.NewBuffer(templateBytes), output)
	if err != nil {
		t.logger.Error("Error during rendering", "error", err)
		return nil, err
	}

	return t.secrets, nil
}

func (t *VaultifyTemplate) RenderToDirectory(templateDir string, outputDir string) (*secrets.Secrets, error) {
	t.logger.Info("Rendering template directory", "directory", templateDir)

	err := filepath.Walk(templateDir, func(templateFile string, info os.FileInfo, err error) error {
		if err != nil {
			t.logger.Error("Error visiting path", "path", templateFile, "error", err)
			return err
		}

		relativePath, err := filepath.Rel(templateDir, templateFile)
		if err != nil {
			t.logger.Error("Path not relative to templateDir", "path", templateFile, "templateDir", templateDir, "error", err)
			return err
		}
		outputPath := path.Join(outputDir, relativePath)

		if info.IsDir() {
			t.logger.Info("Creating directory", "directory", outputPath)
			// TODO: need to be writable while rendering templates but could
			// probably be restored afterwards.
			if err := os.MkdirAll(outputPath, info.Mode().Perm()|0700); err != nil {
				t.logger.Error("Failed to create output directory structure", "outputPath", outputPath)
				return err
			}
			return nil
		}

		_, err = t.RenderToFile(templateFile, outputPath)
		err = nil
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return t.secrets, nil
}

func (t *VaultifyTemplate) render(input io.Reader, output io.Writer) error {
	inputBytes, err := ioutil.ReadAll(input)
	if err != nil {
		return err
	}

	tmpl := template.New(templateName)
	tmpl.Delims("<{", "}>")
	tmpl.Funcs(t.funcMap)

	_, err = tmpl.Parse(string(inputBytes))
	if err != nil {
		return err
	}

	err = tmpl.Execute(output, nil)
	if err != nil {
		return err
	}

	return nil
}
