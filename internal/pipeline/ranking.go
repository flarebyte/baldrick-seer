package pipeline

import "context"

import (
	"math"
	"sort"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

const topsisTieTolerance = 1e-12

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

type scoredAlternative struct {
	Name   string
	Values []float64
	Score  float64
}

func topsisRankAlternatives(
	activeCriteria []string,
	eligible []scoredAlternative,
	weightByCriterion map[string]float64,
	criteriaByName map[string]CriterionConfig,
) ([]domain.RankedAlternative, error) {
	if len(eligible) == 0 {
		return nil, nil
	}
	if len(eligible) == 1 {
		return []domain.RankedAlternative{{
			Name:  eligible[0].Name,
			Rank:  1,
			Score: 0,
		}}, nil
	}

	matrix := make([][]float64, len(eligible))
	for rowIndex, alternative := range eligible {
		matrix[rowIndex] = append([]float64(nil), alternative.Values...)
	}

	weighted, err := buildWeightedNormalizedMatrix(matrix, activeCriteria, weightByCriterion)
	if err != nil {
		return nil, err
	}

	idealBest, idealWorst, err := idealSolutions(activeCriteria, weighted, criteriaByName)
	if err != nil {
		return nil, err
	}

	scored := make([]scoredAlternative, len(eligible))
	for index, alternative := range eligible {
		scored[index] = alternative
		distanceBest := euclideanDistance(weighted[index], idealBest)
		distanceWorst := euclideanDistance(weighted[index], idealWorst)
		denominator := distanceBest + distanceWorst
		if denominator == 0 {
			scored[index].Score = 0
			continue
		}
		scored[index].Score = distanceWorst / denominator
	}

	sort.Slice(scored, func(i int, j int) bool {
		scoreDelta := scored[i].Score - scored[j].Score
		if math.Abs(scoreDelta) > topsisTieTolerance {
			return scored[i].Score > scored[j].Score
		}
		return scored[i].Name < scored[j].Name
	})

	ranked := make([]domain.RankedAlternative, 0, len(scored))
	for index, alternative := range scored {
		ranked = append(ranked, domain.RankedAlternative{
			Name:  alternative.Name,
			Rank:  index + 1,
			Score: alternative.Score,
		})
	}

	return ranked, nil
}

func buildWeightedNormalizedMatrix(
	matrix [][]float64,
	activeCriteria []string,
	weightByCriterion map[string]float64,
) ([][]float64, error) {
	columnNorms := make([]float64, len(activeCriteria))
	for columnIndex := range activeCriteria {
		for rowIndex := range matrix {
			columnNorms[columnIndex] += matrix[rowIndex][columnIndex] * matrix[rowIndex][columnIndex]
		}
		columnNorms[columnIndex] = math.Sqrt(columnNorms[columnIndex])
	}

	weighted := make([][]float64, len(matrix))
	for rowIndex := range matrix {
		weighted[rowIndex] = make([]float64, len(activeCriteria))
		for columnIndex, criterionName := range activeCriteria {
			weight, exists := weightByCriterion[criterionName]
			if !exists {
				return nil, ErrRankingFailed
			}

			normalizedValue := 0.0
			if columnNorms[columnIndex] != 0 {
				normalizedValue = matrix[rowIndex][columnIndex] / columnNorms[columnIndex]
			}
			weighted[rowIndex][columnIndex] = normalizedValue * weight
		}
	}

	return weighted, nil
}

func idealSolutions(
	activeCriteria []string,
	weighted [][]float64,
	criteriaByName map[string]CriterionConfig,
) ([]float64, []float64, error) {
	idealBest := make([]float64, len(activeCriteria))
	idealWorst := make([]float64, len(activeCriteria))
	for columnIndex, criterionName := range activeCriteria {
		criterion, exists := criteriaByName[criterionName]
		if !exists {
			return nil, nil, ErrRankingFailed
		}

		if len(weighted) == 0 {
			return nil, nil, ErrRankingFailed
		}
		best := weighted[0][columnIndex]
		worst := weighted[0][columnIndex]
		for rowIndex := 1; rowIndex < len(weighted); rowIndex++ {
			value := weighted[rowIndex][columnIndex]
			switch criterion.Polarity {
			case "benefit":
				if value > best {
					best = value
				}
				if value < worst {
					worst = value
				}
			case "cost":
				if value < best {
					best = value
				}
				if value > worst {
					worst = value
				}
			default:
				return nil, nil, ErrRankingFailed
			}
		}
		idealBest[columnIndex] = best
		idealWorst[columnIndex] = worst
	}

	return idealBest, idealWorst, nil
}

func euclideanDistance(values []float64, target []float64) float64 {
	sum := 0.0
	for index := range values {
		delta := values[index] - target[index]
		sum += delta * delta
	}
	return math.Sqrt(sum)
}

func normalizeCriterionValue(criterion CriterionConfig, value CriterionValue) (float64, error) {
	switch criterion.ValueType {
	case "number":
		return numericValue(value.Value)
	case "ordinal":
		return numericValue(value.Value)
	case "boolean":
		booleanValue, ok := value.Value.(bool)
		if !ok {
			return 0, ErrRankingFailed
		}
		if booleanValue {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, ErrRankingFailed
	}
}

func numericValue(value any) (float64, error) {
	switch typed := value.(type) {
	case int:
		return float64(typed), nil
	case int8:
		return float64(typed), nil
	case int16:
		return float64(typed), nil
	case int32:
		return float64(typed), nil
	case int64:
		return float64(typed), nil
	case float32:
		return float64(typed), nil
	case float64:
		return typed, nil
	default:
		return 0, ErrRankingFailed
	}
}

func alternativeViolatesConstraints(
	scenario ScenarioConfig,
	alternative AlternativeEvaluationConfig,
	criteriaByName map[string]CriterionConfig,
) (bool, error) {
	for _, constraint := range scenario.Constraints {
		criterion, exists := criteriaByName[constraint.CriterionName]
		if !exists {
			return false, ErrRankingFailed
		}
		value, exists := alternative.Values[constraint.CriterionName]
		if !exists {
			return false, ErrRankingFailed
		}
		matches, err := constraintMatches(criterion, constraint, value)
		if err != nil {
			return false, err
		}
		if !matches {
			return true, nil
		}
	}

	return false, nil
}

func constraintMatches(criterion CriterionConfig, constraint ConstraintConfig, value CriterionValue) (bool, error) {
	switch criterion.ValueType {
	case "number", "ordinal":
		left, err := numericValue(value.Value)
		if err != nil {
			return false, err
		}
		right, err := numericValue(constraint.Value)
		if err != nil {
			return false, err
		}
		switch constraint.Operator {
		case "<=":
			return left <= right, nil
		case ">=":
			return left >= right, nil
		case "=":
			return left == right, nil
		case "!=":
			return left != right, nil
		default:
			return false, ErrRankingFailed
		}
	case "boolean":
		left, ok := value.Value.(bool)
		if !ok {
			return false, ErrRankingFailed
		}
		right, ok := constraint.Value.(bool)
		if !ok {
			return false, ErrRankingFailed
		}
		switch constraint.Operator {
		case "=":
			return left == right, nil
		case "!=":
			return left != right, nil
		default:
			return false, ErrRankingFailed
		}
	default:
		return false, ErrRankingFailed
	}
}

func canonicalAlternativeEvaluations(input []AlternativeEvaluationConfig) []AlternativeEvaluationConfig {
	if len(input) == 0 {
		return nil
	}

	output := append([]AlternativeEvaluationConfig(nil), input...)
	sort.Slice(output, func(i int, j int) bool {
		return output[i].AlternativeName < output[j].AlternativeName
	})
	return output
}
