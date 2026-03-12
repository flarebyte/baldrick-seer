package pipeline

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func TestDefaultScenarioRanker(t *testing.T) {
	t.Parallel()

	mixedConfig, mixedWeights := mixedCriteriaRankingFixture()

	tests := []struct {
		name        string
		config      LoadedConfig
		weights     []ScenarioCriterionWeights
		wantResults []domain.ScenarioRankingResult
		wantErr     error
	}{
		{
			name:    "basic topsis ranking with mixed cost and benefit criteria",
			config:  mixedConfig,
			weights: mixedWeights,
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
			config: oneCriterionRankingConfig(
				CriterionConfig{Name: "priority", Polarity: "benefit", ValueType: "ordinal", ScaleGuidance: []any{"low", "high"}},
				map[string]CriterionValue{"priority": {Kind: "ordinal", Value: 3}},
				map[string]CriterionValue{"priority": {Kind: "ordinal", Value: 1}},
			),
			weights: singleCriterionWeights("baseline", "priority", 1),
			wantResults: singleScenarioResults("baseline", []domain.RankedAlternative{
				{Name: "alpha", Rank: 1, Score: 1},
				{Name: "beta", Rank: 2, Score: 0},
			}),
		},
		{
			name: "ranking with boolean criteria",
			config: oneCriterionRankingConfig(
				CriterionConfig{Name: "approved", Polarity: "benefit", ValueType: "boolean"},
				map[string]CriterionValue{"approved": {Kind: "boolean", Value: true}},
				map[string]CriterionValue{"approved": {Kind: "boolean", Value: false}},
			),
			weights: singleCriterionWeights("baseline", "approved", 1),
			wantResults: singleScenarioResults("baseline", []domain.RankedAlternative{
				{Name: "alpha", Rank: 1, Score: 1},
				{Name: "beta", Rank: 2, Score: 0},
			}),
		},
		{
			name: "scenario with one eligible alternative",
			config: singleAlternativeRankingConfig(
				CriterionConfig{Name: "cost", Polarity: "cost", ValueType: "number"},
				"alpha",
				map[string]CriterionValue{"cost": {Kind: "number", Value: 3}},
			),
			weights: singleCriterionWeights("baseline", "cost", 1),
			wantResults: singleScenarioResults("baseline", []domain.RankedAlternative{
				{Name: "alpha", Rank: 1, Score: 0},
			}),
		},
		{
			name: "scenario where constraints exclude one alternative",
			config: constrainedRankingConfig(
				CriterionConfig{Name: "cost", Polarity: "cost", ValueType: "number"},
				[]AlternativeConfig{{Name: "alpha"}, {Name: "beta"}},
				[]ConstraintConfig{{CriterionName: "cost", Operator: "<=", Value: 5}},
				[]AlternativeEvaluationConfig{
					{AlternativeName: "alpha", Values: map[string]CriterionValue{"cost": {Kind: "number", Value: 3}}},
					{AlternativeName: "beta", Values: map[string]CriterionValue{"cost": {Kind: "number", Value: 8}}},
				},
			),
			weights: singleCriterionWeights("baseline", "cost", 1),
			wantResults: singleScenarioResults("baseline", []domain.RankedAlternative{
				{Name: "alpha", Rank: 1, Score: 0},
				{Name: "beta", Excluded: true},
			}),
		},
		{
			name: "scenario where constraints exclude all alternatives",
			config: constrainedRankingConfig(
				CriterionConfig{Name: "cost", Polarity: "cost", ValueType: "number"},
				[]AlternativeConfig{{Name: "beta"}, {Name: "alpha"}},
				[]ConstraintConfig{{CriterionName: "cost", Operator: "<=", Value: 1}},
				[]AlternativeEvaluationConfig{
					{AlternativeName: "beta", Values: map[string]CriterionValue{"cost": {Kind: "number", Value: 8}}},
					{AlternativeName: "alpha", Values: map[string]CriterionValue{"cost": {Kind: "number", Value: 3}}},
				},
			),
			weights: singleCriterionWeights("baseline", "cost", 1),
			wantResults: singleScenarioResults("baseline", []domain.RankedAlternative{
				{Name: "alpha", Excluded: true},
				{Name: "beta", Excluded: true},
			}),
		},
		{
			name: "stable tie break behavior when scores are equal",
			config: constrainedRankingConfig(
				CriterionConfig{Name: "quality", Polarity: "benefit", ValueType: "number"},
				[]AlternativeConfig{{Name: "beta"}, {Name: "alpha"}},
				nil,
				[]AlternativeEvaluationConfig{
					{AlternativeName: "beta", Values: map[string]CriterionValue{"quality": {Kind: "number", Value: 5}}},
					{AlternativeName: "alpha", Values: map[string]CriterionValue{"quality": {Kind: "number", Value: 5}}},
				},
			),
			weights: singleCriterionWeights("baseline", "quality", 1),
			wantResults: singleScenarioResults("baseline", []domain.RankedAlternative{
				{Name: "alpha", Rank: 1, Score: 0},
				{Name: "beta", Rank: 2, Score: 0},
			}),
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

			assertStageRunResult(t, func() (RankScenariosOutput, error) {
				return ranker.RankScenarios(context.Background(), RankScenariosInput{
					Command:         domain.CommandRequest{CommandName: domain.CommandNameReportGenerate, ConfigPath: tt.config.Path},
					ValidatedModel:  domain.ValidatedModelSummary{ConfigPath: tt.config.Path},
					ScenarioWeights: tt.weights,
					Config:          tt.config,
				})
			}, tt.wantErr, func(got RankScenariosOutput) {
				assertScenarioResults(t, got.ScenarioResults, tt.wantResults, 1e-9)
			})
		})
	}
}

