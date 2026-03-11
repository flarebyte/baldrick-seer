package pipeline

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func TestDefaultScenarioRanker(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		config      LoadedConfig
		weights     []ScenarioCriterionWeights
		wantResults []domain.ScenarioRankingResult
		wantErr     error
	}{
		{
			name: "basic topsis ranking with mixed cost and benefit criteria",
			config: rankingConfig(
				[]CriterionConfig{
					{Name: "cost", Polarity: "cost", ValueType: "number"},
					{Name: "quality", Polarity: "benefit", ValueType: "number"},
				},
				[]AlternativeConfig{{Name: "alpha"}, {Name: "beta"}},
				ScenarioConfig{
					Name: "baseline",
					ActiveCriteria: []ScenarioCriterionRef{
						{CriterionName: "cost"},
						{CriterionName: "quality"},
					},
				},
				[]AlternativeEvaluationConfig{
					{AlternativeName: "alpha", Values: map[string]CriterionValue{"cost": {Kind: "number", Value: 1}, "quality": {Kind: "number", Value: 4}}},
					{AlternativeName: "beta", Values: map[string]CriterionValue{"cost": {Kind: "number", Value: 2}, "quality": {Kind: "number", Value: 5}}},
				},
			),
			weights: []ScenarioCriterionWeights{{
				ScenarioName: "baseline",
				CriterionWeights: []CriterionWeight{
					{CriterionName: "cost", Weight: 0.5},
					{CriterionName: "quality", Weight: 0.5},
				},
			}},
			wantResults: []domain.ScenarioRankingResult{{
				ScenarioName: "baseline",
				RankedAlternatives: []domain.RankedAlternative{
					{Name: "alpha", Rank: 1, Score: 0.7411716371},
					{Name: "beta", Rank: 2, Score: 0.2588283629},
				},
			}},
		},
		{
			name: "ranking with ordinal criteria",
			config: rankingConfig(
				[]CriterionConfig{{Name: "priority", Polarity: "benefit", ValueType: "ordinal", ScaleGuidance: []any{"low", "high"}}},
				[]AlternativeConfig{{Name: "alpha"}, {Name: "beta"}},
				ScenarioConfig{
					Name:           "baseline",
					ActiveCriteria: []ScenarioCriterionRef{{CriterionName: "priority"}},
				},
				[]AlternativeEvaluationConfig{
					{AlternativeName: "alpha", Values: map[string]CriterionValue{"priority": {Kind: "ordinal", Value: 3}}},
					{AlternativeName: "beta", Values: map[string]CriterionValue{"priority": {Kind: "ordinal", Value: 1}}},
				},
			),
			weights: []ScenarioCriterionWeights{{
				ScenarioName: "baseline",
				CriterionWeights: []CriterionWeight{
					{CriterionName: "priority", Weight: 1},
				},
			}},
			wantResults: []domain.ScenarioRankingResult{{
				ScenarioName: "baseline",
				RankedAlternatives: []domain.RankedAlternative{
					{Name: "alpha", Rank: 1, Score: 1},
					{Name: "beta", Rank: 2, Score: 0},
				},
			}},
		},
		{
			name: "ranking with boolean criteria",
			config: rankingConfig(
				[]CriterionConfig{{Name: "approved", Polarity: "benefit", ValueType: "boolean"}},
				[]AlternativeConfig{{Name: "alpha"}, {Name: "beta"}},
				ScenarioConfig{
					Name:           "baseline",
					ActiveCriteria: []ScenarioCriterionRef{{CriterionName: "approved"}},
				},
				[]AlternativeEvaluationConfig{
					{AlternativeName: "alpha", Values: map[string]CriterionValue{"approved": {Kind: "boolean", Value: true}}},
					{AlternativeName: "beta", Values: map[string]CriterionValue{"approved": {Kind: "boolean", Value: false}}},
				},
			),
			weights: []ScenarioCriterionWeights{{
				ScenarioName: "baseline",
				CriterionWeights: []CriterionWeight{
					{CriterionName: "approved", Weight: 1},
				},
			}},
			wantResults: []domain.ScenarioRankingResult{{
				ScenarioName: "baseline",
				RankedAlternatives: []domain.RankedAlternative{
					{Name: "alpha", Rank: 1, Score: 1},
					{Name: "beta", Rank: 2, Score: 0},
				},
			}},
		},
		{
			name: "scenario with one eligible alternative",
			config: rankingConfig(
				[]CriterionConfig{{Name: "cost", Polarity: "cost", ValueType: "number"}},
				[]AlternativeConfig{{Name: "alpha"}},
				ScenarioConfig{
					Name:           "baseline",
					ActiveCriteria: []ScenarioCriterionRef{{CriterionName: "cost"}},
				},
				[]AlternativeEvaluationConfig{
					{AlternativeName: "alpha", Values: map[string]CriterionValue{"cost": {Kind: "number", Value: 3}}},
				},
			),
			weights: []ScenarioCriterionWeights{{
				ScenarioName:     "baseline",
				CriterionWeights: []CriterionWeight{{CriterionName: "cost", Weight: 1}},
			}},
			wantResults: []domain.ScenarioRankingResult{{
				ScenarioName: "baseline",
				RankedAlternatives: []domain.RankedAlternative{
					{Name: "alpha", Rank: 1, Score: 0},
				},
			}},
		},
		{
			name: "scenario where constraints exclude one alternative",
			config: rankingConfig(
				[]CriterionConfig{{Name: "cost", Polarity: "cost", ValueType: "number"}},
				[]AlternativeConfig{{Name: "alpha"}, {Name: "beta"}},
				ScenarioConfig{
					Name:           "baseline",
					ActiveCriteria: []ScenarioCriterionRef{{CriterionName: "cost"}},
					Constraints:    []ConstraintConfig{{CriterionName: "cost", Operator: "<=", Value: 5}},
				},
				[]AlternativeEvaluationConfig{
					{AlternativeName: "alpha", Values: map[string]CriterionValue{"cost": {Kind: "number", Value: 3}}},
					{AlternativeName: "beta", Values: map[string]CriterionValue{"cost": {Kind: "number", Value: 8}}},
				},
			),
			weights: []ScenarioCriterionWeights{{
				ScenarioName:     "baseline",
				CriterionWeights: []CriterionWeight{{CriterionName: "cost", Weight: 1}},
			}},
			wantResults: []domain.ScenarioRankingResult{{
				ScenarioName: "baseline",
				RankedAlternatives: []domain.RankedAlternative{
					{Name: "alpha", Rank: 1, Score: 0},
					{Name: "beta", Excluded: true},
				},
			}},
		},
		{
			name: "scenario where constraints exclude all alternatives",
			config: rankingConfig(
				[]CriterionConfig{{Name: "cost", Polarity: "cost", ValueType: "number"}},
				[]AlternativeConfig{{Name: "beta"}, {Name: "alpha"}},
				ScenarioConfig{
					Name:           "baseline",
					ActiveCriteria: []ScenarioCriterionRef{{CriterionName: "cost"}},
					Constraints:    []ConstraintConfig{{CriterionName: "cost", Operator: "<=", Value: 1}},
				},
				[]AlternativeEvaluationConfig{
					{AlternativeName: "beta", Values: map[string]CriterionValue{"cost": {Kind: "number", Value: 8}}},
					{AlternativeName: "alpha", Values: map[string]CriterionValue{"cost": {Kind: "number", Value: 3}}},
				},
			),
			weights: []ScenarioCriterionWeights{{
				ScenarioName:     "baseline",
				CriterionWeights: []CriterionWeight{{CriterionName: "cost", Weight: 1}},
			}},
			wantResults: []domain.ScenarioRankingResult{{
				ScenarioName: "baseline",
				RankedAlternatives: []domain.RankedAlternative{
					{Name: "alpha", Excluded: true},
					{Name: "beta", Excluded: true},
				},
			}},
		},
		{
			name: "stable tie break behavior when scores are equal",
			config: rankingConfig(
				[]CriterionConfig{{Name: "quality", Polarity: "benefit", ValueType: "number"}},
				[]AlternativeConfig{{Name: "beta"}, {Name: "alpha"}},
				ScenarioConfig{
					Name:           "baseline",
					ActiveCriteria: []ScenarioCriterionRef{{CriterionName: "quality"}},
				},
				[]AlternativeEvaluationConfig{
					{AlternativeName: "beta", Values: map[string]CriterionValue{"quality": {Kind: "number", Value: 5}}},
					{AlternativeName: "alpha", Values: map[string]CriterionValue{"quality": {Kind: "number", Value: 5}}},
				},
			),
			weights: []ScenarioCriterionWeights{{
				ScenarioName:     "baseline",
				CriterionWeights: []CriterionWeight{{CriterionName: "quality", Weight: 1}},
			}},
			wantResults: []domain.ScenarioRankingResult{{
				ScenarioName: "baseline",
				RankedAlternatives: []domain.RankedAlternative{
					{Name: "alpha", Rank: 1, Score: 0},
					{Name: "beta", Rank: 2, Score: 0},
				},
			}},
		},
		{
			name: "internal inconsistency fails deterministically",
			config: rankingConfig(
				[]CriterionConfig{{Name: "cost", Polarity: "", ValueType: "number"}},
				[]AlternativeConfig{{Name: "alpha"}, {Name: "beta"}},
				ScenarioConfig{
					Name:           "baseline",
					ActiveCriteria: []ScenarioCriterionRef{{CriterionName: "cost"}},
				},
				[]AlternativeEvaluationConfig{
					{AlternativeName: "alpha", Values: map[string]CriterionValue{"cost": {Kind: "number", Value: 1}}},
					{AlternativeName: "beta", Values: map[string]CriterionValue{"cost": {Kind: "number", Value: 2}}},
				},
			),
			weights: []ScenarioCriterionWeights{{
				ScenarioName:     "baseline",
				CriterionWeights: []CriterionWeight{{CriterionName: "cost", Weight: 1}},
			}},
			wantErr: ErrRankingFailed,
		},
	}

	ranker := DefaultScenarioRanker{}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ranker.RankScenarios(RankScenariosInput{
				Command:         domain.CommandRequest{CommandName: domain.CommandNameReportGenerate, ConfigPath: tt.config.Path},
				ValidatedModel:  domain.ValidatedModelSummary{ConfigPath: tt.config.Path},
				ScenarioWeights: tt.weights,
				Config:          tt.config,
			})

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("RankScenarios() error = %v", err)
			}

			assertScenarioResults(t, got.ScenarioResults, tt.wantResults, 1e-9)
		})
	}
}

