package domain

import (
	"reflect"
	"testing"
)

func TestCanonicalDiagnostics(t *testing.T) {
	t.Parallel()

	input := []Diagnostic{
		NewDiagnostic(DiagnosticSeverityWarning, "b.code", "z", DiagnosticLocation{Line: 2, Column: 1}, "warn"),
		NewDiagnostic(DiagnosticSeverityError, "b.code", "a", DiagnosticLocation{Line: 3, Column: 1}, "error-b"),
		NewDiagnostic(DiagnosticSeverityError, "a.code", "b", DiagnosticLocation{Line: 1, Column: 2}, "error-a"),
		NewDiagnostic(DiagnosticSeverityError, "b.code", "a", DiagnosticLocation{Line: 1, Column: 5}, "error-c"),
	}

	got := CanonicalDiagnostics(input)
	want := []Diagnostic{
		input[2],
		input[3],
		input[1],
		input[0],
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("CanonicalDiagnostics() = %#v, want %#v", got, want)
	}

	if reflect.DeepEqual(input, got) {
		t.Fatalf("CanonicalDiagnostics() mutated ordering semantics of caller slice: input = %#v", input)
	}
}

func TestCanonicalRankedAlternatives(t *testing.T) {
	t.Parallel()

	input := []RankedAlternative{
		{Name: "beta", Rank: 2, Score: 0.8},
		{Name: "alpha", Rank: 1, Score: 0.7},
		{Name: "aardvark", Rank: 1, Score: 0.9},
		{Name: "beta", Rank: 1, Score: 0.6},
		{Name: "zeta", Excluded: true},
		{Name: "gamma", Excluded: true},
	}

	got := CanonicalRankedAlternatives(input)
	want := []RankedAlternative{
		{Name: "aardvark", Rank: 1, Score: 0.9},
		{Name: "alpha", Rank: 1, Score: 0.7},
		{Name: "beta", Rank: 1, Score: 0.6},
		{Name: "beta", Rank: 2, Score: 0.8},
		{Name: "gamma", Excluded: true},
		{Name: "zeta", Excluded: true},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("CanonicalRankedAlternatives() = %#v, want %#v", got, want)
	}
}

func TestCanonicalNamedCollections(t *testing.T) {
	t.Parallel()

	reports := []ReportDefinition{
		{Name: "zeta", Title: "Zeta", Format: "json"},
		{Name: "alpha", Title: "Alpha", Format: "markdown"},
		{Name: "alpha", Title: "Alpha", Format: "csv"},
	}
	scenarios := []ScenarioRankingResult{
		{
			ScenarioName: "zeta",
			RankedAlternatives: []RankedAlternative{
				{Name: "beta", Rank: 2, Score: 2},
				{Name: "alpha", Rank: 1, Score: 1},
			},
		},
		{
			ScenarioName: "alpha",
			RankedAlternatives: []RankedAlternative{
				{Name: "beta", Rank: 1, Score: 1},
				{Name: "alpha", Rank: 1, Score: 2},
			},
		},
	}

	gotReports := CanonicalReportDefinitions(reports)
	gotScenarios := CanonicalScenarioResults(scenarios)
	gotNames := CanonicalNames([]string{"scenario-b", "scenario-a", "scenario-c"})

	wantReports := []ReportDefinition{
		{Name: "alpha", Title: "Alpha", Format: "csv"},
		{Name: "alpha", Title: "Alpha", Format: "markdown"},
		{Name: "zeta", Title: "Zeta", Format: "json"},
	}
	wantScenarios := []ScenarioRankingResult{
		{
			ScenarioName: "alpha",
			RankedAlternatives: []RankedAlternative{
				{Name: "alpha", Rank: 1, Score: 2},
				{Name: "beta", Rank: 1, Score: 1},
			},
		},
		{
			ScenarioName: "zeta",
			RankedAlternatives: []RankedAlternative{
				{Name: "alpha", Rank: 1, Score: 1},
				{Name: "beta", Rank: 2, Score: 2},
			},
		},
	}
	wantNames := []string{"scenario-a", "scenario-b", "scenario-c"}

	if !reflect.DeepEqual(gotReports, wantReports) {
		t.Fatalf("CanonicalReportDefinitions() = %#v, want %#v", gotReports, wantReports)
	}

	if !reflect.DeepEqual(gotScenarios, wantScenarios) {
		t.Fatalf("CanonicalScenarioResults() = %#v, want %#v", gotScenarios, wantScenarios)
	}

	if !reflect.DeepEqual(gotNames, wantNames) {
		t.Fatalf("CanonicalNames() = %#v, want %#v", gotNames, wantNames)
	}
}

func TestCanonicalCommandResultDeterministic(t *testing.T) {
	t.Parallel()

	input := CommandResult{
		CommandName: CommandNameReportGenerate,
		Diagnostics: []Diagnostic{
			NewDiagnostic(DiagnosticSeverityWarning, "warn", "b", DiagnosticLocation{}, "warn"),
			NewDiagnostic(DiagnosticSeverityError, "err", "a", DiagnosticLocation{}, "err"),
		},
		ValidatedModel: &ValidatedModelSummary{
			ConfigPath: "testdata/config/minimal.cue",
			ReportDefinitions: []ReportDefinition{
				{Name: "zeta", Title: "Zeta", Format: "json"},
				{Name: "alpha", Title: "Alpha", Format: "markdown"},
			},
		},
		ScenarioResults: []ScenarioRankingResult{
			{
				ScenarioName: "zeta",
				RankedAlternatives: []RankedAlternative{
					{Name: "beta", Rank: 2, Score: 2},
					{Name: "alpha", Rank: 1, Score: 1},
				},
			},
		},
		FinalRanking: &AggregatedRankingResult{
			RankedAlternatives: []RankedAlternative{
				{Name: "beta", Rank: 2, Score: 2},
				{Name: "alpha", Rank: 1, Score: 1},
			},
		},
		ReportDefinitions: []ReportDefinition{
			{Name: "zeta", Title: "Zeta", Format: "json"},
			{Name: "alpha", Title: "Alpha", Format: "markdown"},
		},
	}

	first := CanonicalCommandResult(input)
	second := CanonicalCommandResult(input)

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("CanonicalCommandResult() first = %#v, second = %#v", first, second)
	}
}
