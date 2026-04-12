package pipeline

import (
	"fmt"
	"path/filepath"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func validateRequiredSections(config LoadedConfig) []domain.Diagnostic {
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

	diagnostics := make([]domain.Diagnostic, 0, len(requiredSections))
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

	return diagnostics
}

func validateCriteriaCatalog(diagnostics *[]domain.Diagnostic, criteria []CriterionConfig) map[string]CriterionConfig {
	criteriaByName := make(map[string]CriterionConfig, len(criteria))
	for criterionIndex, criterion := range criteria {
		if _, exists := criteriaByName[criterion.Name]; !exists {
			criteriaByName[criterion.Name] = criterion
		}

		if !isSupportedCriterionValueType(criterion.ValueType) {
			*diagnostics = append(*diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.unsupported_criterion_value_type",
				fmt.Sprintf("config.criteriaCatalog[%d].valueType", criterionIndex),
				domain.DiagnosticLocation{},
				fmt.Sprintf("unsupported criterion value type: %s", criterion.ValueType),
			))
		}

		if criterion.ValueType == "ordinal" && len(criterion.ScaleGuidance) == 0 {
			*diagnostics = append(*diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.ordinal_scale_guidance_missing",
				fmt.Sprintf("config.criteriaCatalog[%d].scaleGuidance", criterionIndex),
				domain.DiagnosticLocation{},
				fmt.Sprintf("ordinal criterion is missing scaleGuidance: %s", criterion.Name),
			))
		}
	}

	return criteriaByName
}

func validateScenarios(
	diagnostics *[]domain.Diagnostic,
	scenarios []ScenarioConfig,
	criteriaNames []string,
	criteriaByName map[string]CriterionConfig,
) map[string]scenarioValidationInfo {
	scenarioNameCounts := countScenarioNames(scenarios)
	scenarioInfos := make(map[string]scenarioValidationInfo, len(scenarios))

	for scenarioIndex, scenario := range scenarios {
		activeCriterionNames := make([]string, 0, len(scenario.ActiveCriteria))
		hasUnknownActiveCriteria := false
		for criterionIndex, ref := range scenario.ActiveCriteria {
			activeCriterionNames = append(activeCriterionNames, ref.CriterionName)
			if !hasName(criteriaNames, ref.CriterionName) {
				hasUnknownActiveCriteria = true
				*diagnostics = append(*diagnostics, domain.NewDiagnostic(
					domain.DiagnosticSeverityError,
					"validation.unknown_active_criterion",
					fmt.Sprintf("config.scenarios[%d].activeCriteria[%d].criterionName", scenarioIndex, criterionIndex),
					domain.DiagnosticLocation{},
					fmt.Sprintf("unknown criterion name in active criteria: %s", ref.CriterionName),
				))
			}
		}

		activeCriterionNames = domain.CanonicalNames(activeCriterionNames)
		if scenarioNameCounts[scenario.Name] == 1 && !hasUnknownActiveCriteria {
			scenarioInfos[scenario.Name] = scenarioValidationInfo{
				Index:                scenarioIndex,
				ActiveCriterionNames: activeCriterionNames,
			}
		}

		validateScenarioConstraints(diagnostics, scenarioIndex, scenario, activeCriterionNames, criteriaNames, criteriaByName)
		*diagnostics = append(*diagnostics, validateScenarioPairwiseComparisons(
			scenarioIndex,
			scenario,
			criteriaNames,
			activeCriterionNames,
		)...)
	}

	return scenarioInfos
}

func validateScenarioConstraints(
	diagnostics *[]domain.Diagnostic,
	scenarioIndex int,
	scenario ScenarioConfig,
	activeCriterionNames []string,
	criteriaNames []string,
	criteriaByName map[string]CriterionConfig,
) {
	for constraintIndex, constraint := range scenario.Constraints {
		if !hasName(criteriaNames, constraint.CriterionName) {
			*diagnostics = append(*diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.unknown_constraint_criterion",
				fmt.Sprintf("config.scenarios[%d].constraints[%d].criterionName", scenarioIndex, constraintIndex),
				domain.DiagnosticLocation{},
				fmt.Sprintf("unknown criterion name in constraints: %s", constraint.CriterionName),
			))
			continue
		}

		if !hasName(activeCriterionNames, constraint.CriterionName) {
			*diagnostics = append(*diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.inactive_constraint_criterion",
				fmt.Sprintf("config.scenarios[%d].constraints[%d].criterionName", scenarioIndex, constraintIndex),
				domain.DiagnosticLocation{},
				fmt.Sprintf("constraint references criterion not active in scenario: %s", constraint.CriterionName),
			))
			continue
		}

		criterion := criteriaByName[constraint.CriterionName]
		if !isSupportedConstraintOperator(criterion.ValueType, constraint.Operator) {
			*diagnostics = append(*diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.invalid_constraint_operator",
				fmt.Sprintf("config.scenarios[%d].constraints[%d].operator", scenarioIndex, constraintIndex),
				domain.DiagnosticLocation{},
				fmt.Sprintf("invalid constraint operator for %s criterion %s: %s", criterion.ValueType, constraint.CriterionName, constraint.Operator),
			))
		}

		if !isValidCriterionValue(criterion.ValueType, constraint.Value) {
			*diagnostics = append(*diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				validationConstraintValueCode(criterion.ValueType),
				fmt.Sprintf("config.scenarios[%d].constraints[%d].value", scenarioIndex, constraintIndex),
				domain.DiagnosticLocation{},
				validationConstraintValueMessage(criterion.ValueType, constraint.CriterionName),
			))
		}
	}
}

