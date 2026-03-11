package pipeline

import "github.com/flarebyte/baldrick-seer/internal/domain"

type ConfigLoader interface {
	LoadConfig(LoadConfigInput) (LoadConfigOutput, error)
}

type ModelValidator interface {
	ValidateModel(ValidateModelInput) (ValidateModelOutput, error)
}

type CriteriaWeighter interface {
	WeightCriteria(WeightCriteriaInput) (WeightCriteriaOutput, error)
}

type ScenarioRanker interface {
	RankScenarios(RankScenariosInput) (RankScenariosOutput, error)
}

type ScenarioAggregator interface {
	AggregateScenarios(AggregateScenariosInput) (AggregateScenariosOutput, error)
}

type ReportRenderer interface {
	RenderReports(RenderReportsInput) (RenderReportsOutput, error)
}

type LoadConfigInput struct {
	ConfigPath string
}

type LoadedConfig struct {
	Path           string
	Evaluated      string
	TopLevelFields []string
}

type LoadConfigOutput struct {
	Config LoadedConfig
}

type ValidateModelInput struct {
	Command domain.CommandRequest
	Config  LoadedConfig
}

type ValidateModelOutput struct {
	Diagnostics       []domain.Diagnostic
	ValidatedModel    domain.ValidatedModelSummary
	ReportDefinitions []domain.ReportDefinition
}

type CriterionWeight struct {
	CriterionName string
	Weight        float64
}

type WeightCriteriaInput struct {
	Command        domain.CommandRequest
	ValidatedModel domain.ValidatedModelSummary
}

type WeightCriteriaOutput struct {
	CriterionWeights []CriterionWeight
}

type RankScenariosInput struct {
	Command          domain.CommandRequest
	ValidatedModel   domain.ValidatedModelSummary
	CriterionWeights []CriterionWeight
}

type RankScenariosOutput struct {
	ScenarioResults []domain.ScenarioRankingResult
}

type AggregateScenariosInput struct {
	Command         domain.CommandRequest
	ScenarioResults []domain.ScenarioRankingResult
}

type AggregateScenariosOutput struct {
	FinalRanking domain.AggregatedRankingResult
}

type RenderReportsInput struct {
	Command           domain.CommandRequest
	ValidatedModel    domain.ValidatedModelSummary
	ScenarioResults   []domain.ScenarioRankingResult
	FinalRanking      domain.AggregatedRankingResult
	ReportDefinitions []domain.ReportDefinition
}

type RenderReportsOutput struct {
	ReportDefinitions []domain.ReportDefinition
}
