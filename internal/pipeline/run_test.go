package pipeline

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
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

	got := order
	want := []string{"load", "validate"}
	assertStageOrder(t, got, want)
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

	got := order
	want := []string{"load", "validate", "weight", "rank", "aggregate", "render"}
	assertStageOrder(t, got, want)
}

func TestRunReportGenerateSelectsCurrentV1Strategy(t *testing.T) {
	t.Parallel()

	order := []string{}
	strategy := fakeRankingStrategy{
		recorder: &order,
		output: RankingStrategyOutput{
			ScenarioWeights: []ScenarioCriterionWeights{
				{
					ScenarioName: "baseline",
					CriterionWeights: []CriterionWeight{
						{CriterionName: "cost", Weight: 1},
					},
				},
			},
			ScenarioResults: []domain.ScenarioRankingResult{
				{
					ScenarioName: "baseline",
					RankedAlternatives: []domain.RankedAlternative{
						{Name: "option_a", Rank: 1, Score: 0.9},
					},
				},
			},
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
	weight := &fakeCriteriaWeighter{recorder: &order, err: wantErr}
	rank := &fakeScenarioRanker{recorder: &order}
	aggregate := &fakeScenarioAggregator{recorder: &order}
	render := &fakeReportRenderer{recorder: &order}
	runner := Runner{
		ConfigLoader:       &fakeConfigLoader{recorder: &order},
		ModelValidator:     &fakeModelValidator{recorder: &order},
		CriteriaWeighter:   weight,
		ScenarioRanker:     rank,
		ScenarioAggregator: aggregate,
		ReportRenderer:     render,
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

func TestStageIOContractsCanBeConstructed(t *testing.T) {
	t.Parallel()

	command := domain.CommandRequest{
		CommandName: domain.CommandNameReportGenerate,
		ConfigPath:  fixtureConfigPath(),
	}
	loadOutput := LoadConfigOutput{
		Config: LoadedConfig{
			Path:           fixtureConfigPath(),
			Evaluated:      "config: {\n\tname: \"minimal\"\n}\n",
			TopLevelFields: []string{"config"},
		},
	}
	validateInput := ValidateModelInput{
		Command: command,
		Config:  loadOutput.Config,
	}
	validateOutput := ValidateModelOutput{
		Diagnostics: []domain.Diagnostic{
			domain.NewDiagnostic(domain.DiagnosticSeverityWarning, "stub.warning", loadOutput.Config.Path, domain.DiagnosticLocation{}, "warning"),
		},
		ValidatedModel: domain.ValidatedModelSummary{
			ConfigPath:       loadOutput.Config.Path,
			CriterionCount:   3,
			AlternativeCount: 2,
			ScenarioCount:    1,
		},
		ReportDefinitions: []domain.ReportDefinition{
			{Name: "summary", Title: "Summary", Format: "markdown"},
		},
	}
	weightOutput := WeightCriteriaOutput{
		ScenarioWeights: []ScenarioCriterionWeights{
			{
				ScenarioName: "startup",
				CriterionWeights: []CriterionWeight{
					{CriterionName: "cost", Weight: 0.6},
					{CriterionName: "speed", Weight: 0.4},
				},
			},
		},
	}
	rankOutput := RankScenariosOutput{
		ScenarioResults: []domain.ScenarioRankingResult{
			{
				ScenarioName: "startup",
				RankedAlternatives: []domain.RankedAlternative{
					{Name: "platform-a", Rank: 1, Score: 0.9},
				},
			},
		},
	}
	aggregateOutput := AggregateScenariosOutput{
		FinalRanking: domain.AggregatedRankingResult{
			RankedAlternatives: []domain.RankedAlternative{
				{Name: "platform-a", Rank: 1, Score: 0.9},
			},
		},
	}
	renderInput := RenderReportsInput{
		Command:           command,
		ValidatedModel:    validateOutput.ValidatedModel,
		ScenarioResults:   rankOutput.ScenarioResults,
		FinalRanking:      aggregateOutput.FinalRanking,
		ReportDefinitions: validateOutput.ReportDefinitions,
	}

	if validateInput.Command.CommandName != domain.CommandNameReportGenerate {
		t.Fatalf("CommandName = %q, want %q", validateInput.Command.CommandName, domain.CommandNameReportGenerate)
	}

	if got, want := validateOutput.ValidatedModel.ConfigPath, loadOutput.Config.Path; got != want {
		t.Fatalf("ValidatedModel.ConfigPath = %q, want %q", got, want)
	}

	if got, want := validateInput.Config.TopLevelFields[0], "config"; got != want {
		t.Fatalf("TopLevelFields[0] = %q, want %q", got, want)
	}

	if got, want := weightOutput.ScenarioWeights[0].CriterionWeights[0].CriterionName, "cost"; got != want {
		t.Fatalf("CriterionName = %q, want %q", got, want)
	}

	if got, want := rankOutput.ScenarioResults[0].ScenarioName, "startup"; got != want {
		t.Fatalf("ScenarioName = %q, want %q", got, want)
	}

	if got, want := renderInput.ReportDefinitions[0].Name, "summary"; got != want {
		t.Fatalf("ReportDefinitions[0].Name = %q, want %q", got, want)
	}
}

func TestRunReportGenerateIsDeterministic(t *testing.T) {
	t.Parallel()

	command := domain.CommandRequest{
		CommandName: domain.CommandNameReportGenerate,
		ConfigPath:  filepath.Join("..", "..", "testdata", "config", "minimal.cue"),
	}

	firstOrder := []string{}
	first, err := newFakeRunner(&firstOrder).RunReportGenerate(context.Background(), command)
	if err != nil {
		t.Fatalf("first RunReportGenerate() error = %v", err)
	}

	secondOrder := []string{}
	second, err := newFakeRunner(&secondOrder).RunReportGenerate(context.Background(), command)
	if err != nil {
		t.Fatalf("second RunReportGenerate() error = %v", err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("first result = %#v, second result = %#v", first, second)
	}
}

func TestFixtureDrivenFlowsUseConfigPath(t *testing.T) {
	t.Parallel()

	for _, tt := range fixtureFlowCases(fixtureConfigPath()) {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assertRunnerUsesConfigPath(t, tt.run, tt.command)
		})
	}
}

func TestInvalidCueFailsAtLoadStage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		command domain.CommandRequest
		run     func(Runner, context.Context, domain.CommandRequest) (domain.CommandResult, error)
		wantErr error
	}{
		{
			name: "validate invalid cue",
			command: domain.CommandRequest{
				CommandName: domain.CommandNameValidate,
				ConfigPath:  filepath.Join("..", "..", "testdata", "config", "non_concrete.cue"),
			},
			run:     Runner.RunValidate,
			wantErr: ErrConfigNotConcrete,
		},
		{
			name: "report generate invalid cue",
			command: domain.CommandRequest{
				CommandName: domain.CommandNameReportGenerate,
				ConfigPath:  filepath.Join("..", "..", "testdata", "config", "malformed.cue"),
			},
			run:     Runner.RunReportGenerate,
			wantErr: ErrConfigLoadInvalid,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assertLoadStageFailureStopsPipeline(t, tt.run, tt.command, tt.wantErr)
		})
	}
}

