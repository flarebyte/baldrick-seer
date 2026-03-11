package pipeline

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func validLoadedConfig() LoadedConfig {
	return LoadedConfig{
		Path:           filepath.Clean(fixtureConfigPath()),
		TopLevelFields: []string{"config"},
		ConfigFields: []string{
			"aggregation",
			"alternatives",
			"criteriaCatalog",
			"evaluations",
			"problem",
			"reports",
			"scenarios",
		},
		Config: &ExecutionConfig{
			Problem: &ProblemConfig{Name: "minimal"},
			Reports: []ReportConfig{
				{Name: "summary", Title: "Summary", Format: "markdown"},
			},
			CriteriaCatalog: []CriterionConfig{
				{Name: "cost", Polarity: "cost", ValueType: "number"},
			},
			Alternatives: []AlternativeConfig{
				{Name: "option_a"},
			},
			Scenarios: []ScenarioConfig{
				{
					Name: "baseline",
					ActiveCriteria: []ScenarioCriterionRef{
						{CriterionName: "cost"},
					},
				},
			},
			Evaluations: []EvaluationConfig{
				{
					ScenarioName: "baseline",
					Evaluations: []AlternativeEvaluationConfig{
						{
							AlternativeName: "option_a",
							Values: map[string]CriterionValue{
								"cost": {Kind: "number", Value: 1},
							},
						},
					},
				},
			},
			Aggregation: &AggregationConfig{},
		},
	}
}

func validLoadedConfigWithAHPPairs(activeCriteria []string, comparisons []PairwiseComparison) LoadedConfig {
	config := validLoadedConfig()
	config.Config.CriteriaCatalog = nil
	config.Config.Evaluations[0].Evaluations[0].Values = map[string]CriterionValue{}
	for _, criterionName := range activeCriteria {
		config.Config.CriteriaCatalog = append(config.Config.CriteriaCatalog, CriterionConfig{
			Name:      criterionName,
			Polarity:  defaultPolarityForCriterionName(criterionName),
			ValueType: "number",
		})
		config.Config.Evaluations[0].Evaluations[0].Values[criterionName] = CriterionValue{Kind: "number", Value: 1}
	}

	config.Config.Scenarios = []ScenarioConfig{
		{
			Name: "baseline",
			Preferences: &ScenarioPreferences{
				Method:      "ahp_pairwise",
				Scale:       "saaty_1_9",
				Comparisons: comparisons,
			},
		},
	}
	for _, criterionName := range activeCriteria {
		config.Config.Scenarios[0].ActiveCriteria = append(
			config.Config.Scenarios[0].ActiveCriteria,
			ScenarioCriterionRef{CriterionName: criterionName},
		)
	}

	return config
}

func validLoadedConfigWithScenarioEvaluations(
	criteria []CriterionConfig,
	activeCriteria []string,
	evaluationBlocks []EvaluationConfig,
) LoadedConfig {
	config := validLoadedConfig()
	config.Config.CriteriaCatalog = append([]CriterionConfig(nil), criteria...)
	config.Config.Scenarios = []ScenarioConfig{scenarioWithActiveCriteria("baseline", activeCriteria)}
	config.Config.Evaluations = append([]EvaluationConfig(nil), evaluationBlocks...)
	return config
}

func validLoadedConfigWithConstraints(
	criteria []CriterionConfig,
	activeCriteria []string,
	constraints []ConstraintConfig,
) LoadedConfig {
	config := validLoadedConfig()
	config.Config.CriteriaCatalog = append([]CriterionConfig(nil), criteria...)
	config.Config.Scenarios = []ScenarioConfig{scenarioWithConstraints("baseline", activeCriteria, constraints)}
	config.Config.Evaluations[0].Evaluations[0].Values = map[string]CriterionValue{}
	for _, criterionName := range activeCriteria {
		config.Config.Evaluations[0].Evaluations[0].Values[criterionName] = CriterionValue{
			Kind:  criterionValueTypeForName(criteria, criterionName),
			Value: validCriterionValueForName(criteria, criterionName),
		}
	}
	return config
}

func validLoadedConfigWithReports(reports []ReportConfig) LoadedConfig {
	config := validLoadedConfig()
	config.Config.Reports = append([]ReportConfig(nil), reports...)
	return config
}

