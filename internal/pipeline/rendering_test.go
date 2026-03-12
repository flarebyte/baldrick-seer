package pipeline

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func TestDefaultReportRenderer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		report     ReportConfig
		scenarios  []domain.ScenarioRankingResult
		weights    []ScenarioCriterionWeights
		wantGolden string
	}{
		{
			name: "markdown renderer output shape",
			report: ReportConfig{
				Name:      "summary-markdown",
				Title:     "Summary Markdown",
				Format:    "markdown",
				Arguments: []string{"include-scores=true"},
			},
			scenarios:  reportScenarioResults(),
			weights:    reportScenarioWeights(),
			wantGolden: "report_markdown.out.golden",
		},
		{
			name: "json renderer output shape",
			report: ReportConfig{
				Name:      "summary-json",
				Title:     "Summary JSON",
				Format:    "json",
				Arguments: []string{"include-weights=true", "pretty=true"},
			},
			scenarios:  reportScenarioResults(),
			weights:    reportScenarioWeights(),
			wantGolden: "report_json.out.golden",
		},
		{
			name: "csv renderer output shape",
			report: ReportConfig{
				Name:      "summary-csv",
				Title:     "Summary CSV",
				Format:    "csv",
				Arguments: []string{"columns=scenario,alternative,score,rank", "header=true"},
			},
			scenarios:  reportScenarioResults(),
			wantGolden: "report_csv.out.golden",
		},
		{
			name: "empty final aggregated ranking when all alternatives are ineligible",
			report: ReportConfig{
				Name:   "summary-empty",
				Title:  "Summary Empty",
				Format: "markdown",
			},
			scenarios: []domain.ScenarioRankingResult{
				{
					ScenarioName: "baseline",
					RankedAlternatives: []domain.RankedAlternative{
						{Name: "alpha", Excluded: true},
						{Name: "beta", Excluded: true},
					},
				},
			},
			wantGolden: "report_empty_markdown.out.golden",
		},
		{
			name: "scenario focus aggregates only participating scenarios",
			report: ReportConfig{
				Name:   "summary-focused",
				Title:  "Summary Focused",
				Format: "markdown",
				Focus: &ReportFocus{
					ScenarioNames: []string{"growth"},
				},
			},
			scenarios: []domain.ScenarioRankingResult{
				{
					ScenarioName: "baseline",
					RankedAlternatives: []domain.RankedAlternative{
						{Name: "alpha", Excluded: true, ExclusionReason: "excluded by scenario constraints"},
						{Name: "beta", Rank: 1, Score: 0.7},
					},
				},
				{
					ScenarioName: "growth",
					RankedAlternatives: []domain.RankedAlternative{
						{Name: "alpha", Rank: 1, Score: 0.8},
						{Name: "beta", Rank: 2, Score: 0.4},
					},
				},
			},
			weights: []ScenarioCriterionWeights{
				{
					ScenarioName: "growth",
					CriterionWeights: []CriterionWeight{
						{CriterionName: "cost", Weight: 0.6},
						{CriterionName: "quality", Weight: 0.4},
					},
				},
			},
			wantGolden: "report_focused_markdown.out.golden",
		},
		{
			name: "json report honors alternative focus",
			report: ReportConfig{
				Name:      "summary-json-focused",
				Title:     "Summary JSON Focused",
				Format:    "json",
				Arguments: []string{"include-weights=true", "pretty=true"},
				Focus: &ReportFocus{
					AlternativeNames: []string{"alpha"},
				},
			},
			scenarios:  reportScenarioResults(),
			weights:    reportScenarioWeights(),
			wantGolden: "report_focused_json.out.golden",
		},
		{
			name: "csv report honors criterion focus",
			report: ReportConfig{
				Name:      "summary-csv-focused",
				Title:     "Summary CSV Focused",
				Format:    "csv",
				Arguments: []string{"columns=scenario,alternative,criterion,value,score,rank,excluded,exclusion_reason", "header=true"},
				Focus: &ReportFocus{
					CriterionNames: []string{"cost"},
				},
			},
			scenarios:  reportScenarioResults(),
			wantGolden: "report_focused_csv.out.golden",
		},
	}

	renderer := DefaultReportRenderer{}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			config := reportLoadedConfig(tt.report)
			got, err := renderer.RenderReports(context.Background(), RenderReportsInput{
				Command: domain.CommandRequest{
					CommandName: domain.CommandNameReportGenerate,
					ConfigPath:  config.Path,
				},
				ValidatedModel: domain.ValidatedModelSummary{
					ConfigPath: config.Path,
				},
				ScenarioResults: tt.scenarios,
				FinalRanking:    domain.AggregatedRankingResult{},
				ReportDefinitions: []domain.ReportDefinition{
					{Name: tt.report.Name, Title: tt.report.Title, Format: tt.report.Format},
				},
				ScenarioWeights: tt.weights,
				Config:          config,
			})
			if err != nil {
				t.Fatalf("RenderReports() error = %v", err)
			}

			if got, want := got.RenderedOutput, readPipelineGolden(t, tt.wantGolden); got != want {
				t.Fatalf("RenderedOutput = %q, want %q", got, want)
			}
		})
	}
}