func TestValidationFailureStopsPipeline(t *testing.T) {
	t.Parallel()

	for _, tt := range fixtureFlowCases(fixtureConfigPath()) {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assertValidationStageFailureStopsPipeline(t, tt.run, tt.command)
		})
	}
}

func invalidValidationLoadOutput() LoadConfigOutput {
	config := validLoadedConfig()
	config.Config.Reports = []ReportConfig{{Name: "summary"}}
	config.Config.Scenarios = []ScenarioConfig{{
		Name:           "baseline",
		ActiveCriteria: []ScenarioCriterionRef{{CriterionName: "missing"}},
	}}
	config.Config.Evaluations = scenarioEvaluationBlock(
		"baseline",
		alternativeEvaluation("option_a", map[string]CriterionValue{"cost": {Kind: "number", Value: 1}}),
	)

	return LoadConfigOutput{Config: config}
}

func fixtureFlowCases(configPath string) []struct {
	name    string
	command domain.CommandRequest
	run     func(Runner, context.Context, domain.CommandRequest) (domain.CommandResult, error)
} {
	return []struct {
		name    string
		command domain.CommandRequest
		run     func(Runner, context.Context, domain.CommandRequest) (domain.CommandResult, error)
	}{
		{
			name: "validate",
			command: domain.CommandRequest{
				CommandName: domain.CommandNameValidate,
				ConfigPath:  configPath,
			},
			run: Runner.RunValidate,
		},
		{
			name: "report generate",
			command: domain.CommandRequest{
				CommandName: domain.CommandNameReportGenerate,
				ConfigPath:  configPath,
			},
			run: Runner.RunReportGenerate,
		},
	}
}