func TestDefaultScenarioRankerIsDeterministic(t *testing.T) {
	t.Parallel()

	config, weights := mixedCriteriaRankingFixture()

	input := RankScenariosInput{
		Command:         domain.CommandRequest{CommandName: domain.CommandNameReportGenerate, ConfigPath: config.Path},
		ValidatedModel:  domain.ValidatedModelSummary{ConfigPath: config.Path},
		ScenarioWeights: weights,
		Config:          config,
	}

	assertRepeatedDeepEqual(t, 1, func() (RankScenariosOutput, error) {
		return DefaultScenarioRanker{}.RankScenarios(context.Background(), input)
	})
}

func mixedCriteriaRankingFixture() (LoadedConfig, []ScenarioCriterionWeights) {
	return rankingConfig(
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
		), []ScenarioCriterionWeights{{
			ScenarioName: "baseline",
			CriterionWeights: []CriterionWeight{
				{CriterionName: "cost", Weight: 0.5},
				{CriterionName: "quality", Weight: 0.5},
			},
		}}
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

			got, err := runner.RunReportGenerate(context.Background(), domain.CommandRequest{
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

	_, err := runner.RunReportGenerate(context.Background(), domain.CommandRequest{
		CommandName: domain.CommandNameReportGenerate,
		ConfigPath:  fixtureConfigPath(),
	})
	if !errors.Is(err, ErrRankingFailed) {
		t.Fatalf("error = %v, want %v", err, ErrRankingFailed)
	}

	_ = assertFailureCategory(t, err, ErrRankingFailed, domain.FailureCategoryExecution, "")

	if got, want := order, []string{"load", "validate", "weight", "rank"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("order = %#v, want %#v", got, want)
	}
}

type recordingScenarioRanker struct {
	recorder *[]string
	inner    ScenarioRanker
}

func (r recordingScenarioRanker) RankScenarios(ctx context.Context, input RankScenariosInput) (RankScenariosOutput, error) {
	*r.recorder = append(*r.recorder, "rank")
	return r.inner.RankScenarios(ctx, input)
}

type fixedScenarioWeighter struct {
	recorder        *[]string
	scenarioWeights []ScenarioCriterionWeights
}

func (f *fixedScenarioWeighter) WeightCriteria(_ context.Context, input WeightCriteriaInput) (WeightCriteriaOutput, error) {
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

func oneCriterionRankingConfig(
	criterion CriterionConfig,
	alphaValues map[string]CriterionValue,
	betaValues map[string]CriterionValue,
) LoadedConfig {
	return rankingConfig(
		[]CriterionConfig{criterion},
		[]AlternativeConfig{{Name: "alpha"}, {Name: "beta"}},
		ScenarioConfig{
			Name:           "baseline",
			ActiveCriteria: []ScenarioCriterionRef{{CriterionName: criterion.Name}},
		},
		[]AlternativeEvaluationConfig{
			{AlternativeName: "alpha", Values: alphaValues},
			{AlternativeName: "beta", Values: betaValues},
		},
	)
}

func singleAlternativeRankingConfig(
	criterion CriterionConfig,
	alternativeName string,
	values map[string]CriterionValue,
) LoadedConfig {
	return rankingConfig(
		[]CriterionConfig{criterion},
		[]AlternativeConfig{{Name: alternativeName}},
		ScenarioConfig{
			Name:           "baseline",
			ActiveCriteria: []ScenarioCriterionRef{{CriterionName: criterion.Name}},
		},
		[]AlternativeEvaluationConfig{
			{AlternativeName: alternativeName, Values: values},
		},
	)
}

func constrainedRankingConfig(
	criterion CriterionConfig,
	alternatives []AlternativeConfig,
	constraints []ConstraintConfig,
	evaluations []AlternativeEvaluationConfig,
) LoadedConfig {
	return rankingConfig(
		[]CriterionConfig{criterion},
		alternatives,
		ScenarioConfig{
			Name:           "baseline",
			ActiveCriteria: []ScenarioCriterionRef{{CriterionName: criterion.Name}},
			Constraints:    constraints,
		},
		evaluations,
	)
}
