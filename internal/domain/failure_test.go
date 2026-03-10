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
}
