package pipeline

import (
	"fmt"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

type DefaultModelValidator struct{}

type scenarioValidationInfo struct {
	Index                int
	ActiveCriterionNames []string
}

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
	scenarioNameCounts := countScenarioNames(config.Config.Scenarios)

	criteriaByName := make(map[string]CriterionConfig, len(config.Config.CriteriaCatalog))
	for criterionIndex, criterion := range config.Config.CriteriaCatalog {
		if _, exists := criteriaByName[criterion.Name]; !exists {
			criteriaByName[criterion.Name] = criterion
		}

		if !isSupportedCriterionValueType(criterion.ValueType) {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.unsupported_criterion_value_type",
				fmt.Sprintf("config.criteriaCatalog[%d].valueType", criterionIndex),
				domain.DiagnosticLocation{},
				fmt.Sprintf("unsupported criterion value type: %s", criterion.ValueType),
			))
		}

		if criterion.ValueType == "ordinal" && len(criterion.ScaleGuidance) == 0 {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.ordinal_scale_guidance_missing",
				fmt.Sprintf("config.criteriaCatalog[%d].scaleGuidance", criterionIndex),
				domain.DiagnosticLocation{},
				fmt.Sprintf("ordinal criterion is missing scaleGuidance: %s", criterion.Name),
			))
		}
	}

	scenarioInfos := make(map[string]scenarioValidationInfo, len(config.Config.Scenarios))
	for scenarioIndex, scenario := range config.Config.Scenarios {
		activeCriterionNames := make([]string, 0, len(scenario.ActiveCriteria))
		hasUnknownActiveCriteria := false
		for criterionIndex, ref := range scenario.ActiveCriteria {
			activeCriterionNames = append(activeCriterionNames, ref.CriterionName)
			if !hasName(criteriaNames, ref.CriterionName) {
				hasUnknownActiveCriteria = true
				diagnostics = append(diagnostics, domain.NewDiagnostic(
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

		diagnostics = append(diagnostics, validateScenarioPairwiseComparisons(
			scenarioIndex,
			scenario,
			criteriaNames,
			activeCriterionNames,
		)...)
	}

	seenEvaluationScenarios := make(map[string]struct{}, len(config.Config.Evaluations))
	for evaluationIndex, evaluation := range config.Config.Evaluations {
		if _, exists := seenEvaluationScenarios[evaluation.ScenarioName]; exists {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
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
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.unknown_evaluation_scenario",
				fmt.Sprintf("config.evaluations[%d].scenarioName", evaluationIndex),
				domain.DiagnosticLocation{},
				fmt.Sprintf("unknown scenario name in evaluations: %s", evaluation.ScenarioName),
			))
		}

		seenAlternatives := make(map[string]struct{}, len(evaluation.Evaluations))
		for alternativeIndex, alternative := range evaluation.Evaluations {
			if _, exists := seenAlternatives[alternative.AlternativeName]; exists {
				diagnostics = append(diagnostics, domain.NewDiagnostic(
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
				diagnostics = append(diagnostics, domain.NewDiagnostic(
					domain.DiagnosticSeverityError,
					"validation.unknown_evaluation_alternative",
					fmt.Sprintf("config.evaluations[%d].evaluations[%d].alternativeName", evaluationIndex, alternativeIndex),
					domain.DiagnosticLocation{},
					fmt.Sprintf("unknown alternative name in evaluations: %s", alternative.AlternativeName),
				))
			}

			if hasScenario {
				diagnostics = append(diagnostics, validateAlternativeEvaluationValues(
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

func validateAlternativeEvaluationValues(
	evaluationIndex int,
	alternativeIndex int,
	scenarioInfo scenarioValidationInfo,
	alternative AlternativeEvaluationConfig,
	criteriaNames []string,
	criteriaByName map[string]CriterionConfig,
	scenarioName string,
) []domain.Diagnostic {
	var diagnostics []domain.Diagnostic

	for _, criterionName := range scenarioInfo.ActiveCriterionNames {
		if _, exists := alternative.Values[criterionName]; !exists {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.missing_evaluation_value",
				fmt.Sprintf("config.evaluations[%d].evaluations[%d].values", evaluationIndex, alternativeIndex),
				domain.DiagnosticLocation{},
				fmt.Sprintf("missing value for active criterion in scenario %s: %s", scenarioName, criterionName),
			))
		}
	}

	valueNames := make([]string, 0, len(alternative.Values))
	for criterionName := range alternative.Values {
		valueNames = append(valueNames, criterionName)
	}

	for _, criterionName := range domain.CanonicalNames(valueNames) {
		value := alternative.Values[criterionName]
		valuePath := fmt.Sprintf("config.evaluations[%d].evaluations[%d].values.%s", evaluationIndex, alternativeIndex, criterionName)

		if !hasName(criteriaNames, criterionName) {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.unknown_evaluation_criterion",
				valuePath,
				domain.DiagnosticLocation{},
				fmt.Sprintf("unknown criterion name in evaluation values: %s", criterionName),
			))
			continue
		}

		if !hasName(scenarioInfo.ActiveCriterionNames, criterionName) {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.inactive_evaluation_criterion",
				valuePath,
				domain.DiagnosticLocation{},
				fmt.Sprintf("criterion value is not active in scenario %s: %s", scenarioName, criterionName),
			))
			continue
		}

		criterion := criteriaByName[criterionName]
		if !isSupportedCriterionValueKind(value.Kind) {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.unsupported_evaluation_value_kind",
				valuePath+".kind",
				domain.DiagnosticLocation{},
				fmt.Sprintf("unsupported evaluation value kind for criterion %s: %s", criterionName, value.Kind),
			))
			continue
		}

		if value.Kind != criterion.ValueType {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.evaluation_value_kind_mismatch",
				valuePath+".kind",
				domain.DiagnosticLocation{},
				fmt.Sprintf("evaluation value kind mismatch for criterion %s: want %s, got %s", criterionName, criterion.ValueType, value.Kind),
			))
			continue
		}

		if !isValidCriterionValue(value.Kind, value.Value) {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				validationValueSemanticCode(value.Kind),
				valuePath+".value",
				domain.DiagnosticLocation{},
				validationValueSemanticMessage(criterionName, value.Kind),
			))
		}
	}

	return diagnostics
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