func TestDefaultScenarioRankerIsDeterministic(t *testing.T) {
	t.Parallel()

	config := rankingConfig(
		[]CriterionConfig{
			{Name: "cost", Polarity: "cost", ValueType: "number"},
			{Name: "quality", Polarity: "benefit", ValueType: "number"},
		},
		[]AlternativeConfig{{Name: "alpha"}, {Name: "beta"}},
		ScenarioConfig{
			Name: "baseline",
			ActiveCriteria: []ScenarioCriterionRef{
				{CriterionName: "cost"},
				{CriterionName: "quality"},
			},
		},
		[]AlternativeEvaluationConfig{
			{AlternativeName: "alpha", Values: map[string]CriterionValue{"cost": {Kind: "number", Value: 1}, "quality": {Kind: "number", Value: 4}}},
			{AlternativeName: "beta", Values: map[string]CriterionValue{"cost": {Kind: "number", Value: 2}, "quality": {Kind: "number", Value: 5}}},
		},
	)

	input := RankScenariosInput{
		Command:        domain.CommandRequest{CommandName: domain.CommandNameReportGenerate, ConfigPath: config.Path},
		ValidatedModel: domain.ValidatedModelSummary{ConfigPath: config.Path},
		ScenarioWeights: []ScenarioCriterionWeights{{
			ScenarioName: "baseline",
			CriterionWeights: []CriterionWeight{
				{CriterionName: "cost", Weight: 0.5},
				{CriterionName: "quality", Weight: 0.5},
			},
		}},
		Config: config,
	}

	first, err := DefaultScenarioRanker{}.RankScenarios(input)
	if err != nil {
		t.Fatalf("first RankScenarios() error = %v", err)
	}

	second, err := DefaultScenarioRanker{}.RankScenarios(input)
	if err != nil {
		t.Fatalf("second RankScenarios() error = %v", err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("first = %#v, second = %#v", first, second)
	}
}

func TestRunReportGenerateUsesRealScenarioResults(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
		want []domain.ScenarioRankingResult
	}{
		{
			name: "valid scenario ranking",
			path: filepath.Join("..", "..", "testdata", "config", "topsis_valid.cue"),
			want: []domain.ScenarioRankingResult{{
				ScenarioName: "baseline",
				RankedAlternatives: []domain.RankedAlternative{
					{Name: "alpha", Rank: 1, Score: 0.8691890100},
					{Name: "beta", Rank: 2, Score: 0.1308109900},
				},
			}},
		},
		{
			name: "valid scenario ranking with exclusion status",
			path: filepath.Join("..", "..", "testdata", "config", "topsis_exclusion.cue"),
			want: []domain.ScenarioRankingResult{{
				ScenarioName: "baseline",
				RankedAlternatives: []domain.RankedAlternative{
					{Name: "alpha", Rank: 1, Score: 0},
					{Name: "beta", Excluded: true},
				},
			}},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			runner := Runner{
				ConfigLoader:       DefaultConfigLoader{},
				ModelValidator:     DefaultModelValidator{},
				CriteriaWeighter:   DefaultCriteriaWeighter{},
				ScenarioRanker:     DefaultScenarioRanker{},
				ScenarioAggregator: DefaultScenarioAggregator{},
				ReportRenderer:     DefaultReportRenderer{},
			}

			got, err := runner.RunReportGenerate(domain.CommandRequest{
				CommandName: domain.CommandNameReportGenerate,
				ConfigPath:  tt.path,
			})
			if err != nil {
				t.Fatalf("RunReportGenerate() error = %v", err)
			}

			assertScenarioResults(t, got.ScenarioResults, tt.want, 1e-9)
		})
	}
}

