package pipeline

import "github.com/flarebyte/baldrick-seer/internal/domain"

func filterScenarioResultsForReport(report ReportConfig, scenarioResults []domain.ScenarioRankingResult) []domain.ScenarioRankingResult {
	if report.Focus == nil || len(report.Focus.ScenarioNames) == 0 {
		return domain.CanonicalScenarioResults(scenarioResults)
	}

	allowed := toAllowedNameSet(report.Focus.ScenarioNames)

	var filtered []domain.ScenarioRankingResult
	for _, scenarioResult := range domain.CanonicalScenarioResults(scenarioResults) {
		if _, exists := allowed[scenarioResult.ScenarioName]; exists {
			filtered = append(filtered, scenarioResult)
		}
	}
	return filtered
}

func filterScenarioAlternativesForReport(report ReportConfig, scenarioResults []domain.ScenarioRankingResult) []domain.ScenarioRankingResult {
	if report.Focus == nil || len(report.Focus.AlternativeNames) == 0 {
		return domain.CanonicalScenarioResults(scenarioResults)
	}

	allowed := toAllowedNameSet(report.Focus.AlternativeNames)
	filtered := make([]domain.ScenarioRankingResult, 0, len(scenarioResults))
	for _, scenarioResult := range domain.CanonicalScenarioResults(scenarioResults) {
		alternatives := make([]domain.RankedAlternative, 0, len(scenarioResult.RankedAlternatives))
		for _, alternative := range domain.CanonicalRankedAlternatives(scenarioResult.RankedAlternatives) {
			if _, exists := allowed[alternative.Name]; exists {
				alternatives = append(alternatives, alternative)
			}
		}
		filtered = append(filtered, domain.ScenarioRankingResult{
			ScenarioName:       scenarioResult.ScenarioName,
			RankedAlternatives: alternatives,
		})
	}
	return domain.CanonicalScenarioResults(filtered)
}

func filterFinalRankingForReport(report ReportConfig, finalRanking domain.AggregatedRankingResult) domain.AggregatedRankingResult {
	if report.Focus == nil || len(report.Focus.AlternativeNames) == 0 {
		return domain.CanonicalAggregatedRankingResult(finalRanking)
	}

	allowed := toAllowedNameSet(report.Focus.AlternativeNames)

	var filtered []domain.RankedAlternative
	for _, alternative := range domain.CanonicalAggregatedRankingResult(finalRanking).RankedAlternatives {
		if _, exists := allowed[alternative.Name]; exists {
			filtered = append(filtered, alternative)
		}
	}
	for index := range filtered {
		filtered[index].Rank = index + 1
	}
	return domain.AggregatedRankingResult{RankedAlternatives: filtered}
}

func filterScenarioWeightsForReport(report ReportConfig, scenarioWeights []ScenarioCriterionWeights) []ScenarioCriterionWeights {
	if len(scenarioWeights) == 0 {
		return nil
	}

	allowedScenarios, allowedCriteria := focusedScenarioAndCriterionNames(report)

	filtered := make([]ScenarioCriterionWeights, 0, len(scenarioWeights))
	for _, scenarioWeight := range canonicalScenarioWeights(scenarioWeights) {
		if len(allowedScenarios) > 0 {
			if _, exists := allowedScenarios[scenarioWeight.ScenarioName]; !exists {
				continue
			}
		}

		weights := make([]CriterionWeight, 0, len(scenarioWeight.CriterionWeights))
		for _, weight := range canonicalCriterionWeights(scenarioWeight.CriterionWeights) {
			if len(allowedCriteria) > 0 {
				if _, exists := allowedCriteria[weight.CriterionName]; !exists {
					continue
				}
			}
			weights = append(weights, weight)
		}

		filtered = append(filtered, ScenarioCriterionWeights{
			ScenarioName:     scenarioWeight.ScenarioName,
			CriterionWeights: weights,
		})
	}

	return canonicalScenarioWeights(filtered)
}
