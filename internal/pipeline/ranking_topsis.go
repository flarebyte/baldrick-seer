package pipeline

import (
	"math"
	"sort"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

const topsisTieTolerance = 1e-12

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
