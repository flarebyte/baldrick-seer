package app

import (
	"path/filepath"
	"testing"
)

func TestRunValidate(t *testing.T) {
	t.Parallel()

	got, err := RunValidate(ValidateRequest{
		ConfigPath: filepath.Join("..", "..", "testdata", "config", "minimal.cue"),
	})
	if err != nil {
		t.Fatalf("RunValidate() error = %v", err)
	}

	if want := validateStubOutput; got.Stdout != want {
		t.Fatalf("Stdout = %q, want %q", got.Stdout, want)
	}
}

func TestRunReportGenerate(t *testing.T) {
	t.Parallel()

	got, err := RunReportGenerate(ReportGenerateRequest{
		ConfigPath: filepath.Join("..", "..", "testdata", "config", "minimal.cue"),
	})
	if err != nil {
		t.Fatalf("RunReportGenerate() error = %v", err)
	}

	if want := reportGenerateStubOutput; got.Stdout != want {
		t.Fatalf("Stdout = %q, want %q", got.Stdout, want)
	}
}
