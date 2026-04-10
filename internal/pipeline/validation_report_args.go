package pipeline

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/flarebyte/baldrick-seer/internal/domain"
)

type reportArgumentRule struct {
	AllowedFormats []string
	ValidateValue  func(string) bool
}

var reportArgumentRules = map[string]reportArgumentRule{
	"include-scenarios": {AllowedFormats: []string{"markdown"}, ValidateValue: func(value string) bool {
		return value == "all" || value == "focused"
	}},
	"top-alternatives": {AllowedFormats: []string{"markdown"}, ValidateValue: isPositiveIntegerString},
	"include-scores":   {AllowedFormats: []string{"markdown"}, ValidateValue: isBooleanString},
	"explain":          {AllowedFormats: []string{"markdown"}, ValidateValue: isBooleanString},
	"detail": {AllowedFormats: []string{"markdown"}, ValidateValue: func(value string) bool {
		return value == "brief" || value == "standard" || value == "full"
	}},
	"include-evidence": {AllowedFormats: []string{"json"}, ValidateValue: isBooleanString},
	"include-weights":  {AllowedFormats: []string{"json"}, ValidateValue: isBooleanString},
	"pretty":           {AllowedFormats: []string{"json"}, ValidateValue: isBooleanString},
	"columns":          {AllowedFormats: []string{"csv"}, ValidateValue: isValidCSVColumns},
	"header":           {AllowedFormats: []string{"csv"}, ValidateValue: isBooleanString},
}

func validateReportArguments(reportIndex int, report ReportConfig) []domain.Diagnostic {
	if !isSupportedReportFormat(report.Format) {
		return nil
	}

	var diagnostics []domain.Diagnostic
	seenKeys := make(map[string]struct{}, len(report.Arguments))

	for argumentIndex, argument := range report.Arguments {
		path := fmt.Sprintf("config.reports[%d].arguments[%d]", reportIndex, argumentIndex)
		key, value, ok := strings.Cut(argument, "=")
		if !ok || key == "" {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.malformed_report_argument",
				path,
				domain.DiagnosticLocation{},
				fmt.Sprintf("report argument must use key=value form: %s", argument),
			))
			continue
		}

		rule, exists := reportArgumentRules[key]
		if !exists {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.unknown_report_argument",
				path,
				domain.DiagnosticLocation{},
				fmt.Sprintf("unknown report argument key: %s", key),
			))
			continue
		}

		if _, exists := seenKeys[key]; exists {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.duplicate_report_argument",
				path,
				domain.DiagnosticLocation{},
				fmt.Sprintf("duplicate report argument key: %s", key),
			))
			continue
		}
		seenKeys[key] = struct{}{}

		if !reportArgumentAllowedForFormat(rule, report.Format) {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.incompatible_report_argument",
				path,
				domain.DiagnosticLocation{},
				fmt.Sprintf("report argument key %s is not allowed for format %s", key, report.Format),
			))
			continue
		}

		if !rule.ValidateValue(value) {
			diagnostics = append(diagnostics, domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				"validation.invalid_report_argument_value",
				path,
				domain.DiagnosticLocation{},
				fmt.Sprintf("invalid value for report argument %s: %s", key, value),
			))
		}
	}

	return diagnostics
}

func reportArgumentAllowedForFormat(rule reportArgumentRule, format string) bool {
	for _, allowedFormat := range rule.AllowedFormats {
		if allowedFormat == format {
			return true
		}
	}
	return false
}

func isBooleanString(value string) bool {
	return value == "true" || value == "false"
}

func isPositiveIntegerString(value string) bool {
	number, err := strconv.Atoi(value)
	return err == nil && number > 0
}

func isValidCSVColumns(value string) bool {
	if value == "" {
		return false
	}

	allowedColumns := map[string]struct{}{
		"scenario":         {},
		"alternative":      {},
		"criterion":        {},
		"value":            {},
		"score":            {},
		"rank":             {},
		"excluded":         {},
		"exclusion_reason": {},
	}

	seen := map[string]struct{}{}
	for _, column := range strings.Split(value, ",") {
		if _, exists := allowedColumns[column]; !exists {
			return false
		}
		if _, exists := seen[column]; exists {
			return false
		}
		seen[column] = struct{}{}
	}

	return true
}
