package pipeline

import "github.com/flarebyte/baldrick-seer/internal/domain"

func NewInputFailure(code string, path string, message string, err error) error {
	return newFailure(domain.FailureCategoryInput, code, path, message, err)
}

func NewValidationFailure(code string, path string, message string, err error) error {
	return newFailure(domain.FailureCategoryValidation, code, path, message, err)
}

func NewValidationDiagnosticsFailure(diagnostics []domain.Diagnostic, err error) error {
	return domain.NewFailureWithMessage(
		domain.FailureCategoryValidation,
		"",
		domain.CanonicalDiagnostics(diagnostics),
		err,
	)
}

func NewExecutionFailure(code string, path string, message string, err error) error {
	return newFailure(domain.FailureCategoryExecution, code, path, message, err)
}

func NewRenderingFailure(code string, path string, message string, err error) error {
	return newFailure(domain.FailureCategoryRendering, code, path, message, err)
}

func WrapStageFailure(category domain.FailureCategory, code string, path string, message string, err error) error {
	if domain.AsCommandFailure(err) != nil {
		return err
	}

	switch category {
	case domain.FailureCategoryInput:
		return NewInputFailure(code, path, message, err)
	case domain.FailureCategoryValidation:
		return NewValidationFailure(code, path, message, err)
	case domain.FailureCategoryRendering:
		return NewRenderingFailure(code, path, message, err)
	case domain.FailureCategoryExecution:
		return NewExecutionFailure(code, path, message, err)
	default:
		return domain.NewFailureWithMessage(
			domain.FailureCategoryInternal,
			message,
			[]domain.Diagnostic{
				domain.NewDiagnostic(
					domain.DiagnosticSeverityError,
					code,
					path,
					domain.DiagnosticLocation{},
					message,
				),
			},
			err,
		)
	}
}

func newFailure(category domain.FailureCategory, code string, path string, message string, err error) error {
	return domain.NewFailureWithMessage(
		category,
		message,
		[]domain.Diagnostic{
			domain.NewDiagnostic(
				domain.DiagnosticSeverityError,
				code,
				path,
				domain.DiagnosticLocation{},
				message,
			),
		},
		err,
	)
}
