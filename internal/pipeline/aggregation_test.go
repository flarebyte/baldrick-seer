package pipeline

import (
	"context"
	"testing"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func TestDefaultScenarioAggregator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		config    LoadedConfig
		scenarios []domain.ScenarioRankingResult
		want      domain.AggregatedRankingResult
		wantErr   error
	}{
		{
			name:      "equal average aggregation",
			config:    aggregationConfig("equal_average", nil),
			scenarios: aggregationScenarioResults(0.4, 0.2, 0.2, 0.8),
			want: domain.AggregatedRankingResult{
				RankedAlternatives: []domain.RankedAlternative{
					{Name: "beta", Rank: 1, Score: 0.5},
					{Name: "alpha", Rank: 2, Score: 0.3},
				},
			},
		},
		{
			name: "weighted average aggregation",
			config: aggregationConfig("weighted_average", map[string]float64{
				"baseline": 3,
				"growth":   1,
			}),
			scenarios: aggregationScenarioResults(0.8, 0.2, 0.2, 0.6),
			want: domain.AggregatedRankingResult{
				RankedAlternatives: []domain.RankedAlternative{
					{Name: "beta", Rank: 1, Score: 0.5},
					{Name: "alpha", Rank: 2, Score: 0.35},
				},
			},
		},
		{
			name:   "final ineligibility when excluded in one scenario",
			config: aggregationConfig("equal_average", nil),
			scenarios: []domain.ScenarioRankingResult{
				{
					ScenarioName: "baseline",
					RankedAlternatives: []domain.RankedAlternative{
						{Name: "alpha", Rank: 1, Score: 0.7},
						{Name: "beta", Rank: 2, Score: 0.4},
					},
				},
				{
					ScenarioName: "growth",
					RankedAlternatives: []domain.RankedAlternative{
						{Name: "alpha", Excluded: true},
						{Name: "beta", Rank: 1, Score: 0.6},
					},
				},
			},
			want: domain.AggregatedRankingResult{
				RankedAlternatives: []domain.RankedAlternative{
					{Name: "beta", Rank: 1, Score: 0.5},
				},
			},
		},
		{
			name:   "deterministic final ranking tie break",
			config: aggregationConfig("equal_average", nil),
			scenarios: []domain.ScenarioRankingResult{
				{
					ScenarioName: "baseline",
					RankedAlternatives: []domain.RankedAlternative{
						{Name: "beta", Rank: 1, Score: 0.5},
						{Name: "alpha", Rank: 2, Score: 0.5},
					},
				},
			},
			want: domain.AggregatedRankingResult{
				RankedAlternatives: []domain.RankedAlternative{
					{Name: "alpha", Rank: 1, Score: 0.5},
					{Name: "beta", Rank: 2, Score: 0.5},
				},
			},
		},
		{
			name:   "empty final ranking when all alternatives are ineligible",
			config: aggregationConfig("equal_average", nil),
			scenarios: []domain.ScenarioRankingResult{
				{
					ScenarioName: "baseline",
					RankedAlternatives: []domain.RankedAlternative{
						{Name: "alpha", Excluded: true},
						{Name: "beta", Excluded: true},
					},
				},
			},
			want: domain.AggregatedRankingResult{},
		},
		{
			name:      "invalid aggregation method fails",
			config:    aggregationConfig("unsupported", nil),
			scenarios: []domain.ScenarioRankingResult{{ScenarioName: "baseline"}},
			wantErr:   ErrAggregationFailed,
		},
	}

	aggregator := DefaultScenarioAggregator{}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assertStageRunResult(t, func() (AggregateScenariosOutput, error) {
				return aggregator.AggregateScenarios(context.Background(), AggregateScenariosInput{
					Command:         domain.CommandRequest{CommandName: domain.CommandNameReportGenerate, ConfigPath: tt.config.Path},
					ScenarioResults: tt.scenarios,
					Config:          tt.config,
				})
			}, tt.wantErr, func(got AggregateScenariosOutput) {
				assertAggregatedRanking(t, got.FinalRanking, tt.want, 1e-9)
			})
		})
	}
}

func TestDefaultScenarioAggregatorIsDeterministic(t *testing.T) {
	t.Parallel()

	input := AggregateScenariosInput{
		Command: domain.CommandRequest{
			CommandName: domain.CommandNameReportGenerate,
			ConfigPath:  fixtureConfigPath(),
		},
		Config: aggregationConfig("equal_average", nil),
		ScenarioResults: []domain.ScenarioRankingResult{
			{
				ScenarioName: "growth",
				RankedAlternatives: []domain.RankedAlternative{
					{Name: "alpha", Rank: 1, Score: 0.4},
					{Name: "beta", Rank: 2, Score: 0.2},
				},
			},
			{
				ScenarioName: "baseline",
				RankedAlternatives: []domain.RankedAlternative{
					{Name: "alpha", Rank: 2, Score: 0.2},
					{Name: "beta", Rank: 1, Score: 0.8},
				},
			},
		},
	}

	assertRepeatedDeepEqual(t, 1, func() (AggregateScenariosOutput, error) {
		return DefaultScenarioAggregator{}.AggregateScenarios(context.Background(), input)
	})
}

func aggregationConfig(method string, scenarioWeights map[string]float64) LoadedConfig {
	config := validLoadedConfig()
	config.Config.Aggregation = &AggregationConfig{
		Method:          method,
		ScenarioWeights: scenarioWeights,
	}
	return config
}

func aggregationScenarioResults(growthAlpha float64, growthBeta float64, baselineAlpha float64, baselineBeta float64) []domain.ScenarioRankingResult {
	return []domain.ScenarioRankingResult{
		{
			ScenarioName: "growth",
			RankedAlternatives: []domain.RankedAlternative{
				{Name: "alpha", Rank: 1, Score: growthAlpha},
				{Name: "beta", Rank: 2, Score: growthBeta},
			},
		},
		{
			ScenarioName: "baseline",
			RankedAlternatives: []domain.RankedAlternative{
				{Name: "alpha", Rank: 2, Score: baselineAlpha},
				{Name: "beta", Rank: 1, Score: baselineBeta},
			},
		},
	}
}
