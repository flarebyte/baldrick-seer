package pipeline

import "testing"

func TestDefaultModelValidatorReportDefinitionValidation(t *testing.T) {
	t.Parallel()

	runValidatorFailureCases(t, []validatorFailureCase{
		{
			name: "valid markdown report definition",
			config: validLoadedConfigWithReports([]ReportConfig{{
				Name:      "summary",
				Title:     "Summary",
				Format:    "markdown",
				Arguments: []string{"include-scenarios=all", "top-alternatives=2", "include-scores=true"},
			}}),
		},
		{
			name: "valid json report definition",
			config: validLoadedConfigWithReports([]ReportConfig{{
				Name:      "summary",
				Title:     "Summary",
				Format:    "json",
				Arguments: []string{"include-evidence=false", "pretty=true"},
			}}),
		},
		{
			name: "valid csv report definition",
			config: validLoadedConfigWithReports([]ReportConfig{{
				Name:      "summary",
				Title:     "Summary",
				Format:    "csv",
				Filepath:  "artifacts/summary.csv",
				Arguments: []string{"columns=scenario,alternative,score,rank", "header=true"},
			}}),
		},
		{
			name: "valid relative report filepath with parent directory segments",
			config: validLoadedConfigWithReports([]ReportConfig{{
				Name:     "summary",
				Title:    "Summary",
				Format:   "markdown",
				Filepath: "../decision/data/summary.md",
			}}),
		},
		{
			name: "unsupported format",
			config: validLoadedConfigWithReports([]ReportConfig{{
				Name:   "summary",
				Title:  "Summary",
				Format: "html",
			}}),
			wantCodes:   []string{"validation.unsupported_report_format"},
			wantMessage: "unsupported report format: html",
		},
		{
			name: "unknown scenario in report focus",
			config: validLoadedConfigWithReports([]ReportConfig{{
				Name:   "summary",
				Title:  "Summary",
				Format: "markdown",
				Focus:  &ReportFocus{ScenarioNames: []string{"missing"}},
			}}),
			wantCodes:   []string{"validation.unknown_report_focus_scenario"},
			wantMessage: "unknown scenario name in report focus: missing",
		},
		{
			name: "unknown alternative in report focus",
			config: validLoadedConfigWithReports([]ReportConfig{{
				Name:   "summary",
				Title:  "Summary",
				Format: "markdown",
				Focus:  &ReportFocus{AlternativeNames: []string{"missing"}},
			}}),
			wantCodes:   []string{"validation.unknown_report_focus_alternative"},
			wantMessage: "unknown alternative name in report focus: missing",
		},
		{
			name: "unknown criterion in report focus",
			config: validLoadedConfigWithReports([]ReportConfig{{
				Name:   "summary",
				Title:  "Summary",
				Format: "markdown",
				Focus:  &ReportFocus{CriterionNames: []string{"missing"}},
			}}),
			wantCodes:   []string{"validation.unknown_report_focus_criterion"},
			wantMessage: "unknown criterion name in report focus: missing",
		},
		{
			name: "malformed argument without equals",
			config: validLoadedConfigWithReports([]ReportConfig{{
				Name:      "summary",
				Title:     "Summary",
				Format:    "markdown",
				Arguments: []string{"include-scenarios"},
			}}),
			wantCodes:   []string{"validation.malformed_report_argument"},
			wantMessage: "report argument must use key=value form: include-scenarios",
		},
		{
			name: "unknown argument key",
			config: validLoadedConfigWithReports([]ReportConfig{{
				Name:      "summary",
				Title:     "Summary",
				Format:    "markdown",
				Arguments: []string{"unknown=true"},
			}}),
			wantCodes:   []string{"validation.unknown_report_argument"},
			wantMessage: "unknown report argument key: unknown",
		},
		{
			name: "duplicate argument key",
			config: validLoadedConfigWithReports([]ReportConfig{{
				Name:      "summary",
				Title:     "Summary",
				Format:    "markdown",
				Arguments: []string{"include-scenarios=all", "include-scenarios=focused"},
			}}),
			wantCodes:   []string{"validation.duplicate_report_argument"},
			wantMessage: "duplicate report argument key: include-scenarios",
		},
		{
			name: "format specific key used with wrong report format",
			config: validLoadedConfigWithReports([]ReportConfig{{
				Name:      "summary",
				Title:     "Summary",
				Format:    "json",
				Arguments: []string{"header=true"},
			}}),
			wantCodes:   []string{"validation.incompatible_report_argument"},
			wantMessage: "report argument key header is not allowed for format json",
		},
		{
			name: "invalid argument value",
			config: validLoadedConfigWithReports([]ReportConfig{{
				Name:      "summary",
				Title:     "Summary",
				Format:    "csv",
				Arguments: []string{"header=yes"},
			}}),
			wantCodes:   []string{"validation.invalid_report_argument_value"},
			wantMessage: "invalid value for report argument header: yes",
		},
		{
			name: "absolute report filepath is invalid",
			config: validLoadedConfigWithReports([]ReportConfig{{
				Name:     "summary",
				Title:    "Summary",
				Format:   "markdown",
				Filepath: "/tmp/summary.md",
			}}),
			wantCodes:   []string{"validation.invalid_report_filepath"},
			wantMessage: "report filepath must be relative: /tmp/summary.md",
		},
		{
			name: "filepath resolving to current directory is invalid",
			config: validLoadedConfigWithReports([]ReportConfig{{
				Name:     "summary",
				Title:    "Summary",
				Format:   "markdown",
				Filepath: ".",
			}}),
			wantCodes:   []string{"validation.invalid_report_filepath"},
			wantMessage: "report filepath must name a file, got: .",
		},
		{
			name: "filepath resolving to parent directory only is invalid",
			config: validLoadedConfigWithReports([]ReportConfig{{
				Name:     "summary",
				Title:    "Summary",
				Format:   "markdown",
				Filepath: "..",
			}}),
			wantCodes:   []string{"validation.invalid_report_filepath"},
			wantMessage: "report filepath must name a file, got: ..",
		},
		{
			name: "valid mixed examples using spec argument patterns",
			config: validLoadedConfigWithReports([]ReportConfig{
				{
					Name:      "markdown-summary",
					Title:     "Summary",
					Format:    "markdown",
					Arguments: []string{"include-scenarios=all", "top-alternatives=3", "explain=true"},
				},
				{
					Name:      "json-summary",
					Title:     "Summary JSON",
					Format:    "json",
					Arguments: []string{"include-evidence=true", "include-weights=true", "pretty=true"},
				},
				{
					Name:      "csv-summary",
					Title:     "Summary CSV",
					Format:    "csv",
					Arguments: []string{"columns=scenario,alternative,criterion,value,score,rank", "header=true"},
				},
			}),
		},
	})
}
