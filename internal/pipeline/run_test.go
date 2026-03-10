package pipeline

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

type fakeConfigLoader struct {
	recorder *[]string
	err      error
}

func (f *fakeConfigLoader) LoadConfig(input LoadConfigInput) (LoadConfigOutput, error) {
	*f.recorder = append(*f.recorder, "load")
	if f.err != nil {
		return LoadConfigOutput{}, f.err
	}
	return LoadConfigOutput{ConfigPath: input.ConfigPath}, nil
}

type fakeModelValidator struct {
	recorder *[]string
	err      error
}

func (f *fakeModelValidator) ValidateModel(input ValidateModelInput) (ValidateModelOutput, error) {
	*f.recorder = append(*f.recorder, "validate")
	if f.err != nil {
		return ValidateModelOutput{}, f.err
	}
	return ValidateModelOutput{
		ValidatedModel: domain.ValidatedModelSummary{
			ConfigPath: input.Config.ConfigPath,
		},
	}, nil
}

type fakeCriteriaWeighter struct {
	recorder *[]string
	err      error
}

func (f *fakeCriteriaWeighter) WeightCriteria(WeightCriteriaInput) (WeightCriteriaOutput, error) {
	*f.recorder = append(*f.recorder, "weight")
	if f.err != nil {
		return WeightCriteriaOutput{}, f.err
	}
	return WeightCriteriaOutput{}, nil
}

type fakeScenarioRanker struct {
	recorder *[]string
	err      error
}

func (f *fakeScenarioRanker) RankScenarios(RankScenariosInput) (RankScenariosOutput, error) {
	*f.recorder = append(*f.recorder, "rank")
	if f.err != nil {
		return RankScenariosOutput{}, f.err
	}
	return RankScenariosOutput{}, nil
}

type fakeScenarioAggregator struct {
	recorder *[]string
	err      error
}

func (f *fakeScenarioAggregator) AggregateScenarios(AggregateScenariosInput) (AggregateScenariosOutput, error) {
	*f.recorder = append(*f.recorder, "aggregate")
	if f.err != nil {
		return AggregateScenariosOutput{}, f.err
	}
	return AggregateScenariosOutput{}, nil
}

type fakeReportRenderer struct {
	recorder *[]string
	err      error
}

func (f *fakeReportRenderer) RenderReports(RenderReportsInput) (RenderReportsOutput, error) {
	*f.recorder = append(*f.recorder, "render")
	if f.err != nil {
		return RenderReportsOutput{}, f.err
	}
	return RenderReportsOutput{
		ReportDefinitions: []domain.ReportDefinition{
			{Name: "zeta", Title: "Zeta", Format: "json"},
			{Name: "alpha", Title: "Alpha", Format: "markdown"},
		},
	}, nil
}

