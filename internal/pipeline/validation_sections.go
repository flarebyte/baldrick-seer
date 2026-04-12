package pipeline

import (
	"fmt"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

func validateRequiredSections(config LoadedConfig) []domain.Diagnostic {
	requiredSections := []struct {
		name    string
		present bool
		valid   bool
	}{
		{name: "problem", present: hasName(config.ConfigFields, "problem"), valid: config.Config != nil && config.Config.Problem != nil},
		{name: "reports", present: hasName(config.ConfigFields, "reports"), valid: config.Config != nil && len(config.Config.Reports) > 0},
		{name: "criteriaCatalog", present: hasName(config.ConfigFields, "criteriaCatalog"), valid: config.Config != nil && len(config.Config.CriteriaCatalog) > 0},
		{name: "alternatives", present: hasName(config.ConfigFields, "alternatives"), valid: config.Config != nil && len(config.Config.Alternatives) > 0},
		{name: "scenarios", present: hasName(config.ConfigFields, "scenarios"), valid: config.Config != nil && len(config.Config.Scenarios) > 0},
		{name: "evaluations", present: hasName(config.ConfigFields, "evaluations"), valid: config.Config != nil && len(config.Config.Evaluations) > 0},
		{name: "aggregation", present: hasName(config.ConfigFields, "aggregation"), valid: config.Config != nil && config.Config.Aggregation != nil},
	}

	diagnostics := make([]domain.Diagnostic, 0, len(requiredSections))
	for _, section := range requiredSections {
		if section.present && section.valid {
			continue
		}

		diagnostics = append(diagnostics, domain.NewDiagnostic(
			domain.DiagnosticSeverityError,
			"validation.section_missing",
			fmt.Sprintf("config.%s", section.name),
			domain.DiagnosticLocation{},
			fmt.Sprintf("missing required section: %s", section.name),
		))
	}

	return diagnostics
}

func collectUniqueNames[T any](
	diagnostics *[]domain.Diagnostic,
	pathPrefix string,
	items []T,
	nameOf func(T) string,
	code string,
	messageFormat string,
) []string {
	seen := make(map[string]struct{}, len(items))
	var names []string

	for index, item := range items {
		name := nameOf(item)
		names = append(names, name)
		if _, exists := seen[name]; exists {
			*diagnostics = append(*diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				code,
				fmt.Sprintf("%s[%d].name", pathPrefix, index),
				domain.DiagnosticLocation{},
				fmt.Sprintf(messageFormat, name),
			))
			continue
		}

		seen[name] = struct{}{}
	}

	return domain.CanonicalNames(names)
}

func hasName(names []string, target string) bool {
	for _, name := range names {
		if name == target {
			return true
		}
	}

	return false
}

func isSupportedCriterionValueType(valueType string) bool {
	return valueType == "number" || valueType == "ordinal" || valueType == "boolean"
}

func isSupportedReportFormat(format string) bool {
	return format == "markdown" || format == "json" || format == "csv"
}

func countScenarioNames(scenarios []ScenarioConfig) map[string]int {
	counts := make(map[string]int, len(scenarios))
	for _, scenario := range scenarios {
		counts[scenario.Name]++
	}
	return counts
}
