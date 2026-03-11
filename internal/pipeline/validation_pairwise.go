package pipeline

import (
	"fmt"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

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
		for _, pair := range expectedAHPPairs(activeCriterionNames) {
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
