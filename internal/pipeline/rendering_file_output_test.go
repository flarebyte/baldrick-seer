package pipeline

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func TestDefaultReportRendererWritesFileTargetedReports(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "model.cue")
	report := ReportConfig{
		Name:      "summary-json",
		Title:     "Summary JSON",
		Format:    "json",
		Filepath:  "artifacts/summary.json",
		Arguments: []string{"include-weights=true", "pretty=true"},
	}

	config := reportLoadedConfig(report)
	config.Path = configPath

	renderer := DefaultReportRenderer{}
	got, err := renderer.RenderReports(context.Background(), RenderReportsInput{
		Command: domain.CommandRequest{
			CommandName: domain.CommandNameReportGenerate,
			ConfigPath:  config.Path,
		},
		ValidatedModel:    domain.ValidatedModelSummary{ConfigPath: config.Path},
		ScenarioResults:   reportScenarioResults(),
		FinalRanking:      domain.AggregatedRankingResult{},
		ReportDefinitions: []domain.ReportDefinition{{Name: report.Name, Title: report.Title, Format: report.Format}},
		ScenarioWeights:   reportScenarioWeights(),
		Config:            config,
	})
	if err != nil {
		t.Fatalf("RenderReports() error = %v", err)
	}

	if got.RenderedOutput != "" {
		t.Fatalf("RenderedOutput = %q, want empty stdout output", got.RenderedOutput)
	}

	content, err := os.ReadFile(filepath.Join(tempDir, "artifacts", "summary.json"))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	want, err := renderReport(report, config.Config, reportScenarioResults(), config.Config.Aggregation, reportScenarioWeights())
	if err != nil {
		t.Fatalf("renderReport() error = %v", err)
	}
	if got := string(content); got != want {
		t.Fatalf("file content = %q, want %q", got, want)
	}
}

func TestDefaultReportRendererSplitsMixedStdoutAndFileOutputs(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "model.cue")
	reportFile := ReportConfig{
		Name:      "a-json",
		Title:     "A JSON",
		Format:    "json",
		Filepath:  "artifacts/a.json",
		Arguments: []string{"pretty=true"},
	}
	reportStdout := ReportConfig{
		Name:      "b-markdown",
		Title:     "B Markdown",
		Format:    "markdown",
		Arguments: []string{"include-scores=true"},
	}

	config := reportLoadedConfig(reportFile, reportStdout)
	config.Path = configPath

	renderer := DefaultReportRenderer{}
	got, err := renderer.RenderReports(context.Background(), RenderReportsInput{
		Command: domain.CommandRequest{
			CommandName: domain.CommandNameReportGenerate,
			ConfigPath:  config.Path,
		},
		ValidatedModel: domain.ValidatedModelSummary{
			ConfigPath: config.Path,
		},
		ScenarioResults: reportScenarioResults(),
		FinalRanking:    domain.AggregatedRankingResult{},
		ReportDefinitions: []domain.ReportDefinition{
			{Name: reportStdout.Name, Title: reportStdout.Title, Format: reportStdout.Format},
			{Name: reportFile.Name, Title: reportFile.Title, Format: reportFile.Format},
		},
		ScenarioWeights: reportScenarioWeights(),
		Config:          config,
	})
	if err != nil {
		t.Fatalf("RenderReports() error = %v", err)
	}

	wantStdout, err := renderReport(reportStdout, config.Config, reportScenarioResults(), config.Config.Aggregation, reportScenarioWeights())
	if err != nil {
		t.Fatalf("renderReport(stdout) error = %v", err)
	}
	if got := got.RenderedOutput; got != wantStdout {
		t.Fatalf("RenderedOutput = %q, want %q", got, wantStdout)
	}

	content, err := os.ReadFile(filepath.Join(tempDir, "artifacts", "a.json"))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	wantFile, err := renderReport(reportFile, config.Config, reportScenarioResults(), config.Config.Aggregation, reportScenarioWeights())
	if err != nil {
		t.Fatalf("renderReport(file) error = %v", err)
	}
	if got := string(content); got != wantFile {
		t.Fatalf("file content = %q, want %q", got, wantFile)
	}
}