func assertRunnerUsesConfigPath(
	t *testing.T,
	run func(Runner, context.Context, domain.CommandRequest) (domain.CommandResult, error),
	command domain.CommandRequest,
) {
	t.Helper()

	got, err := run(NewDefaultRunner(), context.Background(), command)
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	assertValidatedModelPath(t, got, filepath.Clean(fixtureConfigPath()))
}

func assertLoadStageFailureStopsPipeline(
	t *testing.T,
	run func(Runner, context.Context, domain.CommandRequest) (domain.CommandResult, error),
	command domain.CommandRequest,
	wantErr error,
) {
	t.Helper()

	order := []string{}
	runner := newFakeRunner(&order)
	runner.ConfigLoader = DefaultConfigLoader{}

	_, err := run(runner, context.Background(), command)
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
	if len(order) != 0 {
		t.Fatalf("order = %#v, want no downstream stage calls", order)
	}
}

func assertValidationStageFailureStopsPipeline(
	t *testing.T,
	run func(Runner, context.Context, domain.CommandRequest) (domain.CommandResult, error),
	command domain.CommandRequest,
) {
	t.Helper()

	order := []string{}
	runner := newFakeRunner(&order)
	runner.ConfigLoader = &fakeConfigLoader{
		recorder: &order,
		output:   invalidValidationLoadOutput(),
	}
	runner.ModelValidator = recordingModelValidator{
		recorder: &order,
		inner:    DefaultModelValidator{},
	}

	_, err := run(runner, context.Background(), command)
	if !errors.Is(err, ErrValidationFailed) {
		t.Fatalf("error = %v, want %v", err, ErrValidationFailed)
	}
	assertStageOrder(t, order, []string{"load", "validate"})
}

func TestRealValidationFailureFromFixture(t *testing.T) {
	t.Parallel()

	runFlowBehaviorTests(t, []flowBehaviorCase{
		{
			name: "validate invalid reference fixture",
			command: domain.CommandRequest{
				CommandName: domain.CommandNameValidate,
				ConfigPath:  filepath.Join("..", "..", "testdata", "config", "invalid_reference.cue"),
			},
			run:         Runner.RunValidate,
			wantErr:     ErrValidationFailed,
			wantMessage: "unknown scenario name in evaluations: missing",
		},
		{
			name: "report generate invalid reference fixture",
			command: domain.CommandRequest{
				CommandName: domain.CommandNameReportGenerate,
				ConfigPath:  filepath.Join("..", "..", "testdata", "config", "invalid_reference.cue"),
			},
			run:         Runner.RunReportGenerate,
			wantErr:     ErrValidationFailed,
			wantMessage: "unknown scenario name in evaluations: missing",
		},
	})
}

func TestPairwiseValidationFlowBehavior(t *testing.T) {
	t.Parallel()

	runFlowBehaviorTests(t, []flowBehaviorCase{
		{
			name: "validate stops on pairwise validation failure",
			command: domain.CommandRequest{
				CommandName: domain.CommandNameValidate,
				ConfigPath:  filepath.Join("..", "..", "testdata", "config", "pairwise_missing.cue"),
			},
			run:         Runner.RunValidate,
			wantErr:     ErrValidationFailed,
			wantMessage: "missing pairwise comparison for pair: reliability/speed",
		},
		{
			name: "report generate stops on pairwise validation failure",
			command: domain.CommandRequest{
				CommandName: domain.CommandNameReportGenerate,
				ConfigPath:  filepath.Join("..", "..", "testdata", "config", "pairwise_missing.cue"),
			},
			run:         Runner.RunReportGenerate,
			wantErr:     ErrValidationFailed,
			wantMessage: "missing pairwise comparison for pair: reliability/speed",
		},
		{
			name: "valid pairwise config proceeds",
			command: domain.CommandRequest{
				CommandName: domain.CommandNameReportGenerate,
				ConfigPath:  filepath.Join("..", "..", "testdata", "config", "pairwise_valid.cue"),
			},
			run: Runner.RunReportGenerate,
		},
	})
}

