package pipeline

import "sort"

type scoredAlternative struct {
	Name   string
	Values []float64
	Score  float64
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
