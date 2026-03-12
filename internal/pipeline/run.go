package pipeline

import "context"

import "github.com/flarebyte/baldrick-seer/internal/domain"

type Runner struct {
	ConfigLoader       ConfigLoader
	ModelValidator     ModelValidator
	CriteriaWeighter   CriteriaWeighter
	ScenarioRanker     ScenarioRanker
	RankingStrategies  RankingStrategySelector
	ScenarioAggregator ScenarioAggregator
	ReportRenderer     ReportRenderer
}

func NewDefaultRunner() Runner {
	weighter := DefaultCriteriaWeighter{}
	ranker := DefaultScenarioRanker{}

	return Runner{
		ConfigLoader:       DefaultConfigLoader{},
		ModelValidator:     DefaultModelValidator{},
		CriteriaWeighter:   weighter,
		ScenarioRanker:     ranker,
		RankingStrategies:  newDefaultRankingStrategySelector(weighter, ranker),
		ScenarioAggregator: DefaultScenarioAggregator{},
		ReportRenderer:     DefaultReportRenderer{},
	}
}

func (r Runner) RunValidate(ctx context.Context, command domain.CommandRequest) (domain.CommandResult, error) {
	if err := checkContext(ctx, command.ConfigPath); err != nil {
		return domain.CommandResult{}, err
	}

	_, validation, err := r.loadAndValidate(ctx, command)
	if err != nil {
		return domain.CommandResult{}, err
	}

	return buildCommandResult(command.CommandName, validation, nil, nil, validation.ReportDefinitions, ""), nil
}

func (r Runner) RunReportGenerate(ctx context.Context, command domain.CommandRequest) (domain.CommandResult, error) {
	if err := checkContext(ctx, command.ConfigPath); err != nil {
		return domain.CommandResult{}, err
	}

	config, validation, err := r.loadAndValidate(ctx, command)
	if err != nil {
		return domain.CommandResult{}, err
	}

	if err := checkContext(ctx, command.ConfigPath); err != nil {
		return domain.CommandResult{}, err
	}

	strategyMethod, err := r.rankingStrategySelector().Select(config, validation)
	if err != nil {
		return domain.CommandResult{}, WrapStageFailure(domain.FailureCategoryExecution, "strategy.select_failed", command.ConfigPath, "command failed", err)
	}

	if err := checkContext(ctx, command.ConfigPath); err != nil {
		return domain.CommandResult{}, err
	}

	strategy, err := r.rankingStrategySelector().Strategy(strategyMethod)
	if err != nil {
		return domain.CommandResult{}, WrapStageFailure(domain.FailureCategoryExecution, "strategy.unsupported", command.ConfigPath, "command failed", err)
	}

	if err := checkContext(ctx, command.ConfigPath); err != nil {
		return domain.CommandResult{}, err
	}

	ranking, err := strategy.Execute(ctx, RankingStrategyInput{
		Command:        command,
		ValidatedModel: validation.ValidatedModel,
		Config:         config,
	})
	if err != nil {
		return domain.CommandResult{}, WrapStageFailure(domain.FailureCategoryExecution, "ranking.failed", command.ConfigPath, "command failed", err)
	}

	aggregated, err := r.ScenarioAggregator.AggregateScenarios(ctx, AggregateScenariosInput{
		Command:         command,
		ScenarioResults: ranking.ScenarioResults,
		Config:          config,
	})
	if err != nil {
		return domain.CommandResult{}, WrapStageFailure(domain.FailureCategoryExecution, "aggregation.failed", command.ConfigPath, "command failed", err)
	}

	if err := checkContext(ctx, command.ConfigPath); err != nil {
		return domain.CommandResult{}, err
	}

	rendered, err := r.ReportRenderer.RenderReports(ctx, RenderReportsInput{
		Command:           command,
		ValidatedModel:    validation.ValidatedModel,
		ScenarioResults:   ranking.ScenarioResults,
		FinalRanking:      aggregated.FinalRanking,
		ReportDefinitions: validation.ReportDefinitions,
		ScenarioWeights:   ranking.ScenarioWeights,
		Config:            config,
	})
	if err != nil {
		return domain.CommandResult{}, WrapStageFailure(domain.FailureCategoryRendering, "rendering.failed", command.ConfigPath, "command failed", err)
	}

	return buildCommandResult(
		command.CommandName,
		validation,
		ranking.ScenarioResults,
		&aggregated.FinalRanking,
		rendered.ReportDefinitions,
		rendered.RenderedOutput,
	), nil
}

func (r Runner) rankingStrategySelector() RankingStrategySelector {
	if r.RankingStrategies != nil {
		return r.RankingStrategies
	}
	return newDefaultRankingStrategySelector(r.CriteriaWeighter, r.ScenarioRanker)
}

func (r Runner) loadAndValidate(ctx context.Context, command domain.CommandRequest) (LoadedConfig, ValidateModelOutput, error) {
	config, err := r.ConfigLoader.LoadConfig(ctx, LoadConfigInput{
		ConfigPath: command.ConfigPath,
	})
	if err != nil {
		return LoadedConfig{}, ValidateModelOutput{}, WrapStageFailure(domain.FailureCategoryInput, "config.load_failed", command.ConfigPath, "command failed", err)
	}

	if err := checkContext(ctx, command.ConfigPath); err != nil {
		return LoadedConfig{}, ValidateModelOutput{}, err
	}

	validation, err := r.ModelValidator.ValidateModel(ctx, ValidateModelInput{
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
