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
	config, err := r.ConfigLoader.LoadConfig(LoadConfigInput{
		ConfigPath: command.ConfigPath,
	})
	if err != nil {
		return domain.CommandResult{}, WrapStageFailure(domain.FailureCategoryInput, "config.load_failed", command.ConfigPath, "command failed", err)
	}

	validation, err := r.ModelValidator.ValidateModel(ValidateModelInput{
		Command: command,
		Config:  config,
	})
	if err != nil {
		return domain.CommandResult{}, WrapStageFailure(domain.FailureCategoryValidation, "validation.failed", command.ConfigPath, "command failed", err)
	}

	return domain.CanonicalCommandResult(domain.CommandResult{
		CommandName:       command.CommandName,
		Diagnostics:       validation.Diagnostics,
		ValidatedModel:    &validation.ValidatedModel,
		ReportDefinitions: validation.ReportDefinitions,
	}), nil
}

func (r Runner) RunReportGenerate(command domain.CommandRequest) (domain.CommandResult, error) {
	config, err := r.ConfigLoader.LoadConfig(LoadConfigInput{
		ConfigPath: command.ConfigPath,
	})
	if err != nil {
		return domain.CommandResult{}, WrapStageFailure(domain.FailureCategoryInput, "config.load_failed", command.ConfigPath, "command failed", err)
	}

	validation, err := r.ModelValidator.ValidateModel(ValidateModelInput{
		Command: command,
		Config:  config,
	})
	if err != nil {
		return domain.CommandResult{}, WrapStageFailure(domain.FailureCategoryValidation, "validation.failed", command.ConfigPath, "command failed", err)
	}

	weights, err := r.CriteriaWeighter.WeightCriteria(WeightCriteriaInput{
		Command:        command,
		ValidatedModel: validation.ValidatedModel,
	})
	if err != nil {
		return domain.CommandResult{}, WrapStageFailure(domain.FailureCategoryExecution, "weighting.failed", command.ConfigPath, "command failed", err)
	}

	scenarios, err := r.ScenarioRanker.RankScenarios(RankScenariosInput{
		Command:          command,
		ValidatedModel:   validation.ValidatedModel,
		CriterionWeights: weights.CriterionWeights,
	})
	if err != nil {
		return domain.CommandResult{}, WrapStageFailure(domain.FailureCategoryExecution, "ranking.failed", command.ConfigPath, "command failed", err)
	}

	aggregated, err := r.ScenarioAggregator.AggregateScenarios(AggregateScenariosInput{
		Command:         command,
		ScenarioResults: scenarios.ScenarioResults,
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
	})
	if err != nil {
		return domain.CommandResult{}, WrapStageFailure(domain.FailureCategoryRendering, "rendering.failed", command.ConfigPath, "command failed", err)
	}

	return domain.CanonicalCommandResult(domain.CommandResult{
		CommandName:       command.CommandName,
		Diagnostics:       validation.Diagnostics,
		ValidatedModel:    &validation.ValidatedModel,
		ScenarioResults:   scenarios.ScenarioResults,
		FinalRanking:      &aggregated.FinalRanking,
		ReportDefinitions: rendered.ReportDefinitions,
	}), nil
}
