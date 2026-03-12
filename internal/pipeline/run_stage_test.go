package pipeline

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

type contextKey string

func TestRunValidateStageOrdering(t *testing.T) {
	t.Parallel()

	order := []string{}
	runner := newFakeRunner(&order)

	_, err := runner.RunValidate(context.Background(), domain.CommandRequest{
		CommandName: domain.CommandNameValidate,
		ConfigPath:  filepath.Join("testdata", "config", "minimal.cue"),
	})
	if err != nil {
		t.Fatalf("RunValidate() error = %v", err)
	}

	assertStageOrder(t, order, []string{"load", "validate"})
}

func TestRunReportGenerateStageOrdering(t *testing.T) {
	t.Parallel()

	order := []string{}
	runner := newFakeRunner(&order)

	_, err := runner.RunReportGenerate(context.Background(), domain.CommandRequest{
		CommandName: domain.CommandNameReportGenerate,
		ConfigPath:  filepath.Join("testdata", "config", "minimal.cue"),
	})
	if err != nil {
		t.Fatalf("RunReportGenerate() error = %v", err)
	}

	assertStageOrder(t, order, []string{"load", "validate", "weight", "rank", "aggregate", "render"})
}

func TestRunReportGenerateSelectsCurrentV1Strategy(t *testing.T) {
	t.Parallel()

	order := []string{}
	strategy := fakeRankingStrategy{
		recorder: &order,
		output: RankingStrategyOutput{
			ScenarioWeights: []ScenarioCriterionWeights{{
				ScenarioName: "baseline",
				CriterionWeights: []CriterionWeight{
					{CriterionName: "cost", Weight: 1},
				},
			}},
			ScenarioResults: []domain.ScenarioRankingResult{{
				ScenarioName: "baseline",
				RankedAlternatives: []domain.RankedAlternative{
					{Name: "option_a", Rank: 1, Score: 0.9},
				},
			}},
		},
	}
	runner := Runner{
		ConfigLoader:       &fakeConfigLoader{recorder: &order},
		ModelValidator:     &fakeModelValidator{recorder: &order},
		RankingStrategies:  fakeRankingStrategySelector{recorder: &order, method: RankingMethodV1AHPTopsis, strategy: strategy},
		ScenarioAggregator: &fakeScenarioAggregator{recorder: &order},
		ReportRenderer:     &fakeReportRenderer{recorder: &order},
	}

	_, err := runReportGenerateForTest(runner)
	if err != nil {
		t.Fatalf("RunReportGenerate() error = %v", err)
	}

	assertStageOrder(t, order, []string{"load", "validate", "select-strategy", "strategy:v1_ahp_topsis", "strategy-execute", "aggregate", "render"})
}

func TestRunReportGenerateUnsupportedStrategyFailsDeterministically(t *testing.T) {
	t.Parallel()

	order := []string{}
	wantErr := ErrRankingFailed
	runner := Runner{
		ConfigLoader:       &fakeConfigLoader{recorder: &order},
		ModelValidator:     &fakeModelValidator{recorder: &order},
		RankingStrategies:  fakeRankingStrategySelector{recorder: &order, method: RankingMethodElectre, strategyErr: wantErr},
		ScenarioAggregator: &fakeScenarioAggregator{recorder: &order},
		ReportRenderer:     &fakeReportRenderer{recorder: &order},
	}

	_, err := runReportGenerateForTest(runner)

	_ = assertFailureCategory(t, err, wantErr, domain.FailureCategoryExecution, "command failed")
	assertStageOrder(t, order, []string{"load", "validate", "select-strategy", "strategy:electre"})
}

func TestDefaultRankingStrategySelectorIsDeterministic(t *testing.T) {
	t.Parallel()

	selector := newDefaultRankingStrategySelector(DefaultCriteriaWeighter{}, DefaultScenarioRanker{})
	config := validLoadedConfig()
	validation := ValidateModelOutput{
		ValidatedModel: domain.ValidatedModelSummary{ConfigPath: config.Path},
	}

	firstMethod, err := selector.Select(config, validation)
	if err != nil {
		t.Fatalf("first Select() error = %v", err)
	}
	secondMethod, err := selector.Select(config, validation)
	if err != nil {
		t.Fatalf("second Select() error = %v", err)
	}
	if firstMethod != secondMethod {
		t.Fatalf("methods differed: %q vs %q", firstMethod, secondMethod)
	}

	if got, want := firstMethod, RankingMethodV1AHPTopsis; got != want {
		t.Fatalf("method = %q, want %q", got, want)
	}
}

