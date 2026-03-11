package pipeline

import "testing"

func TestDefaultModelValidatorConstraintValidation(t *testing.T) {
	t.Parallel()

	runValidatorFailureCases(t, []validatorFailureCase{
		{
			name: "valid numeric constraint",
			config: validLoadedConfigWithConstraints(
				[]CriterionConfig{{Name: "cost", ValueType: "number"}},
				[]string{"cost"},
				[]ConstraintConfig{{CriterionName: "cost", Operator: "<=", Value: 10}},
			),
		},
		{
			name: "valid ordinal constraint",
			config: validLoadedConfigWithConstraints(
				[]CriterionConfig{{Name: "priority", ValueType: "ordinal", ScaleGuidance: []any{"low", "medium", "high"}}},
				[]string{"priority"},
				[]ConstraintConfig{{CriterionName: "priority", Operator: ">=", Value: 2}},
			),
		},
		{
			name: "valid boolean constraint",
			config: validLoadedConfigWithConstraints(
				[]CriterionConfig{{Name: "approved", ValueType: "boolean"}},
				[]string{"approved"},
				[]ConstraintConfig{{CriterionName: "approved", Operator: "=", Value: true}},
			),
		},
		{
			name: "unknown criterion in constraint",
			config: validLoadedConfigWithConstraints(
				[]CriterionConfig{{Name: "cost", ValueType: "number"}},
				[]string{"cost"},
				[]ConstraintConfig{{CriterionName: "missing", Operator: "<=", Value: 10}},
			),
			wantCodes:   []string{"validation.unknown_constraint_criterion"},
			wantMessage: "unknown criterion name in constraints: missing",
		},
		{
			name: "inactive criterion in constraint",
			config: validLoadedConfigWithConstraints(
				[]CriterionConfig{
					{Name: "cost", ValueType: "number"},
					{Name: "speed", ValueType: "number"},
				},
				[]string{"cost"},
				[]ConstraintConfig{{CriterionName: "speed", Operator: "<=", Value: 10}},
			),
			wantCodes:   []string{"validation.inactive_constraint_criterion"},
			wantMessage: "constraint references criterion not active in scenario: speed",
		},
		{
			name: "invalid operator for boolean criterion",
			config: validLoadedConfigWithConstraints(
				[]CriterionConfig{{Name: "approved", ValueType: "boolean"}},
				[]string{"approved"},
				[]ConstraintConfig{{CriterionName: "approved", Operator: "<=", Value: true}},
			),
			wantCodes:   []string{"validation.invalid_constraint_operator"},
			wantMessage: "invalid constraint operator for boolean criterion approved: <=",
		},
		{
			name: "invalid numeric value for number criterion",
			config: validLoadedConfigWithConstraints(
				[]CriterionConfig{{Name: "cost", ValueType: "number"}},
				[]string{"cost"},
				[]ConstraintConfig{{CriterionName: "cost", Operator: "<=", Value: true}},
			),
			wantCodes:   []string{"validation.invalid_constraint_number_value"},
			wantMessage: "constraint value must be numeric for criterion: cost",
		},
		{
			name: "non integer value for ordinal criterion",
			config: validLoadedConfigWithConstraints(
				[]CriterionConfig{{Name: "priority", ValueType: "ordinal", ScaleGuidance: []any{"low", "medium", "high"}}},
				[]string{"priority"},
				[]ConstraintConfig{{CriterionName: "priority", Operator: ">=", Value: 1.5}},
			),
			wantCodes:   []string{"validation.invalid_constraint_ordinal_value"},
			wantMessage: "constraint value must be an integer for criterion: priority",
		},
		{
			name: "wrong value type for boolean criterion",
			config: validLoadedConfigWithConstraints(
				[]CriterionConfig{{Name: "approved", ValueType: "boolean"}},
				[]string{"approved"},
				[]ConstraintConfig{{CriterionName: "approved", Operator: "!=", Value: 1}},
			),
			wantCodes:   []string{"validation.invalid_constraint_boolean_value"},
			wantMessage: "constraint value must be true or false for criterion: approved",
		},
		{
			name: "multiple invalid constraints with deterministic diagnostic order",
			config: validLoadedConfigWithConstraints(
				[]CriterionConfig{
					{Name: "approved", ValueType: "boolean"},
					{Name: "cost", ValueType: "number"},
					{Name: "priority", ValueType: "ordinal", ScaleGuidance: []any{"low", "high"}},
					{Name: "speed", ValueType: "number"},
				},
				[]string{"approved", "cost", "priority"},
				[]ConstraintConfig{
					{CriterionName: "approved", Operator: "<=", Value: 1},
					{CriterionName: "speed", Operator: "<=", Value: 10},
					{CriterionName: "priority", Operator: ">=", Value: 1.5},
					{CriterionName: "missing", Operator: "=", Value: 1},
				},
			),
			wantCodes: []string{
				"validation.inactive_constraint_criterion",
				"validation.invalid_constraint_boolean_value",
				"validation.invalid_constraint_operator",
				"validation.invalid_constraint_ordinal_value",
				"validation.unknown_constraint_criterion",
			},
			wantMessage: "constraint references criterion not active in scenario: speed",
		},
	})
}
