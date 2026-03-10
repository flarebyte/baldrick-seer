package cli

import (
	"github.com/flarebyte/baldrick-seer/internal/app"
	"github.com/flarebyte/baldrick-seer/internal/domain"
	"github.com/spf13/cobra"
)

const configFlagName = "config"

type validateRunner func(domain.CommandRequest) (domain.CommandResult, error)
type reportGenerateRunner func(domain.CommandRequest) (domain.CommandResult, error)

type dependencies struct {
	runValidate       validateRunner
	runReportGenerate reportGenerateRunner
}

func NewRootCmd() *cobra.Command {
	return newRootCmd(dependencies{
		runValidate:       app.RunValidate,
		runReportGenerate: app.RunReportGenerate,
	})
}

func newRootCmd(deps dependencies) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "seer",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.AddCommand(newValidateCmd(deps.runValidate))
	rootCmd.AddCommand(newReportCmd(deps.runReportGenerate))

	return rootCmd
}

func newValidateCmd(run validateRunner) *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate the input model",
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := run(domain.CommandRequest{
				CommandName: domain.CommandNameValidate,
				ConfigPath:  configPath,
			})
			if err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write([]byte(renderValidateSuccess()))
			return err
		},
	}

	cmd.Flags().StringVar(&configPath, configFlagName, "", "Path to the config file")

	return cmd
}

func newReportCmd(run reportGenerateRunner) *cobra.Command {
	reportCmd := &cobra.Command{
		Use:   "report",
		Short: "Manage report commands",
	}

	var configPath string

	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a report",
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := run(domain.CommandRequest{
				CommandName: domain.CommandNameReportGenerate,
				ConfigPath:  configPath,
			})
			if err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write([]byte(renderReportGenerateSuccess()))
			return err
		},
	}

	generateCmd.Flags().StringVar(&configPath, configFlagName, "", "Path to the config file")

	reportCmd.AddCommand(generateCmd)

	return reportCmd
}