func TestDefaultReportRendererResolvesFilepathsRelativeToConfigDirectory(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "config")
	configPath := filepath.Join(configDir, "model.cue")
	report := ReportConfig{
		Name:      "summary-markdown",
		Title:     "Summary Markdown",
		Format:    "markdown",
		Filepath:  "../artifacts/summary.md",
		Arguments: []string{"include-scores=true"},
	}

	config := reportLoadedConfig(report)
	config.Path = configPath

	renderer := DefaultReportRenderer{}
	_, err := renderer.RenderReports(context.Background(), RenderReportsInput{
		Command: domain.CommandRequest{
			CommandName: domain.CommandNameReportGenerate,
			ConfigPath:  config.Path,
		},
		ValidatedModel:    domain.ValidatedModelSummary{ConfigPath: config.Path},
		ScenarioResults:   reportScenarioResults(),
		FinalRanking:      domain.AggregatedRankingResult{},
		ReportDefinitions: []domain.ReportDefinition{{Name: report.Name, Title: report.Title, Format: report.Format}},
		ScenarioWeights:   reportScenarioWeights(),
		Config:            config,
	})
	if err != nil {
		t.Fatalf("RenderReports() error = %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tempDir, "artifacts", "summary.md"))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	want, err := renderReport(report, config.Config, reportScenarioResults(), config.Config.Aggregation, reportScenarioWeights())
	if err != nil {
		t.Fatalf("renderReport() error = %v", err)
	}
	if got := string(content); got != want {
		t.Fatalf("file content = %q, want %q", got, want)
	}
}

func TestDefaultReportRendererOverwritesFileTargetedReportsDeterministically(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "model.cue")
	report := ReportConfig{
		Name:      "summary-json",
		Title:     "Summary JSON",
		Format:    "json",
		Filepath:  "artifacts/summary.json",
		Arguments: []string{"pretty=true"},
	}

	config := reportLoadedConfig(report)
	config.Path = configPath

	input := RenderReportsInput{
		Command: domain.CommandRequest{
			CommandName: domain.CommandNameReportGenerate,
			ConfigPath:  config.Path,
		},
		ValidatedModel:    domain.ValidatedModelSummary{ConfigPath: config.Path},
		ScenarioResults:   reportScenarioResults(),
		FinalRanking:      domain.AggregatedRankingResult{},
		ReportDefinitions: []domain.ReportDefinition{{Name: report.Name, Title: report.Title, Format: report.Format}},
		ScenarioWeights:   reportScenarioWeights(),
		Config:            config,
	}

	renderer := DefaultReportRenderer{}
	if _, err := renderer.RenderReports(context.Background(), input); err != nil {
		t.Fatalf("first RenderReports() error = %v", err)
	}

	target := filepath.Join(tempDir, "artifacts", "summary.json")
	if err := os.WriteFile(target, []byte("drift\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if _, err := renderer.RenderReports(context.Background(), input); err != nil {
		t.Fatalf("second RenderReports() error = %v", err)
	}

	content, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	want, err := renderReport(report, config.Config, reportScenarioResults(), config.Config.Aggregation, reportScenarioWeights())
	if err != nil {
		t.Fatalf("renderReport() error = %v", err)
	}
	if got := string(content); got != want {
		t.Fatalf("file content = %q, want %q", got, want)
	}
}

func TestDefaultReportRendererWritesFileTargetedCSVReports(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "model.cue")
	report := ReportConfig{
		Name:      "summary-csv",
		Title:     "Summary CSV",
		Format:    "csv",
		Filepath:  "artifacts/summary.csv",
		Arguments: []string{"columns=scenario,alternative,criterion,value,score,rank,excluded,exclusion_reason", "header=true"},
	}

	config := reportLoadedConfig(report)
	config.Path = configPath

	renderer := DefaultReportRenderer{}
	got, err := renderer.RenderReports(context.Background(), RenderReportsInput{
		Command: domain.CommandRequest{
			CommandName: domain.CommandNameReportGenerate,
			ConfigPath:  config.Path,
		},
		ValidatedModel:    domain.ValidatedModelSummary{ConfigPath: config.Path},
		ScenarioResults:   reportScenarioResults(),
		FinalRanking:      domain.AggregatedRankingResult{},
		ReportDefinitions: []domain.ReportDefinition{{Name: report.Name, Title: report.Title, Format: report.Format}},
		Config:            config,
	})
	if err != nil {
		t.Fatalf("RenderReports() error = %v", err)
	}
	if got.RenderedOutput != "" {
		t.Fatalf("RenderedOutput = %q, want empty stdout output", got.RenderedOutput)
	}

	content, err := os.ReadFile(filepath.Join(tempDir, "artifacts", "summary.csv"))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	assertCSVSchemaOutput(t, string(content))
}
