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
	return RenderReportsOutput{}, nil
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
