package cli

import (
	"github.com/flarebyte/baldrick-seer/internal/app"
	"github.com/spf13/cobra"
)

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
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate the input model",
		RunE: func(cmd *cobra.Command, _ []string) error {
			response, err := run(app.ValidateRequest{})
			if err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write([]byte(response.Stdout))
			return err
		},
	}
}

func newReportCmd(run reportGenerateRunner) *cobra.Command {
	reportCmd := &cobra.Command{
		Use:   "report",
		Short: "Manage report commands",
	}

	reportCmd.AddCommand(&cobra.Command{
		Use:   "generate",
		Short: "Generate a report",
		RunE: func(cmd *cobra.Command, _ []string) error {
			response, err := run(app.ReportGenerateRequest{})
			if err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write([]byte(response.Stdout))
			return err
		},
	})

	return reportCmd
}