func scenarioWithActiveCriteria(name string, activeCriteria []string) ScenarioConfig {
	scenario := ScenarioConfig{Name: name}
	for _, criterionName := range activeCriteria {
		scenario.ActiveCriteria = append(scenario.ActiveCriteria, ScenarioCriterionRef{CriterionName: criterionName})
	}
	return scenario
}

func scenarioWithConstraints(name string, activeCriteria []string, constraints []ConstraintConfig) ScenarioConfig {
	scenario := scenarioWithActiveCriteria(name, activeCriteria)
	scenario.Constraints = append([]ConstraintConfig(nil), constraints...)
	return scenario
}

func scenarioEvaluationBlock(scenarioName string, evaluations ...AlternativeEvaluationConfig) []EvaluationConfig {
	return []EvaluationConfig{{
		ScenarioName: scenarioName,
		Evaluations:  append([]AlternativeEvaluationConfig(nil), evaluations...),
	}}
}

func alternativeEvaluation(alternativeName string, values map[string]CriterionValue) AlternativeEvaluationConfig {
	return AlternativeEvaluationConfig{
		AlternativeName: alternativeName,
		Values:          values,
	}
}

func criterionValueTypeForName(criteria []CriterionConfig, criterionName string) string {
	for _, criterion := range criteria {
		if criterion.Name == criterionName {
			return criterion.ValueType
		}
	}
	return "number"
}

func validCriterionValueForName(criteria []CriterionConfig, criterionName string) any {
	for _, criterion := range criteria {
		if criterion.Name != criterionName {
			continue
		}

		switch criterion.ValueType {
		case "ordinal":
			return 1
		case "boolean":
			return true
		default:
			return 1
		}
	}

	return 1
}

func defaultPolarityForCriterionName(criterionName string) string {
	if criterionName == "cost" {
		return "cost"
	}
	return "benefit"
}

func validateConfig(t *testing.T, config LoadedConfig) []domain.Diagnostic {
	t.Helper()

	validator := DefaultModelValidator{}
	_, err := validator.ValidateModel(context.Background(), ValidateModelInput{
		Command: domain.CommandRequest{
			CommandName: domain.CommandNameValidate,
			ConfigPath:  config.Path,
		},
		Config: config,
	})
	if err == nil {
		t.Fatal("ValidateModel() error = nil, want error")
	}

	if !errors.Is(err, ErrValidationFailed) {
		t.Fatalf("error = %v, want %v", err, ErrValidationFailed)
	}

	failure := domain.AsCommandFailure(err)
	if failure == nil {
		t.Fatal("AsCommandFailure(err) = nil, want value")
	}

	return failure.Diagnostics
}

type validatorFailureCase struct {
	name        string
	config      LoadedConfig
	wantCodes   []string
	wantMessage string
}

func runValidatorFailureCases(t *testing.T, tests []validatorFailureCase) {
	t.Helper()

	validator := DefaultModelValidator{}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := validator.ValidateModel(context.Background(), ValidateModelInput{
				Command: domain.CommandRequest{
					CommandName: domain.CommandNameValidate,
					ConfigPath:  tt.config.Path,
				},
				Config: tt.config,
			})

			if len(tt.wantCodes) == 0 {
				if err != nil {
					t.Fatalf("ValidateModel() error = %v", err)
				}
				return
			}

			if !errors.Is(err, ErrValidationFailed) {
				t.Fatalf("error = %v, want %v", err, ErrValidationFailed)
			}

			failure := domain.AsCommandFailure(err)
			if failure == nil {
				t.Fatal("AsCommandFailure(err) = nil, want value")
			}

			var gotCodes []string
			for _, diagnostic := range failure.Diagnostics {
				gotCodes = append(gotCodes, diagnostic.Code)
			}

			if !reflect.DeepEqual(gotCodes, tt.wantCodes) {
				t.Fatalf("codes = %#v, want %#v", gotCodes, tt.wantCodes)
			}

			if failure.Diagnostics[0].Message != tt.wantMessage {
				t.Fatalf("message = %q, want %q", failure.Diagnostics[0].Message, tt.wantMessage)
			}
		})
	}
}
