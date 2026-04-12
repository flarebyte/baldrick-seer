package pipeline

import (
	"strconv"
	"strings"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func writeMarkdownReportTitle(builder *strings.Builder, title string) {
	builder.WriteString("# ")
	builder.WriteString(title)
	builder.WriteString("\n")
}

func writeMarkdownProblemSection(builder *strings.Builder, config *ExecutionConfig) {
	if config == nil || config.Problem == nil {
		return
	}

	problem := config.Problem
	hasContent := problem.Title != "" || problem.Name != "" || problem.Goal != "" || problem.Description != "" || len(problem.Notes) > 0
	if !hasContent {
		return
	}

	builder.WriteString("\n## Problem\n")
	if title := firstNonEmpty(problem.Title, problem.Name); title != "" {
		builder.WriteString("\n### Title\n")
		builder.WriteString(title)
		builder.WriteString("\n")
	}
	if problem.Goal != "" {
		builder.WriteString("\n### Goal\n")
		builder.WriteString(problem.Goal)
		builder.WriteString("\n")
	}
	if problem.Description != "" {
		builder.WriteString("\n### Description\n")
		builder.WriteString(problem.Description)
		builder.WriteString("\n")
	}
	if len(problem.Notes) > 0 {
		builder.WriteString("\n### Notes\n")
		for _, note := range problem.Notes {
			builder.WriteString("- ")
			builder.WriteString(note)
			builder.WriteString("\n")
		}
	}
}

func writeMarkdownAlternativesSection(builder *strings.Builder, alternatives []AlternativeConfig, includeDescriptions bool) {
	if len(alternatives) == 0 {
		return
	}

	builder.WriteString("\n## Alternatives\n")
	for _, alternative := range alternatives {
		builder.WriteString("\n### ")
		builder.WriteString(alternativeLabel(alternative))
		builder.WriteString("\n")
		if includeDescriptions && alternative.Description != "" {
			builder.WriteString(alternative.Description)
			builder.WriteString("\n")
		}
	}
}

func writeMarkdownScenariosSection(builder *strings.Builder, scenarios []ScenarioConfig) {
	if len(scenarios) == 0 {
		return
	}

	sectionTitle := "## Scenario\n"
	if len(scenarios) > 1 {
		sectionTitle = "## Scenarios\n"
	}
	builder.WriteString("\n")
	builder.WriteString(sectionTitle)
	for _, scenario := range scenarios {
		builder.WriteString("\n### ")
		builder.WriteString(scenarioLabel(scenario))
		builder.WriteString("\n")
		if scenario.Description != "" {
			builder.WriteString(scenario.Description)
			builder.WriteString("\n")
		}
		if scenario.Narrative != "" {
			if scenario.Description != "" {
				builder.WriteString("\n")
			}
			builder.WriteString(scenario.Narrative)
			builder.WriteString("\n")
		}
	}
}

func writeMarkdownDecisionDriversSection(
	builder *strings.Builder,
	report ReportConfig,
	config *ExecutionConfig,
	criteria []CriterionConfig,
	scenarioWeights []ScenarioCriterionWeights,
	includeWeights bool,
) {
	hasCriteria := len(criteria) > 0
	comparisons := markdownPreferenceComparisons(report, config)
	weights := filterScenarioWeightsForReport(report, scenarioWeights)
	if !hasCriteria && len(comparisons) == 0 && (!includeWeights || len(weights) == 0) {
		return
	}

	builder.WriteString("\n## Decision Drivers\n")
	if hasCriteria {
		builder.WriteString("\n### Criteria\n")
		for _, criterion := range criteria {
			builder.WriteString("\n#### ")
			builder.WriteString(criterionLabel(criterion))
			builder.WriteString("\n")
			if criterion.Description != "" {
				builder.WriteString(criterion.Description)
				builder.WriteString("\n")
			}
		}
	}

	if len(comparisons) > 0 {
		builder.WriteString("\n### Preference Justifications\n")
		for _, comparison := range comparisons {
			builder.WriteString("- ")
			builder.WriteString(comparison)
			builder.WriteString("\n")
		}
	}

	if includeWeights && len(weights) > 0 {
		writeMarkdownScenarioWeightList(builder, "\n### Criteria Weights\n", config, weights)
	}
}

func writeMarkdownScenarioRankingsSection(
	builder *strings.Builder,
	report ReportConfig,
	config *ExecutionConfig,
	scenarioResults []domain.ScenarioRankingResult,
	criteria []CriterionConfig,
	includeScores bool,
	topAlternatives int,
	includeEvaluationNotes bool,
) {
	sectionTitle := "## Scenario Ranking\n"
	if len(scenarioResults) > 1 {
		sectionTitle = "## Scenario Rankings\n"
	}
	builder.WriteString("\n")
	builder.WriteString(sectionTitle)
	if len(scenarioResults) == 0 {
		builder.WriteString("\n(none)\n")
		return
	}

	for _, scenarioResult := range scenarioResults {
		builder.WriteString("\n### ")
		builder.WriteString(scenarioLabelByName(config, scenarioResult.ScenarioName))
		builder.WriteString("\n")
		rows := limitRankedAlternatives(scenarioResult.RankedAlternatives, topAlternatives)
		for _, alternative := range rows {
			writeMarkdownRankedAlternative(builder, config, alternative, includeScores)
		}
		if includeEvaluationNotes {
			writeMarkdownScenarioEvaluationNotes(builder, report, config, scenarioResult.ScenarioName, criteria)
		}
	}
}

func writeMarkdownFinalRankingSection(
	builder *strings.Builder,
	config *ExecutionConfig,
	rankedAlternatives []domain.RankedAlternative,
	includeScores bool,
) {
	builder.WriteString("\n## Final Ranking\n")
	if len(rankedAlternatives) == 0 {
		builder.WriteString("\n(none)\n")
		return
	}
	builder.WriteString("\n")
	for _, alternative := range rankedAlternatives {
		writeMarkdownRankedAlternative(builder, config, alternative, includeScores)
	}
}

func writeMarkdownNotesAndTradeoffs(
	builder *strings.Builder,
	report ReportConfig,
	config *ExecutionConfig,
	aggregation *AggregationConfig,
	scenarioWeights []ScenarioCriterionWeights,
	aggregationScenarioWeights map[string]float64,
	scenarioResults []domain.ScenarioRankingResult,
	options markdownRenderOptions,
) {
	builder.WriteString("\n## Notes and Tradeoffs\n")
	builder.WriteString("\n- Aggregation method: ")
	if aggregation != nil {
		builder.WriteString(aggregation.Method)
	}
	builder.WriteString("\n")

	if options.IncludeWeights {
		writeMarkdownAggregationWeightsSection(builder, config, aggregationScenarioWeights)
		writeMarkdownTradeoffCriteriaWeights(builder, report, config, scenarioWeights)
	}

	if options.IncludeTradeoffs {
		writeMarkdownExclusionSection(builder, config, scenarioResults)
	}
}

func writeMarkdownAggregationWeightsSection(builder *strings.Builder, config *ExecutionConfig, aggregationScenarioWeights map[string]float64) {
	if len(aggregationScenarioWeights) == 0 {
		return
	}
	builder.WriteString("\n### Scenario Weights\n")
	for _, scenarioName := range orderedWeightNames(aggregationScenarioWeights) {
		builder.WriteString("- ")
		builder.WriteString(scenarioLabelByName(config, scenarioName))
		builder.WriteString(": ")
		builder.WriteString(formatScore(aggregationScenarioWeights[scenarioName]))
		builder.WriteString("\n")
	}
}

func writeMarkdownTradeoffCriteriaWeights(
	builder *strings.Builder,
	report ReportConfig,
	config *ExecutionConfig,
	scenarioWeights []ScenarioCriterionWeights,
) {
	weights := filterScenarioWeightsForReport(report, scenarioWeights)
	if len(weights) == 0 {
		return
	}
	writeMarkdownScenarioWeightList(builder, "\n### Scenario Criteria Weights\n", config, weights)
}

func writeMarkdownExclusionSection(builder *strings.Builder, config *ExecutionConfig, scenarioResults []domain.ScenarioRankingResult) {
	found := false
	for _, scenarioResult := range scenarioResults {
		for _, alternative := range scenarioResult.RankedAlternatives {
			if !alternative.Excluded {
				continue
			}
			if !found {
				builder.WriteString("\n### Exclusions\n")
				found = true
			}
			builder.WriteString("- ")
			builder.WriteString(scenarioLabelByName(config, scenarioResult.ScenarioName))
			builder.WriteString(": ")
			builder.WriteString(alternativeLabelByName(config, alternative.Name))
			if alternative.ExclusionReason != "" {
				builder.WriteString(" (")
				builder.WriteString(alternative.ExclusionReason)
				builder.WriteString(")")
			}
			builder.WriteString("\n")
		}
	}
}

func writeMarkdownScenarioEvaluationNotes(
	builder *strings.Builder,
	report ReportConfig,
	config *ExecutionConfig,
	scenarioName string,
	criteria []CriterionConfig,
) {
	evaluation, exists := findScenarioEvaluation(report, config, scenarioName)
	if !exists {
		return
	}

	builder.WriteString("\n#### Evaluation Notes\n")
	if evaluation.Description != "" {
		builder.WriteString(evaluation.Description)
		builder.WriteString("\n")
	}

	for _, alternative := range orderedAlternativeEvaluations(report, config, evaluation.Evaluations) {
		builder.WriteString("\n##### ")
		builder.WriteString(alternativeLabelByName(config, alternative.AlternativeName))
		builder.WriteString("\n")
		if alternative.Description != "" {
			builder.WriteString(alternative.Description)
			builder.WriteString("\n")
		}

		orderedValues := orderedCriterionValueRecords(criteria, alternative.Values)
		if len(orderedValues) == 0 {
			continue
		}

		builder.WriteString("\nScores:\n")
		for _, value := range orderedValues {
			builder.WriteString("- ")
			builder.WriteString(value.Label)
			builder.WriteString(": ")
			builder.WriteString(value.Rendered)
			builder.WriteString("\n")
		}
	}
}

func writeMarkdownRankedAlternative(builder *strings.Builder, config *ExecutionConfig, alternative domain.RankedAlternative, includeScores bool) {
	if alternative.Excluded {
		builder.WriteString("- ")
		builder.WriteString(alternativeLabelByName(config, alternative.Name))
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
	builder.WriteString(alternativeLabelByName(config, alternative.Name))
	if includeScores {
		builder.WriteString(" (")
		builder.WriteString(formatScore(alternative.Score))
		builder.WriteString(")")
	}
	builder.WriteString("\n")
}

func writeMarkdownNamedWeights(builder *strings.Builder, config *ExecutionConfig, weights []CriterionWeight) {
	for index, weight := range weights {
		if index > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(criterionLabelByName(config, weight.CriterionName))
		builder.WriteString("=")
		builder.WriteString(formatScore(weight.Weight))
	}
}

func writeMarkdownScenarioWeightList(
	builder *strings.Builder,
	heading string,
	config *ExecutionConfig,
	scenarioWeights []ScenarioCriterionWeights,
) {
	builder.WriteString(heading)
	for _, scenarioWeight := range scenarioWeights {
		builder.WriteString("- ")
		builder.WriteString(scenarioLabelByName(config, scenarioWeight.ScenarioName))
		builder.WriteString(": ")
		writeMarkdownNamedWeights(builder, config, canonicalCriterionWeights(scenarioWeight.CriterionWeights))
		builder.WriteString("\n")
	}
}