func TestRunReportGenerateStopsOnRankingFailure(t *testing.T) {
	t.Parallel()

	order := []string{}
	runner := Runner{
		ConfigLoader: &fakeConfigLoader{recorder: &order, output: LoadConfigOutput{Config: rankingConfig(
			[]CriterionConfig{{Name: "cost", Polarity: "", ValueType: "number"}},
			[]AlternativeConfig{{Name: "alpha"}, {Name: "beta"}},
			ScenarioConfig{Name: "baseline", ActiveCriteria: []ScenarioCriterionRef{{CriterionName: "cost"}}},
			[]AlternativeEvaluationConfig{
				{AlternativeName: "alpha", Values: map[string]CriterionValue{"cost": {Kind: "number", Value: 1}}},
				{AlternativeName: "beta", Values: map[string]CriterionValue{"cost": {Kind: "number", Value: 2}}},
			},
		)}},
		ModelValidator:     &fakeModelValidator{recorder: &order},
		CriteriaWeighter:   &fixedScenarioWeighter{recorder: &order, scenarioWeights: []ScenarioCriterionWeights{{ScenarioName: "baseline", CriterionWeights: []CriterionWeight{{CriterionName: "cost", Weight: 1}}}}},
		ScenarioRanker:     recordingScenarioRanker{recorder: &order, inner: DefaultScenarioRanker{}},
		ScenarioAggregator: &fakeScenarioAggregator{recorder: &order},
		ReportRenderer:     &fakeReportRenderer{recorder: &order},
	}

	_, err := runner.RunReportGenerate(domain.CommandRequest{
		CommandName: domain.CommandNameReportGenerate,
		ConfigPath:  fixtureConfigPath(),
	})
	if !errors.Is(err, ErrRankingFailed) {
		t.Fatalf("error = %v, want %v", err, ErrRankingFailed)
	}

	failure := domain.AsCommandFailure(err)
	if failure == nil {
		t.Fatal("AsCommandFailure(err) = nil, want value")
	}

	if failure.Category != domain.FailureCategoryExecution {
		t.Fatalf("Category = %q, want %q", failure.Category, domain.FailureCategoryExecution)
	}

	if got, want := order, []string{"load", "validate", "weight", "rank"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("order = %#v, want %#v", got, want)
	}
}

