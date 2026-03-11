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
	return FailurePresentation{
		ExitCode: ExitCodeForError(err),
		Stderr:   "status: error\nmessage: " + UserMessageForError(err) + "\n",
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

func failureMessageFromDiagnostics(diagnostics []Diagnostic) string {
	if len(diagnostics) > 0 && diagnostics[0].Message != "" {
		return diagnostics[0].Message
	}
	return "command failed"
}
