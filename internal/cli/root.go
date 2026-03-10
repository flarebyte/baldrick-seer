package cli

import "github.com/spf13/cobra"

const (
	validateStubOutput       = "validate: ok\n"
	reportGenerateStubOutput = "report generate: ok\n"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "seer",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.AddCommand(newValidateCmd())
	rootCmd.AddCommand(newReportCmd())

	return rootCmd
}

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate the input model",
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := cmd.OutOrStdout().Write([]byte(validateStubOutput))
			return err
		},
	}
}

func newReportCmd() *cobra.Command {
	reportCmd := &cobra.Command{
		Use:   "report",
		Short: "Manage report commands",
	}

	reportCmd.AddCommand(&cobra.Command{
		Use:   "generate",
		Short: "Generate a report",
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := cmd.OutOrStdout().Write([]byte(reportGenerateStubOutput))
			return err
		},
	})

	return reportCmd
}