func TestRunReportGenerateFailsFast(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("weight failed")
	order := []string{}
	runner := Runner{
		ConfigLoader:       &fakeConfigLoader{recorder: &order},
		ModelValidator:     &fakeModelValidator{recorder: &order},
		CriteriaWeighter:   &fakeCriteriaWeighter{recorder: &order, err: wantErr},
		ScenarioRanker:     &fakeScenarioRanker{recorder: &order},
		ScenarioAggregator: &fakeScenarioAggregator{recorder: &order},
		ReportRenderer:     &fakeReportRenderer{recorder: &order},
	}

	_, err := runner.RunReportGenerate(context.Background(), domain.CommandRequest{
		CommandName: domain.CommandNameReportGenerate,
		ConfigPath:  filepath.Join("testdata", "config", "minimal.cue"),
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}

	_ = assertFailureCategory(t, err, wantErr, domain.FailureCategoryExecution, "")
	assertStageOrder(t, order, []string{"load", "validate", "weight"})
}

func TestRunValidatePropagatesContext(t *testing.T) {
	t.Parallel()

	order := []string{}
	contexts := []context.Context{}
	runner := Runner{
		ConfigLoader:   &fakeConfigLoader{recorder: &order, contexts: &contexts},
		ModelValidator: &fakeModelValidator{recorder: &order, contexts: &contexts},
	}

	ctx := context.WithValue(context.Background(), contextKey("validate"), "validate")
	_, err := runner.RunValidate(ctx, domain.CommandRequest{
		CommandName: domain.CommandNameValidate,
		ConfigPath:  fixtureConfigPath(),
	})
	if err != nil {
		t.Fatalf("RunValidate() error = %v", err)
	}

	assertStageOrder(t, order, []string{"load", "validate"})
	if len(contexts) != 2 {
		t.Fatalf("len(contexts) = %d, want 2", len(contexts))
	}
	for index, got := range contexts {
		if got != ctx {
			t.Fatalf("contexts[%d] did not match input context", index)
		}
	}
}

func TestRunReportGeneratePropagatesContext(t *testing.T) {
	t.Parallel()

	order := []string{}
	contexts := []context.Context{}
	runner := Runner{
		ConfigLoader:       &fakeConfigLoader{recorder: &order, contexts: &contexts},
		ModelValidator:     &fakeModelValidator{recorder: &order, contexts: &contexts},
		CriteriaWeighter:   &fakeCriteriaWeighter{recorder: &order, contexts: &contexts},
		ScenarioRanker:     &fakeScenarioRanker{recorder: &order, contexts: &contexts},
		ScenarioAggregator: &fakeScenarioAggregator{recorder: &order, contexts: &contexts},
		ReportRenderer:     &fakeReportRenderer{recorder: &order, contexts: &contexts},
	}

	ctx := context.WithValue(context.Background(), contextKey("report"), "report")
	_, err := runner.RunReportGenerate(ctx, domain.CommandRequest{
		CommandName: domain.CommandNameReportGenerate,
		ConfigPath:  fixtureConfigPath(),
	})
	if err != nil {
		t.Fatalf("RunReportGenerate() error = %v", err)
	}

	assertStageOrder(t, order, []string{"load", "validate", "weight", "rank", "aggregate", "render"})
	if len(contexts) != 6 {
		t.Fatalf("len(contexts) = %d, want 6", len(contexts))
	}
	for index, got := range contexts {
		if got != ctx {
			t.Fatalf("contexts[%d] did not match input context", index)
		}
	}
}

func TestRunValidateCancelsBeforeExecutionStarts(t *testing.T) {
	t.Parallel()

	order := []string{}
	runner := newFakeRunner(&order)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := runner.RunValidate(ctx, domain.CommandRequest{
		CommandName: domain.CommandNameValidate,
		ConfigPath:  fixtureConfigPath(),
	})

	_ = assertFailureCategory(t, err, ErrExecutionCanceled, domain.FailureCategoryExecution, "command canceled")
	if len(order) != 0 {
		t.Fatalf("order = %#v, want no stage calls", order)
	}
}

func TestRunReportGenerateCancelsAtIntermediateStage(t *testing.T) {
	t.Parallel()

	order := []string{}
	ctx, cancel := context.WithCancel(context.Background())
	runner := Runner{
		ConfigLoader:   &fakeConfigLoader{recorder: &order},
		ModelValidator: &fakeModelValidator{recorder: &order},
		CriteriaWeighter: &fakeCriteriaWeighter{
			recorder: &order,
			onCall:   cancel,
		},
		ScenarioRanker:     &fakeScenarioRanker{recorder: &order},
		ScenarioAggregator: &fakeScenarioAggregator{recorder: &order},
		ReportRenderer:     &fakeReportRenderer{recorder: &order},
	}

	_, err := runner.RunReportGenerate(ctx, domain.CommandRequest{
		CommandName: domain.CommandNameReportGenerate,
		ConfigPath:  fixtureConfigPath(),
	})

	_ = assertFailureCategory(t, err, ErrExecutionCanceled, domain.FailureCategoryExecution, "command canceled")
	assertStageOrder(t, order, []string{"load", "validate", "weight"})
}
