package pipeline

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func TestRunValidateStageOrdering(t *testing.T) {
	t.Parallel()

	order := []string{}
	runner := newFakeRunner(&order)

	_, err := runner.RunValidate(domain.CommandRequest{
		CommandName: domain.CommandNameValidate,
		ConfigPath:  filepath.Join("testdata", "config", "minimal.cue"),
	})
	if err != nil {
		t.Fatalf("RunValidate() error = %v", err)
	}

	got := order
	want := []string{"load", "validate"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("order = %#v, want %#v", got, want)
	}
}

func TestRunReportGenerateStageOrdering(t *testing.T) {
	t.Parallel()

	order := []string{}
	runner := newFakeRunner(&order)

	_, err := runner.RunReportGenerate(domain.CommandRequest{
		CommandName: domain.CommandNameReportGenerate,
		ConfigPath:  filepath.Join("testdata", "config", "minimal.cue"),
	})
	if err != nil {
		t.Fatalf("RunReportGenerate() error = %v", err)
	}

	got := order
	want := []string{"load", "validate", "weight", "rank", "aggregate", "render"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("order = %#v, want %#v", got, want)
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

	_, err := runner.RunReportGenerate(domain.CommandRequest{
		CommandName: domain.CommandNameReportGenerate,
		ConfigPath:  filepath.Join("testdata", "config", "minimal.cue"),
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}

	failure := domain.AsCommandFailure(err)
	if failure == nil {
		t.Fatal("AsCommandFailure(err) = nil, want value")
	}

	if failure.Category != domain.FailureCategoryExecution {
		t.Fatalf("Category = %q, want %q", failure.Category, domain.FailureCategoryExecution)
	}

	if got, want := order, []string{"load", "validate", "weight"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("order = %#v, want %#v", got, want)
	}
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
		CriterionWeights: []CriterionWeight{
			{CriterionName: "cost", Weight: 0.6},
			{CriterionName: "speed", Weight: 0.4},
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

	if got, want := weightOutput.CriterionWeights[0].CriterionName, "cost"; got != want {
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
	first, err := newFakeRunner(&firstOrder).RunReportGenerate(command)
	if err != nil {
		t.Fatalf("first RunReportGenerate() error = %v", err)
	}

	secondOrder := []string{}
	second, err := newFakeRunner(&secondOrder).RunReportGenerate(command)
	if err != nil {
		t.Fatalf("second RunReportGenerate() error = %v", err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("first result = %#v, second result = %#v", first, second)
	}
}

func TestFixtureDrivenFlowsUseConfigPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		command domain.CommandRequest
		run     func(Runner, domain.CommandRequest) (domain.CommandResult, error)
	}{
		{
			name: "validate",
			command: domain.CommandRequest{
				CommandName: domain.CommandNameValidate,
				ConfigPath:  fixtureConfigPath(),
			},
			run: Runner.RunValidate,
		},
		{
			name: "report generate",
			command: domain.CommandRequest{
				CommandName: domain.CommandNameReportGenerate,
				ConfigPath:  fixtureConfigPath(),
			},
			run: Runner.RunReportGenerate,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.run(NewDefaultRunner(), tt.command)
			if err != nil {
				t.Fatalf("run() error = %v", err)
			}

			assertValidatedModelPath(
				t,
				got,
				filepath.Clean(fixtureConfigPath()),
			)
		})
	}
}

func TestInvalidCueFailsAtLoadStage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		command domain.CommandRequest
		run     func(Runner, domain.CommandRequest) (domain.CommandResult, error)
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

			order := []string{}
			runner := newFakeRunner(&order)
			runner.ConfigLoader = DefaultConfigLoader{}

			_, err := tt.run(runner, tt.command)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}

			if len(order) != 0 {
				t.Fatalf("order = %#v, want no downstream stage calls", order)
			}
		})
	}
}

