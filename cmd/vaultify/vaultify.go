package main

import (
	"flag"
	"fmt"
	"github.com/ahilsend/vaultify/pkg/run"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"

	"github.com/ahilsend/vaultify/pkg/leases"
	"github.com/ahilsend/vaultify/pkg/template"
)

var (
	logger = hclog.Default()

	flags = struct {
		logLevel           int
		role               string
		vaultAddress       string
		templateOptions    template.Options
		renewLeasesOptions leases.Options
		runOptions         run.Options
	}{}

	rootCmd = &cobra.Command{
		Use:          "vaultify",
		Short:        "TODO",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(0),
	}

	templateCmd = &cobra.Command{
		Use:   "template",
		Short: "TODO",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.templateOptions.Role = flags.role
			flags.templateOptions.VaultAddress = flags.vaultAddress

			if !flags.templateOptions.IsValid() {
				return cmd.Help()
			}

			logger.SetLevel(logLevel())

			if err := template.Run(logger, &flags.templateOptions); err != nil {
				return fmt.Errorf("templating faild: %v", err)
			}
			fmt.Println("OK")
			return nil
		},
	}

	reneawLeasesCmd = &cobra.Command{
		Use:   "renew-leases",
		Short: "TODO",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.renewLeasesOptions.Role = flags.role
			flags.renewLeasesOptions.VaultAddress = flags.vaultAddress

			if !flags.renewLeasesOptions.IsValid() {
				return cmd.Help()
			}

			logger.SetLevel(logLevel())

			if err := leases.Run(logger, &flags.renewLeasesOptions); err != nil {
				return fmt.Errorf("renew-leases faild: %v", err)
			}
			return nil
		},
	}

	runCmd = &cobra.Command{
		Use:   "run",
		Short: "TODO",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.runOptions.Role = flags.role
			flags.runOptions.VaultAddress = flags.vaultAddress

			if !flags.runOptions.IsValid() {
				return cmd.Help()
			}

			logger.SetLevel(logLevel())

			if err := run.Run(logger, &flags.runOptions); err != nil {
				return fmt.Errorf("run faild: %v", err)
			}
			return nil
		},
	}
)

func logLevel() hclog.Level {
	switch flags.logLevel {
	case 0:
		return hclog.Error
	case 1:
		return hclog.Warn
	case 2:
		return hclog.Info
	case 3:
		return hclog.Debug
	}
	return hclog.Trace
}

func init() {
	rootCmd.PersistentFlags().CountVarP(&flags.logLevel, "verbose", "v", "Log level TODO")
	rootCmd.PersistentFlags().StringVar(&flags.vaultAddress, "vault", "", "Vault address")
	rootCmd.PersistentFlags().StringVar(&flags.role, "role", "", "Vault role to assume")
	flag.CommandLine.VisitAll(func(gf *flag.Flag) {
		rootCmd.PersistentFlags().AddGoFlag(gf)
	})

	templateCmd.Flags().StringVar(&flags.templateOptions.TemplateFileName, "template-file", "", "Template file to render")
	templateCmd.Flags().StringVar(&flags.templateOptions.OutputFileName, "output-file", "", "Output file")
	templateCmd.Flags().StringVar(&flags.templateOptions.SecretsOutputFileName, "secrets-output-file", "", "Secrets output file")
	templateCmd.Flags().StringToStringVar(&flags.templateOptions.Variables, "var", map[string]string{}, "TODO")

	reneawLeasesCmd.Flags().StringVar(&flags.renewLeasesOptions.SecretsFileName, "secrets-file", "", "Secrets file")

	runCmd.Flags().StringVar(&flags.runOptions.TemplateFileName, "template-file", "", "Template file to render")
	runCmd.Flags().StringVar(&flags.runOptions.OutputFileName, "output-file", "", "Output file")

	rootCmd.AddCommand(templateCmd)
	rootCmd.AddCommand(reneawLeasesCmd)
	rootCmd.AddCommand(runCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
