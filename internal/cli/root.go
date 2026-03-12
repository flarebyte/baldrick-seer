package cli

import (
	"context"
	"fmt"

	"github.com/flarebyte/baldrick-seer/internal/app"
	"github.com/flarebyte/baldrick-seer/internal/buildinfo"
	"github.com/flarebyte/baldrick-seer/internal/domain"
	"github.com/spf13/cobra"
)

const configFlagName = "config"
const repositoryURL = "https://github.com/flarebyte/baldrick-seer"

type validateRunner func(context.Context, domain.CommandRequest) (domain.CommandResult, error)
type reportGenerateRunner func(context.Context, domain.CommandRequest) (domain.CommandResult, error)

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
		Short:         "Validate decision models and generate scenario-based ranking reports.",
		Long:          rootLongHelp(),
		Example:       rootExamples(),
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       buildinfo.String(),
	}

	rootCmd.AddCommand(newValidateCmd(deps.runValidate))
	rootCmd.AddCommand(newReportCmd(deps.runReportGenerate))

	return rootCmd
}

func newValidateCmd(run validateRunner) *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:     "validate",
		Short:   "Validate a CUE decision model without generating rankings.",
		Long:    validateLongHelp(),
		Example: validateExamples(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := run(cmd.Context(), domain.CommandRequest{
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

	cmd.Flags().StringVar(&configPath, configFlagName, "", configFlagUsage())

	return cmd
}

func newReportCmd(run reportGenerateRunner) *cobra.Command {
	reportCmd := &cobra.Command{
		Use:     "report",
		Short:   "Generate ranking reports from a validated decision model.",
		Long:    reportLongHelp(),
		Example: reportExamples(),
	}

	var configPath string

	generateCmd := &cobra.Command{
		Use:     "generate",
		Short:   "Generate markdown, JSON, and CSV reports from the model.",
		Long:    reportGenerateLongHelp(),
		Example: reportGenerateExamples(),
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := run(cmd.Context(), domain.CommandRequest{
				CommandName: domain.CommandNameReportGenerate,
				ConfigPath:  configPath,
			})
			if err != nil {
				return err
			}

			_, err = cmd.OutOrStdout().Write([]byte(renderReportGenerateSuccess(result)))
			return err
		},
	}

	generateCmd.Flags().StringVar(&configPath, configFlagName, "", configFlagUsage())

	reportCmd.AddCommand(generateCmd)

	return reportCmd
}

func configFlagUsage() string {
	return "Path to a .cue file or a directory containing a CUE package"
}

func rootLongHelp() string {
	return fmt.Sprintf(`seer validates scenario-based MCDA models and generates deterministic reports.

Use this CLI when you want to:
  - check that a CUE decision model is structurally valid
  - compute scenario-local rankings with the current v1 pipeline
  - render markdown, JSON, or CSV outputs from the same model

The --config flag accepts either:
  - a single .cue file
  - a directory containing a CUE package

Repository:
  %s`, repositoryURL)
}

func rootExamples() string {
	return `  seer validate --config testdata/config/minimal.cue
  seer validate --config testdata/config_split
  seer report generate --config testdata/config/valid_report.cue
  seer help report generate`
}

func validateLongHelp() string {
	return fmt.Sprintf(`validate loads a CUE model, checks the implemented v1 validation rules,
and exits without generating reports.

Typical use:
  - run it before report generation
  - inspect deterministic validation failures and guidance
  - verify either a single file or a directory-based CUE package

Repository:
  %s`, repositoryURL)
}

func validateExamples() string {
	return `  seer validate --config testdata/config/minimal.cue
  seer validate --config testdata/config_split
  seer validate --config path/to/model.cue`
}

func reportLongHelp() string {
	return fmt.Sprintf(`report contains commands for generating end-user and machine-readable
outputs from a validated model.

Repository:
  %s`, repositoryURL)
}

func reportExamples() string {
	return `  seer report generate --config testdata/config/minimal.cue
  seer help report generate`
}

func reportGenerateLongHelp() string {
	return fmt.Sprintf(`generate runs the current v1 pipeline end to end:
  - CUE loading
  - validation
  - AHP weighting
  - TOPSIS scenario ranking
  - aggregation
  - report rendering

Use this when you want the actual rendered output instead of validation only.

Repository:
  %s`, repositoryURL)
}

func reportGenerateExamples() string {
	return `  seer report generate --config testdata/config/minimal.cue
  seer report generate --config testdata/config/valid_report.cue
  seer report generate --config path/to/modeldir`
}
