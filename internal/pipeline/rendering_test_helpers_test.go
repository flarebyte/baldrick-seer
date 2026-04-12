package pipeline

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func assertMarkdownStandardOutput(t *testing.T, got string) {
	t.Helper()

	patterns := []string{
		"## Problem",
		"### Title\nDecision Demo Title",
		"### Goal\nChoose the most robust option",
		"### Description\nCompare options across baseline and growth scenarios.",
		"### Notes",
		"- context note",
		"## Alternatives",
		"### Alpha",
		"Lower cost option",
		"### Beta",
		"Higher feature depth",
		"## Scenarios",
		"### Baseline",
		"Current operating constraints.",
		"### Growth",
		"Expansion-oriented scenario.",
		"## Decision Drivers",
		"### Criteria",
		"#### Cost",
		"Lower is better",
		"### Preference Justifications",
		"- Growth: Cost over Quality (strength 3.000000)",
		"### Criteria Weights",
		"- Baseline: Cost=1.000000",
		"- Growth: Cost=0.600000, Quality=0.400000",
		"## Scenario Rankings",
		"### Baseline",
		"1. Alpha (0.900000)",
		"- Beta: excluded (excluded by scenario constraints)",
		"#### Evaluation Notes",
		"Observed baseline measurements.",
		"##### Alpha",
		"Scores:",
		"- Cost: 10",
		"### Growth",
		"1. Alpha (0.800000)",
		"2. Beta (0.400000)",
		"## Notes and Tradeoffs",
		"- Aggregation method: equal_average",
		"### Scenario Weights",
		"- Baseline: 0.500000",
		"### Exclusions",
		"- Baseline: Beta (excluded by scenario constraints)",
	}
	for _, pattern := range patterns {
		if !strings.Contains(got, pattern) {
			t.Fatalf("output missing %q in %q", pattern, got)
		}
	}
}

func assertMarkdownFullOutput(t *testing.T, got string) {
	t.Helper()

	assertMarkdownStandardOutput(t, got)
	patterns := []string{
		"# Summary Markdown Full",
		"##### Beta",
		"Beta growth projection.",
		"- Quality: 2",
		"## Final Ranking",
	}
	for _, pattern := range patterns {
		if !strings.Contains(got, pattern) {
			t.Fatalf("output missing %q in %q", pattern, got)
		}
	}
}

func assertMarkdownFlagsOverrideOutput(t *testing.T, got string) {
	t.Helper()

	present := []string{
		"# Summary Markdown Flags",
		"## Problem",
		"## Decision Drivers",
		"### Criteria Weights",
		"## Notes and Tradeoffs",
		"#### Evaluation Notes",
	}
	for _, pattern := range present {
		if !strings.Contains(got, pattern) {
			t.Fatalf("output missing %q in %q", pattern, got)
		}
	}

	absent := []string{
		"Lower cost option",
		"Higher feature depth",
	}
	for _, pattern := range absent {
		if strings.Contains(got, pattern) {
			t.Fatalf("output unexpectedly contained %q in %q", pattern, got)
		}
	}
}

func assertMarkdownFlagsSuppressedOutput(t *testing.T, got string) {
	t.Helper()

	present := []string{
		"# Summary Markdown Flags Off",
		"## Scenario Rankings",
		"## Final Ranking",
	}
	for _, pattern := range present {
		if !strings.Contains(got, pattern) {
			t.Fatalf("output missing %q in %q", pattern, got)
		}
	}

	absent := []string{
		"## Problem",
		"## Alternatives",
		"## Scenarios",
		"## Decision Drivers",
		"## Notes and Tradeoffs",
		"#### Evaluation Notes",
	}
	for _, pattern := range absent {
		if strings.Contains(got, pattern) {
			t.Fatalf("output unexpectedly contained %q in %q", pattern, got)
		}
	}
}

