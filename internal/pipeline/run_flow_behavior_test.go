package pipeline

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

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
