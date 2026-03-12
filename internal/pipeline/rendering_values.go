package pipeline

import (
	"sort"
	"strconv"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

type criterionValueRecord struct {
	Name     string
	Rendered string
}

func buildEvaluationValueLookup(config *ExecutionConfig, report ReportConfig) map[string]map[string][]criterionValueRecord {
	if config == nil || len(config.Evaluations) == 0 {
		return nil
	}

	allowedScenarios := allowedFocusedNames(report.Focus, func(focus *ReportFocus) []string {
		return focus.ScenarioNames
	})
	allowedAlternatives := allowedFocusedNames(report.Focus, func(focus *ReportFocus) []string {
		return focus.AlternativeNames
	})
	allowedCriteria := allowedFocusedNames(report.Focus, func(focus *ReportFocus) []string {
		return focus.CriterionNames
	})

	output := make(map[string]map[string][]criterionValueRecord, len(config.Evaluations))
	for _, evaluation := range canonicalEvaluations(config.Evaluations) {
		if len(allowedScenarios) > 0 {
			if _, exists := allowedScenarios[evaluation.ScenarioName]; !exists {
				continue
			}
		}

		alternatives := make(map[string][]criterionValueRecord, len(evaluation.Evaluations))
		for _, alternative := range canonicalAlternativeEvaluations(evaluation.Evaluations) {
			if len(allowedAlternatives) > 0 {
				if _, exists := allowedAlternatives[alternative.AlternativeName]; !exists {
					continue
				}
			}
			values := make([]criterionValueRecord, 0, len(alternative.Values))
			for _, criterionName := range domain.CanonicalNames(valueNames(alternative.Values)) {
				if len(allowedCriteria) > 0 {
					if _, exists := allowedCriteria[criterionName]; !exists {
						continue
					}
				}
				values = append(values, criterionValueRecord{
					Name:     criterionName,
					Rendered: renderCriterionValue(alternative.Values[criterionName].Value),
				})
			}
			alternatives[alternative.AlternativeName] = values
		}
		output[evaluation.ScenarioName] = alternatives
	}

	return output
}

func canonicalEvaluations(input []EvaluationConfig) []EvaluationConfig {
	if len(input) == 0 {
		return nil
	}

	output := append([]EvaluationConfig(nil), input...)
	for i := range output {
		output[i].Evaluations = canonicalAlternativeEvaluations(output[i].Evaluations)
	}
	sort.Slice(output, func(i int, j int) bool {
		return output[i].ScenarioName < output[j].ScenarioName
	})
	return output
}

func valueNames(values map[string]CriterionValue) []string {
	names := make([]string, 0, len(values))
	for name := range values {
		names = append(names, name)
	}
	return names
}

func renderCriterionValue(value any) string {
	switch typed := value.(type) {
	case bool:
		if typed {
			return "true"
		}
		return "false"
	case int:
		return strconv.Itoa(typed)
	case int8, int16, int32, int64, float32, float64:
		number, err := numericValue(typed)
		if err != nil {
			return ""
		}
		return strconv.FormatFloat(number, 'f', -1, 64)
	default:
		return ""
	}
}