type recordingScenarioRanker struct {
	recorder *[]string
	inner    ScenarioRanker
}

func (r recordingScenarioRanker) RankScenarios(input RankScenariosInput) (RankScenariosOutput, error) {
	*r.recorder = append(*r.recorder, "rank")
	return r.inner.RankScenarios(input)
}

type fixedScenarioWeighter struct {
	recorder        *[]string
	scenarioWeights []ScenarioCriterionWeights
}

func (f *fixedScenarioWeighter) WeightCriteria(WeightCriteriaInput) (WeightCriteriaOutput, error) {
	*f.recorder = append(*f.recorder, "weight")
	return WeightCriteriaOutput{ScenarioWeights: f.scenarioWeights}, nil
}

func rankingConfig(
	criteria []CriterionConfig,
	alternatives []AlternativeConfig,
	scenario ScenarioConfig,
	evaluations []AlternativeEvaluationConfig,
) LoadedConfig {
	config := validLoadedConfig()
	config.Config.CriteriaCatalog = append([]CriterionConfig(nil), criteria...)
	config.Config.Alternatives = append([]AlternativeConfig(nil), alternatives...)
	config.Config.Scenarios = []ScenarioConfig{scenario}
	config.Config.Evaluations = []EvaluationConfig{{
		ScenarioName: scenario.Name,
		Evaluations:  append([]AlternativeEvaluationConfig(nil), evaluations...),
	}}
	return config
}
