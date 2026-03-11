package pipeline

import "context"

import (
	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func (DefaultScenarioRanker) RankScenarios(ctx context.Context, input RankScenariosInput) (RankScenariosOutput, error) {
	if err := checkContext(ctx, input.Command.ConfigPath); err != nil {
		return RankScenariosOutput{}, err
	}

	if input.Config.Config == nil {
		return RankScenariosOutput{}, NewExecutionFailure("ranking.config_missing", input.Command.ConfigPath, "scenario ranking could not be computed", ErrRankingFailed)
	}

	criteriaByName := make(map[string]CriterionConfig, len(input.Config.Config.CriteriaCatalog))
	for _, criterion := range input.Config.Config.CriteriaCatalog {
		criteriaByName[criterion.Name] = criterion
	}

	evaluationsByScenario := make(map[string]EvaluationConfig, len(input.Config.Config.Evaluations))
	for _, evaluation := range input.Config.Config.Evaluations {
		if _, exists := evaluationsByScenario[evaluation.ScenarioName]; exists {
			return RankScenariosOutput{}, NewExecutionFailure("ranking.duplicate_evaluation", input.Command.ConfigPath, "scenario ranking could not be computed", ErrRankingFailed)
		}
		evaluationsByScenario[evaluation.ScenarioName] = evaluation
	}

	weightsByScenario := make(map[string]ScenarioCriterionWeights, len(input.ScenarioWeights))
	for _, scenarioWeight := range input.ScenarioWeights {
		if _, exists := weightsByScenario[scenarioWeight.ScenarioName]; exists {
			return RankScenariosOutput{}, NewExecutionFailure("ranking.duplicate_weight", input.Command.ConfigPath, "scenario ranking could not be computed", ErrRankingFailed)
		}
		weightsByScenario[scenarioWeight.ScenarioName] = scenarioWeight
	}

	scenarioByName := make(map[string]ScenarioConfig, len(input.Config.Config.Scenarios))
	var scenarioNames []string
	for _, scenario := range input.Config.Config.Scenarios {
		if _, exists := scenarioByName[scenario.Name]; exists {
			return RankScenariosOutput{}, NewExecutionFailure("ranking.duplicate_scenario", input.Command.ConfigPath, "scenario ranking could not be computed", ErrRankingFailed)
		}
		scenarioByName[scenario.Name] = scenario
		scenarioNames = append(scenarioNames, scenario.Name)
	}

	results := make([]domain.ScenarioRankingResult, 0, len(scenarioNames))
	for _, scenarioName := range domain.CanonicalNames(scenarioNames) {
		evaluation, hasEvaluation := evaluationsByScenario[scenarioName]
		if !hasEvaluation {
			return RankScenariosOutput{}, NewExecutionFailure("ranking.missing_evaluation", input.Command.ConfigPath, "scenario ranking could not be computed", ErrRankingFailed)
		}
		scenarioWeight, hasWeights := weightsByScenario[scenarioName]
		if !hasWeights {
			return RankScenariosOutput{}, NewExecutionFailure("ranking.missing_weights", input.Command.ConfigPath, "scenario ranking could not be computed", ErrRankingFailed)
		}

		result, err := rankScenario(
			scenarioByName[scenarioName],
			evaluation,
			scenarioWeight,
			criteriaByName,
		)
		if err != nil {
			return RankScenariosOutput{}, NewExecutionFailure("ranking.invalid_input", input.Command.ConfigPath, "scenario ranking could not be computed", ErrRankingFailed)
		}
		results = append(results, result)
	}

	return RankScenariosOutput{
		ScenarioResults: domain.CanonicalScenarioResults(results),
	}, nil
}

func rankScenario(
	scenario ScenarioConfig,
	evaluation EvaluationConfig,
	scenarioWeights ScenarioCriterionWeights,
	criteriaByName map[string]CriterionConfig,
) (domain.ScenarioRankingResult, error) {
	activeCriteria, err := activeCriterionNames(scenario)
	if err != nil {
		return domain.ScenarioRankingResult{}, err
	}

	weightByCriterion := make(map[string]float64, len(scenarioWeights.CriterionWeights))
	for _, weight := range scenarioWeights.CriterionWeights {
		if _, exists := weightByCriterion[weight.CriterionName]; exists {
			return domain.ScenarioRankingResult{}, ErrRankingFailed
		}
		weightByCriterion[weight.CriterionName] = weight.Weight
	}

	eligible := make([]scoredAlternative, 0, len(evaluation.Evaluations))
	excluded := make([]domain.RankedAlternative, 0, len(evaluation.Evaluations))
	orderedEvaluations := canonicalAlternativeEvaluations(evaluation.Evaluations)
	for _, alternativeEvaluation := range orderedEvaluations {
		violates, err := alternativeViolatesConstraints(scenario, alternativeEvaluation, criteriaByName)
		if err != nil {
			return domain.ScenarioRankingResult{}, err
		}
		if violates {
			excluded = append(excluded, domain.RankedAlternative{
				Name:     alternativeEvaluation.AlternativeName,
				Excluded: true,
			})
			continue
		}

		values := make([]float64, 0, len(activeCriteria))
		for _, criterionName := range activeCriteria {
			criterion, exists := criteriaByName[criterionName]
			if !exists {
				return domain.ScenarioRankingResult{}, ErrRankingFailed
			}

			criterionValue, exists := alternativeEvaluation.Values[criterionName]
			if !exists {
				return domain.ScenarioRankingResult{}, ErrRankingFailed
			}

			normalizedValue, err := normalizeCriterionValue(criterion, criterionValue)
			if err != nil {
				return domain.ScenarioRankingResult{}, err
			}
			values = append(values, normalizedValue)
		}

		eligible = append(eligible, scoredAlternative{
			Name:   alternativeEvaluation.AlternativeName,
			Values: values,
		})
	}

	ranked, err := topsisRankAlternatives(activeCriteria, eligible, weightByCriterion, criteriaByName)
	if err != nil {
		return domain.ScenarioRankingResult{}, err
	}

	return domain.ScenarioRankingResult{
		ScenarioName:       scenario.Name,
		RankedAlternatives: append(ranked, excluded...),
	}, nil
}
