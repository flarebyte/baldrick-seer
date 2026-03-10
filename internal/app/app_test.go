package app

import "testing"

func TestRunValidate(t *testing.T) {
	t.Parallel()

	got, err := RunValidate(ValidateRequest{})
	if err != nil {
		t.Fatalf("RunValidate() error = %v", err)
	}

	if want := validateStubOutput; got.Stdout != want {
		t.Fatalf("Stdout = %q, want %q", got.Stdout, want)
	}
}

func TestRunReportGenerate(t *testing.T) {
	t.Parallel()

	got, err := RunReportGenerate(ReportGenerateRequest{})
	if err != nil {
		t.Fatalf("RunReportGenerate() error = %v", err)
	}

	if want := reportGenerateStubOutput; got.Stdout != want {
		t.Fatalf("Stdout = %q, want %q", got.Stdout, want)
	}
}
