package pipeline

import (
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
	options := resolveMarkdownRenderOptions(report)
	return renderMarkdownRichReport(report, config, scenarioResults, finalRanking, aggregation, scenarioWeights, aggregationScenarioWeights, options)
}

type markdownRenderOptions struct {
	Detail                 string
	IncludeContext         bool
	IncludeWeights         bool
	IncludeAltDescriptions bool
	IncludeEvaluationNotes bool
	IncludeTradeoffs       bool
}

func resolveMarkdownRenderOptions(report ReportConfig) markdownRenderOptions {
	detail := reportArgumentValue(report.Arguments, "detail", "standard")
	options := markdownRenderOptions{
		Detail:                 detail,
		IncludeContext:         true,
		IncludeWeights:         false,
		IncludeAltDescriptions: true,
		IncludeEvaluationNotes: true,
		IncludeTradeoffs:       true,
	}

	switch detail {
	case "full":
		options.IncludeWeights = true
	case "standard":
		options.IncludeWeights = true
	case "brief":
		options.IncludeWeights = false
		options.IncludeTradeoffs = false
	default:
		options.IncludeWeights = true
	}

	if reportArgumentPresent(report.Arguments, "explain") && reportArgumentValue(report.Arguments, "explain", "true") == "false" {
		options.IncludeWeights = false
		options.IncludeTradeoffs = false
	}

	if reportArgumentPresent(report.Arguments, "include-context") {
		options.IncludeContext = reportArgumentValue(report.Arguments, "include-context", "false") == "true"
	}
	if reportArgumentPresent(report.Arguments, "include-weights") {
		options.IncludeWeights = reportArgumentValue(report.Arguments, "include-weights", "false") == "true"
	}
	if reportArgumentPresent(report.Arguments, "include-alternative-descriptions") {
		options.IncludeAltDescriptions = reportArgumentValue(report.Arguments, "include-alternative-descriptions", "false") == "true"
	}
	if reportArgumentPresent(report.Arguments, "include-evaluation-notes") {
		options.IncludeEvaluationNotes = reportArgumentValue(report.Arguments, "include-evaluation-notes", "false") == "true"
	}
	if reportArgumentPresent(report.Arguments, "include-tradeoffs") {
		options.IncludeTradeoffs = reportArgumentValue(report.Arguments, "include-tradeoffs", "false") == "true"
	}

	return options
}

func renderMarkdownRichReport(
	report ReportConfig,
	config *ExecutionConfig,
	scenarioResults []domain.ScenarioRankingResult,
	finalRanking domain.AggregatedRankingResult,
	aggregation *AggregationConfig,
	scenarioWeights []ScenarioCriterionWeights,
	aggregationScenarioWeights map[string]float64,
	options markdownRenderOptions,
) string {
	includeScores, topAlternatives := markdownReportSettings(report)
	orderedResults := orderedScenarioResultsForMarkdown(report, config, scenarioResults)
	orderedFinalRanking := limitRankedAlternatives(finalRanking.RankedAlternatives, topAlternatives)
	alternatives := filteredAlternativesForMarkdown(report, config)
	scenarios := filteredScenariosForMarkdown(report, config)
	criteria := filteredCriteriaForMarkdown(report, config)
	var builder strings.Builder

	writeMarkdownReportTitle(&builder, markdownDocumentTitle(report, config))

	if options.IncludeContext {
		writeMarkdownProblemSection(&builder, config)
		writeMarkdownAlternativesSection(&builder, alternatives, options.IncludeAltDescriptions)
		writeMarkdownScenariosSection(&builder, scenarios)
		writeMarkdownDecisionDriversSection(&builder, report, config, criteria, scenarioWeights, options.IncludeWeights)
	}

	writeMarkdownScenarioRankingsSection(&builder, report, config, orderedResults, criteria, includeScores, topAlternatives, options.IncludeEvaluationNotes)
	writeMarkdownFinalRankingSection(&builder, config, orderedFinalRanking, includeScores)

	if options.IncludeTradeoffs || options.IncludeWeights {
		writeMarkdownNotesAndTradeoffs(&builder, report, config, aggregation, scenarioWeights, aggregationScenarioWeights, orderedResults, options)
	}

	return builder.String()
}

func markdownReportSettings(report ReportConfig) (bool, int) {
	includeScores := reportArgumentValue(report.Arguments, "include-scores", "true") == "true"
	topAlternatives := reportArgumentInt(report.Arguments, "top-alternatives")
	return includeScores, topAlternatives
}

func markdownDocumentTitle(report ReportConfig, config *ExecutionConfig) string {
	if report.Title != "" {
		return report.Title
	}
	if config != nil && config.Problem != nil {
		if config.Problem.Title != "" {
			return config.Problem.Title
		}
		if config.Problem.Name != "" {
			return config.Problem.Name
		}
	}
	return ""
}