func validateEvaluations(
	diagnostics *[]domain.Diagnostic,
	evaluations []EvaluationConfig,
	scenarioInfos map[string]scenarioValidationInfo,
	scenarioNames []string,
	alternativeNames []string,
	criteriaNames []string,
	criteriaByName map[string]CriterionConfig,
) {
	seenEvaluationScenarios := make(map[string]struct{}, len(evaluations))
	for evaluationIndex, evaluation := range evaluations {
		if _, exists := seenEvaluationScenarios[evaluation.ScenarioName]; exists {
			*diagnostics = append(*diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.duplicate_evaluation_scenario",
				fmt.Sprintf("config.evaluations[%d].scenarioName", evaluationIndex),
				domain.DiagnosticLocation{},
				fmt.Sprintf("duplicate evaluation block for scenario: %s", evaluation.ScenarioName),
			))
		} else {
			seenEvaluationScenarios[evaluation.ScenarioName] = struct{}{}
		}

		scenarioInfo, hasScenario := scenarioInfos[evaluation.ScenarioName]
		if !hasName(scenarioNames, evaluation.ScenarioName) {
			*diagnostics = append(*diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.unknown_evaluation_scenario",
				fmt.Sprintf("config.evaluations[%d].scenarioName", evaluationIndex),
				domain.DiagnosticLocation{},
				fmt.Sprintf("unknown scenario name in evaluations: %s", evaluation.ScenarioName),
			))
		}

		validateScenarioEvaluationAlternatives(
			diagnostics,
			evaluationIndex,
			evaluation,
			hasScenario,
			scenarioInfo,
			alternativeNames,
			criteriaNames,
			criteriaByName,
		)
	}
}

func validateScenarioEvaluationAlternatives(
	diagnostics *[]domain.Diagnostic,
	evaluationIndex int,
	evaluation EvaluationConfig,
	hasScenario bool,
	scenarioInfo scenarioValidationInfo,
	alternativeNames []string,
	criteriaNames []string,
	criteriaByName map[string]CriterionConfig,
) {
	seenAlternatives := make(map[string]struct{}, len(evaluation.Evaluations))
	for alternativeIndex, alternative := range evaluation.Evaluations {
		if _, exists := seenAlternatives[alternative.AlternativeName]; exists {
			*diagnostics = append(*diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.duplicate_evaluation_alternative",
				fmt.Sprintf("config.evaluations[%d].evaluations[%d].alternativeName", evaluationIndex, alternativeIndex),
				domain.DiagnosticLocation{},
				fmt.Sprintf("duplicate alternative evaluation in scenario %s: %s", evaluation.ScenarioName, alternative.AlternativeName),
			))
		} else {
			seenAlternatives[alternative.AlternativeName] = struct{}{}
		}

		if !hasName(alternativeNames, alternative.AlternativeName) {
			*diagnostics = append(*diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.unknown_evaluation_alternative",
				fmt.Sprintf("config.evaluations[%d].evaluations[%d].alternativeName", evaluationIndex, alternativeIndex),
				domain.DiagnosticLocation{},
				fmt.Sprintf("unknown alternative name in evaluations: %s", alternative.AlternativeName),
			))
		}

		if hasScenario {
			*diagnostics = append(*diagnostics, validateAlternativeEvaluationValues(
				evaluationIndex,
				alternativeIndex,
				scenarioInfo,
				alternative,
				criteriaNames,
				criteriaByName,
				evaluation.ScenarioName,
			)...)
		}
	}
}

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

func isSupportedCriterionValueType(valueType string) bool {
	return valueType == "number" || valueType == "ordinal" || valueType == "boolean"
}

func isSupportedReportFormat(format string) bool {
	return format == "markdown" || format == "json" || format == "csv"
}

func countScenarioNames(scenarios []ScenarioConfig) map[string]int {
	counts := make(map[string]int, len(scenarios))
	for _, scenario := range scenarios {
		counts[scenario.Name]++
	}
	return counts
}
