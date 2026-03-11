package pipeline

import "context"

import "github.com/flarebyte/baldrick-seer/internal/domain"

func (DefaultCriteriaWeighter) WeightCriteria(ctx context.Context, input WeightCriteriaInput) (WeightCriteriaOutput, error) {
	if err := checkContext(ctx, input.Command.ConfigPath); err != nil {
		return WeightCriteriaOutput{}, err
	}

	if input.Config.Config == nil {
		return WeightCriteriaOutput{}, NewExecutionFailure("weighting.config_missing", input.Command.ConfigPath, "criteria weights could not be computed", ErrWeightingFailed)
	}

	scenarioByName := make(map[string]ScenarioConfig, len(input.Config.Config.Scenarios))
	var scenarioNames []string
	for _, scenario := range input.Config.Config.Scenarios {
		if _, exists := scenarioByName[scenario.Name]; exists {
			return WeightCriteriaOutput{}, NewExecutionFailure("weighting.duplicate_scenario", input.Command.ConfigPath, "criteria weights could not be computed", ErrWeightingFailed)
		}
		scenarioByName[scenario.Name] = scenario
		scenarioNames = append(scenarioNames, scenario.Name)
	}

	orderedScenarioNames := domain.CanonicalNames(scenarioNames)
	scenarioWeights := make([]ScenarioCriterionWeights, 0, len(orderedScenarioNames))
	for _, scenarioName := range orderedScenarioNames {
		weights, err := computeScenarioCriterionWeights(scenarioByName[scenarioName])
		if err != nil {
			return WeightCriteriaOutput{}, NewExecutionFailure("weighting.invalid_preferences", input.Command.ConfigPath, "criteria weights could not be computed", ErrWeightingFailed)
		}

		scenarioWeights = append(scenarioWeights, ScenarioCriterionWeights{
			ScenarioName:     scenarioName,
			CriterionWeights: weights,
		})
	}

	return WeightCriteriaOutput{
		ScenarioWeights: canonicalScenarioCriterionWeights(scenarioWeights),
	}, nil
}

func computeScenarioCriterionWeights(scenario ScenarioConfig) ([]CriterionWeight, error) {
	activeCriterionNames, err := activeCriterionNames(scenario)
	if err != nil {
		return nil, err
	}
	switch len(activeCriterionNames) {
	case 0:
		return nil, nil
	case 1:
		return []CriterionWeight{{CriterionName: activeCriterionNames[0], Weight: 1}}, nil
	}

	if scenario.Preferences == nil || scenario.Preferences.Method != "ahp_pairwise" {
		return nil, ErrWeightingFailed
	}

	matrix, err := buildPairwiseMatrix(activeCriterionNames, scenario.Preferences.Comparisons)
	if err != nil {
		return nil, err
	}

	return normalizePairwiseMatrix(activeCriterionNames, matrix)
}

func activeCriterionNames(scenario ScenarioConfig) ([]string, error) {
	names := make([]string, 0, len(scenario.ActiveCriteria))
	seen := make(map[string]struct{}, len(scenario.ActiveCriteria))
	for _, activeCriterion := range scenario.ActiveCriteria {
		if _, exists := seen[activeCriterion.CriterionName]; exists {
			return nil, ErrWeightingFailed
		}
		seen[activeCriterion.CriterionName] = struct{}{}
		names = append(names, activeCriterion.CriterionName)
	}
	return domain.CanonicalNames(names), nil
}

func buildPairwiseMatrix(activeCriterionNames []string, comparisons []PairwiseComparison) ([][]float64, error) {
	size := len(activeCriterionNames)
	matrix := make([][]float64, size)
	indexByName := make(map[string]int, size)
	for index, criterionName := range activeCriterionNames {
		matrix[index] = make([]float64, size)
		matrix[index][index] = 1
		indexByName[criterionName] = index
	}

	for _, comparison := range comparisons {
		if comparison.Strength <= 0 {
			return nil, ErrWeightingFailed
		}

		moreIndex, moreExists := indexByName[comparison.MoreImportantCriterionName]
		lessIndex, lessExists := indexByName[comparison.LessImportantCriterionName]
		if !moreExists || !lessExists {
			return nil, ErrWeightingFailed
		}

		matrix[moreIndex][lessIndex] = comparison.Strength
		matrix[lessIndex][moreIndex] = 1 / comparison.Strength
	}

	for rowIndex := range matrix {
		for columnIndex := range matrix[rowIndex] {
			if matrix[rowIndex][columnIndex] <= 0 {
				return nil, ErrWeightingFailed
			}
		}
	}

	return matrix, nil
}

func normalizePairwiseMatrix(activeCriterionNames []string, matrix [][]float64) ([]CriterionWeight, error) {
	size := len(matrix)
	columnSums := make([]float64, size)
	for columnIndex := 0; columnIndex < size; columnIndex++ {
		for rowIndex := 0; rowIndex < size; rowIndex++ {
			columnSums[columnIndex] += matrix[rowIndex][columnIndex]
		}
		if columnSums[columnIndex] <= 0 {
			return nil, ErrWeightingFailed
		}
	}

	weights := make([]CriterionWeight, size)
	totalWeight := 0.0
	for rowIndex := 0; rowIndex < size; rowIndex++ {
		rowAverage := 0.0
		for columnIndex := 0; columnIndex < size; columnIndex++ {
			rowAverage += matrix[rowIndex][columnIndex] / columnSums[columnIndex]
		}
		rowAverage /= float64(size)
		totalWeight += rowAverage
		weights[rowIndex] = CriterionWeight{
			CriterionName: activeCriterionNames[rowIndex],
			Weight:        rowAverage,
		}
	}

	if totalWeight <= 0 {
		return nil, ErrWeightingFailed
	}

	for index := range weights {
		weights[index].Weight /= totalWeight
	}

	return weights, nil
}

func canonicalScenarioCriterionWeights(input []ScenarioCriterionWeights) []ScenarioCriterionWeights {
	if len(input) == 0 {
		return nil
	}

	output := append([]ScenarioCriterionWeights(nil), input...)
	for index := range output {
		output[index].CriterionWeights = canonicalCriterionWeights(output[index].CriterionWeights)
	}
	return output
}

func canonicalCriterionWeights(input []CriterionWeight) []CriterionWeight {
	if len(input) == 0 {
		return nil
	}

	orderedNames := make([]string, 0, len(input))
	weightByName := make(map[string]float64, len(input))
	for _, weight := range input {
		orderedNames = append(orderedNames, weight.CriterionName)
		weightByName[weight.CriterionName] = weight.Weight
	}

	canonicalNames := domain.CanonicalNames(orderedNames)
	output := make([]CriterionWeight, 0, len(canonicalNames))
	for _, criterionName := range canonicalNames {
		output = append(output, CriterionWeight{
			CriterionName: criterionName,
			Weight:        weightByName[criterionName],
		})
	}

	return output
}