func countScenarioNames(scenarios []ScenarioConfig) map[string]int {
	counts := make(map[string]int, len(scenarios))
	for _, scenario := range scenarios {
		counts[scenario.Name]++
	}
	return counts
}

func isSupportedCriterionValueKind(kind string) bool {
	return isSupportedCriterionValueType(kind)
}

func isValidCriterionValue(kind string, value any) bool {
	switch kind {
	case "number":
		return isNumericValue(value)
	case "ordinal":
		return isIntegerValue(value)
	case "boolean":
		_, ok := value.(bool)
		return ok
	default:
		return false
	}
}

func isNumericValue(value any) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, float32, float64:
		return true
	default:
		return false
	}
}

func isIntegerValue(value any) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64:
		return true
	default:
		return false
	}
}

func validationValueSemanticCode(kind string) string {
	switch kind {
	case "number":
		return "validation.invalid_number_value"
	case "ordinal":
		return "validation.invalid_ordinal_value"
	case "boolean":
		return "validation.invalid_boolean_value"
	default:
		return "validation.invalid_value"
	}
}

func validationValueSemanticMessage(criterionName string, kind string) string {
	switch kind {
	case "number":
		return fmt.Sprintf("number criterion value must be numeric: %s", criterionName)
	case "ordinal":
		return fmt.Sprintf("ordinal criterion value must be an integer: %s", criterionName)
	case "boolean":
		return fmt.Sprintf("boolean criterion value must be true or false: %s", criterionName)
	default:
		return fmt.Sprintf("invalid criterion value: %s", criterionName)
	}
}