func TestEvaluationValidationFlowBehavior(t *testing.T) {
	t.Parallel()

	runFlowBehaviorTests(t, []flowBehaviorCase{
		{
			name: "validate stops on evaluation validation failure",
			command: domain.CommandRequest{
				CommandName: domain.CommandNameValidate,
				ConfigPath:  filepath.Join("..", "..", "testdata", "config", "invalid_evaluation.cue"),
			},
			run:         Runner.RunValidate,
			wantErr:     ErrValidationFailed,
			wantMessage: "evaluation value kind mismatch for criterion cost: want number, got boolean",
		},
		{
			name: "report generate stops on evaluation validation failure",
			command: domain.CommandRequest{
				CommandName: domain.CommandNameReportGenerate,
				ConfigPath:  filepath.Join("..", "..", "testdata", "config", "invalid_evaluation.cue"),
			},
			run:         Runner.RunReportGenerate,
			wantErr:     ErrValidationFailed,
			wantMessage: "evaluation value kind mismatch for criterion cost: want number, got boolean",
		},
		{
			name: "valid evaluation input proceeds",
			command: domain.CommandRequest{
				CommandName: domain.CommandNameValidate,
				ConfigPath:  filepath.Join("..", "..", "testdata", "config", "minimal.cue"),
			},
			run: Runner.RunValidate,
		},
	})
}

func TestConstraintValidationFlowBehavior(t *testing.T) {
	t.Parallel()

	runFlowBehaviorTests(t, []flowBehaviorCase{
		{
			name: "validate stops on constraint validation failure",
			command: domain.CommandRequest{
				CommandName: domain.CommandNameValidate,
				ConfigPath:  filepath.Join("..", "..", "testdata", "config", "invalid_constraint.cue"),
			},
			run:         Runner.RunValidate,
			wantErr:     ErrValidationFailed,
			wantMessage: "invalid constraint operator for boolean criterion approved: <=",
		},
		{
			name: "report generate stops on constraint validation failure",
			command: domain.CommandRequest{
				CommandName: domain.CommandNameReportGenerate,
				ConfigPath:  filepath.Join("..", "..", "testdata", "config", "invalid_constraint.cue"),
			},
			run:         Runner.RunReportGenerate,
			wantErr:     ErrValidationFailed,
			wantMessage: "invalid constraint operator for boolean criterion approved: <=",
		},
		{
			name: "valid constraint input proceeds",
			command: domain.CommandRequest{
				CommandName: domain.CommandNameValidate,
				ConfigPath:  filepath.Join("..", "..", "testdata", "config", "valid_constraint.cue"),
			},
			run: Runner.RunValidate,
		},
	})
}

func TestReportDefinitionValidationFlowBehavior(t *testing.T) {
	t.Parallel()

	runFlowBehaviorTests(t, []flowBehaviorCase{
		{
			name: "validate stops on report definition validation failure",
			command: domain.CommandRequest{
				CommandName: domain.CommandNameValidate,
				ConfigPath:  filepath.Join("..", "..", "testdata", "config", "invalid_report.cue"),
			},
			run:         Runner.RunValidate,
			wantErr:     ErrValidationFailed,
			wantMessage: "report argument key header is not allowed for format json",
		},
		{
			name: "report generate stops on report definition validation failure",
			command: domain.CommandRequest{
				CommandName: domain.CommandNameReportGenerate,
				ConfigPath:  filepath.Join("..", "..", "testdata", "config", "invalid_report.cue"),
			},
			run:         Runner.RunReportGenerate,
			wantErr:     ErrValidationFailed,
			wantMessage: "report argument key header is not allowed for format json",
		},
		{
			name: "valid report definitions proceed",
			command: domain.CommandRequest{
				CommandName: domain.CommandNameValidate,
				ConfigPath:  filepath.Join("..", "..", "testdata", "config", "valid_report.cue"),
			},
			run: Runner.RunValidate,
		},
	})
}

func TestRunReportGenerateProducesRealRenderedReport(t *testing.T) {
	t.Parallel()

	got, err := NewDefaultRunner().RunReportGenerate(context.Background(), domain.CommandRequest{
		CommandName: domain.CommandNameReportGenerate,
		ConfigPath:  fixtureConfigPath(),
	})
	if err != nil {
		t.Fatalf("RunReportGenerate() error = %v", err)
	}

	if got, want := got.RenderedOutput, readPipelineGolden(t, "report_generate_success.stdout.golden"); got != want {
		t.Fatalf("RenderedOutput = %q, want %q", got, want)
	}
}

