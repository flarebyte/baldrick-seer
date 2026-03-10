package domain

import "testing"

func TestContractsCanBeConstructed(t *testing.T) {
	t.Parallel()

	reportDefinitions := []ReportDefinition{
		{Name: "summary", Title: "Summary", Format: "markdown"},
	}
	diagnostics := []Diagnostic{
		{Severity: DiagnosticSeverityWarning, Code: "stub.warning", Message: "warning"},
	}
	scenarioResults := []ScenarioRankingResult{
		{
			ScenarioName: "startup",
			RankedAlternatives: []RankedAlternative{
				{Name: "platform-a", Rank: 1, Score: 1},
			},
		},
	}

	request := CommandRequest{
		CommandName: CommandNameValidate,
		ConfigPath:  "testdata/config/minimal.cue",
	}
	result := CommandResult{
		CommandName: CommandNameValidate,
		Diagnostics: diagnostics,
		ValidatedModel: &ValidatedModelSummary{
			ConfigPath:        request.ConfigPath,
			CriterionCount:    3,
			AlternativeCount:  2,
			ScenarioCount:     1,
			ReportDefinitions: reportDefinitions,
		},
		ScenarioResults: scenarioResults,
		FinalRanking: &AggregatedRankingResult{
			RankedAlternatives: []RankedAlternative{
				{Name: "platform-a", Rank: 1, Score: 1},
			},
		},
		ReportDefinitions: reportDefinitions,
	}

	if request.CommandName != CommandNameValidate {
		t.Fatalf("CommandName = %q, want %q", request.CommandName, CommandNameValidate)
	}

	if result.ValidatedModel == nil {
		t.Fatal("ValidatedModel = nil, want value")
	}

	if got, want := result.ValidatedModel.ConfigPath, request.ConfigPath; got != want {
		t.Fatalf("ConfigPath = %q, want %q", got, want)
	}

	if got, want := result.ScenarioResults[0].RankedAlternatives[0].Name, "platform-a"; got != want {
		t.Fatalf("RankedAlternative.Name = %q, want %q", got, want)
	}

	if got, want := result.ReportDefinitions[0].Format, "markdown"; got != want {
		t.Fatalf("ReportDefinitions[0].Format = %q, want %q", got, want)
	}
}