func validateScenarioPairwiseComparisons(
	scenarioIndex int,
	scenario ScenarioConfig,
	criteriaNames []string,
	activeCriterionNames []string,
) []domain.Diagnostic {
	if scenario.Preferences == nil || scenario.Preferences.Method != "ahp_pairwise" {
		return nil
	}

	var diagnostics []domain.Diagnostic
	expectedPairs := expectedAHPPairs(activeCriterionNames)
	seenCanonical := map[string]int{}
	seenDirectional := map[string]int{}

	for comparisonIndex, comparison := range scenario.Preferences.Comparisons {
		morePath := fmt.Sprintf("config.scenarios[%d].preferences.comparisons[%d].moreImportantCriterionName", scenarioIndex, comparisonIndex)
		lessPath := fmt.Sprintf("config.scenarios[%d].preferences.comparisons[%d].lessImportantCriterionName", scenarioIndex, comparisonIndex)

		moreKnown := hasName(criteriaNames, comparison.MoreImportantCriterionName)
		lessKnown := hasName(criteriaNames, comparison.LessImportantCriterionName)
		if !moreKnown {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.unknown_pairwise_criterion",
				morePath,
				domain.DiagnosticLocation{},
				fmt.Sprintf("unknown criterion name in pairwise comparison: %s", comparison.MoreImportantCriterionName),
			))
		}
		if !lessKnown {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.unknown_pairwise_criterion",
				lessPath,
				domain.DiagnosticLocation{},
				fmt.Sprintf("unknown criterion name in pairwise comparison: %s", comparison.LessImportantCriterionName),
			))
		}

		if moreKnown && !hasName(activeCriterionNames, comparison.MoreImportantCriterionName) {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.inactive_pairwise_criterion",
				morePath,
				domain.DiagnosticLocation{},
				fmt.Sprintf("pairwise comparison references criterion not active in scenario: %s", comparison.MoreImportantCriterionName),
			))
		}
		if lessKnown && !hasName(activeCriterionNames, comparison.LessImportantCriterionName) {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.inactive_pairwise_criterion",
				lessPath,
				domain.DiagnosticLocation{},
				fmt.Sprintf("pairwise comparison references criterion not active in scenario: %s", comparison.LessImportantCriterionName),
			))
		}

		if comparison.MoreImportantCriterionName == comparison.LessImportantCriterionName {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.pairwise_self_comparison",
				morePath,
				domain.DiagnosticLocation{},
				fmt.Sprintf("pairwise comparison cannot compare criterion with itself: %s", comparison.MoreImportantCriterionName),
			))
			continue
		}

		if !moreKnown || !lessKnown ||
			!hasName(activeCriterionNames, comparison.MoreImportantCriterionName) ||
			!hasName(activeCriterionNames, comparison.LessImportantCriterionName) {
			continue
		}

		directionalKey := comparison.MoreImportantCriterionName + ">" + comparison.LessImportantCriterionName
		canonicalKey := canonicalPairKey(comparison.MoreImportantCriterionName, comparison.LessImportantCriterionName)
		if previousIndex, exists := seenDirectional[directionalKey]; exists {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.duplicate_pairwise_comparison",
				morePath,
				domain.DiagnosticLocation{},
				fmt.Sprintf("duplicate pairwise comparison for pair: %s (already defined at comparison %d)", canonicalKey, previousIndex),
			))
			continue
		}

		inverseKey := comparison.LessImportantCriterionName + ">" + comparison.MoreImportantCriterionName
		if previousIndex, exists := seenDirectional[inverseKey]; exists {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.inverse_duplicate_pairwise_comparison",
				morePath,
				domain.DiagnosticLocation{},
				fmt.Sprintf("inverse duplicate pairwise comparison for pair: %s (already defined at comparison %d)", canonicalKey, previousIndex),
			))
			continue
		}

		seenDirectional[directionalKey] = comparisonIndex
		seenCanonical[canonicalKey]++
	}

	if len(activeCriterionNames) > 1 {
		for _, pair := range expectedPairs {
			if seenCanonical[pair] == 0 {
				diagnostics = append(diagnostics, domain.NewDiagnostic(
					domain.DiagnosticSeverityError,
					"validation.missing_pairwise_comparison",
					fmt.Sprintf("config.scenarios[%d].preferences.comparisons", scenarioIndex),
					domain.DiagnosticLocation{},
					fmt.Sprintf("missing pairwise comparison for pair: %s", pair),
				))
			}
		}
	}

	return diagnostics
}

func expectedAHPPairs(activeCriterionNames []string) []string {
	if len(activeCriterionNames) < 2 {
		return nil
	}

	var pairs []string
	for leftIndex := 0; leftIndex < len(activeCriterionNames); leftIndex++ {
		for rightIndex := leftIndex + 1; rightIndex < len(activeCriterionNames); rightIndex++ {
			pairs = append(pairs, canonicalPairKey(activeCriterionNames[leftIndex], activeCriterionNames[rightIndex]))
		}
	}

	return pairs
}

func canonicalPairKey(left string, right string) string {
	if left < right {
		return left + "/" + right
	}
	return right + "/" + left
}
