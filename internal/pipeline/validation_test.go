package pipeline

import (
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
				{Name: "cost"},
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
						{AlternativeName: "option_a"},
					},
				},
			},
			Aggregation: &AggregationConfig{},
		},
	}
}

func validateConfig(t *testing.T, config LoadedConfig) []domain.Diagnostic {
	t.Helper()

	validator := DefaultModelValidator{}
	_, err := validator.ValidateModel(ValidateModelInput{
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

func TestDefaultModelValidatorFailures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		mutate      func(*LoadedConfig)
		wantCodes   []string
		wantMessage string
	}{
		{
			name: "missing required top level sections",
			mutate: func(config *LoadedConfig) {
				config.ConfigFields = nil
				config.Config = &ExecutionConfig{}
			},
			wantCodes: []string{
				"validation.section_missing",
				"validation.section_missing",
				"validation.section_missing",
				"validation.section_missing",
				"validation.section_missing",
				"validation.section_missing",
				"validation.section_missing",
			},
			wantMessage: "missing required section: aggregation",
		},
		{
			name: "duplicate criterion names",
			mutate: func(config *LoadedConfig) {
				config.Config.CriteriaCatalog = []CriterionConfig{{Name: "cost"}, {Name: "cost"}}
			},
			wantCodes:   []string{"validation.duplicate_criterion_name"},
			wantMessage: "duplicate criterion name: cost",
		},
		{
			name: "duplicate alternative names",
			mutate: func(config *LoadedConfig) {
				config.Config.Alternatives = []AlternativeConfig{{Name: "option_a"}, {Name: "option_a"}}
			},
			wantCodes:   []string{"validation.duplicate_alternative_name"},
			wantMessage: "duplicate alternative name: option_a",
		},
		{
			name: "duplicate scenario names",
			mutate: func(config *LoadedConfig) {
				config.Config.Scenarios = []ScenarioConfig{{Name: "baseline"}, {Name: "baseline"}}
			},
			wantCodes:   []string{"validation.duplicate_scenario_name"},
			wantMessage: "duplicate scenario name: baseline",
		},
		{
			name: "duplicate report names",
			mutate: func(config *LoadedConfig) {
				config.Config.Reports = []ReportConfig{{Name: "summary"}, {Name: "summary"}}
			},
			wantCodes:   []string{"validation.duplicate_report_name"},
			wantMessage: "duplicate report name: summary",
		},
		{
			name: "unknown scenario reference in evaluations",
			mutate: func(config *LoadedConfig) {
				config.Config.Evaluations[0].ScenarioName = "missing"
			},
			wantCodes:   []string{"validation.unknown_evaluation_scenario"},
			wantMessage: "unknown scenario name in evaluations: missing",
		},
		{
			name: "unknown alternative reference in scenario evaluations",
			mutate: func(config *LoadedConfig) {
				config.Config.Evaluations[0].Evaluations[0].AlternativeName = "missing"
			},
			wantCodes:   []string{"validation.unknown_evaluation_alternative"},
			wantMessage: "unknown alternative name in evaluations: missing",
		},
		{
			name: "unknown criterion reference in active criteria",
			mutate: func(config *LoadedConfig) {
				config.Config.Scenarios[0].ActiveCriteria[0].CriterionName = "missing"
			},
			wantCodes:   []string{"validation.unknown_active_criterion"},
			wantMessage: "unknown criterion name in active criteria: missing",
		},
		{
			name: "unknown criterion reference in constraints",
			mutate: func(config *LoadedConfig) {
				config.Config.Scenarios[0].Constraints = []ConstraintConfig{{CriterionName: "missing"}}
			},
			wantCodes:   []string{"validation.unknown_constraint_criterion"},
			wantMessage: "unknown criterion name in constraints: missing",
		},
		{
			name: "unknown names in report focus selectors",
			mutate: func(config *LoadedConfig) {
				config.Config.Reports[0].Focus = &ReportFocus{
					ScenarioNames:    []string{"missing-scenario"},
					AlternativeNames: []string{"missing-alternative"},
					CriterionNames:   []string{"missing-criterion"},
				}
			},
			wantCodes: []string{
				"validation.unknown_report_focus_alternative",
				"validation.unknown_report_focus_criterion",
				"validation.unknown_report_focus_scenario",
			},
			wantMessage: "unknown alternative name in report focus: missing-alternative",
		},
		{
			name: "unknown scenario names in aggregation weights",
			mutate: func(config *LoadedConfig) {
				config.Config.Aggregation = &AggregationConfig{
					ScenarioWeights: map[string]float64{"missing": 1},
				}
			},
			wantCodes:   []string{"validation.unknown_aggregation_scenario"},
			wantMessage: "unknown scenario name in aggregation weights: missing",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			config := validLoadedConfig()
			tt.mutate(&config)

			diagnostics := validateConfig(t, config)
			var gotCodes []string
			for _, diagnostic := range diagnostics {
				gotCodes = append(gotCodes, diagnostic.Code)
			}

			if !reflect.DeepEqual(gotCodes, tt.wantCodes) {
				t.Fatalf("codes = %#v, want %#v", gotCodes, tt.wantCodes)
			}

			if diagnostics[0].Message != tt.wantMessage {
				t.Fatalf("message = %q, want %q", diagnostics[0].Message, tt.wantMessage)
			}
		})
	}
}

func TestDefaultModelValidatorValidConfig(t *testing.T) {
	t.Parallel()

	validator := DefaultModelValidator{}
	config := validLoadedConfig()

	got, err := validator.ValidateModel(ValidateModelInput{
		Command: domain.CommandRequest{
			CommandName: domain.CommandNameValidate,
			ConfigPath:  config.Path,
		},
		Config: config,
	})
	if err != nil {
		t.Fatalf("ValidateModel() error = %v", err)
	}

	if got.ValidatedModel.ConfigPath != config.Path {
		t.Fatalf("ConfigPath = %q, want %q", got.ValidatedModel.ConfigPath, config.Path)
	}

	if got.ValidatedModel.CriterionCount != 1 || got.ValidatedModel.AlternativeCount != 1 || got.ValidatedModel.ScenarioCount != 1 {
		t.Fatalf("summary counts = %#v, want all ones", got.ValidatedModel)
	}
}
