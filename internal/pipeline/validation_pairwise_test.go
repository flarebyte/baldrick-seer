package pipeline

import "testing"

func TestDefaultModelValidatorAHPPairwiseValidation(t *testing.T) {
	t.Parallel()

	runValidatorFailureCases(t, []validatorFailureCase{
		{
			name: "valid complete pairwise set for 2 criteria",
			config: validLoadedConfigWithAHPPairs(
				[]string{"cost", "speed"},
				[]PairwiseComparison{
					{MoreImportantCriterionName: "cost", LessImportantCriterionName: "speed"},
				},
			),
		},
		{
			name: "valid complete pairwise set for 3 criteria",
			config: validLoadedConfigWithAHPPairs(
				[]string{"cost", "speed", "reliability"},
				[]PairwiseComparison{
					{MoreImportantCriterionName: "cost", LessImportantCriterionName: "speed"},
					{MoreImportantCriterionName: "cost", LessImportantCriterionName: "reliability"},
					{MoreImportantCriterionName: "speed", LessImportantCriterionName: "reliability"},
				},
			),
		},
		{
			name: "missing pair",
			config: validLoadedConfigWithAHPPairs(
				[]string{"cost", "speed", "reliability"},
				[]PairwiseComparison{
					{MoreImportantCriterionName: "cost", LessImportantCriterionName: "speed"},
					{MoreImportantCriterionName: "cost", LessImportantCriterionName: "reliability"},
				},
			),
			wantCodes:   []string{"validation.missing_pairwise_comparison"},
			wantMessage: "missing pairwise comparison for pair: reliability/speed",
		},
		{
			name: "duplicate canonical pair",
			config: validLoadedConfigWithAHPPairs(
				[]string{"cost", "speed"},
				[]PairwiseComparison{
					{MoreImportantCriterionName: "cost", LessImportantCriterionName: "speed"},
					{MoreImportantCriterionName: "cost", LessImportantCriterionName: "speed"},
				},
			),
			wantCodes:   []string{"validation.duplicate_pairwise_comparison"},
			wantMessage: "duplicate pairwise comparison for pair: cost/speed (already defined at comparison 0)",
		},
		{
			name: "inverse duplicate",
			config: validLoadedConfigWithAHPPairs(
				[]string{"cost", "speed"},
				[]PairwiseComparison{
					{MoreImportantCriterionName: "cost", LessImportantCriterionName: "speed"},
					{MoreImportantCriterionName: "speed", LessImportantCriterionName: "cost"},
				},
			),
			wantCodes:   []string{"validation.inverse_duplicate_pairwise_comparison"},
			wantMessage: "inverse duplicate pairwise comparison for pair: cost/speed (already defined at comparison 0)",
		},
		{
			name: "self comparison",
			config: validLoadedConfigWithAHPPairs(
				[]string{"cost", "speed"},
				[]PairwiseComparison{
					{MoreImportantCriterionName: "cost", LessImportantCriterionName: "cost"},
				},
			),
			wantCodes: []string{
				"validation.missing_pairwise_comparison",
				"validation.pairwise_self_comparison",
			},
			wantMessage: "missing pairwise comparison for pair: cost/speed",
		},
		{
			name: "reference to unknown criterion",
			config: validLoadedConfigWithAHPPairs(
				[]string{"cost", "speed"},
				[]PairwiseComparison{
					{MoreImportantCriterionName: "missing", LessImportantCriterionName: "speed"},
				},
			),
			wantCodes: []string{
				"validation.missing_pairwise_comparison",
				"validation.unknown_pairwise_criterion",
			},
			wantMessage: "missing pairwise comparison for pair: cost/speed",
		},
		{
			name: "reference to criterion not active in scenario",
			config: func() LoadedConfig {
				config := validLoadedConfigWithAHPPairs(
					[]string{"cost", "speed"},
					[]PairwiseComparison{
						{MoreImportantCriterionName: "cost", LessImportantCriterionName: "reliability"},
					},
				)
				config.Config.CriteriaCatalog = append(config.Config.CriteriaCatalog, CriterionConfig{Name: "reliability", ValueType: "number"})
				return config
			}(),
			wantCodes: []string{
				"validation.inactive_pairwise_criterion",
				"validation.missing_pairwise_comparison",
			},
			wantMessage: "pairwise comparison references criterion not active in scenario: reliability",
		},
		{
			name: "scenario with 0 active criteria",
			config: func() LoadedConfig {
				config := validLoadedConfig()
				config.Config.Scenarios[0].ActiveCriteria = nil
				config.Config.Evaluations[0].Evaluations[0].Values = nil
				config.Config.Scenarios[0].Preferences = &ScenarioPreferences{
					Method: "ahp_pairwise",
				}
				return config
			}(),
		},
		{
			name: "scenario with 1 active criterion",
			config: func() LoadedConfig {
				config := validLoadedConfig()
				config.Config.Scenarios[0].Preferences = &ScenarioPreferences{
					Method: "ahp_pairwise",
				}
				return config
			}(),
		},
	})
}
