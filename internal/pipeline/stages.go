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
	ConfigFields   []string
	Config         *ExecutionConfig
}

type ExecutionConfig struct {
	Problem         *ProblemConfig      `json:"problem"`
	Reports         []ReportConfig      `json:"reports"`
	CriteriaCatalog []CriterionConfig   `json:"criteriaCatalog"`
	Alternatives    []AlternativeConfig `json:"alternatives"`
	Scenarios       []ScenarioConfig    `json:"scenarios"`
	Evaluations     []EvaluationConfig  `json:"evaluations"`
	Aggregation     *AggregationConfig  `json:"aggregation"`
}

type ProblemConfig struct {
	Name string `json:"name"`
}

type ReportConfig struct {
	Name   string       `json:"name"`
	Title  string       `json:"title"`
	Format string       `json:"format"`
	Focus  *ReportFocus `json:"focus"`
}

type ReportFocus struct {
	ScenarioNames    []string `json:"scenarioNames"`
	AlternativeNames []string `json:"alternativeNames"`
	CriterionNames   []string `json:"criterionNames"`
}

type CriterionConfig struct {
	Name string `json:"name"`
}

type AlternativeConfig struct {
	Name string `json:"name"`
}

type ScenarioConfig struct {
	Name           string                 `json:"name"`
	ActiveCriteria []ScenarioCriterionRef `json:"activeCriteria"`
	Constraints    []ConstraintConfig     `json:"constraints"`
}

type ScenarioCriterionRef struct {
	CriterionName string `json:"criterionName"`
}

type ConstraintConfig struct {
	CriterionName string `json:"criterionName"`
}

type EvaluationConfig struct {
	ScenarioName string                        `json:"scenarioName"`
	Evaluations  []AlternativeEvaluationConfig `json:"evaluations"`
}

type AlternativeEvaluationConfig struct {
	AlternativeName string `json:"alternativeName"`
}

type AggregationConfig struct {
	ScenarioWeights map[string]float64 `json:"scenarioWeights"`
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
