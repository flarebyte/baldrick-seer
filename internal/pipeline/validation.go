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
		activeCriterionNames := make([]string, 0, len(scenario.ActiveCriteria))
		for criterionIndex, ref := range scenario.ActiveCriteria {
			activeCriterionNames = append(activeCriterionNames, ref.CriterionName)
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
		activeCriterionNames = domain.CanonicalNames(activeCriterionNames)

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