func TestDefaultReportRendererRepeatedRunDeterminism(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		report ReportConfig
	}{
		{
			name: "markdown",
			report: ReportConfig{
				Name:      "summary-markdown",
				Title:     "Summary Markdown",
				Format:    "markdown",
				Arguments: []string{"include-scores=true"},
			},
		},
		{
			name: "json",
			report: ReportConfig{
				Name:      "summary-json",
				Title:     "Summary JSON",
				Format:    "json",
				Arguments: []string{"include-weights=true", "pretty=true"},
			},
		},
		{
			name: "csv",
			report: ReportConfig{
				Name:      "summary-csv",
				Title:     "Summary CSV",
				Format:    "csv",
				Arguments: []string{"columns=scenario,alternative,score,rank", "header=true"},
			},
		},
		{
			name: "focused json",
			report: ReportConfig{
				Name:      "summary-json-focused",
				Title:     "Summary JSON Focused",
				Format:    "json",
				Arguments: []string{"include-weights=true", "pretty=true"},
				Focus: &ReportFocus{
					AlternativeNames: []string{"alpha"},
				},
			},
		},
		{
			name: "focused csv",
			report: ReportConfig{
				Name:      "summary-csv-focused",
				Title:     "Summary CSV Focused",
				Format:    "csv",
				Arguments: []string{"columns=scenario,alternative,criterion,value,score,rank,excluded,exclusion_reason", "header=true"},
				Focus: &ReportFocus{
					CriterionNames: []string{"cost"},
				},
			},
		},
	}

	renderer := DefaultReportRenderer{}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			input := RenderReportsInput{
				Command: domain.CommandRequest{
					CommandName: domain.CommandNameReportGenerate,
					ConfigPath:  fixtureConfigPath(),
				},
				ValidatedModel:  domain.ValidatedModelSummary{ConfigPath: fixtureConfigPath()},
				ScenarioResults: reportScenarioResults(),
				FinalRanking:    domain.AggregatedRankingResult{},
				ReportDefinitions: []domain.ReportDefinition{
					{Name: tt.report.Name, Title: tt.report.Title, Format: tt.report.Format},
				},
				ScenarioWeights: reportScenarioWeights(),
				Config:          reportLoadedConfig(tt.report),
			}

			first, err := renderer.RenderReports(context.Background(), input)
			if err != nil {
				t.Fatalf("first RenderReports() error = %v", err)
			}

			second, err := renderer.RenderReports(context.Background(), input)
			if err != nil {
				t.Fatalf("second RenderReports() error = %v", err)
			}

			if !reflect.DeepEqual(first, second) {
				t.Fatalf("first = %#v, second = %#v", first, second)
			}
		})
	}
}