func TestValidationFailureStopsPipeline(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		command domain.CommandRequest
		run     func(Runner, domain.CommandRequest) (domain.CommandResult, error)
	}{
		{
			name: "validate stops at validation",
			command: domain.CommandRequest{
				CommandName: domain.CommandNameValidate,
				ConfigPath:  fixtureConfigPath(),
			},
			run: Runner.RunValidate,
		},
		{
			name: "report generate stops at validation",
			command: domain.CommandRequest{
				CommandName: domain.CommandNameReportGenerate,
				ConfigPath:  fixtureConfigPath(),
			},
			run: Runner.RunReportGenerate,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			order := []string{}
			runner := newFakeRunner(&order)
			runner.ConfigLoader = &fakeConfigLoader{
				recorder: &order,
				output: LoadConfigOutput{
					Config: LoadedConfig{
						Path:           fixtureConfigPath(),
						TopLevelFields: []string{"config"},
						ConfigFields: []string{
							"aggregation",
							"alternatives",
							"criteriaCatalog",
							"evaluations",
							"problem",
							"reports",
							"scenarios",
						},
						Config: &ExecutionConfig{
							Problem:         &ProblemConfig{Name: "minimal"},
							Reports:         []ReportConfig{{Name: "summary"}},
							CriteriaCatalog: []CriterionConfig{{Name: "cost"}},
							Alternatives:    []AlternativeConfig{{Name: "option_a"}},
							Scenarios:       []ScenarioConfig{{Name: "baseline", ActiveCriteria: []ScenarioCriterionRef{{CriterionName: "missing"}}}},
							Evaluations:     []EvaluationConfig{{ScenarioName: "baseline", Evaluations: []AlternativeEvaluationConfig{{AlternativeName: "option_a"}}}},
							Aggregation:     &AggregationConfig{},
						},
					},
				},
			}
			runner.ModelValidator = recordingModelValidator{
				recorder: &order,
				inner:    DefaultModelValidator{},
			}

			_, err := tt.run(runner, tt.command)
			if !errors.Is(err, ErrValidationFailed) {
				t.Fatalf("error = %v, want %v", err, ErrValidationFailed)
			}

			if got, want := order, []string{"load", "validate"}; !reflect.DeepEqual(got, want) {
				t.Fatalf("order = %#v, want %#v", got, want)
			}
		})
	}
}

func TestRealValidationFailureFromFixture(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		command domain.CommandRequest
		run     func(Runner, domain.CommandRequest) (domain.CommandResult, error)
	}{
		{
			name: "validate invalid reference fixture",
			command: domain.CommandRequest{
				CommandName: domain.CommandNameValidate,
				ConfigPath:  filepath.Join("..", "..", "testdata", "config", "invalid_reference.cue"),
			},
			run: Runner.RunValidate,
		},
		{
			name: "report generate invalid reference fixture",
			command: domain.CommandRequest{
				CommandName: domain.CommandNameReportGenerate,
				ConfigPath:  filepath.Join("..", "..", "testdata", "config", "invalid_reference.cue"),
			},
			run: Runner.RunReportGenerate,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := tt.run(NewDefaultRunner(), tt.command)
			if !errors.Is(err, ErrValidationFailed) {
				t.Fatalf("error = %v, want %v", err, ErrValidationFailed)
			}

			failure := domain.AsCommandFailure(err)
			if failure == nil {
				t.Fatal("AsCommandFailure(err) = nil, want value")
			}

			if got, want := failure.Diagnostics[0].Message, "unknown scenario name in evaluations: missing"; got != want {
				t.Fatalf("message = %q, want %q", got, want)
			}
		})
	}
}

func TestDefaultConfigLoader(t *testing.T) {
	t.Parallel()

	loader := DefaultConfigLoader{}

	got, err := loader.LoadConfig(LoadConfigInput{
		ConfigPath: fixtureConfigPath(),
	})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	want := filepath.Clean(fixtureConfigPath())
	if got.Config.Path != want {
		t.Fatalf("ConfigPath = %q, want %q", got.Config.Path, want)
	}
}

func TestDefaultConfigLoaderMissingFile(t *testing.T) {
	t.Parallel()

	loader := DefaultConfigLoader{}

	_, err := loader.LoadConfig(LoadConfigInput{
		ConfigPath: filepath.Join("..", "..", "testdata", "config", "missing.cue"),
	})
	if !errors.Is(err, ErrConfigPathDoesNotExist) {
		t.Fatalf("error = %v, want %v", err, ErrConfigPathDoesNotExist)
	}

	failure := domain.AsCommandFailure(err)
	if failure == nil {
		t.Fatal("AsCommandFailure(err) = nil, want value")
	}

	if failure.Category != domain.FailureCategoryInput {
		t.Fatalf("Category = %q, want %q", failure.Category, domain.FailureCategoryInput)
	}
}

func TestDefaultConfigLoaderDirectoryPath(t *testing.T) {
	t.Parallel()

	loader := DefaultConfigLoader{}

	_, err := loader.LoadConfig(LoadConfigInput{
		ConfigPath: filepath.Join("..", "..", "testdata", "config"),
	})
	if !errors.Is(err, ErrConfigPathIsDirectory) {
		t.Fatalf("error = %v, want %v", err, ErrConfigPathIsDirectory)
	}

	failure := domain.AsCommandFailure(err)
	if failure == nil {
		t.Fatal("AsCommandFailure(err) = nil, want value")
	}

	if failure.Category != domain.FailureCategoryInput {
		t.Fatalf("Category = %q, want %q", failure.Category, domain.FailureCategoryInput)
	}
}
