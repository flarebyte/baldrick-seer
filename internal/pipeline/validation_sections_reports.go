package pipeline

import (
	"fmt"
	"path/filepath"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func validateReports(
	diagnostics *[]domain.Diagnostic,
	reports []ReportConfig,
	scenarioNames []string,
	alternativeNames []string,
	criteriaNames []string,
) {
	for reportIndex, report := range reports {
		if !isSupportedReportFormat(report.Format) {
			*diagnostics = append(*diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.unsupported_report_format",
				fmt.Sprintf("config.reports[%d].format", reportIndex),
				domain.DiagnosticLocation{},
				fmt.Sprintf("unsupported report format: %s", report.Format),
			))
		}

		if diagnostic, ok := validateReportFilepath(reportIndex, report.Filepath); ok {
			*diagnostics = append(*diagnostics, diagnostic)
		}

		if report.Focus != nil {
			validateReportFocusNames(diagnostics, reportIndex, "scenarioNames", report.Focus.ScenarioNames, scenarioNames, "validation.unknown_report_focus_scenario", "unknown scenario name in report focus: %s")
			validateReportFocusNames(diagnostics, reportIndex, "alternativeNames", report.Focus.AlternativeNames, alternativeNames, "validation.unknown_report_focus_alternative", "unknown alternative name in report focus: %s")
			validateReportFocusNames(diagnostics, reportIndex, "criterionNames", report.Focus.CriterionNames, criteriaNames, "validation.unknown_report_focus_criterion", "unknown criterion name in report focus: %s")
		}

		*diagnostics = append(*diagnostics, validateReportArguments(reportIndex, report)...)
	}
}

func validateReportFilepath(reportIndex int, value string) (domain.Diagnostic, bool) {
	if value == "" {
		return domain.Diagnostic{}, false
	}

	if filepath.IsAbs(value) {
		return domain.NewDiagnostic(
			domain.DiagnosticSeverityError,
			"validation.invalid_report_filepath",
			fmt.Sprintf("config.reports[%d].filepath", reportIndex),
			domain.DiagnosticLocation{},
			fmt.Sprintf("report filepath must be relative: %s", value),
		), true
	}

	cleaned := filepath.Clean(value)
	if cleaned == "." || cleaned == ".." {
		return domain.NewDiagnostic(
			domain.DiagnosticSeverityError,
			"validation.invalid_report_filepath",
			fmt.Sprintf("config.reports[%d].filepath", reportIndex),
			domain.DiagnosticLocation{},
			fmt.Sprintf("report filepath must name a file, got: %s", value),
		), true
	}

	return domain.Diagnostic{}, false
}

func validateAggregation(diagnostics *[]domain.Diagnostic, aggregation *AggregationConfig, scenarioNames []string) {
	if aggregation == nil || len(aggregation.ScenarioWeights) == 0 {
		return
	}

	aggregationScenarioNames := make([]string, 0, len(aggregation.ScenarioWeights))
	for name := range aggregation.ScenarioWeights {
		aggregationScenarioNames = append(aggregationScenarioNames, name)
	}

	for _, name := range domain.CanonicalNames(aggregationScenarioNames) {
		if !hasName(scenarioNames, name) {
			*diagnostics = append(*diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.unknown_aggregation_scenario",
				fmt.Sprintf("config.aggregation.scenarioWeights.%s", name),
				domain.DiagnosticLocation{},
				fmt.Sprintf("unknown scenario name in aggregation weights: %s", name),
			))
		}
	}
}

func validateReportFocusNames(
	diagnostics *[]domain.Diagnostic,
	reportIndex int,
	selectorName string,
	values []string,
	allowed []string,
	code string,
	messageFormat string,
) {
	for valueIndex, value := range values {
		if hasName(allowed, value) {
			continue
		}

		*diagnostics = append(*diagnostics, domain.NewDiagnostic(
			domain.DiagnosticSeverityError,
			code,
			fmt.Sprintf("config.reports[%d].focus.%s[%d]", reportIndex, selectorName, valueIndex),
			domain.DiagnosticLocation{},
			fmt.Sprintf(messageFormat, value),
		))
	}
}