func TestDefaultReportRendererHonorsSelectedReportDefinitions(t *testing.T) {
	t.Parallel()

	renderer := DefaultReportRenderer{}
	reportA := ReportConfig{
		Name:   "a-markdown",
		Title:  "A Markdown",
		Format: "markdown",
	}
	reportB := ReportConfig{
		Name:      "b-json",
		Title:     "B JSON",
		Format:    "json",
		Arguments: []string{"pretty=true"},
	}

	config := reportLoadedConfig(reportA, reportB)
	got, err := renderer.RenderReports(context.Background(), RenderReportsInput{
		Command: domain.CommandRequest{
			CommandName: domain.CommandNameReportGenerate,
			ConfigPath:  config.Path,
		},
		ValidatedModel: domain.ValidatedModelSummary{
			ConfigPath: config.Path,
		},
		ScenarioResults: reportScenarioResults(),
		ReportDefinitions: []domain.ReportDefinition{
			{Name: reportB.Name, Title: reportB.Title, Format: reportB.Format},
		},
		Config: config,
	})
	if err != nil {
		t.Fatalf("RenderReports() error = %v", err)
	}

	if got, want := len(got.ReportDefinitions), 1; got != want {
		t.Fatalf("len(ReportDefinitions) = %d, want %d", got, want)
	}
	if got, want := got.ReportDefinitions[0].Name, reportB.Name; got != want {
		t.Fatalf("ReportDefinitions[0].Name = %q, want %q", got, want)
	}
	if got.RenderedOutput != readPipelineGolden(t, "report_selected_json.out.golden") {
		t.Fatalf("RenderedOutput = %q, want selected json golden", got.RenderedOutput)
	}
}

func TestDefaultReportRendererCanonicalizesShuffledInput(t *testing.T) {
	t.Parallel()

	report := ReportConfig{
		Name:      "summary-json",
		Title:     "Summary JSON",
		Format:    "json",
		Arguments: []string{"include-weights=true", "pretty=true"},
	}
	config := reportLoadedConfig(report)
	definitions := []domain.ReportDefinition{{Name: report.Name, Title: report.Title, Format: report.Format}}

	canonicalInput := RenderReportsInput{
		Command: domain.CommandRequest{
			CommandName: domain.CommandNameReportGenerate,
			ConfigPath:  fixtureConfigPath(),
		},
		ValidatedModel:    domain.ValidatedModelSummary{ConfigPath: fixtureConfigPath()},
		ScenarioResults:   reportScenarioResults(),
		FinalRanking:      domain.AggregatedRankingResult{},
		ReportDefinitions: definitions,
		ScenarioWeights:   reportScenarioWeights(),
		Config:            config,
	}
	shuffledInput := RenderReportsInput{
		Command:        canonicalInput.Command,
		ValidatedModel: canonicalInput.ValidatedModel,
		FinalRanking:   canonicalInput.FinalRanking,
		ReportDefinitions: []domain.ReportDefinition{
			definitions[0],
		},
		ScenarioResults: []domain.ScenarioRankingResult{
			{
				ScenarioName: "baseline",
				RankedAlternatives: []domain.RankedAlternative{
					{Name: "beta", Excluded: true, ExclusionReason: "excluded by scenario constraints"},
					{Name: "alpha", Rank: 1, Score: 0.9},
				},
			},
			{
				ScenarioName: "growth",
				RankedAlternatives: []domain.RankedAlternative{
					{Name: "beta", Rank: 2, Score: 0.4},
					{Name: "alpha", Rank: 1, Score: 0.8},
				},
			},
		},
		ScenarioWeights: []ScenarioCriterionWeights{
			{
				ScenarioName: "baseline",
				CriterionWeights: []CriterionWeight{
					{CriterionName: "cost", Weight: 1},
				},
			},
			{
				ScenarioName: "growth",
				CriterionWeights: []CriterionWeight{
					{CriterionName: "quality", Weight: 0.4},
					{CriterionName: "cost", Weight: 0.6},
				},
			},
		},
		Config: config,
	}

	renderer := DefaultReportRenderer{}
	canonical, err := renderer.RenderReports(context.Background(), canonicalInput)
	if err != nil {
		t.Fatalf("canonical RenderReports() error = %v", err)
	}

	shuffled, err := renderer.RenderReports(context.Background(), shuffledInput)
	if err != nil {
		t.Fatalf("shuffled RenderReports() error = %v", err)
	}

	if canonical.RenderedOutput != shuffled.RenderedOutput {
		t.Fatalf("canonical output = %q, shuffled output = %q", canonical.RenderedOutput, shuffled.RenderedOutput)
	}
}

