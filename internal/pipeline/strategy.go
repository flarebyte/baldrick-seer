package pipeline

import (
	"context"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

// RankingMethod identifies the internal ranking pipeline variant.
type RankingMethod string

const (
	RankingMethodV1AHPTopsis RankingMethod = "v1_ahp_topsis"
	RankingMethodElectre     RankingMethod = "electre"
)

type RankingStrategyInput struct {
	Command        domain.CommandRequest
	ValidatedModel domain.ValidatedModelSummary
	Config         LoadedConfig
}

type RankingStrategyOutput struct {
	ScenarioWeights []ScenarioCriterionWeights
	ScenarioResults []domain.ScenarioRankingResult
}

type RankingStrategy interface {
	Execute(context.Context, RankingStrategyInput) (RankingStrategyOutput, error)
}

type RankingStrategySelector interface {
	Select(LoadedConfig, ValidateModelOutput) (RankingMethod, error)
	Strategy(RankingMethod) (RankingStrategy, error)
}

type defaultRankingStrategySelector struct {
	v1 RankingStrategy
}

func newDefaultRankingStrategySelector(weighter CriteriaWeighter, ranker ScenarioRanker) RankingStrategySelector {
	return defaultRankingStrategySelector{
		v1: v1RankingStrategy{
			weighter: weighter,
			ranker:   ranker,
		},
	}
}

func (s defaultRankingStrategySelector) Select(_ LoadedConfig, _ ValidateModelOutput) (RankingMethod, error) {
	return RankingMethodV1AHPTopsis, nil
}

func (s defaultRankingStrategySelector) Strategy(method RankingMethod) (RankingStrategy, error) {
	switch method {
	case RankingMethodV1AHPTopsis:
		return s.v1, nil
	default:
		return nil, ErrRankingFailed
	}
}

type v1RankingStrategy struct {
	weighter CriteriaWeighter
	ranker   ScenarioRanker
}

func (s v1RankingStrategy) Execute(ctx context.Context, input RankingStrategyInput) (RankingStrategyOutput, error) {
	weights, err := s.weighter.WeightCriteria(ctx, WeightCriteriaInput(input))
	if err != nil {
		return RankingStrategyOutput{}, err
	}

	if err := checkContext(ctx, input.Command.ConfigPath); err != nil {
		return RankingStrategyOutput{}, err
	}

	scenarios, err := s.ranker.RankScenarios(ctx, RankScenariosInput{
		Command:         input.Command,
		ValidatedModel:  input.ValidatedModel,
		ScenarioWeights: weights.ScenarioWeights,
		Config:          input.Config,
	})
	if err != nil {
		return RankingStrategyOutput{}, err
	}

	return RankingStrategyOutput{
		ScenarioWeights: weights.ScenarioWeights,
		ScenarioResults: scenarios.ScenarioResults,
	}, nil
}
