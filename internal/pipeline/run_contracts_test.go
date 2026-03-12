package pipeline

import (
	"context"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

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
		ScenarioWeights: []ScenarioCriterionWeights{{
			ScenarioName: "startup",
			CriterionWeights: []CriterionWeight{
				{CriterionName: "cost", Weight: 0.6},
				{CriterionName: "speed", Weight: 0.4},
			},
		}},
	}
	rankOutput := RankScenariosOutput{
		ScenarioResults: []domain.ScenarioRankingResult{{
			ScenarioName: "startup",
			RankedAlternatives: []domain.RankedAlternative{
				{Name: "platform-a", Rank: 1, Score: 0.9},
			},
		}},
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