func assertJSONContextOutput(t *testing.T, got string) {
	t.Helper()

	var payload map[string]any
	if err := json.Unmarshal([]byte(got), &payload); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if _, exists := payload["problem"]; !exists {
		t.Fatalf("problem context missing in %q", got)
	}
	if _, exists := payload["report"]; !exists {
		t.Fatalf("report context missing in %q", got)
	}
	if _, exists := payload["alternatives"]; !exists {
		t.Fatalf("alternatives context missing in %q", got)
	}
	if _, exists := payload["criteria"]; !exists {
		t.Fatalf("criteria context missing in %q", got)
	}
	if _, exists := payload["scenarios"]; !exists {
		t.Fatalf("scenarios context missing in %q", got)
	}
	if _, exists := payload["evaluations"]; !exists {
		t.Fatalf("evaluations context missing in %q", got)
	}

	problem := payload["problem"].(map[string]any)
	if got, want := problem["title"], "Decision Demo Title"; got != want {
		t.Fatalf("problem.title = %#v, want %#v", got, want)
	}
	if got, want := problem["goal"], "Choose the most robust option"; got != want {
		t.Fatalf("problem.goal = %#v, want %#v", got, want)
	}

	report := payload["report"].(map[string]any)
	if got, want := report["title"], "Summary JSON Context"; got != want {
		t.Fatalf("report.title = %#v, want %#v", got, want)
	}
	if got, want := len(report["arguments"].([]any)), 3; got != want {
		t.Fatalf("len(report.arguments) = %d, want %d", got, want)
	}

	scenarioResults := payload["scenarioResults"].([]any)
	firstScenario := scenarioResults[0].(map[string]any)
	if got, want := firstScenario["scenarioTitle"], "Baseline"; got != want {
		t.Fatalf("scenarioResults[0].scenarioTitle = %#v, want %#v", got, want)
	}

	finalRanking := payload["finalRanking"].([]any)
	firstRank := finalRanking[0].(map[string]any)
	if got, want := firstRank["alternativeTitle"], "Alpha"; got != want {
		t.Fatalf("finalRanking[0].alternativeTitle = %#v, want %#v", got, want)
	}

	evaluations := payload["evaluations"].([]any)
	firstEvaluation := evaluations[0].(map[string]any)
	entries := firstEvaluation["evaluations"].([]any)
	firstAlternative := entries[0].(map[string]any)
	values := firstAlternative["values"].([]any)
	firstValue := values[0].(map[string]any)
	if got, want := firstValue["criterionTitle"], "Cost"; got != want {
		t.Fatalf("evaluations[0].evaluations[0].values[0].criterionTitle = %#v, want %#v", got, want)
	}
}

