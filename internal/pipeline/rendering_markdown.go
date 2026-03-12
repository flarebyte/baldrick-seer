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

func writeMarkdownScenarioExplanation(builder *strings.Builder, scenarioName string, scenarioWeights []ScenarioCriterionWeights) {
	for _, scenarioWeight := range canonicalScenarioWeights(scenarioWeights) {
		if scenarioWeight.ScenarioName != scenarioName {
			continue
		}
		builder.WriteString("Weights: ")
		for index, weight := range canonicalCriterionWeights(scenarioWeight.CriterionWeights) {
			if index > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(weight.CriterionName)
			builder.WriteString("=")
			builder.WriteString(formatScore(weight.Weight))
		}
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
