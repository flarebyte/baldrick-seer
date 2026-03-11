package pipeline

import (
	"errors"
	"math"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func TestDefaultCriteriaWeighter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		config      LoadedConfig
		wantWeights []ScenarioCriterionWeights
		wantErr     error
	}{
		{
			name: "scenario with 0 active criteria",
			config: func() LoadedConfig {
				config := validLoadedConfig()
				config.Config.Scenarios[0].ActiveCriteria = nil
				config.Config.Evaluations[0].Evaluations[0].Values = nil
				return config
			}(),
			wantWeights: []ScenarioCriterionWeights{{ScenarioName: "baseline"}},
		},
		{
			name:   "scenario with 1 active criterion",
			config: validLoadedConfig(),
			wantWeights: []ScenarioCriterionWeights{{
				ScenarioName: "baseline",
				CriterionWeights: []CriterionWeight{
					{CriterionName: "cost", Weight: 1},
				},
			}},
		},
		{
			name: "scenario with 2 criteria and one comparison",
			config: validLoadedConfigWithAHPPairs(
				[]string{"cost", "speed"},
				[]PairwiseComparison{
					{MoreImportantCriterionName: "cost", LessImportantCriterionName: "speed", Strength: 3},
				},
			),
			wantWeights: []ScenarioCriterionWeights{{
				ScenarioName: "baseline",
				CriterionWeights: []CriterionWeight{
					{CriterionName: "cost", Weight: 0.75},
					{CriterionName: "speed", Weight: 0.25},
				},
			}},
		},
		{
			name: "scenario with 3 criteria and full pairwise set",
			config: validLoadedConfigWithAHPPairs(
				[]string{"cost", "speed", "reliability"},
				[]PairwiseComparison{
					{MoreImportantCriterionName: "cost", LessImportantCriterionName: "speed", Strength: 3},
					{MoreImportantCriterionName: "cost", LessImportantCriterionName: "reliability", Strength: 5},
					{MoreImportantCriterionName: "speed", LessImportantCriterionName: "reliability", Strength: 2},
				},
			),
			wantWeights: []ScenarioCriterionWeights{{
				ScenarioName: "baseline",
				CriterionWeights: []CriterionWeight{
					{CriterionName: "cost", Weight: 0.6479468599},
					{CriterionName: "reliability", Weight: 0.1221819646},
					{CriterionName: "speed", Weight: 0.2298711755},
				},
			}},
		},
		{
			name: "internal inconsistency fails deterministically",
			config: func() LoadedConfig {
				config := validLoadedConfigWithAHPPairs(
					[]string{"cost", "speed"},
					[]PairwiseComparison{
						{MoreImportantCriterionName: "cost", LessImportantCriterionName: "speed", Strength: 0},
					},
				)
				return config
			}(),
			wantErr: ErrWeightingFailed,
		},
	}

	weighter := DefaultCriteriaWeighter{}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := weighter.WeightCriteria(WeightCriteriaInput{
				Command: domain.CommandRequest{
					CommandName: domain.CommandNameReportGenerate,
					ConfigPath:  tt.config.Path,
				},
				ValidatedModel: domain.ValidatedModelSummary{ConfigPath: tt.config.Path},
				Config:         tt.config,
			})

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("WeightCriteria() error = %v", err)
			}

			assertScenarioWeights(t, got.ScenarioWeights, tt.wantWeights, 1e-9)
		})
	}
}

func TestDefaultCriteriaWeighterNormalizesWeights(t *testing.T) {
	t.Parallel()

	config := validThreeCriteriaAHPConfig()

	got, err := DefaultCriteriaWeighter{}.WeightCriteria(WeightCriteriaInput{
		Command:        domain.CommandRequest{CommandName: domain.CommandNameReportGenerate, ConfigPath: config.Path},
		ValidatedModel: domain.ValidatedModelSummary{ConfigPath: config.Path},
		Config:         config,
	})
	if err != nil {
		t.Fatalf("WeightCriteria() error = %v", err)
	}

	sum := 0.0
	for _, weight := range got.ScenarioWeights[0].CriterionWeights {
		sum += weight.Weight
	}

	if math.Abs(sum-1) > 1e-9 {
		t.Fatalf("sum = %0.12f, want 1", sum)
	}
}

