package pipeline

import "context"

import (
	"os"
	"path/filepath"
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

		if reportConfig.Filepath != "" {
			if err := writeRenderedReport(input.Config.Path, reportConfig.Filepath, rendered); err != nil {
				return RenderReportsOutput{}, NewRenderingFailure("rendering.report_write_failed", reportConfig.Filepath, "reports could not be rendered", err)
			}
			continue
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

func writeRenderedReport(configPath string, reportFilepath string, content string) error {
	targetPath, err := resolveReportOutputPath(configPath, reportFilepath)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return err
	}

	return os.WriteFile(targetPath, []byte(content), 0o644)
}

func resolveReportOutputPath(configPath string, reportFilepath string) (string, error) {
	basePath := configPath
	if filepath.Ext(configPath) == ".cue" {
		basePath = filepath.Dir(configPath)
	}

	return filepath.Abs(filepath.Join(basePath, reportFilepath))
}
