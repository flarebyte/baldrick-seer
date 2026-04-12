package pipeline

import (
	"strings"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func markdownPreferenceComparisons(report ReportConfig, config *ExecutionConfig) []string {
	if config == nil {
		return nil
	}

	allowedScenarios, allowedCriteria := focusedScenarioAndCriterionNames(report)

	var output []string
	for _, scenario := range config.Scenarios {
		if len(allowedScenarios) > 0 {
			if _, exists := allowedScenarios[scenario.Name]; !exists {
				continue
			}
		}
		if scenario.Preferences == nil || len(scenario.Preferences.Comparisons) == 0 {
			continue
		}
		for _, comparison := range scenario.Preferences.Comparisons {
			if len(allowedCriteria) > 0 {
				if _, exists := allowedCriteria[comparison.MoreImportantCriterionName]; !exists {
					continue
				}
				if _, exists := allowedCriteria[comparison.LessImportantCriterionName]; !exists {
					continue
				}
			}
			var builder strings.Builder
			builder.WriteString(scenarioLabel(scenario))
			builder.WriteString(": ")
			builder.WriteString(criterionLabelByName(config, comparison.MoreImportantCriterionName))
			builder.WriteString(" over ")
			builder.WriteString(criterionLabelByName(config, comparison.LessImportantCriterionName))
			builder.WriteString(" (strength ")
			builder.WriteString(formatScore(comparison.Strength))
			builder.WriteString(")")
			output = append(output, builder.String())
		}
	}
	return output
}

func focusedScenarioAndCriterionNames(report ReportConfig) (map[string]struct{}, map[string]struct{}) {
	return allowedFocusedNames(report.Focus, func(focus *ReportFocus) []string {
			return focus.ScenarioNames
		}),
		allowedFocusedNames(report.Focus, func(focus *ReportFocus) []string {
			return focus.CriterionNames
		})
}

func orderedScenarioResultsForMarkdown(report ReportConfig, config *ExecutionConfig, scenarioResults []domain.ScenarioRankingResult) []domain.ScenarioRankingResult {
	resultByName := make(map[string]domain.ScenarioRankingResult, len(scenarioResults))
	for _, result := range scenarioResults {
		resultByName[result.ScenarioName] = result
	}

	ordered := make([]domain.ScenarioRankingResult, 0, len(scenarioResults))
	for _, scenario := range filteredScenariosForMarkdown(report, config) {
		result, exists := resultByName[scenario.Name]
		if !exists {
			continue
		}
		ordered = append(ordered, result)
		delete(resultByName, scenario.Name)
	}

	for _, result := range domain.CanonicalScenarioResults(scenarioResults) {
		if _, exists := resultByName[result.ScenarioName]; !exists {
			continue
		}
		ordered = append(ordered, result)
	}

	return ordered
}

func filteredAlternativesForMarkdown(report ReportConfig, config *ExecutionConfig) []AlternativeConfig {
	if config == nil {
		return nil
	}
	allowedAlternatives := allowedFocusedNames(report.Focus, func(focus *ReportFocus) []string {
		return focus.AlternativeNames
	})

	output := make([]AlternativeConfig, 0, len(config.Alternatives))
	for _, alternative := range config.Alternatives {
		if len(allowedAlternatives) > 0 {
			if _, exists := allowedAlternatives[alternative.Name]; !exists {
				continue
			}
		}
		output = append(output, alternative)
	}
	return output
}

func filteredScenariosForMarkdown(report ReportConfig, config *ExecutionConfig) []ScenarioConfig {
	if config == nil {
		return nil
	}
	allowedScenarios := allowedFocusedNames(report.Focus, func(focus *ReportFocus) []string {
		return focus.ScenarioNames
	})

	output := make([]ScenarioConfig, 0, len(config.Scenarios))
	for _, scenario := range config.Scenarios {
		if len(allowedScenarios) > 0 {
			if _, exists := allowedScenarios[scenario.Name]; !exists {
				continue
			}
		}
		output = append(output, scenario)
	}
	return output
}

func filteredCriteriaForMarkdown(report ReportConfig, config *ExecutionConfig) []CriterionConfig {
	if config == nil {
		return nil
	}
	allowedCriteria := allowedFocusedNames(report.Focus, func(focus *ReportFocus) []string {
		return focus.CriterionNames
	})

	output := make([]CriterionConfig, 0, len(config.CriteriaCatalog))
	for _, criterion := range config.CriteriaCatalog {
		if len(allowedCriteria) > 0 {
			if _, exists := allowedCriteria[criterion.Name]; !exists {
				continue
			}
		}
		output = append(output, criterion)
	}
	return output
}

func findScenarioEvaluation(report ReportConfig, config *ExecutionConfig, scenarioName string) (EvaluationConfig, bool) {
	if config == nil {
		return EvaluationConfig{}, false
	}
	allowedScenarios := allowedFocusedNames(report.Focus, func(focus *ReportFocus) []string {
		return focus.ScenarioNames
	})
	for _, evaluation := range config.Evaluations {
		if evaluation.ScenarioName != scenarioName {
			continue
		}
		if len(allowedScenarios) > 0 {
			if _, exists := allowedScenarios[evaluation.ScenarioName]; !exists {
				return EvaluationConfig{}, false
			}
		}
		return evaluation, true
	}
	return EvaluationConfig{}, false
}

func orderedAlternativeEvaluations(report ReportConfig, config *ExecutionConfig, evaluations []AlternativeEvaluationConfig) []AlternativeEvaluationConfig {
	allowedAlternatives := allowedFocusedNames(report.Focus, func(focus *ReportFocus) []string {
		return focus.AlternativeNames
	})
	evaluationByName := make(map[string]AlternativeEvaluationConfig, len(evaluations))
	for _, evaluation := range evaluations {
		if len(allowedAlternatives) > 0 {
			if _, exists := allowedAlternatives[evaluation.AlternativeName]; !exists {
				continue
			}
		}
		evaluationByName[evaluation.AlternativeName] = evaluation
	}

	ordered := make([]AlternativeEvaluationConfig, 0, len(evaluationByName))
	for _, alternative := range filteredAlternativesForMarkdown(report, config) {
		evaluation, exists := evaluationByName[alternative.Name]
		if !exists {
			continue
		}
		ordered = append(ordered, evaluation)
		delete(evaluationByName, alternative.Name)
	}

	for _, evaluation := range canonicalAlternativeEvaluations(evaluations) {
		if _, exists := evaluationByName[evaluation.AlternativeName]; !exists {
			continue
		}
		ordered = append(ordered, evaluation)
	}
	return ordered
}

type markdownCriterionValueRecord struct {
	Name     string
	Label    string
	Rendered string
}

func orderedCriterionValueRecords(criteria []CriterionConfig, values map[string]CriterionValue) []markdownCriterionValueRecord {
	if len(values) == 0 {
		return nil
	}

	ordered := make([]markdownCriterionValueRecord, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, criterion := range criteria {
		value, exists := values[criterion.Name]
		if !exists {
			continue
		}
		ordered = append(ordered, markdownCriterionValueRecord{
			Name:     criterion.Name,
			Label:    criterionLabel(criterion),
			Rendered: renderCriterionValue(value.Value),
		})
		seen[criterion.Name] = struct{}{}
	}

	for _, criterionName := range domain.CanonicalNames(valueNames(values)) {
		if _, exists := seen[criterionName]; exists {
			continue
		}
		ordered = append(ordered, markdownCriterionValueRecord{
			Name:     criterionName,
			Label:    criterionName,
			Rendered: renderCriterionValue(values[criterionName].Value),
		})
	}
	return ordered
}
