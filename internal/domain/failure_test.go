package domain

import (
	"errors"
	"testing"
)

func TestNewDiagnostic(t *testing.T) {
	t.Parallel()

	got := NewDiagnostic(
		DiagnosticSeverityError,
		"config.not_found",
		"testdata/config/missing.cue",
		DiagnosticLocation{Line: 3, Column: 7},
		"config path does not exist",
	)

	if got.Code != "config.not_found" {
		t.Fatalf("Code = %q, want %q", got.Code, "config.not_found")
	}

	if got.Severity != DiagnosticSeverityError {
		t.Fatalf("Severity = %q, want %q", got.Severity, DiagnosticSeverityError)
	}

	if got.Path != "testdata/config/missing.cue" {
		t.Fatalf("Path = %q, want %q", got.Path, "testdata/config/missing.cue")
	}

	if got.Location.Line != 3 || got.Location.Column != 7 {
		t.Fatalf("Location = %#v, want line=3 column=7", got.Location)
	}

	if got.Guidance != "" {
		t.Fatalf("Guidance = %q, want empty", got.Guidance)
	}
}

func TestNewDiagnosticAddsValidationGuidance(t *testing.T) {
	t.Parallel()

	got := NewDiagnostic(
		DiagnosticSeverityError,
		"validation.missing_pairwise_comparison",
		"config",
		DiagnosticLocation{},
		"missing pairwise comparison for pair: reliability/speed",
	)

	if got.Guidance != "Add exactly one pairwise comparison for the missing criterion pair." {
		t.Fatalf("Guidance = %q, want pairwise guidance", got.Guidance)
	}
}

func TestExitCodeForCategory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		category FailureCategory
		want     ExitCode
	}{
		{category: FailureCategoryInput, want: ExitCodeInvalidInput},
		{category: FailureCategoryValidation, want: ExitCodeValidation},
		{category: FailureCategoryExecution, want: ExitCodeExecution},
		{category: FailureCategoryRendering, want: ExitCodeExecution},
		{category: FailureCategoryInternal, want: ExitCodeInternal},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(string(tt.category), func(t *testing.T) {
			t.Parallel()

			if got := ExitCodeForCategory(tt.category); got != tt.want {
				t.Fatalf("ExitCodeForCategory(%q) = %d, want %d", tt.category, got, tt.want)
			}
		})
	}
}

func TestFailureCarriesWrappedError(t *testing.T) {
	t.Parallel()

	cause := errors.New("boom")
	failure := NewFailure(
		FailureCategoryExecution,
		[]Diagnostic{
			NewDiagnostic(DiagnosticSeverityError, "execution.failed", "", DiagnosticLocation{}, "execution failed"),
		},
		cause,
	)

	if !errors.Is(failure, cause) {
		t.Fatalf("errors.Is(failure, cause) = false, want true")
	}

	if got, want := ExitCodeForError(failure), ExitCodeExecution; got != want {
		t.Fatalf("ExitCodeForError(failure) = %d, want %d", got, want)
	}

	if got, want := failure.Error(), "execution failed"; got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}

	if got, want := failure.Message, "execution failed"; got != want {
		t.Fatalf("Message = %q, want %q", got, want)
	}
}

func TestAsCommandFailureReturnsNilForUnknownError(t *testing.T) {
	t.Parallel()

	if got := AsCommandFailure(errors.New("plain")); got != nil {
		t.Fatalf("AsCommandFailure() = %#v, want nil", got)
	}
}

func TestPresentError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		err        error
		wantExit   ExitCode
		wantStderr string
	}{
		{
			name: "validation failure uses diagnostic message",
			err: NewFailure(
				FailureCategoryValidation,
				[]Diagnostic{
					NewDiagnostic(DiagnosticSeverityError, "validation.failed", "config", DiagnosticLocation{}, "validation failed"),
				},
				errors.New("boom"),
			),
			wantExit:   ExitCodeValidation,
			wantStderr: "status: error\nmessage: validation failed\n",
		},
		{
			name: "validation failure includes guidance when available",
			err: NewFailure(
				FailureCategoryValidation,
				[]Diagnostic{
					NewDiagnostic(DiagnosticSeverityError, "validation.unknown_evaluation_scenario", "config", DiagnosticLocation{}, "unknown scenario name in evaluations: missing"),
				},
				errors.New("boom"),
			),
			wantExit:   ExitCodeValidation,
			wantStderr: "status: error\nmessage: unknown scenario name in evaluations: missing\nguidance: Define the referenced name or fix the reference to match an existing entry.\n",
		},
		{
			name: "explicit message is preserved",
			err: NewFailureWithMessage(
				FailureCategoryExecution,
				"command canceled",
				nil,
				errors.New("boom"),
			),
			wantExit:   ExitCodeExecution,
			wantStderr: "status: error\nmessage: command canceled\n",
		},
		{
			name:       "plain error uses fallback message",
			err:        errors.New("plain"),
			wantExit:   ExitCodeInternal,
			wantStderr: "status: error\nmessage: command failed\n",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := PresentError(tt.err)
			if got.ExitCode != tt.wantExit {
				t.Fatalf("ExitCode = %d, want %d", got.ExitCode, tt.wantExit)
			}
			if got.Stderr != tt.wantStderr {
				t.Fatalf("Stderr = %q, want %q", got.Stderr, tt.wantStderr)
			}
		})
	}
}