func TestDefaultCriteriaWeighterIsDeterministic(t *testing.T) {
	t.Parallel()

	config := validThreeCriteriaAHPConfig()

	input := WeightCriteriaInput{
		Command:        domain.CommandRequest{CommandName: domain.CommandNameReportGenerate, ConfigPath: config.Path},
		ValidatedModel: domain.ValidatedModelSummary{ConfigPath: config.Path},
		Config:         config,
	}

	first, err := DefaultCriteriaWeighter{}.WeightCriteria(input)
	if err != nil {
		t.Fatalf("first WeightCriteria() error = %v", err)
	}

	second, err := DefaultCriteriaWeighter{}.WeightCriteria(input)
	if err != nil {
		t.Fatalf("second WeightCriteria() error = %v", err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("first = %#v, second = %#v", first, second)
	}
}

func TestRunReportGenerateUsesRealScenarioWeights(t *testing.T) {
	t.Parallel()

	order := []string{}
	ranker := &capturingScenarioRanker{recorder: &order}
	runner := Runner{
		ConfigLoader:       DefaultConfigLoader{},
		ModelValidator:     DefaultModelValidator{},
		CriteriaWeighter:   DefaultCriteriaWeighter{},
		ScenarioRanker:     ranker,
		ScenarioAggregator: &fakeScenarioAggregator{recorder: &order},
		ReportRenderer:     &fakeReportRenderer{recorder: &order},
	}

	_, err := runner.RunReportGenerate(domain.CommandRequest{
		CommandName: domain.CommandNameReportGenerate,
		ConfigPath:  filepath.Join("..", "..", "testdata", "config", "pairwise_valid.cue"),
	})
	if err != nil {
		t.Fatalf("RunReportGenerate() error = %v", err)
	}

	want := []ScenarioCriterionWeights{{
		ScenarioName: "baseline",
		CriterionWeights: []CriterionWeight{
			{CriterionName: "cost", Weight: 0.6479468599},
			{CriterionName: "reliability", Weight: 0.1221819646},
			{CriterionName: "speed", Weight: 0.2298711755},
		},
	}}
	assertScenarioWeights(t, ranker.scenarioWeights, want, 1e-9)
}

func TestRunReportGenerateStopsOnRealWeightingFailure(t *testing.T) {
	t.Parallel()

	order := []string{}
	runner := Runner{
		ConfigLoader: &fakeConfigLoader{
			recorder: &order,
			output: LoadConfigOutput{
				Config: LoadedConfig{
					Path: fixtureConfigPath(),
					Config: &ExecutionConfig{
						Problem:         &ProblemConfig{Name: "minimal"},
						Reports:         []ReportConfig{{Name: "summary", Title: "Summary", Format: "markdown"}},
						CriteriaCatalog: []CriterionConfig{{Name: "cost", ValueType: "number"}, {Name: "speed", ValueType: "number"}},
						Alternatives:    []AlternativeConfig{{Name: "option_a"}},
						Scenarios: []ScenarioConfig{{
							Name: "baseline",
							ActiveCriteria: []ScenarioCriterionRef{
								{CriterionName: "cost"},
								{CriterionName: "speed"},
							},
						}},
						Evaluations: []EvaluationConfig{{
							ScenarioName: "baseline",
							Evaluations: []AlternativeEvaluationConfig{{
								AlternativeName: "option_a",
								Values: map[string]CriterionValue{
									"cost":  {Kind: "number", Value: 1},
									"speed": {Kind: "number", Value: 2},
								},
							}},
						}},
						Aggregation: &AggregationConfig{},
					},
				},
			},
		},
		ModelValidator:     &fakeModelValidator{recorder: &order},
		CriteriaWeighter:   DefaultCriteriaWeighter{},
		ScenarioRanker:     &fakeScenarioRanker{recorder: &order},
		ScenarioAggregator: &fakeScenarioAggregator{recorder: &order},
		ReportRenderer:     &fakeReportRenderer{recorder: &order},
	}

	_, err := runner.RunReportGenerate(domain.CommandRequest{
		CommandName: domain.CommandNameReportGenerate,
		ConfigPath:  fixtureConfigPath(),
	})
	if !errors.Is(err, ErrWeightingFailed) {
		t.Fatalf("error = %v, want %v", err, ErrWeightingFailed)
	}

	failure := domain.AsCommandFailure(err)
	if failure == nil {
		t.Fatal("AsCommandFailure(err) = nil, want value")
	}

	if failure.Category != domain.FailureCategoryExecution {
		t.Fatalf("Category = %q, want %q", failure.Category, domain.FailureCategoryExecution)
	}

	if got, want := order, []string{"load", "validate"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("order = %#v, want %#v", got, want)
	}
}

type capturingScenarioRanker struct {
	recorder        *[]string
	scenarioWeights []ScenarioCriterionWeights
}

func (c *capturingScenarioRanker) RankScenarios(input RankScenariosInput) (RankScenariosOutput, error) {
	*c.recorder = append(*c.recorder, "rank")
	c.scenarioWeights = append([]ScenarioCriterionWeights(nil), input.ScenarioWeights...)
	return RankScenariosOutput{}, nil
}

func validThreeCriteriaAHPConfig() LoadedConfig {
	return validLoadedConfigWithAHPPairs(
		[]string{"cost", "speed", "reliability"},
		[]PairwiseComparison{
			{MoreImportantCriterionName: "cost", LessImportantCriterionName: "speed", Strength: 3},
			{MoreImportantCriterionName: "cost", LessImportantCriterionName: "reliability", Strength: 5},
			{MoreImportantCriterionName: "speed", LessImportantCriterionName: "reliability", Strength: 2},
		},
	)
}
