package pipeline

import (
	"strconv"
	"strings"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func renderMarkdownReport(
	report ReportConfig,
	config *ExecutionConfig,
	scenarioResults []domain.ScenarioRankingResult,
	finalRanking domain.AggregatedRankingResult,
	aggregation *AggregationConfig,
	scenarioWeights []ScenarioCriterionWeights,
	aggregationScenarioWeights map[string]float64,
) string {
	switch reportArgumentValue(report.Arguments, "detail", "brief") {
	case "standard":
		return renderMarkdownStandardReport(report, config, scenarioResults, finalRanking, aggregation, scenarioWeights, aggregationScenarioWeights)
	case "full":
		return renderMarkdownFullReport(report, config, scenarioResults, finalRanking, aggregation, scenarioWeights, aggregationScenarioWeights)
	default:
		return renderMarkdownBriefReport(report, config, scenarioResults, finalRanking, aggregation, scenarioWeights, aggregationScenarioWeights)
	}
}

func renderMarkdownBriefReport(
	report ReportConfig,
	config *ExecutionConfig,
	scenarioResults []domain.ScenarioRankingResult,
	finalRanking domain.AggregatedRankingResult,
	aggregation *AggregationConfig,
	scenarioWeights []ScenarioCriterionWeights,
	aggregationScenarioWeights map[string]float64,
) string {
	includeScores := reportArgumentValue(report.Arguments, "include-scores", "true") == "true"
	topAlternatives := reportArgumentInt(report.Arguments, "top-alternatives")
	includeExplain := reportArgumentValue(report.Arguments, "explain", "true") == "true"

	var builder strings.Builder
	builder.WriteString("# ")
	builder.WriteString(report.Title)
	builder.WriteString("\n\n")
	builder.WriteString("Problem: ")
	builder.WriteString(problemName(config))
	builder.WriteString("\n\n")

	builder.WriteString("## Scenarios\n")
	for _, scenarioResult := range scenarioResults {
		builder.WriteString("\n### ")
		builder.WriteString(scenarioResult.ScenarioName)
		builder.WriteString("\n")
		if includeExplain {
			writeMarkdownScenarioExplanation(&builder, scenarioResult.ScenarioName, scenarioWeights)
		}
		rows := limitRankedAlternatives(scenarioResult.RankedAlternatives, topAlternatives)
		for _, alternative := range rows {
			writeMarkdownAlternative(&builder, alternative, includeScores)
		}
	}

	if includeExplain {
		writeMarkdownAggregationExplanation(&builder, aggregation, aggregationScenarioWeights)
	}

	builder.WriteString("\n## Final Ranking\n")
	if len(finalRanking.RankedAlternatives) == 0 {
		builder.WriteString("\n(none)\n")
		return builder.String()
	}

	for _, alternative := range limitRankedAlternatives(finalRanking.RankedAlternatives, topAlternatives) {
		writeMarkdownAlternative(&builder, alternative, includeScores)
	}

	return builder.String()
}

func renderMarkdownStandardReport(
	report ReportConfig,
	config *ExecutionConfig,
	scenarioResults []domain.ScenarioRankingResult,
	finalRanking domain.AggregatedRankingResult,
	aggregation *AggregationConfig,
	scenarioWeights []ScenarioCriterionWeights,
	aggregationScenarioWeights map[string]float64,
) string {
	includeScores := reportArgumentValue(report.Arguments, "include-scores", "true") == "true"
	topAlternatives := reportArgumentInt(report.Arguments, "top-alternatives")

	var builder strings.Builder
	builder.WriteString("# ")
	builder.WriteString(report.Title)
	builder.WriteString("\n\n")

	builder.WriteString("## Problem\n\n")
	builder.WriteString("- Name: ")
	builder.WriteString(problemName(config))
	builder.WriteString("\n")
	builder.WriteString("- Report: ")
	builder.WriteString(report.Name)
	builder.WriteString("\n")

	builder.WriteString("\n## Alternatives\n")
	for _, alternative := range canonicalAlternatives(config.Alternatives) {
		builder.WriteString("- ")
		builder.WriteString(alternative.Name)
		builder.WriteString("\n")
	}

	builder.WriteString("\n## Scenarios\n")
	for _, scenario := range canonicalScenarios(config.Scenarios) {
		builder.WriteString("- ")
		builder.WriteString(scenario.Name)
		builder.WriteString("\n")
	}

	builder.WriteString("\n## Criteria Weights\n")
	if len(scenarioWeights) == 0 {
		builder.WriteString("- (none)\n")
	} else {
		for _, scenarioWeight := range canonicalScenarioWeights(scenarioWeights) {
			builder.WriteString("- ")
			builder.WriteString(scenarioWeight.ScenarioName)
			builder.WriteString(": ")
			writeMarkdownInlineWeights(&builder, canonicalCriterionWeights(scenarioWeight.CriterionWeights))
			builder.WriteString("\n")
		}
	}

	builder.WriteString("\n## Scenario Rankings\n")
	for _, scenarioResult := range scenarioResults {
		builder.WriteString("\n### ")
		builder.WriteString(scenarioResult.ScenarioName)
		builder.WriteString("\n")
		for _, alternative := range limitRankedAlternatives(scenarioResult.RankedAlternatives, topAlternatives) {
			writeMarkdownAlternative(&builder, alternative, includeScores)
		}
	}

	builder.WriteString("\n## Final Ranking\n")
	if len(finalRanking.RankedAlternatives) == 0 {
		builder.WriteString("\n(none)\n")
	} else {
		for _, alternative := range limitRankedAlternatives(finalRanking.RankedAlternatives, topAlternatives) {
			writeMarkdownAlternative(&builder, alternative, includeScores)
		}
	}

	builder.WriteString("\n## Notes and Tradeoffs\n")
	writeMarkdownAggregationNotes(&builder, aggregation, aggregationScenarioWeights)
	writeMarkdownExclusionNotes(&builder, scenarioResults)

	return builder.String()
}

func renderMarkdownFullReport(
	report ReportConfig,
	config *ExecutionConfig,
	scenarioResults []domain.ScenarioRankingResult,
	finalRanking domain.AggregatedRankingResult,
	aggregation *AggregationConfig,
	scenarioWeights []ScenarioCriterionWeights,
	aggregationScenarioWeights map[string]float64,
) string {
	var builder strings.Builder
	builder.WriteString(renderMarkdownStandardReport(report, config, scenarioResults, finalRanking, aggregation, scenarioWeights, aggregationScenarioWeights))
	builder.WriteString("\n## Detailed Scenario Notes\n")
	for _, scenarioResult := range scenarioResults {
		builder.WriteString("\n### ")
		builder.WriteString(scenarioResult.ScenarioName)
		builder.WriteString("\n")
		writeMarkdownScenarioDetailNotes(&builder, scenarioResult)
	}
	builder.WriteString("\n## Aggregation Notes\n")
	writeMarkdownAggregationDetailNotes(&builder, aggregation, aggregationScenarioWeights, finalRanking)
	return builder.String()
}

func writeMarkdownScenarioExplanation(builder *strings.Builder, scenarioName string, scenarioWeights []ScenarioCriterionWeights) {
	for _, scenarioWeight := range canonicalScenarioWeights(scenarioWeights) {
		if scenarioWeight.ScenarioName != scenarioName {
			continue
		}
		builder.WriteString("Weights: ")
		writeMarkdownInlineWeights(builder, canonicalCriterionWeights(scenarioWeight.CriterionWeights))
		builder.WriteString("\n")
		return
	}
}

func writeMarkdownAggregationExplanation(builder *strings.Builder, aggregation *AggregationConfig, aggregationScenarioWeights map[string]float64) {
	builder.WriteString("\n## Aggregation\n")
	if aggregation == nil {
		builder.WriteString("\nMethod: \n")
		return
	}

	builder.WriteString("\nMethod: ")
	builder.WriteString(aggregation.Method)
	builder.WriteString("\n")
	if len(aggregationScenarioWeights) == 0 {
		return
	}
	builder.WriteString("Scenario weights:\n")
	for _, scenarioName := range orderedWeightNames(aggregationScenarioWeights) {
		builder.WriteString("- ")
		builder.WriteString(scenarioName)
		builder.WriteString(": ")
		builder.WriteString(formatScore(aggregationScenarioWeights[scenarioName]))
		builder.WriteString("\n")
	}
}

func writeMarkdownAlternative(builder *strings.Builder, alternative domain.RankedAlternative, includeScores bool) {
	builder.WriteString("- ")
	if alternative.Excluded {
		builder.WriteString(alternative.Name)
		builder.WriteString(": excluded")
		if alternative.ExclusionReason != "" {
			builder.WriteString(" (")
			builder.WriteString(alternative.ExclusionReason)
			builder.WriteString(")")
		}
		builder.WriteString("\n")
		return
	}

	builder.WriteString(strconv.Itoa(alternative.Rank))
	builder.WriteString(". ")
	builder.WriteString(alternative.Name)
	if includeScores {
		builder.WriteString(" (")
		builder.WriteString(formatScore(alternative.Score))
		builder.WriteString(")")
	}
	builder.WriteString("\n")
}

func writeMarkdownInlineWeights(builder *strings.Builder, weights []CriterionWeight) {
	for index, weight := range weights {
		if index > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(weight.CriterionName)
		builder.WriteString("=")
		builder.WriteString(formatScore(weight.Weight))
	}
}

func writeMarkdownAggregationNotes(builder *strings.Builder, aggregation *AggregationConfig, aggregationScenarioWeights map[string]float64) {
	if aggregation == nil {
		builder.WriteString("- Aggregation method: \n")
		return
	}

	builder.WriteString("- Aggregation method: ")
	builder.WriteString(aggregation.Method)
	builder.WriteString("\n")
	if len(aggregationScenarioWeights) == 0 {
		builder.WriteString("- Scenario weights: none\n")
		return
	}
	builder.WriteString("- Scenario weights: ")
	for index, scenarioName := range orderedWeightNames(aggregationScenarioWeights) {
		if index > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(scenarioName)
		builder.WriteString("=")
		builder.WriteString(formatScore(aggregationScenarioWeights[scenarioName]))
	}
	builder.WriteString("\n")
}

func writeMarkdownExclusionNotes(builder *strings.Builder, scenarioResults []domain.ScenarioRankingResult) {
	found := false
	for _, scenarioResult := range scenarioResults {
		for _, alternative := range scenarioResult.RankedAlternatives {
			if !alternative.Excluded {
				continue
			}
			if !found {
				builder.WriteString("- Exclusions:\n")
				found = true
			}
			builder.WriteString("  ")
			builder.WriteString(scenarioResult.ScenarioName)
			builder.WriteString(": ")
			builder.WriteString(alternative.Name)
			if alternative.ExclusionReason != "" {
				builder.WriteString(" (")
				builder.WriteString(alternative.ExclusionReason)
				builder.WriteString(")")
			}
			builder.WriteString("\n")
		}
	}
	if !found {
		builder.WriteString("- Exclusions: none\n")
	}
}

func writeMarkdownScenarioDetailNotes(builder *strings.Builder, scenarioResult domain.ScenarioRankingResult) {
	rankedCount := 0
	excludedCount := 0
	for _, alternative := range scenarioResult.RankedAlternatives {
		if alternative.Excluded {
			excludedCount++
			continue
		}
		rankedCount++
	}
	builder.WriteString("- Ranked alternatives: ")
	builder.WriteString(strconv.Itoa(rankedCount))
	builder.WriteString("\n")
	builder.WriteString("- Excluded alternatives: ")
	builder.WriteString(strconv.Itoa(excludedCount))
	builder.WriteString("\n")
	if len(scenarioResult.RankedAlternatives) > 0 && !scenarioResult.RankedAlternatives[0].Excluded {
		builder.WriteString("- Leading alternative: ")
		builder.WriteString(scenarioResult.RankedAlternatives[0].Name)
		builder.WriteString("\n")
	}
}

func writeMarkdownAggregationDetailNotes(builder *strings.Builder, aggregation *AggregationConfig, aggregationScenarioWeights map[string]float64, finalRanking domain.AggregatedRankingResult) {
	if aggregation == nil {
		builder.WriteString("- Aggregation method unavailable\n")
	} else {
		builder.WriteString("- Aggregation method: ")
		builder.WriteString(aggregation.Method)
		builder.WriteString("\n")
	}
	builder.WriteString("- Participating scenarios: ")
	builder.WriteString(strconv.Itoa(len(aggregationScenarioWeights)))
	builder.WriteString("\n")
	builder.WriteString("- Final eligible alternatives: ")
	builder.WriteString(strconv.Itoa(len(finalRanking.RankedAlternatives)))
	builder.WriteString("\n")
}

func canonicalAlternatives(alternatives []AlternativeConfig) []AlternativeConfig {
	output := append([]AlternativeConfig(nil), alternatives...)
	sortAlternatives(output)
	return output
}

func canonicalScenarios(scenarios []ScenarioConfig) []ScenarioConfig {
	output := append([]ScenarioConfig(nil), scenarios...)
	sortScenarios(output)
	return output
}

func sortAlternatives(alternatives []AlternativeConfig) {
	for index := 0; index < len(alternatives); index++ {
		for next := index + 1; next < len(alternatives); next++ {
			if alternatives[next].Name < alternatives[index].Name {
				alternatives[index], alternatives[next] = alternatives[next], alternatives[index]
			}
		}
	}
}

func sortScenarios(scenarios []ScenarioConfig) {
	for index := 0; index < len(scenarios); index++ {
		for next := index + 1; next < len(scenarios); next++ {
			if scenarios[next].Name < scenarios[index].Name {
				scenarios[index], scenarios[next] = scenarios[next], scenarios[index]
			}
		}
	}
}
