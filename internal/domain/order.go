package domain

import "sort"

// CanonicalDiagnostics returns diagnostics ordered by severity, code, path,
// location, and message so repeated runs stay stable.
func CanonicalDiagnostics(input []Diagnostic) []Diagnostic {
	if len(input) == 0 {
		return nil
	}

	output := append([]Diagnostic(nil), input...)
	sort.Slice(output, func(i int, j int) bool {
		left := output[i]
		right := output[j]

		leftSeverity := diagnosticSeverityOrder(left.Severity)
		rightSeverity := diagnosticSeverityOrder(right.Severity)
		if leftSeverity != rightSeverity {
			return leftSeverity < rightSeverity
		}
		if left.Code != right.Code {
			return left.Code < right.Code
		}
		if left.Path != right.Path {
			return left.Path < right.Path
		}
		if left.Location.Line != right.Location.Line {
			return left.Location.Line < right.Location.Line
		}
		if left.Location.Column != right.Location.Column {
			return left.Location.Column < right.Location.Column
		}
		return left.Message < right.Message
	})

	return output
}

// CanonicalNames returns names sorted lexicographically.
func CanonicalNames(input []string) []string {
	if len(input) == 0 {
		return nil
	}

	output := append([]string(nil), input...)
	sort.Strings(output)
	return output
}

// CanonicalReportDefinitions returns report definitions ordered by name, title,
// and format.
func CanonicalReportDefinitions(input []ReportDefinition) []ReportDefinition {
	if len(input) == 0 {
		return nil
	}

	output := append([]ReportDefinition(nil), input...)
	sort.Slice(output, func(i int, j int) bool {
		left := output[i]
		right := output[j]

		if left.Name != right.Name {
			return left.Name < right.Name
		}
		if left.Title != right.Title {
			return left.Title < right.Title
		}
		return left.Format < right.Format
	})

	return output
}

// CanonicalRankedAlternatives returns ranking rows ordered by rank ascending,
// then name ascending, then score ascending as a final stable tie-breaker.
func CanonicalRankedAlternatives(input []RankedAlternative) []RankedAlternative {
	if len(input) == 0 {
		return nil
	}

	output := append([]RankedAlternative(nil), input...)
	sort.Slice(output, func(i int, j int) bool {
		left := output[i]
		right := output[j]

		if left.Rank != right.Rank {
			return left.Rank < right.Rank
		}
		if left.Name != right.Name {
			return left.Name < right.Name
		}
		return left.Score < right.Score
	})

	return output
}

// CanonicalScenarioResults returns scenarios ordered by scenario name, with
// each scenario-local ranking ordered by rank and name.
func CanonicalScenarioResults(input []ScenarioRankingResult) []ScenarioRankingResult {
	if len(input) == 0 {
		return nil
	}

	output := append([]ScenarioRankingResult(nil), input...)
	for i := range output {
		output[i].RankedAlternatives = CanonicalRankedAlternatives(output[i].RankedAlternatives)
	}

	sort.Slice(output, func(i int, j int) bool {
		return output[i].ScenarioName < output[j].ScenarioName
	})

	return output
}

func CanonicalValidatedModelSummary(input ValidatedModelSummary) ValidatedModelSummary {
	output := input
	output.ReportDefinitions = CanonicalReportDefinitions(input.ReportDefinitions)
	return output
}

func CanonicalAggregatedRankingResult(input AggregatedRankingResult) AggregatedRankingResult {
	output := input
	output.RankedAlternatives = CanonicalRankedAlternatives(input.RankedAlternatives)
	return output
}

func CanonicalCommandResult(input CommandResult) CommandResult {
	output := input
	output.Diagnostics = CanonicalDiagnostics(input.Diagnostics)
	output.ReportDefinitions = CanonicalReportDefinitions(input.ReportDefinitions)
	output.ScenarioResults = CanonicalScenarioResults(input.ScenarioResults)
	if input.ValidatedModel != nil {
		validatedModel := CanonicalValidatedModelSummary(*input.ValidatedModel)
		output.ValidatedModel = &validatedModel
	}
	if input.FinalRanking != nil {
		finalRanking := CanonicalAggregatedRankingResult(*input.FinalRanking)
		output.FinalRanking = &finalRanking
	}
	return output
}

func diagnosticSeverityOrder(severity DiagnosticSeverity) int {
	switch severity {
	case DiagnosticSeverityError:
		return 0
	case DiagnosticSeverityWarning:
		return 1
	default:
		return 2
	}
}
