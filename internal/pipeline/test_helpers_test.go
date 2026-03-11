package pipeline

import (
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

type fakeConfigLoader struct {
	recorder *[]string
	err      error
	output   LoadConfigOutput
}

func (f *fakeConfigLoader) LoadConfig(input LoadConfigInput) (LoadConfigOutput, error) {
	*f.recorder = append(*f.recorder, "load")
	if f.err != nil {
		return LoadConfigOutput{}, f.err
	}
	if f.output.Config.Path != "" || f.output.Config.Config != nil {
		return f.output, nil
	}
	return LoadConfigOutput{
		Config: LoadedConfig{
			Path: input.ConfigPath,
		},
	}, nil
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
			ConfigPath: input.Config.Path,
		},
	}, nil
}

type recordingModelValidator struct {
	recorder *[]string
	inner    ModelValidator
}

func (r recordingModelValidator) ValidateModel(input ValidateModelInput) (ValidateModelOutput, error) {
	*r.recorder = append(*r.recorder, "validate")
	return r.inner.ValidateModel(input)
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

func newFakeRunner(order *[]string) Runner {
	return Runner{
		ConfigLoader:       &fakeConfigLoader{recorder: order},
		ModelValidator:     &fakeModelValidator{recorder: order},
		CriteriaWeighter:   &fakeCriteriaWeighter{recorder: order},
		ScenarioRanker:     &fakeScenarioRanker{recorder: order},
		ScenarioAggregator: &fakeScenarioAggregator{recorder: order},
		ReportRenderer:     &fakeReportRenderer{recorder: order},
	}
}

func fixtureConfigPath() string {
	return filepath.Join("..", "..", "testdata", "config", "minimal.cue")
}

func assertValidatedModelPath(t *testing.T, result domain.CommandResult, wantPath string) {
	t.Helper()

	if result.ValidatedModel == nil {
		t.Fatal("ValidatedModel = nil, want value")
	}

	if result.ValidatedModel.ConfigPath != wantPath {
		t.Fatalf("ConfigPath = %q, want %q", result.ValidatedModel.ConfigPath, wantPath)
	}
}
