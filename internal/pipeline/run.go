package pipeline

import "github.com/flarebyte/baldrick-seer/internal/domain"

type Runner struct {
	ConfigLoader       ConfigLoader
	ModelValidator     ModelValidator
	CriteriaWeighter   CriteriaWeighter
	ScenarioRanker     ScenarioRanker
	ScenarioAggregator ScenarioAggregator
	ReportRenderer     ReportRenderer
}

func NewDefaultRunner() Runner {
	return Runner{
		ConfigLoader:       DefaultConfigLoader{},
		ModelValidator:     DefaultModelValidator{},
		CriteriaWeighter:   DefaultCriteriaWeighter{},
		ScenarioRanker:     DefaultScenarioRanker{},
		ScenarioAggregator: DefaultScenarioAggregator{},
		ReportRenderer:     DefaultReportRenderer{},
	}
}

func (r Runner) RunValidate(command domain.CommandRequest) (domain.CommandResult, error) {
	_, validation, err := r.loadAndValidate(command)
	if err != nil {
		return domain.CommandResult{}, err
	}

	return buildCommandResult(command.CommandName, validation, nil, nil, validation.ReportDefinitions, ""), nil
}

func (r Runner) RunReportGenerate(command domain.CommandRequest) (domain.CommandResult, error) {
	config, validation, err := r.loadAndValidate(command)
	if err != nil {
		return domain.CommandResult{}, err
	}

	weights, err := r.CriteriaWeighter.WeightCriteria(WeightCriteriaInput{
		Command:        command,
		ValidatedModel: validation.ValidatedModel,
		Config:         config,
	})
	if err != nil {
		return domain.CommandResult{}, WrapStageFailure(domain.FailureCategoryExecution, "weighting.failed", command.ConfigPath, "command failed", err)
	}

	scenarios, err := r.ScenarioRanker.RankScenarios(RankScenariosInput{
		Command:         command,
		ValidatedModel:  validation.ValidatedModel,
		ScenarioWeights: weights.ScenarioWeights,
		Config:          config,
	})
	if err != nil {
		return domain.CommandResult{}, WrapStageFailure(domain.FailureCategoryExecution, "ranking.failed", command.ConfigPath, "command failed", err)
	}

	aggregated, err := r.ScenarioAggregator.AggregateScenarios(AggregateScenariosInput{
		Command:         command,
		ScenarioResults: scenarios.ScenarioResults,
		Config:          config,
	})
	if err != nil {
		return domain.CommandResult{}, WrapStageFailure(domain.FailureCategoryExecution, "aggregation.failed", command.ConfigPath, "command failed", err)
	}

	rendered, err := r.ReportRenderer.RenderReports(RenderReportsInput{
		Command:           command,
		ValidatedModel:    validation.ValidatedModel,
		ScenarioResults:   scenarios.ScenarioResults,
		FinalRanking:      aggregated.FinalRanking,
		ReportDefinitions: validation.ReportDefinitions,
		ScenarioWeights:   weights.ScenarioWeights,
		Config:            config,
	})
	if err != nil {
		return domain.CommandResult{}, WrapStageFailure(domain.FailureCategoryRendering, "rendering.failed", command.ConfigPath, "command failed", err)
	}

	return buildCommandResult(
		command.CommandName,
		validation,
		scenarios.ScenarioResults,
		&aggregated.FinalRanking,
		rendered.ReportDefinitions,
		rendered.RenderedOutput,
	), nil
}

func (r Runner) loadAndValidate(command domain.CommandRequest) (LoadedConfig, ValidateModelOutput, error) {
	config, err := r.ConfigLoader.LoadConfig(LoadConfigInput{
		ConfigPath: command.ConfigPath,
	})
	if err != nil {
		return LoadedConfig{}, ValidateModelOutput{}, WrapStageFailure(domain.FailureCategoryInput, "config.load_failed", command.ConfigPath, "command failed", err)
	}

	validation, err := r.ModelValidator.ValidateModel(ValidateModelInput{
		Command: command,
		Config:  config.Config,
	})
	if err != nil {
		return LoadedConfig{}, ValidateModelOutput{}, WrapStageFailure(domain.FailureCategoryValidation, "validation.failed", command.ConfigPath, "command failed", err)
	}

	return config.Config, validation, nil
}

func buildCommandResult(
	commandName domain.CommandName,
	validation ValidateModelOutput,
	scenarioResults []domain.ScenarioRankingResult,
	finalRanking *domain.AggregatedRankingResult,
	reportDefinitions []domain.ReportDefinition,
	renderedOutput string,
) domain.CommandResult {
	return domain.CanonicalCommandResult(domain.CommandResult{
		CommandName:       commandName,
		Diagnostics:       validation.Diagnostics,
		ValidatedModel:    &validation.ValidatedModel,
		ScenarioResults:   scenarioResults,
		FinalRanking:      finalRanking,
		ReportDefinitions: reportDefinitions,
		RenderedOutput:    renderedOutput,
	})
}
