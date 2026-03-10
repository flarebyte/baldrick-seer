package pipeline

import (
	"os"
	"path/filepath"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

type DefaultConfigLoader struct{}

func (DefaultConfigLoader) LoadConfig(input LoadConfigInput) (LoadConfigOutput, error) {
	if input.ConfigPath == "" {
		return LoadConfigOutput{}, ErrConfigPathRequired
	}

	info, err := os.Stat(input.ConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return LoadConfigOutput{}, ErrConfigPathDoesNotExist
		}

		return LoadConfigOutput{}, err
	}

	if info.IsDir() {
		return LoadConfigOutput{}, ErrConfigPathIsDirectory
	}

	return LoadConfigOutput{
		ConfigPath: filepath.Clean(input.ConfigPath),
	}, nil
}

type DefaultModelValidator struct{}

func (DefaultModelValidator) ValidateModel(input ValidateModelInput) (ValidateModelOutput, error) {
	return ValidateModelOutput{
		ValidatedModel: domain.ValidatedModelSummary{
			ConfigPath: input.Config.ConfigPath,
		},
	}, nil
}

type DefaultCriteriaWeighter struct{}

func (DefaultCriteriaWeighter) WeightCriteria(WeightCriteriaInput) (WeightCriteriaOutput, error) {
	return WeightCriteriaOutput{}, nil
}

type DefaultScenarioRanker struct{}

func (DefaultScenarioRanker) RankScenarios(RankScenariosInput) (RankScenariosOutput, error) {
	return RankScenariosOutput{}, nil
}

type DefaultScenarioAggregator struct{}

func (DefaultScenarioAggregator) AggregateScenarios(AggregateScenariosInput) (AggregateScenariosOutput, error) {
	return AggregateScenariosOutput{}, nil
}

type DefaultReportRenderer struct{}

func (DefaultReportRenderer) RenderReports(input RenderReportsInput) (RenderReportsOutput, error) {
	return RenderReportsOutput{
		ReportDefinitions: input.ReportDefinitions,
	}, nil
}