func TestRunValidateStageOrdering(t *testing.T) {
	t.Parallel()

	order := []string{}
	runner := Runner{
		ConfigLoader:       &fakeConfigLoader{recorder: &order},
		ModelValidator:     &fakeModelValidator{recorder: &order},
		CriteriaWeighter:   &fakeCriteriaWeighter{recorder: &order},
		ScenarioRanker:     &fakeScenarioRanker{recorder: &order},
		ScenarioAggregator: &fakeScenarioAggregator{recorder: &order},
		ReportRenderer:     &fakeReportRenderer{recorder: &order},
	}

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
	runner := Runner{
		ConfigLoader:       &fakeConfigLoader{recorder: &order},
		ModelValidator:     &fakeModelValidator{recorder: &order},
		CriteriaWeighter:   &fakeCriteriaWeighter{recorder: &order},
		ScenarioRanker:     &fakeScenarioRanker{recorder: &order},
		ScenarioAggregator: &fakeScenarioAggregator{recorder: &order},
		ReportRenderer:     &fakeReportRenderer{recorder: &order},
	}

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
		ConfigPath:  filepath.Join("..", "..", "testdata", "config", "minimal.cue"),
	}
	loadOutput := LoadConfigOutput{
		ConfigPath: filepath.Join("..", "..", "testdata", "config", "minimal.cue"),
	}
	validateInput := ValidateModelInput{
		Command: command,
		Config:  loadOutput,
	}
	validateOutput := ValidateModelOutput{
		Diagnostics: []domain.Diagnostic{
			domain.NewDiagnostic(domain.DiagnosticSeverityWarning, "stub.warning", loadOutput.ConfigPath, domain.DiagnosticLocation{}, "warning"),
		},
		ValidatedModel: domain.ValidatedModelSummary{
			ConfigPath:       loadOutput.ConfigPath,
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

	if got, want := validateOutput.ValidatedModel.ConfigPath, loadOutput.ConfigPath; got != want {
		t.Fatalf("ValidatedModel.ConfigPath = %q, want %q", got, want)
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

	buildRunner := func() Runner {
		order := []string{}
		return Runner{
			ConfigLoader:       &fakeConfigLoader{recorder: &order},
			ModelValidator:     &fakeModelValidator{recorder: &order},
			CriteriaWeighter:   &fakeCriteriaWeighter{recorder: &order},
			ScenarioRanker:     &fakeScenarioRanker{recorder: &order},
			ScenarioAggregator: &fakeScenarioAggregator{recorder: &order},
			ReportRenderer:     &fakeReportRenderer{recorder: &order},
		}
	}

	command := domain.CommandRequest{
		CommandName: domain.CommandNameReportGenerate,
		ConfigPath:  filepath.Join("..", "..", "testdata", "config", "minimal.cue"),
	}

	first, err := buildRunner().RunReportGenerate(command)
	if err != nil {
		t.Fatalf("first RunReportGenerate() error = %v", err)
	}

	second, err := buildRunner().RunReportGenerate(command)
	if err != nil {
		t.Fatalf("second RunReportGenerate() error = %v", err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("first result = %#v, second result = %#v", first, second)
	}
}

func TestRunValidateUsesFixtureDrivenConfigPath(t *testing.T) {
	t.Parallel()

	runner := NewDefaultRunner()
	command := domain.CommandRequest{
		CommandName: domain.CommandNameValidate,
		ConfigPath:  filepath.Join("..", "..", "testdata", "config", "minimal.cue"),
	}

	got, err := runner.RunValidate(command)
	if err != nil {
		t.Fatalf("RunValidate() error = %v", err)
	}

	if got.ValidatedModel == nil {
		t.Fatal("ValidatedModel = nil, want value")
	}

	if want := filepath.Clean(filepath.Join("..", "..", "testdata", "config", "minimal.cue")); got.ValidatedModel.ConfigPath != want {
		t.Fatalf("ConfigPath = %q, want %q", got.ValidatedModel.ConfigPath, want)
	}
}

func TestRunReportGenerateUsesFixtureDrivenConfigPath(t *testing.T) {
	t.Parallel()

	runner := NewDefaultRunner()
	command := domain.CommandRequest{
		CommandName: domain.CommandNameReportGenerate,
		ConfigPath:  filepath.Join("..", "..", "testdata", "config", "minimal.cue"),
	}

	got, err := runner.RunReportGenerate(command)
	if err != nil {
		t.Fatalf("RunReportGenerate() error = %v", err)
	}

	if got.ValidatedModel == nil {
		t.Fatal("ValidatedModel = nil, want value")
	}

	if want := filepath.Clean(filepath.Join("..", "..", "testdata", "config", "minimal.cue")); got.ValidatedModel.ConfigPath != want {
		t.Fatalf("ConfigPath = %q, want %q", got.ValidatedModel.ConfigPath, want)
	}
}

func TestDefaultConfigLoader(t *testing.T) {
	t.Parallel()

	loader := DefaultConfigLoader{}

	got, err := loader.LoadConfig(LoadConfigInput{
		ConfigPath: filepath.Join("..", "..", "testdata", "config", "minimal.cue"),
	})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	want := filepath.Clean(filepath.Join("..", "..", "testdata", "config", "minimal.cue"))
	if got.ConfigPath != want {
		t.Fatalf("ConfigPath = %q, want %q", got.ConfigPath, want)
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
