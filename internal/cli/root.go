package cli

import (
	"github.com/flarebyte/baldrick-seer/internal/app"
	"github.com/spf13/cobra"
)

const configFlagName = "config"

type validateRunner func(app.ValidateRequest) (app.ValidateResponse, error)
type reportGenerateRunner func(app.ReportGenerateRequest) (app.ReportGenerateResponse, error)

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
			response, err := run(app.ValidateRequest{ConfigPath: configPath})
			if err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write([]byte(response.Stdout))
			return err
		},
	}

	cmd.Flags().StringVar(&configPath, configFlagName, "", "Path to the config file")
	_ = cmd.MarkFlagRequired(configFlagName)

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
			response, err := run(app.ReportGenerateRequest{ConfigPath: configPath})
			if err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write([]byte(response.Stdout))
			return err
		},
	}

	generateCmd.Flags().StringVar(&configPath, configFlagName, "", "Path to the config file")
	_ = generateCmd.MarkFlagRequired(configFlagName)

	reportCmd.AddCommand(generateCmd)

	return reportCmd
}
