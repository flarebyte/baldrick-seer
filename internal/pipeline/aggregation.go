package pipeline

import "context"

import (
	"sort"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func (DefaultScenarioAggregator) AggregateScenarios(ctx context.Context, input AggregateScenariosInput) (AggregateScenariosOutput, error) {
	if err := checkContext(ctx, input.Command.ConfigPath); err != nil {
		return AggregateScenariosOutput{}, err
	}

	if input.Config.Config == nil || input.Config.Config.Aggregation == nil {
		return AggregateScenariosOutput{}, NewExecutionFailure("aggregation.config_missing", input.Command.ConfigPath, "final ranking could not be computed", ErrAggregationFailed)
	}

	finalRanking, err := aggregateScenarioResults(input.Config.Config.Aggregation, input.ScenarioResults)
	if err != nil {
		return AggregateScenariosOutput{}, NewExecutionFailure("aggregation.invalid_method", input.Command.ConfigPath, "final ranking could not be computed", ErrAggregationFailed)
	}

	return AggregateScenariosOutput{
		FinalRanking: finalRanking,
	}, nil
}

func aggregateScenarioResults(config *AggregationConfig, scenarioResults []domain.ScenarioRankingResult) (domain.AggregatedRankingResult, error) {
	if config == nil {
		return domain.AggregatedRankingResult{}, ErrAggregationFailed
	}

	orderedScenarioResults := domain.CanonicalScenarioResults(scenarioResults)
	scoreByAlternative := map[string]float64{}
	eligibility := map[string]bool{}
	presence := map[string]int{}

	scenarioWeights, err := aggregationWeights(config, orderedScenarioResults)
	if err != nil {
		return domain.AggregatedRankingResult{}, err
	}

	for _, scenarioResult := range orderedScenarioResults {
		scenarioWeight, exists := scenarioWeights[scenarioResult.ScenarioName]
		if !exists {
			return domain.AggregatedRankingResult{}, ErrAggregationFailed
		}

		for _, rankedAlternative := range scenarioResult.RankedAlternatives {
			if _, exists := eligibility[rankedAlternative.Name]; !exists {
				eligibility[rankedAlternative.Name] = true
			}
			if rankedAlternative.Excluded {
				eligibility[rankedAlternative.Name] = false
				continue
			}

			scoreByAlternative[rankedAlternative.Name] += rankedAlternative.Score * scenarioWeight
			presence[rankedAlternative.Name]++
		}
	}

	participatingScenarioCount := len(orderedScenarioResults)
	var ranked []domain.RankedAlternative
	for alternativeName, score := range scoreByAlternative {
		if !eligibility[alternativeName] {
			continue
		}
		if presence[alternativeName] != participatingScenarioCount {
			continue
		}
		ranked = append(ranked, domain.RankedAlternative{
			Name:  alternativeName,
			Score: score,
		})
	}

	sort.Slice(ranked, func(i int, j int) bool {
		if ranked[i].Score != ranked[j].Score {
			return ranked[i].Score > ranked[j].Score
		}
		return ranked[i].Name < ranked[j].Name
	})

	for index := range ranked {
		ranked[index].Rank = index + 1
	}

	return domain.CanonicalAggregatedRankingResult(domain.AggregatedRankingResult{
		RankedAlternatives: ranked,
	}), nil
}

func aggregationWeights(config *AggregationConfig, scenarioResults []domain.ScenarioRankingResult) (map[string]float64, error) {
	switch config.Method {
	case "equal_average":
		if len(scenarioResults) == 0 {
			return map[string]float64{}, nil
		}
		weight := 1 / float64(len(scenarioResults))
		weights := make(map[string]float64, len(scenarioResults))
		for _, scenarioResult := range scenarioResults {
			weights[scenarioResult.ScenarioName] = weight
		}
		return weights, nil
	case "weighted_average":
		if len(scenarioResults) == 0 {
			return map[string]float64{}, nil
		}
		total := 0.0
		weights := make(map[string]float64, len(scenarioResults))
		for _, scenarioResult := range scenarioResults {
			weight, exists := config.ScenarioWeights[scenarioResult.ScenarioName]
			if !exists || weight <= 0 {
				return nil, ErrAggregationFailed
			}
			weights[scenarioResult.ScenarioName] = weight
			total += weight
		}
		if total <= 0 {
			return nil, ErrAggregationFailed
		}
		for scenarioName := range weights {
			weights[scenarioName] /= total
		}
		return weights, nil
	default:
		return nil, ErrAggregationFailed
	}
}
