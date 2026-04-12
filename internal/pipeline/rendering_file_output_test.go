package pipeline

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func tempConfigPath(tempDir string) string {
	return filepath.Join(tempDir, "model.cue")
}

func readArtifactFile(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	return string(content)
}

func TestDefaultReportRendererWritesFileTargetedReports(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	report := ReportConfig{
		Name:      "summary-json",
		Title:     "Summary JSON",
		Format:    "json",
		Filepath:  "artifacts/summary.json",
		Arguments: []string{"include-weights=true", "pretty=true"},
	}

	config := reportLoadedConfig(report)
	config.Path = tempConfigPath(tempDir)

	got := renderReportsForTest(t, config, reportScenarioResults(), domain.AggregatedRankingResult{}, singleReportDefinitions(report), reportScenarioWeights())
	if got.RenderedOutput != "" {
		t.Fatalf("RenderedOutput = %q, want empty stdout output", got.RenderedOutput)
	}

	if got, want := readArtifactFile(t, filepath.Join(tempDir, "artifacts", "summary.json")), expectedRenderedReport(t, report, config, reportScenarioResults(), reportScenarioWeights()); got != want {
		t.Fatalf("file content = %q, want %q", got, want)
	}
}

func TestDefaultReportRendererSplitsMixedStdoutAndFileOutputs(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
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
	config.Path = tempConfigPath(tempDir)

	got := renderReportsForTest(t, config, reportScenarioResults(), domain.AggregatedRankingResult{}, []domain.ReportDefinition{
		{Name: reportStdout.Name, Title: reportStdout.Title, Format: reportStdout.Format},
		{Name: reportFile.Name, Title: reportFile.Title, Format: reportFile.Format},
	}, reportScenarioWeights())

	if got, want := got.RenderedOutput, expectedRenderedReport(t, reportStdout, config, reportScenarioResults(), reportScenarioWeights()); got != want {
		t.Fatalf("RenderedOutput = %q, want %q", got, want)
	}

	if got, want := readArtifactFile(t, filepath.Join(tempDir, "artifacts", "a.json")), expectedRenderedReport(t, reportFile, config, reportScenarioResults(), reportScenarioWeights()); got != want {
		t.Fatalf("file content = %q, want %q", got, want)
	}
}

func TestDefaultReportRendererResolvesFilepathsRelativeToConfigDirectory(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	report := ReportConfig{
		Name:      "summary-markdown",
		Title:     "Summary Markdown",
		Format:    "markdown",
		Filepath:  "../artifacts/summary.md",
		Arguments: []string{"include-scores=true"},
	}

	config := reportLoadedConfig(report)
	config.Path = filepath.Join(tempDir, "config", "model.cue")

	_ = renderReportsForTest(t, config, reportScenarioResults(), domain.AggregatedRankingResult{}, singleReportDefinitions(report), reportScenarioWeights())

	if got, want := readArtifactFile(t, filepath.Join(tempDir, "artifacts", "summary.md")), expectedRenderedReport(t, report, config, reportScenarioResults(), reportScenarioWeights()); got != want {
		t.Fatalf("file content = %q, want %q", got, want)
	}
}

func TestDefaultReportRendererOverwritesFileTargetedReportsDeterministically(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	report := ReportConfig{
		Name:      "summary-json",
		Title:     "Summary JSON",
		Format:    "json",
		Filepath:  "artifacts/summary.json",
		Arguments: []string{"pretty=true"},
	}

	config := reportLoadedConfig(report)
	config.Path = tempConfigPath(tempDir)

	renderReportsForTest(t, config, reportScenarioResults(), domain.AggregatedRankingResult{}, singleReportDefinitions(report), reportScenarioWeights())

	target := filepath.Join(tempDir, "artifacts", "summary.json")
	if err := os.WriteFile(target, []byte("drift\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	renderReportsForTest(t, config, reportScenarioResults(), domain.AggregatedRankingResult{}, singleReportDefinitions(report), reportScenarioWeights())

	if got, want := readArtifactFile(t, target), expectedRenderedReport(t, report, config, reportScenarioResults(), reportScenarioWeights()); got != want {
		t.Fatalf("file content = %q, want %q", got, want)
	}
}

func TestDefaultReportRendererWritesFileTargetedCSVReports(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	report := ReportConfig{
		Name:      "summary-csv",
		Title:     "Summary CSV",
		Format:    "csv",
		Filepath:  "artifacts/summary.csv",
		Arguments: []string{"columns=scenario,alternative,criterion,value,score,rank,excluded,exclusion_reason", "header=true"},
	}

	config := reportLoadedConfig(report)
	config.Path = tempConfigPath(tempDir)

	got := renderReportsForTest(t, config, reportScenarioResults(), domain.AggregatedRankingResult{}, singleReportDefinitions(report), nil)
	if got.RenderedOutput != "" {
		t.Fatalf("RenderedOutput = %q, want empty stdout output", got.RenderedOutput)
	}

	assertCSVSchemaOutput(t, readArtifactFile(t, filepath.Join(tempDir, "artifacts", "summary.csv")))
}