func assertCSVSchemaOutput(t *testing.T, got string) {
	t.Helper()

	reader := csv.NewReader(strings.NewReader(got))
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("csv.ReadAll() error = %v", err)
	}
	if len(records) < 2 {
		t.Fatalf("len(records) = %d, want at least 2", len(records))
	}

	wantHeader := []string{"scenario", "alternative", "criterion", "value", "score", "rank", "excluded", "exclusion_reason"}
	if !reflect.DeepEqual(records[0], wantHeader) {
		t.Fatalf("header = %#v, want %#v", records[0], wantHeader)
	}

	for index, record := range records[1:] {
		if len(record) != len(wantHeader) {
			t.Fatalf("record %d width = %d, want %d", index+1, len(record), len(wantHeader))
		}
	}

	schema := csvSchemaDescriptions()
	for _, column := range wantHeader {
		if schema[column] == "" {
			t.Fatalf("missing schema description for %q", column)
		}
	}

	if !strings.Contains(got, "overall,alpha,,,0.850000,1,false,") {
		t.Fatalf("overall final ranking row missing in %q", got)
	}
	if !strings.Contains(got, "baseline,beta,cost,20,,,true,excluded by scenario constraints") {
		t.Fatalf("excluded scenario row missing in %q", got)
	}
	if !strings.Contains(got, "growth,alpha,quality,3,0.800000,1,false,") {
		t.Fatalf("criterion-level row missing in %q", got)
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

func singleReportDefinitions(report ReportConfig) []domain.ReportDefinition {
	return []domain.ReportDefinition{{Name: report.Name, Title: report.Title, Format: report.Format}}
}

func renderReportsForTest(
	t *testing.T,
	config LoadedConfig,
	scenarioResults []domain.ScenarioRankingResult,
	finalRanking domain.AggregatedRankingResult,
	reportDefinitions []domain.ReportDefinition,
	scenarioWeights []ScenarioCriterionWeights,
) RenderReportsOutput {
	t.Helper()

	renderer := DefaultReportRenderer{}
	got, err := renderer.RenderReports(context.Background(), RenderReportsInput{
		Command: domain.CommandRequest{
			CommandName: domain.CommandNameReportGenerate,
			ConfigPath:  config.Path,
		},
		ValidatedModel:    domain.ValidatedModelSummary{ConfigPath: config.Path},
		ScenarioResults:   scenarioResults,
		FinalRanking:      finalRanking,
		ReportDefinitions: reportDefinitions,
		ScenarioWeights:   scenarioWeights,
		Config:            config,
	})
	if err != nil {
		t.Fatalf("RenderReports() error = %v", err)
	}
	return got
}

func expectedRenderedReport(
	t *testing.T,
	report ReportConfig,
	config LoadedConfig,
	scenarioResults []domain.ScenarioRankingResult,
	scenarioWeights []ScenarioCriterionWeights,
) string {
	t.Helper()

	want, err := renderReport(report, config.Config, scenarioResults, config.Config.Aggregation, scenarioWeights)
	if err != nil {
		t.Fatalf("renderReport() error = %v", err)
	}
	return want
}

func reportLoadedConfig(reports ...ReportConfig) LoadedConfig {
	config := validLoadedConfig()
	config.Config.Problem = &ProblemConfig{
		Name:        "Decision Demo",
		Title:       "Decision Demo Title",
		Goal:        "Choose the most robust option",
		Description: "Compare options across baseline and growth scenarios.",
		Owner:       "platform-team",
		Notes:       []string{"context note", "reviewable artifact"},
	}
	config.Config.Reports = append([]ReportConfig(nil), reports...)
	config.Config.CriteriaCatalog = []CriterionConfig{
		{Name: "cost", Title: "Cost", Description: "Lower is better", Polarity: "cost", ValueType: "number"},
		{Name: "quality", Title: "Quality", Description: "Higher is better", Polarity: "benefit", ValueType: "ordinal", ScaleGuidance: []any{"poor", "good"}},
	}
	config.Config.Alternatives = []AlternativeConfig{
		{Name: "alpha", Title: "Alpha", Description: "Lower cost option"},
		{Name: "beta", Title: "Beta", Description: "Higher feature depth"},
	}
	config.Config.Scenarios = []ScenarioConfig{
		{
			Name:        "baseline",
			Title:       "Baseline",
			Description: "Current operating constraints.",
			Narrative:   "Steady-state execution with tight delivery pressure.",
			ActiveCriteria: []ScenarioCriterionRef{
				{CriterionName: "cost"},
			},
		},
		{
			Name:        "growth",
			Title:       "Growth",
			Description: "Expansion-oriented scenario.",
			Narrative:   "Rapid scale-up with higher demand volatility.",
			ActiveCriteria: []ScenarioCriterionRef{
				{CriterionName: "cost"},
				{CriterionName: "quality"},
			},
			Preferences: &ScenarioPreferences{
				Method: "ahp",
				Scale:  "saaty",
				Comparisons: []PairwiseComparison{
					{
						MoreImportantCriterionName: "cost",
						LessImportantCriterionName: "quality",
						Strength:                   3,
					},
				},
			},
		},
	}
	config.Config.Evaluations = []EvaluationConfig{
		{
			ScenarioName: "baseline",
			Description:  "Observed baseline measurements.",
			Evaluations: []AlternativeEvaluationConfig{
				{
					AlternativeName: "alpha",
					Description:     "Alpha baseline observation.",
					Values: map[string]CriterionValue{
						"cost": {Kind: "number", Value: 10},
					},
				},
				{
					AlternativeName: "beta",
					Description:     "Beta baseline observation.",
					Values: map[string]CriterionValue{
						"cost": {Kind: "number", Value: 20},
					},
				},
			},
		},
		{
			ScenarioName: "growth",
			Description:  "Projected growth measurements.",
			Evaluations: []AlternativeEvaluationConfig{
				{
					AlternativeName: "alpha",
					Description:     "Alpha growth projection.",
					Values: map[string]CriterionValue{
						"cost":    {Kind: "number", Value: 12},
						"quality": {Kind: "ordinal", Value: 3},
					},
				},
				{
					AlternativeName: "beta",
					Description:     "Beta growth projection.",
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
