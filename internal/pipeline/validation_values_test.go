package pipeline

import "testing"

func TestDefaultModelValidatorEvaluationValueValidation(t *testing.T) {
	t.Parallel()

	runValidatorFailureCases(t, []validatorFailureCase{
		{
			name: "valid evaluation coverage for active criteria",
			config: validLoadedConfigWithScenarioEvaluations(
				[]CriterionConfig{
					{Name: "cost", ValueType: "number"},
					{Name: "approved", ValueType: "boolean"},
				},
				[]string{"cost", "approved"},
				[]EvaluationConfig{
					{
						ScenarioName: "baseline",
						Evaluations: []AlternativeEvaluationConfig{
							{
								AlternativeName: "option_a",
								Values: map[string]CriterionValue{
									"approved": {Kind: "boolean", Value: true},
									"cost":     {Kind: "number", Value: 10},
								},
							},
						},
					},
				},
			),
		},
		{
			name: "duplicate evaluation block for the same scenario",
			config: validLoadedConfigWithScenarioEvaluations(
				[]CriterionConfig{{Name: "cost", ValueType: "number"}},
				[]string{"cost"},
				[]EvaluationConfig{
					{
						ScenarioName: "baseline",
						Evaluations:  []AlternativeEvaluationConfig{{AlternativeName: "option_a", Values: map[string]CriterionValue{"cost": {Kind: "number", Value: 1}}}},
					},
					{
						ScenarioName: "baseline",
						Evaluations:  []AlternativeEvaluationConfig{{AlternativeName: "option_a", Values: map[string]CriterionValue{"cost": {Kind: "number", Value: 1}}}},
					},
				},
			),
			wantCodes:   []string{"validation.duplicate_evaluation_scenario"},
			wantMessage: "duplicate evaluation block for scenario: baseline",
		},
		{
			name: "unknown scenario in evaluation block",
			config: func() LoadedConfig {
				config := validLoadedConfig()
				config.Config.Evaluations[0].ScenarioName = "missing"
				return config
			}(),
			wantCodes:   []string{"validation.unknown_evaluation_scenario"},
			wantMessage: "unknown scenario name in evaluations: missing",
		},
		{
			name: "duplicate alternative evaluation within a scenario",
			config: validLoadedConfigWithScenarioEvaluations(
				[]CriterionConfig{{Name: "cost", ValueType: "number"}},
				[]string{"cost"},
				[]EvaluationConfig{
					{
						ScenarioName: "baseline",
						Evaluations: []AlternativeEvaluationConfig{
							{AlternativeName: "option_a", Values: map[string]CriterionValue{"cost": {Kind: "number", Value: 1}}},
							{AlternativeName: "option_a", Values: map[string]CriterionValue{"cost": {Kind: "number", Value: 2}}},
						},
					},
				},
			),
			wantCodes:   []string{"validation.duplicate_evaluation_alternative"},
			wantMessage: "duplicate alternative evaluation in scenario baseline: option_a",
		},
		{
			name: "unknown alternative in a scenario evaluation",
			config: func() LoadedConfig {
				config := validLoadedConfig()
				config.Config.Evaluations[0].Evaluations[0].AlternativeName = "missing"
				return config
			}(),
			wantCodes:   []string{"validation.unknown_evaluation_alternative"},
			wantMessage: "unknown alternative name in evaluations: missing",
		},
		{
			name: "missing value for an active criterion",
			config: validLoadedConfigWithScenarioEvaluations(
				[]CriterionConfig{
					{Name: "cost", ValueType: "number"},
					{Name: "speed", ValueType: "number"},
				},
				[]string{"cost", "speed"},
				[]EvaluationConfig{
					{
						ScenarioName: "baseline",
						Evaluations: []AlternativeEvaluationConfig{
							{AlternativeName: "option_a", Values: map[string]CriterionValue{"cost": {Kind: "number", Value: 1}}},
						},
					},
				},
			),
			wantCodes:   []string{"validation.missing_evaluation_value"},
			wantMessage: "missing value for active criterion in scenario baseline: speed",
		},
		{
			name: "unknown criterion in values",
			config: validLoadedConfigWithScenarioEvaluations(
				[]CriterionConfig{{Name: "cost", ValueType: "number"}},
				[]string{"cost"},
				[]EvaluationConfig{
					{
						ScenarioName: "baseline",
						Evaluations: []AlternativeEvaluationConfig{
							{
								AlternativeName: "option_a",
								Values: map[string]CriterionValue{
									"cost":    {Kind: "number", Value: 1},
									"missing": {Kind: "number", Value: 2},
								},
							},
						},
					},
				},
			),
			wantCodes:   []string{"validation.unknown_evaluation_criterion"},
			wantMessage: "unknown criterion name in evaluation values: missing",
		},
		{
			name: "inactive criterion in values",
			config: validLoadedConfigWithScenarioEvaluations(
				[]CriterionConfig{
					{Name: "cost", ValueType: "number"},
					{Name: "speed", ValueType: "number"},
				},
				[]string{"cost"},
				[]EvaluationConfig{
					{
						ScenarioName: "baseline",
						Evaluations: []AlternativeEvaluationConfig{
							{
								AlternativeName: "option_a",
								Values: map[string]CriterionValue{
									"cost":  {Kind: "number", Value: 1},
									"speed": {Kind: "number", Value: 2},
								},
							},
						},
					},
				},
			),
			wantCodes:   []string{"validation.inactive_evaluation_criterion"},
			wantMessage: "criterion value is not active in scenario baseline: speed",
		},
		{
			name: "number criterion with wrong kind",
			config: validLoadedConfigWithScenarioEvaluations(
				[]CriterionConfig{{Name: "cost", ValueType: "number"}},
				[]string{"cost"},
				[]EvaluationConfig{
					{
						ScenarioName: "baseline",
						Evaluations: []AlternativeEvaluationConfig{
							{AlternativeName: "option_a", Values: map[string]CriterionValue{"cost": {Kind: "boolean", Value: true}}},
						},
					},
				},
			),
			wantCodes:   []string{"validation.evaluation_value_kind_mismatch"},
			wantMessage: "evaluation value kind mismatch for criterion cost: want number, got boolean",
		},
		{
			name: "ordinal criterion with wrong kind",
			config: validLoadedConfigWithScenarioEvaluations(
				[]CriterionConfig{{Name: "priority", ValueType: "ordinal", ScaleGuidance: []any{"low", "high"}}},
				[]string{"priority"},
				[]EvaluationConfig{
					{
						ScenarioName: "baseline",
						Evaluations: []AlternativeEvaluationConfig{
							{AlternativeName: "option_a", Values: map[string]CriterionValue{"priority": {Kind: "number", Value: 1}}},
						},
					},
				},
			),
			wantCodes:   []string{"validation.evaluation_value_kind_mismatch"},
			wantMessage: "evaluation value kind mismatch for criterion priority: want ordinal, got number",
		},
		{
			name: "boolean criterion with wrong kind",
			config: validLoadedConfigWithScenarioEvaluations(
				[]CriterionConfig{{Name: "approved", ValueType: "boolean"}},
				[]string{"approved"},
				[]EvaluationConfig{
					{
						ScenarioName: "baseline",
						Evaluations: []AlternativeEvaluationConfig{
							{AlternativeName: "option_a", Values: map[string]CriterionValue{"approved": {Kind: "number", Value: 1}}},
						},
					},
				},
			),
			wantCodes:   []string{"validation.evaluation_value_kind_mismatch"},
			wantMessage: "evaluation value kind mismatch for criterion approved: want boolean, got number",
		},
		{
			name: "ordinal value that is not an integer",
			config: validLoadedConfigWithScenarioEvaluations(
				[]CriterionConfig{{Name: "priority", ValueType: "ordinal", ScaleGuidance: []any{"low", "high"}}},
				[]string{"priority"},
				[]EvaluationConfig{
					{
						ScenarioName: "baseline",
						Evaluations: []AlternativeEvaluationConfig{
							{AlternativeName: "option_a", Values: map[string]CriterionValue{"priority": {Kind: "ordinal", Value: 1.5}}},
						},
					},
				},
			),
			wantCodes:   []string{"validation.invalid_ordinal_value"},
			wantMessage: "ordinal criterion value must be an integer: priority",
		},
		{
			name: "ordinal criterion missing scale guidance",
			config: validLoadedConfigWithScenarioEvaluations(
				[]CriterionConfig{{Name: "priority", ValueType: "ordinal"}},
				[]string{"priority"},
				[]EvaluationConfig{
					{
						ScenarioName: "baseline",
						Evaluations: []AlternativeEvaluationConfig{
							{AlternativeName: "option_a", Values: map[string]CriterionValue{"priority": {Kind: "ordinal", Value: 1}}},
						},
					},
				},
			),
			wantCodes:   []string{"validation.ordinal_scale_guidance_missing"},
			wantMessage: "ordinal criterion is missing scaleGuidance: priority",
		},
		{
			name: "valid boolean and valid ordinal cases",
			config: validLoadedConfigWithScenarioEvaluations(
				[]CriterionConfig{
					{Name: "approved", ValueType: "boolean"},
					{Name: "priority", ValueType: "ordinal", ScaleGuidance: []any{"low", "medium", "high"}},
				},
				[]string{"approved", "priority"},
				[]EvaluationConfig{
					{
						ScenarioName: "baseline",
						Evaluations: []AlternativeEvaluationConfig{
							{
								AlternativeName: "option_a",
								Values: map[string]CriterionValue{
									"approved": {Kind: "boolean", Value: true},
									"priority": {Kind: "ordinal", Value: 2},
								},
							},
						},
					},
				},
			),
		},
	})
}
