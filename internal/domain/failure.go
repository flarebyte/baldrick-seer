package domain

import "errors"

type FailureCategory string

const (
	FailureCategoryInput      FailureCategory = "input"
	FailureCategoryValidation FailureCategory = "validation"
	FailureCategoryExecution  FailureCategory = "execution"
	FailureCategoryRendering  FailureCategory = "rendering"
	FailureCategoryInternal   FailureCategory = "internal"
)

type ExitCode int

const (
	ExitCodeSuccess      ExitCode = 0
	ExitCodeInvalidInput ExitCode = 1
	ExitCodeValidation   ExitCode = 1
	ExitCodeExecution    ExitCode = 1
	ExitCodeInternal     ExitCode = 1
)

type CommandFailure struct {
	Category    FailureCategory
	ExitCode    ExitCode
	Message     string
	Diagnostics []Diagnostic
	Err         error
}

func NewDiagnostic(
	severity DiagnosticSeverity,
	code string,
	path string,
	location DiagnosticLocation,
	message string,
) Diagnostic {
	return Diagnostic{
		Severity: severity,
		Code:     code,
		Path:     path,
		Location: location,
		Message:  message,
		Guidance: guidanceForDiagnosticCode(code),
	}
}

func NewFailure(category FailureCategory, diagnostics []Diagnostic, err error) *CommandFailure {
	return NewFailureWithMessage(category, "", diagnostics, err)
}

func NewFailureWithMessage(category FailureCategory, message string, diagnostics []Diagnostic, err error) *CommandFailure {
	if message == "" {
		message = failureMessageFromDiagnostics(diagnostics)
	}

	return &CommandFailure{
		Category:    category,
		ExitCode:    ExitCodeForCategory(category),
		Message:     message,
		Diagnostics: diagnostics,
		Err:         err,
	}
}

func (f *CommandFailure) Error() string {
	if f.Message != "" {
		return f.Message
	}
	return "command failed"
}

func (f *CommandFailure) Unwrap() error {
	return f.Err
}

func ExitCodeForCategory(category FailureCategory) ExitCode {
	switch category {
	case FailureCategoryInput:
		return ExitCodeInvalidInput
	case FailureCategoryValidation:
		return ExitCodeValidation
	case FailureCategoryExecution, FailureCategoryRendering:
		return ExitCodeExecution
	case FailureCategoryInternal:
		return ExitCodeInternal
	default:
		return ExitCodeInternal
	}
}

func ExitCodeForError(err error) ExitCode {
	var failure *CommandFailure
	if errors.As(err, &failure) {
		return failure.ExitCode
	}
	return ExitCodeInternal
}

func AsCommandFailure(err error) *CommandFailure {
	var failure *CommandFailure
	if errors.As(err, &failure) {
		return failure
	}
	return nil
}

type FailurePresentation struct {
	ExitCode ExitCode
	Stderr   string
}

func PresentError(err error) FailurePresentation {
	stderr := "status: error\nmessage: " + UserMessageForError(err) + "\n"
	if guidance := UserGuidanceForError(err); guidance != "" {
		stderr += "guidance: " + guidance + "\n"
	}
	return FailurePresentation{
		ExitCode: ExitCodeForError(err),
		Stderr:   stderr,
	}
}

func UserMessageForError(err error) string {
	if failure := AsCommandFailure(err); failure != nil {
		if failure.Message != "" {
			return failure.Message
		}
	}

	return "command failed"
}

func UserGuidanceForError(err error) string {
	if failure := AsCommandFailure(err); failure != nil {
		if failure.Category != FailureCategoryValidation {
			return ""
		}
		for _, diagnostic := range failure.Diagnostics {
			if diagnostic.Guidance != "" {
				return diagnostic.Guidance
			}
		}
	}
	return ""
}

func failureMessageFromDiagnostics(diagnostics []Diagnostic) string {
	if len(diagnostics) > 0 && diagnostics[0].Message != "" {
		return diagnostics[0].Message
	}
	return "command failed"
}

func guidanceForDiagnosticCode(code string) string {
	switch code {
	case "validation.section_missing":
		return "Add the missing top-level section under config."
	case "validation.duplicate_criterion_name", "validation.duplicate_alternative_name", "validation.duplicate_scenario_name", "validation.duplicate_report_name":
		return "Rename or remove the duplicate entry so each name is unique."
	case "validation.unknown_active_criterion", "validation.unknown_constraint_criterion", "validation.unknown_evaluation_scenario", "validation.unknown_evaluation_alternative", "validation.unknown_evaluation_criterion", "validation.unknown_report_focus_scenario", "validation.unknown_report_focus_alternative", "validation.unknown_report_focus_criterion", "validation.unknown_aggregation_scenario", "validation.unknown_pairwise_criterion":
		return "Define the referenced name or fix the reference to match an existing entry."
	case "validation.missing_pairwise_comparison":
		return "Add exactly one pairwise comparison for the missing criterion pair."
	case "validation.duplicate_pairwise_comparison", "validation.inverse_duplicate_pairwise_comparison":
		return "Keep only one canonical pairwise comparison for each unordered criterion pair."
	case "validation.pairwise_self_comparison":
		return "Compare two distinct active criteria in each pairwise comparison."
	case "validation.missing_evaluation_value":
		return "Add a value for every active criterion in the scenario evaluation."
	case "validation.evaluation_value_kind_mismatch", "validation.unsupported_evaluation_value_kind", "validation.invalid_number_value", "validation.invalid_ordinal_value", "validation.invalid_boolean_value":
		return "Use the value kind and concrete value required by the referenced criterion type."
	case "validation.ordinal_scale_guidance_missing":
		return "Add scaleGuidance to each ordinal criterion definition."
	case "validation.invalid_constraint_operator":
		return "Use an operator allowed by the constrained criterion value type."
	case "validation.invalid_constraint_number_value", "validation.invalid_constraint_ordinal_value", "validation.invalid_constraint_boolean_value", "validation.invalid_constraint_value":
		return "Use a constraint value compatible with the constrained criterion value type."
	case "validation.malformed_report_argument":
		return "Write each report argument as key=value."
	case "validation.unknown_report_argument":
		return "Use only supported report argument keys for the selected format."
	case "validation.duplicate_report_argument":
		return "Provide each report argument key at most once."
	case "validation.incompatible_report_argument":
		return "Use report arguments that are allowed for the selected format."
	case "validation.invalid_report_argument_value":
		return "Use a supported value for the selected report argument."
	default:
		return ""
	}
}
