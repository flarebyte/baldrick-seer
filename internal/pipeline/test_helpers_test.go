package pipeline

import (
	"context"
	"errors"
	"math"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

type fakeConfigLoader struct {
	recorder *[]string
	err      error
	output   LoadConfigOutput
	contexts *[]context.Context
}

func (f *fakeConfigLoader) LoadConfig(ctx context.Context, input LoadConfigInput) (LoadConfigOutput, error) {
	*f.recorder = append(*f.recorder, "load")
	if f.contexts != nil {
		*f.contexts = append(*f.contexts, ctx)
	}
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
	contexts *[]context.Context
}

func (f *fakeModelValidator) ValidateModel(ctx context.Context, input ValidateModelInput) (ValidateModelOutput, error) {
	*f.recorder = append(*f.recorder, "validate")
	if f.contexts != nil {
		*f.contexts = append(*f.contexts, ctx)
	}
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

func (r recordingModelValidator) ValidateModel(ctx context.Context, input ValidateModelInput) (ValidateModelOutput, error) {
	*r.recorder = append(*r.recorder, "validate")
	return r.inner.ValidateModel(ctx, input)
}

type fakeCriteriaWeighter struct {
	recorder *[]string
	err      error
	contexts *[]context.Context
	onCall   func()
}

func (f *fakeCriteriaWeighter) WeightCriteria(ctx context.Context, input WeightCriteriaInput) (WeightCriteriaOutput, error) {
	*f.recorder = append(*f.recorder, "weight")
	if f.contexts != nil {
		*f.contexts = append(*f.contexts, ctx)
	}
	if f.onCall != nil {
		f.onCall()
	}
	if f.err != nil {
		return WeightCriteriaOutput{}, f.err
	}
	return WeightCriteriaOutput{}, nil
}

type fakeScenarioRanker struct {
	recorder *[]string
	err      error
	contexts *[]context.Context
}

func (f *fakeScenarioRanker) RankScenarios(ctx context.Context, input RankScenariosInput) (RankScenariosOutput, error) {
	*f.recorder = append(*f.recorder, "rank")
	if f.contexts != nil {
		*f.contexts = append(*f.contexts, ctx)
	}
	if f.err != nil {
		return RankScenariosOutput{}, f.err
	}
	return RankScenariosOutput{}, nil
}

type fakeScenarioAggregator struct {
	recorder *[]string
	err      error
	contexts *[]context.Context
}

func (f *fakeScenarioAggregator) AggregateScenarios(ctx context.Context, input AggregateScenariosInput) (AggregateScenariosOutput, error) {
	*f.recorder = append(*f.recorder, "aggregate")
	if f.contexts != nil {
		*f.contexts = append(*f.contexts, ctx)
	}
	if f.err != nil {
		return AggregateScenariosOutput{}, f.err
	}
	return AggregateScenariosOutput{}, nil
}

type fakeReportRenderer struct {
	recorder *[]string
	err      error
	contexts *[]context.Context
}

func (f *fakeReportRenderer) RenderReports(ctx context.Context, input RenderReportsInput) (RenderReportsOutput, error) {
	*f.recorder = append(*f.recorder, "render")
	if f.contexts != nil {
		*f.contexts = append(*f.contexts, ctx)
	}
	if f.err != nil {
		return RenderReportsOutput{}, f.err
	}
	return RenderReportsOutput{
		ReportDefinitions: []domain.ReportDefinition{
			{Name: "zeta", Title: "Zeta", Format: "json"},
			{Name: "alpha", Title: "Alpha", Format: "markdown"},
		},
		RenderedOutput: "rendered output\n",
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

type flowBehaviorCase struct {
	name        string
	command     domain.CommandRequest
	run         func(Runner, context.Context, domain.CommandRequest) (domain.CommandResult, error)
	wantErr     error
	wantMessage string
}

func runFlowBehaviorTests(t *testing.T, tests []flowBehaviorCase) {
	t.Helper()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.run(NewDefaultRunner(), context.Background(), tt.command)
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("run() error = %v", err)
				}
				if got.ValidatedModel == nil {
					t.Fatal("ValidatedModel = nil, want value")
				}
				return
			}

			_ = assertCommandFailure(t, err, tt.wantErr, tt.wantMessage)
		})
	}
}

func assertCommandFailure(t *testing.T, err error, wantErr error, wantMessage string) *domain.CommandFailure {
	t.Helper()

	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}

	failure := domain.AsCommandFailure(err)
	if failure == nil {
		t.Fatal("AsCommandFailure(err) = nil, want value")
	}

	if wantMessage != "" {
		if got := failure.Diagnostics[0].Message; got != wantMessage {
			t.Fatalf("message = %q, want %q", got, wantMessage)
		}
	}

	return failure
}

func assertStageOrder(t *testing.T, got []string, want []string) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("order = %#v, want %#v", got, want)
	}
}

