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

func TestExecutionContractsForReportFlow(t *testing.T) {
	t.Parallel()

	reportDefinitions := []ReportDefinition{
		{Name: "summary", Title: "Summary", Format: "markdown"},
		{Name: "scores", Title: "Scores", Format: "json"},
	}
	result := CommandResult{
		CommandName: CommandNameReportGenerate,
		Diagnostics: []Diagnostic{
			NewDiagnostic(DiagnosticSeverityWarning, "stub.warning", "testdata/config/minimal.cue", DiagnosticLocation{}, "warning"),
		},
		ValidatedModel: &ValidatedModelSummary{
			ConfigPath:        "testdata/config/minimal.cue",
			CriterionCount:    3,
			AlternativeCount:  2,
			ScenarioCount:     2,
			ReportDefinitions: reportDefinitions,
		},
		ScenarioResults: []ScenarioRankingResult{
			{
				ScenarioName: "startup",
				RankedAlternatives: []RankedAlternative{
					{Name: "platform-a", Rank: 1, Score: 0.9},
					{Name: "platform-b", Rank: 2, Score: 0.7},
				},
			},
		},
		FinalRanking: &AggregatedRankingResult{
			RankedAlternatives: []RankedAlternative{
				{Name: "platform-a", Rank: 1, Score: 0.9},
				{Name: "platform-b", Rank: 2, Score: 0.7},
			},
		},
		ReportDefinitions: reportDefinitions,
	}

	if result.CommandName != CommandNameReportGenerate {
		t.Fatalf("CommandName = %q, want %q", result.CommandName, CommandNameReportGenerate)
	}

	if got, want := len(result.Diagnostics), 1; got != want {
		t.Fatalf("len(Diagnostics) = %d, want %d", got, want)
	}

	if got, want := result.ValidatedModel.ScenarioCount, 2; got != want {
		t.Fatalf("ScenarioCount = %d, want %d", got, want)
	}

	if got, want := result.ScenarioResults[0].ScenarioName, "startup"; got != want {
		t.Fatalf("ScenarioName = %q, want %q", got, want)
	}

	if got, want := result.FinalRanking.RankedAlternatives[0].Rank, 1; got != want {
		t.Fatalf("FinalRanking rank = %d, want %d", got, want)
	}

	if got, want := result.ReportDefinitions[1].Name, "scores"; got != want {
		t.Fatalf("ReportDefinitions[1].Name = %q, want %q", got, want)
	}
}
