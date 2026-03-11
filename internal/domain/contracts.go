// Package domain defines the internal execution contracts used by the CLI pipeline.
package domain

type CommandName string

const (
	CommandNameValidate       CommandName = "validate"
	CommandNameReportGenerate CommandName = "report generate"
)

type DiagnosticSeverity string

const (
	DiagnosticSeverityError   DiagnosticSeverity = "error"
	DiagnosticSeverityWarning DiagnosticSeverity = "warning"
)

// CommandRequest carries the minimal command input needed by the execution pipeline.
type CommandRequest struct {
	CommandName CommandName
	ConfigPath  string
}

// CommandResult carries ordered execution outputs for later presentation stages.
type CommandResult struct {
	CommandName       CommandName
	Diagnostics       []Diagnostic
	ValidatedModel    *ValidatedModelSummary
	ScenarioResults   []ScenarioRankingResult
	FinalRanking      *AggregatedRankingResult
	ReportDefinitions []ReportDefinition
}

type Diagnostic struct {
	Severity DiagnosticSeverity
	Code     string
	Path     string
	Location DiagnosticLocation
	Message  string
}

type DiagnosticLocation struct {
	Line   int
	Column int
}

type ReportDefinition struct {
	Name   string
	Title  string
	Format string
}

type ValidatedModelSummary struct {
	ConfigPath        string
	CriterionCount    int
	AlternativeCount  int
	ScenarioCount     int
	ReportDefinitions []ReportDefinition
}

// RankedAlternatives is ordered for deterministic rendering.
type RankedAlternative struct {
	Name     string
	Rank     int
	Score    float64
	Excluded bool
}

// ScenarioRankingResult carries ordered scenario-local rankings.
type ScenarioRankingResult struct {
	ScenarioName       string
	RankedAlternatives []RankedAlternative
}

// AggregatedRankingResult carries the final ordered ranking across scenarios.
type AggregatedRankingResult struct {
	RankedAlternatives []RankedAlternative
}