func assertScenarioWeights(t *testing.T, got []ScenarioCriterionWeights, want []ScenarioCriterionWeights, tolerance float64) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("len(got) = %d, want %d", len(got), len(want))
	}

	for scenarioIndex := range want {
		if got[scenarioIndex].ScenarioName != want[scenarioIndex].ScenarioName {
			t.Fatalf("ScenarioName[%d] = %q, want %q", scenarioIndex, got[scenarioIndex].ScenarioName, want[scenarioIndex].ScenarioName)
		}
		if len(got[scenarioIndex].CriterionWeights) != len(want[scenarioIndex].CriterionWeights) {
			t.Fatalf("len(CriterionWeights[%d]) = %d, want %d", scenarioIndex, len(got[scenarioIndex].CriterionWeights), len(want[scenarioIndex].CriterionWeights))
		}
		for weightIndex := range want[scenarioIndex].CriterionWeights {
			gotWeight := got[scenarioIndex].CriterionWeights[weightIndex]
			wantWeight := want[scenarioIndex].CriterionWeights[weightIndex]
			if gotWeight.CriterionName != wantWeight.CriterionName {
				t.Fatalf("CriterionName[%d][%d] = %q, want %q", scenarioIndex, weightIndex, gotWeight.CriterionName, wantWeight.CriterionName)
			}
			if math.Abs(gotWeight.Weight-wantWeight.Weight) > tolerance {
				t.Fatalf("Weight[%d][%d] = %0.12f, want %0.12f", scenarioIndex, weightIndex, gotWeight.Weight, wantWeight.Weight)
			}
		}
	}
}

func assertScenarioResults(t *testing.T, got []domain.ScenarioRankingResult, want []domain.ScenarioRankingResult, tolerance float64) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("len(got) = %d, want %d", len(got), len(want))
	}

	for scenarioIndex := range want {
		if got[scenarioIndex].ScenarioName != want[scenarioIndex].ScenarioName {
			t.Fatalf("ScenarioName[%d] = %q, want %q", scenarioIndex, got[scenarioIndex].ScenarioName, want[scenarioIndex].ScenarioName)
		}
		if len(got[scenarioIndex].RankedAlternatives) != len(want[scenarioIndex].RankedAlternatives) {
			t.Fatalf("len(RankedAlternatives[%d]) = %d, want %d", scenarioIndex, len(got[scenarioIndex].RankedAlternatives), len(want[scenarioIndex].RankedAlternatives))
		}
		for alternativeIndex := range want[scenarioIndex].RankedAlternatives {
			gotAlternative := got[scenarioIndex].RankedAlternatives[alternativeIndex]
			wantAlternative := want[scenarioIndex].RankedAlternatives[alternativeIndex]
			if gotAlternative.Name != wantAlternative.Name {
				t.Fatalf("Name[%d][%d] = %q, want %q", scenarioIndex, alternativeIndex, gotAlternative.Name, wantAlternative.Name)
			}
			if gotAlternative.Rank != wantAlternative.Rank {
				t.Fatalf("Rank[%d][%d] = %d, want %d", scenarioIndex, alternativeIndex, gotAlternative.Rank, wantAlternative.Rank)
			}
			if gotAlternative.Excluded != wantAlternative.Excluded {
				t.Fatalf("Excluded[%d][%d] = %t, want %t", scenarioIndex, alternativeIndex, gotAlternative.Excluded, wantAlternative.Excluded)
			}
			if math.Abs(gotAlternative.Score-wantAlternative.Score) > tolerance {
				t.Fatalf("Score[%d][%d] = %0.12f, want %0.12f", scenarioIndex, alternativeIndex, gotAlternative.Score, wantAlternative.Score)
			}
		}
	}
}

func assertAggregatedRanking(t *testing.T, got domain.AggregatedRankingResult, want domain.AggregatedRankingResult, tolerance float64) {
	t.Helper()

	gotRows := domain.CanonicalRankedAlternatives(got.RankedAlternatives)
	wantRows := domain.CanonicalRankedAlternatives(want.RankedAlternatives)
	if len(gotRows) != len(wantRows) {
		t.Fatalf("len(RankedAlternatives) = %d, want %d", len(gotRows), len(wantRows))
	}

	for index := range gotRows {
		if gotRows[index].Name != wantRows[index].Name {
			t.Fatalf("RankedAlternatives[%d].Name = %q, want %q", index, gotRows[index].Name, wantRows[index].Name)
		}
		if gotRows[index].Rank != wantRows[index].Rank {
			t.Fatalf("RankedAlternatives[%d].Rank = %d, want %d", index, gotRows[index].Rank, wantRows[index].Rank)
		}
		if gotRows[index].Excluded != wantRows[index].Excluded {
			t.Fatalf("RankedAlternatives[%d].Excluded = %t, want %t", index, gotRows[index].Excluded, wantRows[index].Excluded)
		}
		if math.Abs(gotRows[index].Score-wantRows[index].Score) > tolerance {
			t.Fatalf("RankedAlternatives[%d].Score = %.10f, want %.10f", index, gotRows[index].Score, wantRows[index].Score)
		}
	}
}
