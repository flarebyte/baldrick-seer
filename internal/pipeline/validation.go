package pipeline

import (
	"fmt"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

type DefaultModelValidator struct{}

func (DefaultModelValidator) ValidateModel(input ValidateModelInput) (ValidateModelOutput, error) {
	diagnostics := validateLoadedConfig(input.Config)
	if len(diagnostics) > 0 {
		return ValidateModelOutput{}, NewValidationDiagnosticsFailure(diagnostics, ErrValidationFailed)
	}

	reportDefinitions := make([]domain.ReportDefinition, 0, len(input.Config.Config.Reports))
	for _, report := range input.Config.Config.Reports {
		reportDefinitions = append(reportDefinitions, domain.ReportDefinition{
			Name:   report.Name,
			Title:  report.Title,
			Format: report.Format,
		})
	}

	return ValidateModelOutput{
		ValidatedModel: domain.CanonicalValidatedModelSummary(domain.ValidatedModelSummary{
			ConfigPath:        input.Config.Path,
			CriterionCount:    len(input.Config.Config.CriteriaCatalog),
			AlternativeCount:  len(input.Config.Config.Alternatives),
			ScenarioCount:     len(input.Config.Config.Scenarios),
			ReportDefinitions: reportDefinitions,
		}),
		ReportDefinitions: domain.CanonicalReportDefinitions(reportDefinitions),
	}, nil
}

func validateLoadedConfig(config LoadedConfig) []domain.Diagnostic {
	var diagnostics []domain.Diagnostic

	requiredSections := []struct {
		name    string
		present bool
		valid   bool
	}{
		{name: "problem", present: hasName(config.ConfigFields, "problem"), valid: config.Config != nil && config.Config.Problem != nil},
		{name: "reports", present: hasName(config.ConfigFields, "reports"), valid: config.Config != nil && len(config.Config.Reports) > 0},
		{name: "criteriaCatalog", present: hasName(config.ConfigFields, "criteriaCatalog"), valid: config.Config != nil && len(config.Config.CriteriaCatalog) > 0},
		{name: "alternatives", present: hasName(config.ConfigFields, "alternatives"), valid: config.Config != nil && len(config.Config.Alternatives) > 0},
		{name: "scenarios", present: hasName(config.ConfigFields, "scenarios"), valid: config.Config != nil && len(config.Config.Scenarios) > 0},
		{name: "evaluations", present: hasName(config.ConfigFields, "evaluations"), valid: config.Config != nil && len(config.Config.Evaluations) > 0},
		{name: "aggregation", present: hasName(config.ConfigFields, "aggregation"), valid: config.Config != nil && config.Config.Aggregation != nil},
	}

	for _, section := range requiredSections {
		if section.present && section.valid {
			continue
		}

		diagnostics = append(diagnostics, domain.NewDiagnostic(
			domain.DiagnosticSeverityError,
			"validation.section_missing",
			fmt.Sprintf("config.%s", section.name),
			domain.DiagnosticLocation{},
			fmt.Sprintf("missing required section: %s", section.name),
		))
	}

	if config.Config == nil {
		return domain.CanonicalDiagnostics(diagnostics)
	}

	criteriaNames := collectUniqueNames(
		&diagnostics,
		"config.criteriaCatalog",
		config.Config.CriteriaCatalog,
		func(item CriterionConfig) string { return item.Name },
		"validation.duplicate_criterion_name",
		"duplicate criterion name: %s",
	)
	alternativeNames := collectUniqueNames(
		&diagnostics,
		"config.alternatives",
		config.Config.Alternatives,
		func(item AlternativeConfig) string { return item.Name },
		"validation.duplicate_alternative_name",
		"duplicate alternative name: %s",
	)
	scenarioNames := collectUniqueNames(
		&diagnostics,
		"config.scenarios",
		config.Config.Scenarios,
		func(item ScenarioConfig) string { return item.Name },
		"validation.duplicate_scenario_name",
		"duplicate scenario name: %s",
	)
	collectUniqueNames(
		&diagnostics,
		"config.reports",
		config.Config.Reports,
		func(item ReportConfig) string { return item.Name },
		"validation.duplicate_report_name",
		"duplicate report name: %s",
	)

	for scenarioIndex, scenario := range config.Config.Scenarios {
		for criterionIndex, ref := range scenario.ActiveCriteria {
			if !hasName(criteriaNames, ref.CriterionName) {
				diagnostics = append(diagnostics, domain.NewDiagnostic(
					domain.DiagnosticSeverityError,
					"validation.unknown_active_criterion",
					fmt.Sprintf("config.scenarios[%d].activeCriteria[%d].criterionName", scenarioIndex, criterionIndex),
					domain.DiagnosticLocation{},
					fmt.Sprintf("unknown criterion name in active criteria: %s", ref.CriterionName),
				))
			}
		}

		for constraintIndex, constraint := range scenario.Constraints {
			if !hasName(criteriaNames, constraint.CriterionName) {
				diagnostics = append(diagnostics, domain.NewDiagnostic(
					domain.DiagnosticSeverityError,
					"validation.unknown_constraint_criterion",
					fmt.Sprintf("config.scenarios[%d].constraints[%d].criterionName", scenarioIndex, constraintIndex),
					domain.DiagnosticLocation{},
					fmt.Sprintf("unknown criterion name in constraints: %s", constraint.CriterionName),
				))
			}
		}
	}

	for evaluationIndex, evaluation := range config.Config.Evaluations {
		if !hasName(scenarioNames, evaluation.ScenarioName) {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.unknown_evaluation_scenario",
				fmt.Sprintf("config.evaluations[%d].scenarioName", evaluationIndex),
				domain.DiagnosticLocation{},
				fmt.Sprintf("unknown scenario name in evaluations: %s", evaluation.ScenarioName),
			))
		}

		for alternativeIndex, alternative := range evaluation.Evaluations {
			if !hasName(alternativeNames, alternative.AlternativeName) {
				diagnostics = append(diagnostics, domain.NewDiagnostic(
					domain.DiagnosticSeverityError,
					"validation.unknown_evaluation_alternative",
					fmt.Sprintf("config.evaluations[%d].evaluations[%d].alternativeName", evaluationIndex, alternativeIndex),
					domain.DiagnosticLocation{},
					fmt.Sprintf("unknown alternative name in evaluations: %s", alternative.AlternativeName),
				))
			}
		}
	}

	for reportIndex, report := range config.Config.Reports {
		if report.Focus == nil {
			continue
		}

		validateReportFocusNames(&diagnostics, reportIndex, "scenarioNames", report.Focus.ScenarioNames, scenarioNames, "validation.unknown_report_focus_scenario", "unknown scenario name in report focus: %s")
		validateReportFocusNames(&diagnostics, reportIndex, "alternativeNames", report.Focus.AlternativeNames, alternativeNames, "validation.unknown_report_focus_alternative", "unknown alternative name in report focus: %s")
		validateReportFocusNames(&diagnostics, reportIndex, "criterionNames", report.Focus.CriterionNames, criteriaNames, "validation.unknown_report_focus_criterion", "unknown criterion name in report focus: %s")
	}

	if config.Config.Aggregation != nil && len(config.Config.Aggregation.ScenarioWeights) > 0 {
		aggregationScenarioNames := make([]string, 0, len(config.Config.Aggregation.ScenarioWeights))
		for name := range config.Config.Aggregation.ScenarioWeights {
			aggregationScenarioNames = append(aggregationScenarioNames, name)
		}
		aggregationScenarioNames = domain.CanonicalNames(aggregationScenarioNames)

		for _, name := range aggregationScenarioNames {
			if !hasName(scenarioNames, name) {
				diagnostics = append(diagnostics, domain.NewDiagnostic(
					domain.DiagnosticSeverityError,
					"validation.unknown_aggregation_scenario",
					fmt.Sprintf("config.aggregation.scenarioWeights.%s", name),
					domain.DiagnosticLocation{},
					fmt.Sprintf("unknown scenario name in aggregation weights: %s", name),
				))
			}
		}
	}
	return domain.CanonicalDiagnostics(diagnostics)
}

func collectUniqueNames[T any](
	diagnostics *[]domain.Diagnostic,
	pathPrefix string,
	items []T,
	nameOf func(T) string,
	code string,
	messageFormat string,
) []string {
	seen := make(map[string]struct{}, len(items))
	var names []string

	for index, item := range items {
		name := nameOf(item)
		names = append(names, name)
		if _, exists := seen[name]; exists {
			*diagnostics = append(*diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				code,
				fmt.Sprintf("%s[%d].name", pathPrefix, index),
				domain.DiagnosticLocation{},
				fmt.Sprintf(messageFormat, name),
			))
			continue
		}

		seen[name] = struct{}{}
	}

	return domain.CanonicalNames(names)
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

func hasName(names []string, target string) bool {
	for _, name := range names {
		if name == target {
			return true
		}
	}

	return false
}