func TestRunValidateRepeatedRunDeterminism(t *testing.T) {
	t.Parallel()

	command := domain.CommandRequest{
		CommandName: domain.CommandNameValidate,
		ConfigPath:  fixtureConfigPath(),
	}

	assertRepeatedDeepEqual(t, 1, func() (domain.CommandResult, error) {
		return NewDefaultRunner().RunValidate(context.Background(), command)
	})
}

func TestRunReportGenerateRepeatedRunDeterminism(t *testing.T) {
	t.Parallel()

	command := domain.CommandRequest{
		CommandName: domain.CommandNameReportGenerate,
		ConfigPath:  filepath.Join("..", "..", "testdata", "config", "valid_report.cue"),
	}

	assertRepeatedDeepEqual(t, 1, func() (domain.CommandResult, error) {
		return NewDefaultRunner().RunReportGenerate(context.Background(), command)
	})
}

func TestRunValidateValidationFailureDeterminism(t *testing.T) {
	t.Parallel()

	command := domain.CommandRequest{
		CommandName: domain.CommandNameValidate,
		ConfigPath:  filepath.Join("..", "..", "testdata", "config", "invalid_reference.cue"),
	}

	_, firstErr := NewDefaultRunner().RunValidate(context.Background(), command)
	_, secondErr := NewDefaultRunner().RunValidate(context.Background(), command)

	first := domain.PresentError(firstErr)
	second := domain.PresentError(secondErr)

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("first = %#v, second = %#v", first, second)
	}
}

func TestRunReportGenerateStopsOnAggregationFailure(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("aggregate failed")

	assertReportGenerateStageFailure(t, wantErr, domain.FailureCategoryExecution, []string{"load", "validate", "weight", "rank", "aggregate"}, func(order *[]string) Runner {
		return Runner{
			ConfigLoader:       &fakeConfigLoader{recorder: order},
			ModelValidator:     &fakeModelValidator{recorder: order},
			CriteriaWeighter:   &fakeCriteriaWeighter{recorder: order},
			ScenarioRanker:     &fakeScenarioRanker{recorder: order},
			ScenarioAggregator: &fakeScenarioAggregator{recorder: order, err: wantErr},
			ReportRenderer:     &fakeReportRenderer{recorder: order},
		}
	})
}

func TestRunReportGenerateStopsOnRenderingFailure(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("render failed")
	assertReportGenerateStageFailure(t, wantErr, domain.FailureCategoryRendering, []string{"load", "validate", "weight", "rank", "aggregate", "render"}, func(order *[]string) Runner {
		return Runner{
			ConfigLoader:       &fakeConfigLoader{recorder: order},
			ModelValidator:     &fakeModelValidator{recorder: order},
			CriteriaWeighter:   &fakeCriteriaWeighter{recorder: order},
			ScenarioRanker:     &fakeScenarioRanker{recorder: order},
			ScenarioAggregator: &fakeScenarioAggregator{recorder: order},
			ReportRenderer:     &fakeReportRenderer{recorder: order, err: wantErr},
		}
	})
}

func TestDefaultConfigLoader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		configPath string
	}{
		{
			name:       "single file",
			configPath: fixtureConfigPath(),
		},
		{
			name:       "directory package",
			configPath: filepath.Join("..", "..", "testdata", "config_split"),
		},
	}

	loader := DefaultConfigLoader{}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := mustLoadConfig(t, loader, tt.configPath)
			if got.Config.Path != filepath.Clean(tt.configPath) {
				t.Fatalf("ConfigPath = %q, want %q", got.Config.Path, filepath.Clean(tt.configPath))
			}
		})
	}
}

func TestDefaultConfigLoaderMissingFile(t *testing.T) {
	t.Parallel()

	loader := DefaultConfigLoader{}

	_, err := loader.LoadConfig(context.Background(), LoadConfigInput{
		ConfigPath: filepath.Join("..", "..", "testdata", "config", "missing.cue"),
	})
	assertLoaderFailure(t, err, ErrConfigPathDoesNotExist, domain.FailureCategoryInput, "config.not_found", "config path does not exist")
}

func TestDefaultConfigLoaderEmptyDirectoryPath(t *testing.T) {
	t.Parallel()

	loader := DefaultConfigLoader{}

	_, err := loader.LoadConfig(context.Background(), LoadConfigInput{
		ConfigPath: filepath.Join("..", "..", "testdata", "config_empty"),
	})
	assertLoaderFailure(t, err, ErrConfigDirectoryEmpty, domain.FailureCategoryInput, "config.directory_empty", "config directory does not contain any .cue files")
}
