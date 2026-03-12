package pipeline

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

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
