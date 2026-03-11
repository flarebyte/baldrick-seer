package pipeline

import (
	"errors"
	"reflect"
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
			name:   "equal average aggregation",
			config: aggregationConfig("equal_average", nil),
			scenarios: []domain.ScenarioRankingResult{
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
			scenarios: []domain.ScenarioRankingResult{
				{
					ScenarioName: "growth",
					RankedAlternatives: []domain.RankedAlternative{
						{Name: "alpha", Rank: 1, Score: 0.8},
						{Name: "beta", Rank: 2, Score: 0.2},
					},
				},
				{
					ScenarioName: "baseline",
					RankedAlternatives: []domain.RankedAlternative{
						{Name: "alpha", Rank: 2, Score: 0.2},
						{Name: "beta", Rank: 1, Score: 0.6},
					},
				},
			},
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

			got, err := aggregator.AggregateScenarios(AggregateScenariosInput{
				Command:         domain.CommandRequest{CommandName: domain.CommandNameReportGenerate, ConfigPath: tt.config.Path},
				ScenarioResults: tt.scenarios,
				Config:          tt.config,
			})
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("AggregateScenarios() error = %v", err)
			}

			assertAggregatedRanking(t, got.FinalRanking, tt.want, 1e-9)
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

	first, err := DefaultScenarioAggregator{}.AggregateScenarios(input)
	if err != nil {
		t.Fatalf("first AggregateScenarios() error = %v", err)
	}

	second, err := DefaultScenarioAggregator{}.AggregateScenarios(input)
	if err != nil {
		t.Fatalf("second AggregateScenarios() error = %v", err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("first = %#v, second = %#v", first, second)
	}
}

func aggregationConfig(method string, scenarioWeights map[string]float64) LoadedConfig {
	config := validLoadedConfig()
	config.Config.Aggregation = &AggregationConfig{
		Method:          method,
		ScenarioWeights: scenarioWeights,
	}
	return config
}