func readPipelineGolden(t *testing.T, name string) string {
	t.Helper()

	content, err := os.ReadFile(filepath.Join("..", "..", "testdata", "golden", name))
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", name, err)
	}
	return string(content)
}

func reportLoadedConfig(reports ...ReportConfig) LoadedConfig {
	config := validLoadedConfig()
	config.Config.Problem = &ProblemConfig{Name: "Decision Demo"}
	config.Config.Reports = append([]ReportConfig(nil), reports...)
	config.Config.CriteriaCatalog = []CriterionConfig{
		{Name: "cost", Polarity: "cost", ValueType: "number"},
		{Name: "quality", Polarity: "benefit", ValueType: "ordinal", ScaleGuidance: []any{"poor", "good"}},
	}
	config.Config.Alternatives = []AlternativeConfig{
		{Name: "alpha"},
		{Name: "beta"},
	}
	config.Config.Scenarios = []ScenarioConfig{
		{
			Name: "baseline",
			ActiveCriteria: []ScenarioCriterionRef{
				{CriterionName: "cost"},
			},
		},
		{
			Name: "growth",
			ActiveCriteria: []ScenarioCriterionRef{
				{CriterionName: "cost"},
				{CriterionName: "quality"},
			},
		},
	}
	config.Config.Evaluations = []EvaluationConfig{
		{
			ScenarioName: "baseline",
			Evaluations: []AlternativeEvaluationConfig{
				{
					AlternativeName: "alpha",
					Values: map[string]CriterionValue{
						"cost": {Kind: "number", Value: 10},
					},
				},
				{
					AlternativeName: "beta",
					Values: map[string]CriterionValue{
						"cost": {Kind: "number", Value: 20},
					},
				},
			},
		},
		{
			ScenarioName: "growth",
			Evaluations: []AlternativeEvaluationConfig{
				{
					AlternativeName: "alpha",
					Values: map[string]CriterionValue{
						"cost":    {Kind: "number", Value: 12},
						"quality": {Kind: "ordinal", Value: 3},
					},
				},
				{
					AlternativeName: "beta",
					Values: map[string]CriterionValue{
						"cost":    {Kind: "number", Value: 18},
						"quality": {Kind: "ordinal", Value: 2},
					},
				},
			},
		},
	}
	config.Config.Aggregation = &AggregationConfig{Method: "equal_average"}
	return config
}

func reportScenarioResults() []domain.ScenarioRankingResult {
	return []domain.ScenarioRankingResult{
		{
			ScenarioName: "growth",
			RankedAlternatives: []domain.RankedAlternative{
				{Name: "alpha", Rank: 1, Score: 0.8},
				{Name: "beta", Rank: 2, Score: 0.4},
			},
		},
		{
			ScenarioName: "baseline",
			RankedAlternatives: []domain.RankedAlternative{
				{Name: "alpha", Rank: 1, Score: 0.9},
				{Name: "beta", Excluded: true, ExclusionReason: "excluded by scenario constraints"},
			},
		},
	}
}

func reportScenarioWeights() []ScenarioCriterionWeights {
	return []ScenarioCriterionWeights{
		{
			ScenarioName: "growth",
			CriterionWeights: []CriterionWeight{
				{CriterionName: "quality", Weight: 0.4},
				{CriterionName: "cost", Weight: 0.6},
			},
		},
		{
			ScenarioName: "baseline",
			CriterionWeights: []CriterionWeight{
				{CriterionName: "cost", Weight: 1},
			},
		},
	}
}
