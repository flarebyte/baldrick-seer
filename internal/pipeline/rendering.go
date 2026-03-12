package pipeline

import "context"

import (
	"strings"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func (DefaultReportRenderer) RenderReports(ctx context.Context, input RenderReportsInput) (RenderReportsOutput, error) {
	if err := checkContext(ctx, input.Command.ConfigPath); err != nil {
		return RenderReportsOutput{}, err
	}

	if input.Config.Config == nil {
		return RenderReportsOutput{}, NewRenderingFailure("rendering.config_missing", input.Command.ConfigPath, "reports could not be rendered", ErrRenderingFailed)
	}

	reportByName := make(map[string]ReportConfig, len(input.Config.Config.Reports))
	for _, report := range input.Config.Config.Reports {
		reportByName[report.Name] = report
	}

	orderedDefinitions := domain.CanonicalReportDefinitions(input.ReportDefinitions)
	renderedParts := make([]string, 0, len(orderedDefinitions))
	for _, reportDefinition := range orderedDefinitions {
		reportConfig, exists := reportByName[reportDefinition.Name]
		if !exists {
			return RenderReportsOutput{}, NewRenderingFailure("rendering.report_missing", input.Command.ConfigPath, "reports could not be rendered", ErrRenderingFailed)
		}

		rendered, err := renderReport(
			reportConfig,
			input.Config.Config,
			input.ScenarioResults,
			input.Config.Config.Aggregation,
			input.ScenarioWeights,
		)
		if err != nil {
			return RenderReportsOutput{}, NewRenderingFailure("rendering.report_failed", input.Command.ConfigPath, "reports could not be rendered", ErrRenderingFailed)
		}
		renderedParts = append(renderedParts, rendered)
	}

	return RenderReportsOutput{
		ReportDefinitions: orderedDefinitions,
		RenderedOutput:    strings.Join(renderedParts, "\n"),
	}, nil
}

func renderReport(
	report ReportConfig,
	config *ExecutionConfig,
	scenarioResults []domain.ScenarioRankingResult,
	aggregation *AggregationConfig,
	scenarioWeights []ScenarioCriterionWeights,
) (string, error) {
	filteredScenarioResults := filterScenarioResultsForReport(report, scenarioResults)
	displayScenarioResults := filterScenarioAlternativesForReport(report, filteredScenarioResults)
	filteredFinalRanking, err := aggregateScenarioResults(aggregation, filteredScenarioResults)
	if err != nil {
		return "", err
	}
	filteredFinalRanking = filterFinalRankingForReport(report, filteredFinalRanking)
	filteredScenarioWeights := filterScenarioWeightsForReport(report, scenarioWeights)
	aggregationScenarioWeights, err := aggregationWeights(aggregation, filteredScenarioResults)
	if err != nil {
		return "", err
	}

	switch report.Format {
	case "markdown":
		return renderMarkdownReport(report, config, displayScenarioResults, filteredFinalRanking, aggregation, filteredScenarioWeights, aggregationScenarioWeights), nil
	case "json":
		return renderJSONReport(report, config, displayScenarioResults, filteredFinalRanking, aggregation, filteredScenarioWeights, aggregationScenarioWeights)
	case "csv":
		return renderCSVReport(report, config, displayScenarioResults, filteredFinalRanking)
	default:
		return "", ErrRenderingFailed
	}
}
