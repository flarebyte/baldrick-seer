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
	return &CommandFailure{
		Category:    category,
		ExitCode:    ExitCodeForCategory(category),
		Diagnostics: diagnostics,
		Err:         err,
	}
}

func (f *CommandFailure) Error() string {
	if len(f.Diagnostics) > 0 && f.Diagnostics[0].Message != "" {
		return f.Diagnostics[0].Message
	}
	if f.Err != nil {
		return f.Err.Error()
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
