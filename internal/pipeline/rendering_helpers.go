package pipeline

import (
	"strconv"
	"strings"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func allowedFocusedNames(focus *ReportFocus, selectNames func(*ReportFocus) []string) map[string]struct{} {
	if focus == nil {
		return nil
	}
	names := selectNames(focus)
	if len(names) == 0 {
		return nil
	}
	return toAllowedNameSet(names)
}

func toAllowedNameSet(names []string) map[string]struct{} {
	allowed := make(map[string]struct{}, len(names))
	for _, name := range names {
		allowed[name] = struct{}{}
	}
	return allowed
}

func reportArgumentValue(arguments []string, key string, fallback string) string {
	for _, argument := range arguments {
		parsedKey, parsedValue, ok := strings.Cut(argument, "=")
		if ok && parsedKey == key {
			return parsedValue
		}
	}
	return fallback
}

func reportArgumentPresent(arguments []string, key string) bool {
	for _, argument := range arguments {
		parsedKey, _, ok := strings.Cut(argument, "=")
		if ok && parsedKey == key {
			return true
		}
	}
	return false
}

func reportArgumentInt(arguments []string, key string) int {
	value := reportArgumentValue(arguments, key, "")
	if value == "" {
		return 0
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return parsed
}

func limitRankedAlternatives(input []domain.RankedAlternative, top int) []domain.RankedAlternative {
	if top <= 0 {
		return domain.CanonicalRankedAlternatives(input)
	}

	var ranked []domain.RankedAlternative
	var excluded []domain.RankedAlternative
	for _, alternative := range domain.CanonicalRankedAlternatives(input) {
		if alternative.Excluded {
			excluded = append(excluded, alternative)
			continue
		}
		if len(ranked) < top {
			ranked = append(ranked, alternative)
		}
	}
	return append(ranked, excluded...)
}

func formatScore(score float64) string {
	return strconv.FormatFloat(score, 'f', 6, 64)
}

func intPointer(value int) *int {
	return &value
}

func floatPointer(value float64) *float64 {
	return &value
}

func canonicalScenarioWeights(input []ScenarioCriterionWeights) []ScenarioCriterionWeights {
	if len(input) == 0 {
		return nil
	}

	output := append([]ScenarioCriterionWeights(nil), input...)
	for index := range output {
		output[index].CriterionWeights = canonicalCriterionWeights(output[index].CriterionWeights)
	}
	sortScenarioWeights(output)
	return output
}

func sortScenarioWeights(input []ScenarioCriterionWeights) {
	if len(input) < 2 {
		return
	}
	for i := 1; i < len(input); i++ {
		current := input[i]
		j := i - 1
		for ; j >= 0 && input[j].ScenarioName > current.ScenarioName; j-- {
			input[j+1] = input[j]
		}
		input[j+1] = current
	}
}

func problemName(config *ExecutionConfig) string {
	if config == nil || config.Problem == nil || config.Problem.Name == "" {
		return ""
	}
	return config.Problem.Name
}

func orderedWeightNames(weights map[string]float64) []string {
	names := make([]string, 0, len(weights))
	for name := range weights {
		names = append(names, name)
	}
	return domain.CanonicalNames(names)
}
